CREATE TABLE IF NOT EXISTS events (
    time DATETIME,
    node_name VARCHAR(255),
    node_ip VARCHAR(45),
    old_node_status TINYINT,
    new_node_status TINYINT,
    service_name VARCHAR(255),
    service_id VARCHAR(255),
    old_service_status TINYINT,
    new_service_status TINYINT,
    old_instance_count INT,
    new_instance_count INT,
    check_name  VARCHAR(255),
    old_check_status TINYINT,
    new_check_status TINYINT,
    check_output VARCHAR(2048)
) CHARSET=utf8;

CREATE INDEX IF NOT EXISTS time_idx ON events (`time` DESC);
CREATE INDEX IF NOT EXISTS time_service_idx ON events (`time` DESC, `service_name`);
CREATE INDEX IF NOT EXISTS time_node_idx ON events (`time` DESC, `node_name`);
CREATE INDEX IF NOT EXISTS time_node_service_idx ON events (`time` DESC, `service_name`, `node_name`);