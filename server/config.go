package server

import "flag"

type Config struct {
	ListenAddr string `json:"listen"`
}

var DefaultConfig = Config{
	ListenAddr: ":8888",
}

var flagConfig Config

func init() {
	flag.StringVar(&flagConfig.ListenAddr, "listen", DefaultConfig.ListenAddr, "Server listen address")
}

func ConfigFromFlags() Config {
	return flagConfig
}
