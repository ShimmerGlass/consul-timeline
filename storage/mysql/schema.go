package mysql

import (
	"fmt"
	"strings"
)

var Schema = []string{
	"CREATE TABLE IF NOT EXISTS events (\n" +
		"    time DATETIME,\n" +
		"    datacenter VARCHAR(50),\n" +
		"    node_name VARCHAR(255),\n" +
		"    node_ip VARCHAR(45),\n" +
		"    old_node_status TINYINT,\n" +
		"    new_node_status TINYINT,\n" +
		"    service_name VARCHAR(255),\n" +
		"    service_id VARCHAR(255),\n" +
		"    old_service_status TINYINT,\n" +
		"    new_service_status TINYINT,\n" +
		"    old_instance_count INT,\n" +
		"    new_instance_count INT,\n" +
		"    check_name  VARCHAR(255),\n" +
		"    old_check_status TINYINT,\n" +
		"    new_check_status TINYINT,\n" +
		"    check_output VARCHAR(2048)\n" +
		") CHARSET=utf8;\n",
	"CREATE INDEX IF NOT EXISTS time_idx ON events (`time` DESC);",
	"CREATE INDEX IF NOT EXISTS time_service_idx ON events (`time` DESC, `service_name`);",
	"CREATE INDEX IF NOT EXISTS time_node_idx ON events (`time` DESC, `node_name`);",
	"CREATE INDEX IF NOT EXISTS time_node_service_idx ON events (`time` DESC, `service_name`, `node_name`);",
}

func PrintSchema() {
	fmt.Println(strings.Join(Schema, "\n"))
}
