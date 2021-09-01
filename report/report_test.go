package report

import (
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestReport(t *testing.T) {
	r := New(os.Stdout, 100)
	t1 := time.Now()
	defer func() {
		r.Finalize(time.Since(t1))
	}()

	for i := 0; i < 500; i++ {
		r.Add(&Result{
			StatusCode: 200,
			Duration:   time.Millisecond * time.Duration(1+rand.Intn(20)*100),
		})
	}
	for i := 0; i < 10; i++ {
		r.Add(&Result{
			StatusCode: 100,
		})
	}
}
