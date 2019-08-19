package consul

import (
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/hashicorp/consul/api"

	"github.com/hashicorp/consul/agent/pool"
	"github.com/hashicorp/consul/agent/structs"
)

type Consul struct {
	config      Config
	connPool    *pool.ConnPool
	client      *api.Client
	dc          string
	servers     []net.Addr
	serverIndex uint64

	ready sync.WaitGroup
}

func New(cfg Config) *Consul {
	client, err := api.NewClient(&api.Config{
		Address: cfg.Address,
		Token:   cfg.Token,
	})
	if err != nil {
		log.Fatal(err)
	}
	c := &Consul{
		config: cfg,
		connPool: &pool.ConnPool{
			LogOutput:  os.Stderr,
			MaxTime:    30 * time.Second,
			MaxStreams: 50,
		},
		client: client,
	}

	c.ready.Add(1)

	go c.watchServers()

	return c
}

func (c *Consul) Services(idx uint64) (*structs.IndexedServices, error) {
	c.ready.Wait()

	out := &structs.IndexedServices{}
	err := c.rpc("Catalog.ListServices", &structs.DCSpecificRequest{
		Datacenter: c.dc,
		QueryOptions: structs.QueryOptions{
			MinQueryIndex: idx,
			MaxQueryTime:  10 * time.Minute,
		},
	}, out)

	return out, err
}

func (c *Consul) Service(idx uint64, name string) (*structs.IndexedCheckServiceNodes, error) {
	c.ready.Wait()

	out := &structs.IndexedCheckServiceNodes{}
	err := c.rpc("Health.ServiceNodes", &structs.ServiceSpecificRequest{
		Datacenter:  c.dc,
		ServiceName: name,
		QueryOptions: structs.QueryOptions{
			MinQueryIndex: idx,
			MaxQueryTime:  10 * time.Minute,
		},
	}, out)

	return out, err
}

func (c *Consul) Nodes(idx uint64) (*structs.IndexedNodes, error) {
	c.ready.Wait()

	out := &structs.IndexedNodes{}
	err := c.rpc("Catalog.ListNodes", &structs.DCSpecificRequest{
		Datacenter: c.dc,
		QueryOptions: structs.QueryOptions{
			MinQueryIndex: idx,
			MaxQueryTime:  10 * time.Minute,
		},
	}, out)

	return out, err
}
func (c *Consul) Node(idx uint64, name string) (*structs.IndexedHealthChecks, error) {
	c.ready.Wait()

	out := &structs.IndexedHealthChecks{}
	err := c.rpc("Health.NodeChecks", &structs.NodeSpecificRequest{
		Datacenter: c.dc,
		Node:       name,
		QueryOptions: structs.QueryOptions{
			MinQueryIndex: idx,
			MaxQueryTime:  10 * time.Minute,
		},
	}, out)

	return out, err
}

func (c *Consul) Lock() (*api.Lock, error) {
	return c.client.LockOpts(&api.LockOptions{
		SessionTTL: (10 * time.Second).String(),
		Key:        c.config.LockPath,
	})
}

func (c *Consul) watchServers() {
	idx := uint64(0)
	for {
		instances, meta, err := c.client.Health().Service("consul", "", true, &api.QueryOptions{
			WaitIndex: idx,
			WaitTime:  10 * time.Minute,
		})
		if err != nil {
			log.Errorf("error retrieving consul servers: %s", err)
			time.Sleep(time.Second)
			continue
		}

		servers := []net.Addr{}
		for _, i := range instances {
			servers = append(servers, &net.TCPAddr{
				IP:   net.ParseIP(i.Node.Address),
				Port: i.Service.Port,
			})
		}
		c.servers = servers
		c.dc = instances[0].Node.Datacenter

		if idx == 0 {
			c.ready.Done()
		}

		idx = meta.LastIndex
	}
}

func (c *Consul) rpc(method string, in, out interface{}) error {
	idx := atomic.AddUint64(&c.serverIndex, 1)
	server := c.servers[int(idx)%len(c.servers)]
	return c.connPool.RPC(c.dc, server, 3, method, false, in, out)
}
