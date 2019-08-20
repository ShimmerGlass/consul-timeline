package storage

import (
	"context"
	"time"

	tl "github.com/aestek/consul-timeline/timeline"
)

type Query struct {
	Start  time.Time
	Filter string
	Limit  int
}

type Storage interface {
	Store(evt tl.Event) error
	Query(ctx context.Context, q Query) ([]tl.Event, error)
}
