package consul

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	rpcSuccessCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "consul_timeline_rpc_success_total",
		Help: "The total number of sent RPCs",
	})
	rpcErrorCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "consul_timeline_rpc_error_total",
		Help: "The total number of sent RPCs",
	})
)
