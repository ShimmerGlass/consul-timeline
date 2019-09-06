package watch

import (
	"github.com/hashicorp/consul/agent/structs"
)

type consul interface {
	Services(idx uint64) (*structs.IndexedServices, error)
	Service(idx uint64, name string) (*structs.IndexedCheckServiceNodes, error)

	Nodes(idx uint64) (*structs.IndexedNodes, error)
	Node(idx uint64, name string) (*structs.IndexedHealthChecks, error)

	Datacenter() string
}
