package kim

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var channelTotalGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "kim",
	Name:      "channel_total",
	Help:      "网关并发数",
}, []string{"serviceId", "serviceName"})
