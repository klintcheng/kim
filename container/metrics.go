package container

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var messageOutFlowBytes = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: "kim",
	Name:      "message_out_flow_bytes",
	Help:      "网关下发的消息字节数",
}, []string{"command"})
