package consul

import "flag"

type Config struct {
	Address string `json:"address"`
}

var flagConfig Config

func init() {
	flag.StringVar(&flagConfig.Address, "consul", "localhost:8500", "Consul agent address")
}

func ConfigFromFlags() Config {
	return flagConfig
}
