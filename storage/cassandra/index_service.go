package cass

import (
	"time"

	"github.com/aestek/consul-timeline/timeline"
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
)

type idxService struct {
	session *gocql.Session
}

func (i *idxService) Store(evt tl.Event) error {
	if evt.ServiceName == "" {
		return nil
	}

	stmt, names := qb.Insert("events_service").Columns(
		"id",
		"time",
		"node_name",
		"node_ip",
		"old_node_status",
		"new_node_status",
		"service_name",
		"service_id",
		"old_service_status",
		"new_service_status",
		"old_instance_count",
		"new_instance_count",
		"check_name",
		"old_check_status",
		"new_check_status",
		"check_output",
	).ToCql()
	q := gocqlx.Query(i.session.Query(stmt), names).BindStruct(evt)
	return errors.Wrap(q.ExecRelease(), "cassandra event insert")
}

func (i *idxService) FetchBefore(service string, start time.Time, limit int) ([]tl.Event, error) {
	stmt, names := qb.Select("events_service").Where(
		qb.Eq("service_name"),
		qb.Lt("time"),
	).Limit(uint(limit)).OrderBy("time", qb.DESC).ToCql()
	q := gocqlx.Query(i.session.Query(stmt), names).BindMap(qb.M{
		"service_name": service,
		"time":         start,
	})

	events := make([]tl.Event, 0)

	err := q.SelectRelease(&events)
	if err != nil {
		return nil, errors.Wrap(err, "cassandra events fetch")
	}

	return events, nil
}
