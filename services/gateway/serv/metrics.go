package serv

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var messageInTotal = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: "kim",
	Name:      "message_in_total",
	Help:      "The total number of message received",
}, []string{"serviceId", "serviceName", "command"})

var messageInFlowBytes = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: "kim",
	Name:      "message_in_flow_bytes",
	Help:      "The gateway receives message bytes",
}, []string{"serviceId", "serviceName", "command"})

var noServerFoundErrorTotal = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: "kim",
	Name:      "no_server_found_error_total",
	Help:      "the total number of query failures",
}, []string{"zone"})
