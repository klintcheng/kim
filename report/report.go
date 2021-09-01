package report

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"sort"
	"strings"
	"time"
)

type Result struct {
	StatusCode    int
	Err           error
	Duration      time.Duration
	ContentLength int64
}

const mapcap = 100000

type Report struct {
	fastest  float64
	slowest  float64
	average  float64
	rps      float64
	avgTotal float64

	lats        []float64
	errorDist   map[string]int
	statusCodes []int
	results     chan *Result
	SizeTotal   int64
	numRes      int
	total       time.Duration
	sizeTotal   int64
	w           io.Writer
}

func New(w io.Writer, n int) *Report {
	if n < 500 {
		n = 500
	}
	cap := min(n, mapcap)

	r := &Report{
		lats:        make([]float64, 0, cap),
		w:           w,
		results:     make(chan *Result, 10),
		errorDist:   make(map[string]int),
		statusCodes: make([]int, 0, cap),
	}
	go r.start()

	return r
}

func (r *Report) Add(res *Result) {
	r.results <- res
}

func (r *Report) Finalize(total time.Duration) {
	close(r.results)

	r.total = total
	r.rps = float64(r.numRes) / r.total.Seconds()
	r.average = r.avgTotal / float64(len(r.lats))
	r.print()
}

func (r *Report) start() {
	// Loop will continue until channel is closed
	for res := range r.results {
		r.numRes++
		if res.Err != nil {
			r.errorDist[res.Err.Error()]++
		} else {
			r.avgTotal += res.Duration.Seconds()
			if len(r.lats) < mapcap {
				r.lats = append(r.lats, res.Duration.Seconds())
				r.statusCodes = append(r.statusCodes, res.StatusCode)
			}
			if res.ContentLength > 0 {
				r.sizeTotal += res.ContentLength
			}
		}
	}
}

func (r *Report) histogram() []Bucket {
	bc := 4
	buckets := make([]float64, bc+1)
	counts := make([]int, bc+1)
	bs := (r.slowest - r.fastest) / float64(bc)
	for i := 0; i < bc; i++ {
		buckets[i] = r.fastest + bs*float64(i)
	}
	buckets[bc] = r.slowest
	var bi int
	var max int
	for i := 0; i < len(r.lats); {
		if r.lats[i] <= buckets[bi] {
			i++
			counts[bi]++
			if max < counts[bi] {
				max = counts[bi]
			}
		} else if bi < len(buckets)-1 {
			bi++
		}
	}
	res := make([]Bucket, len(buckets))
	for i := 0; i < len(buckets); i++ {
		res[i] = Bucket{
			Mark:      buckets[i],
			Count:     counts[i],
			Frequency: float64(counts[i]) / float64(len(r.lats)),
		}
	}
	return res
}

func (r *Report) snapshot() Snapshot {
	snapshot := Snapshot{
		Average:     r.average,
		Rps:         r.rps,
		Total:       r.total,
		SizeTotal:   r.sizeTotal,
		ErrorDist:   r.errorDist,
		Lats:        make([]float64, len(r.lats)),
		StatusCodes: make([]int, len(r.lats)),
	}

	if len(r.lats) == 0 {
		return snapshot
	}

	copy(snapshot.Lats, r.lats)
	copy(snapshot.StatusCodes, r.statusCodes)

	sort.Float64s(r.lats)
	r.fastest = r.lats[0]
	r.slowest = r.lats[len(r.lats)-1]

	snapshot.Histogram = r.histogram()
	snapshot.LatencyDistribution = r.latencies()

	snapshot.Fastest = r.fastest
	snapshot.Slowest = r.slowest

	statusCodeDist := make(map[int]int, len(snapshot.StatusCodes))
	for _, statusCode := range snapshot.StatusCodes {
		statusCodeDist[statusCode]++
	}
	snapshot.StatusCodeDist = statusCodeDist

	return snapshot
}

func (r *Report) latencies() []LatencyDistribution {
	pctls := []int{10, 50, 75, 90, 99}
	data := make([]float64, len(pctls))
	j := 0
	for i := 0; i < len(r.lats) && j < len(pctls); i++ {
		current := i * 100 / len(r.lats)
		if current >= pctls[j] {
			data[j] = r.lats[i]
			j++
		}
	}
	res := make([]LatencyDistribution, len(pctls))
	for i := 0; i < len(pctls); i++ {
		if data[i] > 0 {
			res[i] = LatencyDistribution{Percentage: pctls[i], Latency: data[i]}
		}
	}
	return res
}

func (r *Report) print() {
	buf := &bytes.Buffer{}
	if err := newTemplate().Execute(buf, r.snapshot()); err != nil {
		log.Println("error:", err.Error())
		return
	}
	r.printf(buf.String())

	r.printf("\n")
}

func (r *Report) printf(s string, v ...interface{}) {
	fmt.Fprintf(r.w, s, v...)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type Bucket struct {
	Mark      float64
	Count     int
	Frequency float64
}

func histogram(buckets []Bucket) string {
	max := 0
	for _, b := range buckets {
		if v := b.Count; v > max {
			max = v
		}
	}
	res := new(bytes.Buffer)
	for i := 0; i < len(buckets); i++ {
		// Normalize bar lengths.
		var barLen int
		if max > 0 {
			barLen = (buckets[i].Count*40 + max/2) / max
		}
		res.WriteString(fmt.Sprintf("  %4.3f [%v]\t|%v\n", buckets[i].Mark, buckets[i].Count, strings.Repeat(barChar, barLen)))
	}
	return res.String()
}

type Snapshot struct {
	AvgTotal float64
	Fastest  float64
	Slowest  float64
	Average  float64
	Rps      float64

	AvgDelay float64
	DelayMax float64
	DelayMin float64

	Lats        []float64
	StatusCodes []int

	Total time.Duration

	ErrorDist      map[string]int
	StatusCodeDist map[int]int
	SizeTotal      int64
	SizeReq        int64
	NumRes         int64

	LatencyDistribution []LatencyDistribution
	Histogram           []Bucket
}

type LatencyDistribution struct {
	Percentage int
	Latency    float64
}
