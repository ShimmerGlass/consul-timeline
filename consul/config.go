package consul

import "flag"

type Config struct {
	Address               string `json:"address"`
	Token                 string `json:"token"`
	EnableDistributedLock bool   `json:"enable_distributed_lock"`
	LockPath              string `json:"lock_path"`
}

var DefaultConfig = Config{
	Address:               "localhost:8500",
	Token:                 "",
	EnableDistributedLock: false,
	LockPath:              "consul_timeline/lock",
}

var flagConfig Config

func init() {
	flag.StringVar(&flagConfig.Address, "consul", DefaultConfig.Address, "Consul agent address")
	flag.StringVar(&flagConfig.Token, "consul-token", DefaultConfig.Token, "Consul ACL token")
	flag.BoolVar(&flagConfig.EnableDistributedLock, "consul-enable-distributed-lcok", DefaultConfig.EnableDistributedLock, "Multi timeline instance lock for storage")
	flag.StringVar(&flagConfig.LockPath, "consul-lock-path", DefaultConfig.LockPath, "Consul lock path")
}

func ConfigFromFlags() Config {
	return flagConfig
}
