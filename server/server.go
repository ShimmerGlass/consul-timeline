package server

import (
	"encoding/json"
	"net/http"

	"github.com/NYTimes/gziphandler"
	"github.com/aestek/consul-timeline/storage"
	tl "github.com/aestek/consul-timeline/timeline"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
)

type FilterEntriesProvider interface {
	FilterEntries() []string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Server struct {
	listenAddr string

	storage storage.Storage
	router  *httprouter.Router

	events   <-chan tl.Event
	services FilterEntriesProvider

	ws *ws
}

func New(cfg Config, storage storage.Storage, services FilterEntriesProvider, events <-chan tl.Event) *Server {
	return &Server{
		listenAddr: cfg.ListenAddr,
		storage:    storage,
		router:     httprouter.New(),
		events:     events,
		services:   services,
		ws:         newWs(),
	}
}

func (s *Server) Serve() error {
	s.serveStatic()

	s.router.GET("/", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		http.Redirect(w, r, "/web/", 301)
	})

	s.router.GET("/events", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		filter, err := filterFromQuery(r.URL.Query())
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		q := storage.Query{
			Start:  filter.Start,
			Filter: filter.Filter,
			Limit:  filter.Limit,
		}

		events, err := s.storage.Query(r.Context(), q)
		if err != nil {
			log.Errorf("query error %s", err)
			http.Error(w, err.Error(), 500)
			return
		}

		json.NewEncoder(w).Encode(events)
	})

	s.router.GET("/filter-entries", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		json.NewEncoder(w).Encode(s.services.FilterEntries())
	})

	s.router.GET("/ws", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		filter, err := filterFromQuery(r.URL.Query())
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Errorf("ws connection upgrade: %s", err)
			return
		}

		s.ws.Add(conn, filter)
	})

	s.router.GET("/status", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Write([]byte("OK"))
	})

	go func() {
		for e := range s.events {
			s.ws.Send(e)
		}
	}()

	return http.ListenAndServe(s.listenAddr, gziphandler.GzipHandler(s.router))
}
