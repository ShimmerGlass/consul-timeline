package memory

import (
	"context"

	"github.com/aestek/consul-timeline/storage"
	tl "github.com/aestek/consul-timeline/timeline"
)

type Storage struct {
	cfg Config

	events []tl.Event
	next   int
	size   int
}

func New(cfg Config) *Storage {
	return &Storage{events: make([]tl.Event, cfg.MaxSize)}
}

func (s *Storage) Store(evt tl.Event) error {
	s.events[s.next] = evt
	s.next = (s.next + 1) % len(s.events)
	if s.size < len(s.events) {
		s.size++
		sizeGauge.Set(float64(s.size))
	}

	return nil
}

func (s *Storage) Query(_ context.Context, q storage.Query) ([]tl.Event, error) {
	if len(s.events) == 0 {
		return nil, nil
	}
	res := []tl.Event{}

	i := s.next
	for j := 0; j < s.size; j++ {
		i--
		if i < 0 {
			i = len(s.events) - 1
		}

		evt := s.events[i]

		if evt.Time.After(q.Start) {
			continue
		}
		if q.Filter != "" && evt.ServiceName != q.Filter && evt.NodeName != q.Filter {
			continue
		}

		res = append(res, evt)

		if len(res) >= q.Limit {
			break
		}
	}

	return res, nil
}
