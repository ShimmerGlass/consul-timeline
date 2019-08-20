package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/aestek/consul-timeline/storage"
	tl "github.com/aestek/consul-timeline/timeline"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

var _ storage.Storage = (*Storage)(nil)

type Storage struct {
	db *sql.DB
}

func New(cfg Config) (*Storage, error) {
	db, err := sql.Open(
		"mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?parseTime=true",
			cfg.User,
			cfg.Password,
			cfg.Host,
			cfg.Port,
			cfg.Database,
		),
	)
	if err != nil {
		return nil, err
	}

	return &Storage{db}, nil
}

func (s *Storage) Store(evt tl.Event) error {
	if len(evt.CheckOutput) > 2048 {
		evt.CheckOutput = evt.CheckOutput[:2048]
	}
	_, err := s.db.Exec(
		`
			INSERT INTO events (
				time,
				node_name,
				node_ip,
				old_node_status,
				new_node_status,
				service_name,
				service_id,
				old_service_status,
				new_service_status,
				old_instance_count,
				new_instance_count,
				check_name,
				old_check_status,
				new_check_status,
				check_output
			) VALUES (
				?, ?, ?, ?, ?, ?, ?, ?,
				?, ?, ?, ?, ?, ?, ?
			)
		`,
		evt.Time,
		evt.NodeName,
		evt.NodeIP,
		evt.OldNodeStatus,
		evt.NewNodeStatus,
		evt.ServiceName,
		evt.ServiceID,
		evt.OldServiceStatus,
		evt.NewServiceStatus,
		evt.OldInstanceCount,
		evt.NewInstanceCount,
		evt.CheckName,
		evt.OldCheckStatus,
		evt.NewCheckStatus,
		evt.CheckOutput,
	)
	return errors.Wrap(err, "mysql event insert")
}

func (s *Storage) Query(ctx context.Context, q storage.Query) ([]tl.Event, error) {
	args := []interface{}{}
	qs := `
		SELECT
			time,
			node_name,
			node_ip,
			old_node_status,
			new_node_status,
			service_name,
			service_id,
			old_service_status,
			new_service_status,
			old_instance_count,
			new_instance_count,
			check_name,
			old_check_status,
			new_check_status,
			check_output
		FROM events
		WHERE time <= ?
	`
	args = append(args, q.Start)
	if q.Filter != "" {
		qs += "&& (service_name = ? || node_name = ?) \n"
		args = append(args, q.Filter, q.Filter)
	}
	qs += "ORDER BY `time` DESC\n"
	qs += "LIMIT 0, ?\n"
	args = append(args, q.Limit)

	rows, err := s.db.QueryContext(ctx, qs, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := []tl.Event{}
	for rows.Next() {
		evt := tl.Event{}
		err := rows.Scan(
			&evt.Time,
			&evt.NodeName,
			&evt.NodeIP,
			&evt.OldNodeStatus,
			&evt.NewNodeStatus,
			&evt.ServiceName,
			&evt.ServiceID,
			&evt.OldServiceStatus,
			&evt.NewServiceStatus,
			&evt.OldInstanceCount,
			&evt.NewInstanceCount,
			&evt.CheckName,
			&evt.OldCheckStatus,
			&evt.NewCheckStatus,
			&evt.CheckOutput,
		)
		if err != nil {
			return nil, err
		}
		res = append(res, evt)
	}

	return res, nil
}
