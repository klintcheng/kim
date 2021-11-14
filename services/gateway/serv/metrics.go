package serv

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var messageInTotal = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: "kim",
	Name:      "message_in_total",
	Help:      "网关接收消息总数",
}, []string{"serviceId", "serviceName", "command"})

var messageInFlowBytes = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: "kim",
	Name:      "message_in_flow_bytes",
	Help:      "网关接收消息字节数",
}, []string{"serviceId", "serviceName", "command"})

var noServerFoundErrorTotal = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: "kim",
	Name:      "no_server_found_error_total",
	Help:      "查找zone分区中服务失败的次数",
}, []string{"zone"})
