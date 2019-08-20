package cass

import (
	"context"
	"time"

	"github.com/aestek/consul-timeline/storage"
	tl "github.com/aestek/consul-timeline/timeline"
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
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

func New(cfg Config) (*Storage, error) {
	cluster := gocql.NewCluster(cfg.Addresses...)
	cluster.Keyspace = cfg.Keyspace
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	s := &Storage{}

	gi := &idxGlobal{session}
	s.Global = gi
	s.indexes = append(s.indexes, gi)

	gs := &idxService{session}
	s.Services = gs
	s.indexes = append(s.indexes, gs)

	return s, nil
}

func (s *Storage) Store(evt tl.Event) error {
	for _, i := range s.indexes {
		if err := i.Store(evt); err != nil {
			return errors.Wrapf(err, "cass event insert (%T)", i)
		}
	}
	return nil
}

func (s *Storage) Query(_ context.Context, q storage.Query) ([]tl.Event, error) {
	if q.Filter != "" {
		return s.Services.FetchBefore(q.Filter, q.Start, q.Limit)
	}

	return s.Global.FetchBefore(q.Start, q.Limit)
}
