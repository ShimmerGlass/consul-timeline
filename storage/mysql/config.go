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

	PurgeFrequency   int `json:"purge_frequency"`
	PurgeMaxAgeHours int `json:"purge_max_age_hours"`

	SetupSchema bool `json:"setup_schema"`
	PrintSchema bool `json:"-"`
}

var flagConfig Config

func init() {
	flag.BoolVar(&flagConfig.SetupSchema, "mysql-setup-schema", false, "Automatically setup MySQL schema")
	flag.BoolVar(&flagConfig.PrintSchema, "mysql-print-schema", false, "Print MySQL schema")

	flag.StringVar(&flagConfig.Host, "mysql-host", "localhost", "MySQL server host")
	flag.IntVar(&flagConfig.Port, "mysql-port", 3306, "MySQL server port")
	flag.StringVar(&flagConfig.User, "mysql-user", "root", "MySQL user")
	flag.StringVar(&flagConfig.Password, "mysql-password", "", "MySQL server password")
	flag.StringVar(&flagConfig.Database, "mysql-db", "consul_timeline", "MySQL database name")

	flag.IntVar(&flagConfig.PurgeMaxAgeHours, "mysql-purge-max-age-hours", 2*7*24, "Periodically delete events older than this duration")
	flag.IntVar(&flagConfig.PurgeFrequency, "mysql-purge-frequency", 10000, "Purge events every n writes")
}

func ConfigFromFlags() Config {
	return flagConfig
}
