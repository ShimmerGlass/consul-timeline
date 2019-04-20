package cass

import (
	"time"

	"github.com/aestek/consul-timeline/storage"
	"github.com/aestek/consul-timeline/timeline"
	"github.com/gocql/gocql"
)

var _ storage.Storage = (*Storage)(nil)

const (
	timeBlockResolution = 5 * time.Minute
)

type cassEvt struct {
	tl.Event
	TimeBlock time.Time
}

type index interface {
	Store(tl.Event) error
}

type Storage struct {
	indexes  []index
	Services *idxService
	Global   *idxGlobal
}

func New(session *gocql.Session) *Storage {
	s := &Storage{}

	gi := &idxGlobal{session}
	s.Global = gi
	s.indexes = append(s.indexes, gi)

	gs := &idxService{session}
	s.Services = gs
	s.indexes = append(s.indexes, gs)

	return s
}

func (s *Storage) Store(evt tl.Event) error {
	for _, i := range s.indexes {
		if err := i.Store(evt); err != nil {
			return err
		}
	}
	return nil
}

func (s *Storage) Query(q storage.Query) ([]tl.Event, error) {
	if q.Service != "" {
		return s.Services.FetchBefore(q.Service, q.Start, q.Limit)
	}

	return s.Global.FetchBefore(q.Start, q.Limit)
}
