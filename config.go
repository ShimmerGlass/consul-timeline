package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/ghodss/yaml"

	"github.com/aestek/consul-timeline/consul"
	"github.com/aestek/consul-timeline/server"
	"github.com/aestek/consul-timeline/storage/cassandra"
	"github.com/aestek/consul-timeline/storage/mysql"
)

type Config struct {
	Consul    consul.Config `json:"consul"`
	Server    server.Config `json:"server"`
	Mysql     *mysql.Config `json:"mysql"`
	Cassandra *cass.Config  `json:"cassandra"`
}

var configFileFlag = flag.String("config", "", "Config file path (yaml, json)")
var storageFlag = flag.String("storage", "mysql", "Storage backend (mysql, cassandra)")

func FromFlags() Config {
	cfg := Config{
		Consul: consul.ConfigFromFlags(),
		Server: server.ConfigFromFlags(),
	}

	switch *storageFlag {
	case "mysql":
		c := mysql.ConfigFromFlags()
		cfg.Mysql = &c
	case "cassandra":
		c := cass.ConfigFromFlags()
		cfg.Cassandra = &c
	default:
		log.Fatalf("unknown storage %s", *storageFlag)
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

	var cfg Config
	err = yaml.Unmarshal(f, &cfg)
	if err != nil {
		log.Fatal(err)
	}

	return cfg
}
