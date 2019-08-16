CREATE TABLE events (
    time datetime,
    node_name varchar(255),
    node_ip varchar(45),
    old_node_status int,
    new_node_status int,
    service_name varchar(255),
    service_id varchar(255),
    old_service_status int,
    new_service_status int,
    old_instance_count int,
    new_instance_count int,
    check_name  varchar(255),
    old_check_status int,
    new_check_status int,
    check_output varchar(2048)
)