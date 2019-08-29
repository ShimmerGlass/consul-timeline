package watch

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	eventsCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "consul_timeline_events_total",
		Help: "Total number of events sent",
	})
)
