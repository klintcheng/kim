package report

import (
	"encoding/json"
	"fmt"
	"text/template"
)

var (
	defaultTmpl = `
Summary:
  Total:	{{ formatNumber .Total.Seconds }} secs
  Slowest:	{{ formatNumber .Slowest }} secs
  Fastest:	{{ formatNumber .Fastest }} secs
  Average:	{{ formatNumber .Average }} secs
  Requests/sec:	{{ formatNumber .Rps }}
  {{ if gt .SizeTotal 0 }}
  Total data:	{{ .SizeTotal }} bytes{{ end }}
Response time histogram:
{{ histogram .Histogram }}
Latency distribution:{{ range .LatencyDistribution }}
  {{ .Percentage }}%% in {{ formatNumber .Latency }} secs{{ end }}

Status code distribution:{{ range $code, $num := .StatusCodeDist }}
  [{{ $code }}]	{{ $num }} responses{{ end }}
{{ if gt (len .ErrorDist) 0 }}Error distribution:{{ range $err, $num := .ErrorDist }}
  [{{ $num }}]	{{ $err }}{{ end }}{{ end }}
`
)

const (
	barChar = "■"
)

func newTemplate() *template.Template {
	return template.Must(template.New("tmpl").Funcs(tmplFuncMap).Parse(defaultTmpl))
}

var tmplFuncMap = template.FuncMap{
	"formatNumber":    formatNumber,
	"formatNumberInt": formatNumberInt,
	"histogram":       histogram,
	"jsonify":         jsonify,
}

func jsonify(v interface{}) string {
	d, _ := json.Marshal(v)
	return string(d)
}

func formatNumber(duration float64) string {
	return fmt.Sprintf("%4.4f", duration)
}

func formatNumberInt(duration int) string {
	return fmt.Sprintf("%d", duration)
}
