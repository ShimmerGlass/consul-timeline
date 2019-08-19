package cass

import (
	"flag"
	"strings"
)

const Name = "cassandra"

type Config struct {
	Addresses []string
	Keyspace  string
}

var (
	addr     = flag.String("cassandra", "", "Cassandra addresses, comma separated")
	keyspace = flag.String("cassandra-keyspace", "consul_timeline", "Cassandra keyspace")
)

func ConfigFromFlags() Config {
	return Config{
		Addresses: strings.Split(*addr, ","),
		Keyspace:  *keyspace,
	}
}
