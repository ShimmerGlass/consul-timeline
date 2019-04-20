package watch

import (
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/aestek/consul-timeline/timeline"
	"github.com/hashicorp/consul/agent/structs"
	log "github.com/sirupsen/logrus"
)

const (
	waitOnErr = 5 * time.Second
)

type serviceWatch struct {
	Stop  uint32
	State []structs.CheckServiceNode
}

type nodeWatch struct {
	Stop   uint32
	Node   *structs.Node
	Checks structs.HealthChecks
}

type Watcher struct {
	c consul

	lock sync.Mutex

	services map[string]*serviceWatch
	nodes    map[string]*nodeWatch

	ready   bool
	readyWg sync.WaitGroup

	evtID int32
	out   chan tl.Event
}

func New(c consul, bufferSize int) *Watcher {
	return &Watcher{
		c:        c,
		services: make(map[string]*serviceWatch),
		nodes:    make(map[string]*nodeWatch),
		out:      make(chan tl.Event, bufferSize),
	}
}

func (w *Watcher) Run() <-chan tl.Event {
	log.Info("watch: starting")
	w.readyWg.Add(2)
	go w.watchServices()
	go w.watchNodes()
	w.readyWg.Wait()
	w.ready = true
	log.Info("watch ready")
	return w.out
}

func (w *Watcher) Services() []string {
	services := []string{}

	w.lock.Lock()
	for s := range w.services {
		services = append(services, s)
	}
	w.lock.Unlock()

	sort.Strings(services)
	return services
}

func (w *Watcher) watchServices() {
	var idx uint64

	for {
		if idx > 0 {
			w.readyWg.Wait()
		}
		res, err := w.c.Services(idx)
		if err != nil {
			log.Errorf("error getting service list: %s", err)
			time.Sleep(waitOnErr)
			continue
		}

		// init
		if idx == 0 {
			w.readyWg.Add(len(res.Services))
			w.readyWg.Done()
		}

		w.handleServicesChanged(res.Services)
		idx = res.Index
	}
}

func (w *Watcher) handleServicesChanged(services map[string][]string) {
	w.lock.Lock()
	defer w.lock.Unlock()

	for s := range services {
		if _, ok := w.services[s]; ok {
			atomic.StoreUint32(&w.services[s].Stop, 0)
			continue
		}

		w.watchService(s)
	}

	for s := range w.services {
		if _, ok := services[s]; !ok {
			atomic.StoreUint32(&w.services[s].Stop, 1)
		}
	}
}

func (w *Watcher) watchService(name string) {
	log.Printf("watching service %s", name)

	state := &serviceWatch{}
	w.services[name] = state

	go func() {
		var idx uint64

		for {
			if idx > 0 {
				w.readyWg.Wait()
			}

			if atomic.LoadUint32(&state.Stop) == 1 {
				w.lock.Lock()
				delete(w.services, name)
				w.lock.Unlock()
				log.Printf("stopping watching service %s", name)
				return
			}

			res, err := w.c.Service(idx, name)

			if err != nil {
				log.Errorf("error getting service %s: %s", name, err)
				time.Sleep(waitOnErr)
				continue
			}

			if idx == 0 && !w.ready {
				w.readyWg.Done()
			}

			w.lock.Lock()
			if w.ready {
				w.handleServiceChanged(name, time.Now(), res.Nodes)
			}
			w.services[name].State = res.Nodes
			w.lock.Unlock()

			idx = res.Index
		}
	}()
}

func (w *Watcher) handleServiceChanged(name string, at time.Time, state []structs.CheckServiceNode) {
	if _, ok := w.services[name]; !ok {
		log.Printf("no entry for %s", name)
		return
	}

	w.compareServiceStates(at, w.services[name].State, state)
}

func (w *Watcher) watchNodes() {
	var idx uint64

	for {
		if idx > 0 {
			w.readyWg.Wait()
		}
		res, err := w.c.Nodes(idx)
		if err != nil {
			log.Errorf("error getting node list: %s", err)
			time.Sleep(waitOnErr)
			continue
		}

		if idx == 0 && !w.ready {
			w.readyWg.Add(len(res.Nodes))
			w.readyWg.Done()
		}

		idx = res.Index

		w.handleNodesChanged(res.Nodes)
	}
}

func (w *Watcher) handleNodesChanged(nodes []*structs.Node) {
	w.lock.Lock()
	defer w.lock.Unlock()

	new := map[string]*structs.Node{}
	for _, n := range nodes {
		new[n.Node] = n
	}

	for n, node := range new {
		if _, ok := w.nodes[n]; ok {
			w.nodes[n].Node = node
			atomic.StoreUint32(&w.nodes[n].Stop, 0)
			continue
		}

		w.watchNode(node)
	}

	for n := range w.nodes {
		if _, ok := new[n]; !ok {
			atomic.StoreUint32(&w.nodes[n].Stop, 1)
		}
	}
}

func (w *Watcher) watchNode(node *structs.Node) {
	log.Printf("watching node %s", node.Node)

	state := &nodeWatch{Node: node}
	w.nodes[node.Node] = state

	go func() {
		var idx uint64

		for {
			if idx > 0 {
				w.readyWg.Wait()
			}

			if atomic.LoadUint32(&state.Stop) == 1 {
				w.lock.Lock()
				delete(w.nodes, node.Node)
				w.lock.Unlock()
				log.Printf("stopping watching node %s", node.Node)
				return
			}

			res, err := w.c.Node(idx, node.Node)

			if err != nil {
				log.Errorf("error getting node %s: %s", node.Node, err)
				time.Sleep(waitOnErr)
				continue
			}

			filteredChecks := res.HealthChecks[:0]
			for _, c := range res.HealthChecks {
				if c.ServiceID == "" {
					filteredChecks = append(filteredChecks, c)
				}
			}

			if idx == 0 && !w.ready {
				w.readyWg.Done()
			}

			w.lock.Lock()
			if w.ready {
				w.handleNodeChanged(node.Node, time.Now(), filteredChecks)
			}
			state.Checks = filteredChecks
			w.lock.Unlock()

			idx = res.Index
		}
	}()
}

func (w *Watcher) handleNodeChanged(name string, at time.Time, checks structs.HealthChecks) {
	if _, ok := w.nodes[name]; !ok {
		log.Printf("no entry for %s", name)
		return
	}

	state := w.nodes[name]

	oldStatus, newStatus := tl.StatusMissing, tl.StatusMissing

	if len(state.Checks) > 0 {
		oldStatus = aggregatedStatus(state.Checks)
	}
	if len(checks) > 0 {
		newStatus = aggregatedStatus(checks)
	}

	base := tl.Event{
		Time:          at,
		NodeName:      state.Node.Node,
		NodeIP:        state.Node.Address,
		OldNodeStatus: oldStatus,
		NewNodeStatus: newStatus,
	}

	w.compareChecks(base, state.Checks, checks)
}

func (w *Watcher) nextEventID() int32 {
	return atomic.AddInt32(&w.evtID, 1)
}
