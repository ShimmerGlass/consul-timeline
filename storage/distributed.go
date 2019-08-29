package storage

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/aestek/consul-timeline/consul"
	tl "github.com/aestek/consul-timeline/timeline"
	"github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
)

var _ Storage = (*Distributed)(nil)

type Distributed struct {
	consul *consul.Consul
	inner  Storage

	stop chan struct{}
	done sync.WaitGroup

	enabled uint32
}

func NewDistributed(consul *consul.Consul, inner Storage) *Distributed {
	s := &Distributed{
		consul: consul,
		inner:  inner,
		stop:   make(chan struct{}),
	}

	go s.lockLoop()

	return s
}

func (s *Distributed) lockLoop() {

	var lock *api.Lock

	for {
		l, err := s.consul.Lock()
		if err != nil {
			log.Error(err)
			time.Sleep(time.Second)
			continue
		}
		lock = l
		break
	}

	for {
		isLeaderGauge.Set(0)

		log.Info("aquiring lock")
		lockChan, err := lock.Lock(nil)

		if err != nil {
			log.Error(err)
			time.Sleep(time.Second)
			continue
		}

		log.Info("storage lock aquired")
		isLeaderGauge.Set(1)
		s.done.Add(1)
		atomic.StoreUint32(&s.enabled, 1)
		select {
		case <-lockChan:
			atomic.StoreUint32(&s.enabled, 0)
			log.Info("storage lock lost")
			s.done.Done()

		case <-s.stop:
			atomic.StoreUint32(&s.enabled, 0)
			log.Info("unlocking storage lock")
			lock.Unlock()
			s.done.Done()
			return
		}
	}
}

func (s *Distributed) Stop() {
	close(s.stop)
	s.done.Wait()
}

func (s *Distributed) Store(evt tl.Event) error {
	if atomic.LoadUint32(&s.enabled) == 0 {
		// we are not leader, frop the event
		return nil
	}

	return s.inner.Store(evt)
}

func (s *Distributed) Query(ctx context.Context, q Query) ([]tl.Event, error) {
	return s.inner.Query(ctx, q)
}
