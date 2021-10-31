package kim

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var channelTotalGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
	Namespace: "kim",
	Name:      "channel_total",
	Help:      "The total number of channel logined",
}, []string{"serviceId", "serviceName"})
