package main

import (
	"flag"
	"strings"

	"github.com/aestek/consul-timeline/consul"
	"github.com/aestek/consul-timeline/server"
	cass "github.com/aestek/consul-timeline/storage/cassandra"
	"github.com/aestek/consul-timeline/timeline"
	"github.com/aestek/consul-timeline/watch"
	"github.com/gocql/gocql"
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

	consulAddr := flag.String("consul", "localhost:8500", "Consul agent address")
	listenAddr := flag.String("listen", ":8888", "Listen address")
	cassandraAddr := flag.String("cassandra", "127.0.0.1", "Cassandra addresses, comma separated")
	cassandraKeyspace := flag.String("cassandra-keyspace", "consul_timeline", "Cassandra keyspace")
	flag.Parse()

	// consul client
	consul := consul.New(*consulAddr)

	// cassandra client
	cluster := gocql.NewCluster(strings.Split(*cassandraAddr, ",")...)
	cluster.Keyspace = *cassandraKeyspace
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}

	// consul watch
	w := watch.New(consul, eventsBuffer)
	events := w.Run()

	storageEvents, apiEvents := dupEvents(events)

	// storage
	storage := cass.New(session)

	// HTTP api
	api := server.New(storage, w, apiEvents)
	go func() {
		err = api.Serve(*listenAddr)
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
