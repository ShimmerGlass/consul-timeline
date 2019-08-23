package memory

import (
	"flag"
)

const Name = "memory"

type Config struct {
	MaxSize int `json:"max_size"`
}

var DefaultConfig = Config{
	MaxSize: 10000,
}

var flagConfig Config

func init() {
	flag.IntVar(&flagConfig.MaxSize, "storage-memory-max-size", DefaultConfig.MaxSize, "Max events to store")
}

func ConfigFromFlags() Config {
	return flagConfig
}
