package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/ghodss/yaml"

	"github.com/aestek/consul-timeline/consul"
	"github.com/aestek/consul-timeline/server"
	cass "github.com/aestek/consul-timeline/storage/cassandra"
	"github.com/aestek/consul-timeline/storage/mysql"
)

type Config struct {
	LogLevel  string        `json:"log_level"`
	Storage   string        `json:"storage"`
	Consul    consul.Config `json:"consul"`
	Server    server.Config `json:"server"`
	Mysql     mysql.Config  `json:"mysql"`
	Cassandra cass.Config   `json:"cassandra"`
}

var (
	logLevelFlag   = flag.String("log-level", "info", "(debug, info, warning, error, fatal)")
	configFileFlag = flag.String("config", "", "Config file path (yaml, json)")
	storageFlag    = flag.String("storage", "mysql", "Storage backend (mysql, cassandra)")
)

func FromFlags() Config {
	cfg := Config{
		LogLevel:  *logLevelFlag,
		Storage:   *storageFlag,
		Consul:    consul.ConfigFromFlags(),
		Server:    server.ConfigFromFlags(),
		Mysql:     mysql.ConfigFromFlags(),
		Cassandra: cass.ConfigFromFlags(),
	}

	return cfg
}

func GetConfig() Config {
	flag.Parse()

	if *configFileFlag == "" {
		return FromFlags()
	}

	f, err := ioutil.ReadFile(*configFileFlag)
	if err != nil {
		log.Fatal(err)
	}

	cfg := Config{
		LogLevel: "info",
	}

	err = yaml.Unmarshal(f, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}
