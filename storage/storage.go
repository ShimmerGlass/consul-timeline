package storage

import (
	"context"
	"time"

	"github.com/aestek/consul-timeline/timeline"
)

type Query struct {
	Start   time.Time
	Service string
	Limit   int
}

type Storage interface {
	Store(evt tl.Event) error
	Query(ctx context.Context, q Query) ([]tl.Event, error)
}
