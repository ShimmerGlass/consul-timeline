package server

import "flag"

type Config struct {
	ListenAddr string `json:"listen"`
}

var flagConfig Config

func init() {
	flag.StringVar(&flagConfig.ListenAddr, "listen", ":8888", "Server listen address")
}

func ConfigFromFlags() Config {
	return flagConfig
}
