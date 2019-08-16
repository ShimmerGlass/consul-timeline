package main

import (
	"github.com/aestek/consul-timeline/consul"
	"github.com/aestek/consul-timeline/storage"

	_ "github.com/go-sql-driver/mysql"

	"github.com/aestek/consul-timeline/server"
	cass "github.com/aestek/consul-timeline/storage/cassandra"
	"github.com/aestek/consul-timeline/storage/mysql"
	"github.com/aestek/consul-timeline/timeline"
	"github.com/aestek/consul-timeline/watch"
	log "github.com/sirupsen/logrus"
)

const (
	eventsBuffer = 200
)

func dupEvents(in <-chan tl.Event) (<-chan tl.Event, <-chan tl.Event) {
	o1 := make(chan tl.Event, eventsBuffer)
	o2 := make(chan tl.Event, eventsBuffer)
	go func() {
		for e := range in {
			o1 <- e
			o2 <- e
		}
	}()

	return o1, o2
}

func main() {
	log.SetLevel(log.DebugLevel)

	cfg := GetConfig()

	// storage
	var storage storage.Storage
	var err error

	if cfg.Mysql != nil {
		storage, err = mysql.New(*cfg.Mysql)
		if err != nil {
			log.Fatal(err)
		}
	} else if cfg.Cassandra != nil {
		storage, err = cass.New(*cfg.Cassandra)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Fatal("no storage provided")
	}

	// consul client
	consul := consul.New(cfg.Consul)

	// consul watch
	w := watch.New(consul, eventsBuffer)
	events := w.Run()

	storageEvents, apiEvents := dupEvents(events)

	// HTTP api
	api := server.New(cfg.Server, storage, w, apiEvents)
	go func() {
		err := api.Serve()
		if err != nil {
			log.Fatal(err)
		}
	}()

	for e := range storageEvents {
		err := storage.Store(e)
		if err != nil {
			log.Error(err)
		}
	}
}
