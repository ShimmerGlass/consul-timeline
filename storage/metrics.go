package storage

import (
	"context"
	"time"

	tl "github.com/aestek/consul-timeline/timeline"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	isLeaderGauge = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "consul_timeline_storage_is_leader",
		Help: "Is leader",
	})

	writeHistogram = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "consul_timeline_storage_write_seconds",
		Help:    "consul_timeline_storage_write_seconds.",
		Buckets: prometheus.ExponentialBuckets(0.0001, 4, 10),
	})

	readHistogram = promauto.NewHistogram(prometheus.HistogramOpts{
		Name:    "consul_timeline_storage_read_seconds",
		Help:    "consul_timeline_storage_read_seconds.",
		Buckets: prometheus.ExponentialBuckets(0.0001, 4, 10),
	})
)

type Metrics struct {
	inner Storage
}

func NewMetrics(inner Storage) *Metrics {
	return &Metrics{inner}
}

func (s *Metrics) Store(evt tl.Event) error {
	start := time.Now()
	defer func() {
		writeHistogram.Observe(time.Since(start).Seconds())
	}()
	return s.inner.Store(evt)
}

func (s *Metrics) Query(ctx context.Context, q Query) ([]tl.Event, error) {
	start := time.Now()
	defer func() {
		readHistogram.Observe(time.Since(start).Seconds())
	}()
	return s.inner.Query(ctx, q)
}
