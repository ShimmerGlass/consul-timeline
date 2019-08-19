package noop

import (
	"context"

	"github.com/aestek/consul-timeline/storage"
	tl "github.com/aestek/consul-timeline/timeline"
)

var _ storage.Storage = (*Storage)(nil)

const Name = "noop"

type Storage struct {
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Store(evt tl.Event) error {
	return nil
}

func (s *Storage) Query(_ context.Context, q storage.Query) ([]tl.Event, error) {
	return nil, nil
}
