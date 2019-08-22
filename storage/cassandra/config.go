package cass

import (
	"flag"
	"strings"
)

const Name = "cassandra"

type Config struct {
	Addresses []string `json:"addresses"`
	Keyspace  string   `json:"keyspace"`
}

var DefaultConfig = Config{
	Keyspace: "consul_timeline",
}

var (
	addr     = flag.String("cassandra", strings.Join(DefaultConfig.Addresses, ", "), "Cassandra addresses, comma separated")
	keyspace = flag.String("cassandra-keyspace", DefaultConfig.Keyspace, "Cassandra keyspace")
)

func ConfigFromFlags() Config {
	return Config{
		Addresses: strings.Split(*addr, ","),
		Keyspace:  *keyspace,
	}
}
