package storage

import (
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
	Query(q Query) ([]tl.Event, error)
}
