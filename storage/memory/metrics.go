package memory

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	sizeGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "consul_timeline_storage_memory_size",
		Help: "Storage memory size",
	})
)
