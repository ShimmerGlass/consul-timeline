package mysql

import (
	"flag"
)

const Name = "mysql"

type Config struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

var flagConfig Config

func init() {
	flag.StringVar(&flagConfig.Host, "mysql-host", "localhost", "MySQL server host")
	flag.IntVar(&flagConfig.Port, "mysql-port", 3306, "MySQL server port")
	flag.StringVar(&flagConfig.User, "mysql-user", "root", "MySQL user")
	flag.StringVar(&flagConfig.Password, "mysql-password", "", "MySQL server password")
	flag.StringVar(&flagConfig.Database, "mysql-db", "consul_timeline", "MySQL database name")
}

func ConfigFromFlags() Config {
	return flagConfig
}
