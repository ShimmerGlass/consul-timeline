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

var DefaultConfig = Config{
	Host:             "localhost",
	Port:             3306,
	User:             "",
	Password:         "",
	Database:         "consul_timeline",
	PurgeFrequency:   10000,
	PurgeMaxAgeHours: 2 * 7 * 24,
	SetupSchema:      false,
	PrintSchema:      false,
}

var flagConfig Config

func init() {
	flag.BoolVar(&flagConfig.SetupSchema, "mysql-setup-schema", DefaultConfig.SetupSchema, "Automatically setup MySQL schema")
	flag.BoolVar(&flagConfig.PrintSchema, "mysql-print-schema", DefaultConfig.PrintSchema, "Print MySQL schema")

	flag.StringVar(&flagConfig.Host, "mysql-host", DefaultConfig.Host, "MySQL server host")
	flag.IntVar(&flagConfig.Port, "mysql-port", DefaultConfig.Port, "MySQL server port")
	flag.StringVar(&flagConfig.User, "mysql-user", DefaultConfig.User, "MySQL user")
	flag.StringVar(&flagConfig.Password, "mysql-password", DefaultConfig.Password, "MySQL server password")
	flag.StringVar(&flagConfig.Database, "mysql-db", DefaultConfig.Database, "MySQL database name")

	flag.IntVar(&flagConfig.PurgeMaxAgeHours, "mysql-purge-max-age-hours", DefaultConfig.PurgeMaxAgeHours, "Periodically delete events older than this duration")
	flag.IntVar(&flagConfig.PurgeFrequency, "mysql-purge-frequency", DefaultConfig.PurgeFrequency, "Purge events every n writes")
}

func ConfigFromFlags() Config {
	return flagConfig
}
