CREATE TABLE events_global (
    id int,
    time_block timestamp,
    time timestamp,
    node_name text,
    node_ip text,
    old_node_status int,
    new_node_status int,
    service_name text,
    service_id text,
    old_service_status int,
    new_service_status int,
    old_instance_count int,
    new_instance_count int,
    check_name  text,
    old_check_status int,
    new_check_status int,
    check_output text,
    PRIMARY KEY (
        time_block,
        time,
        id
    )
) WITH CLUSTERING ORDER BY (time DESC);

CREATE TABLE events_service (
    id int,
    time timestamp,
    node_name text,
    node_ip text,
    old_node_status int,
    new_node_status int,
    service_name text,
    service_id text,
    old_service_status int,
    new_service_status int,
    old_instance_count int,
    new_instance_count int,
    check_name  text,
    old_check_status int,
    new_check_status int,
    check_output text,
    PRIMARY KEY (
        service_name,
        time,
        id
    )
) WITH CLUSTERING ORDER BY (time DESC);