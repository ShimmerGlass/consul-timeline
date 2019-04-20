package cass

import (
	"time"

	"github.com/aestek/consul-timeline/timeline"
	"github.com/gocql/gocql"
	"github.com/pkg/errors"
	"github.com/scylladb/gocqlx"
	"github.com/scylladb/gocqlx/qb"
	log "github.com/sirupsen/logrus"
)

type idxGlobal struct {
	session *gocql.Session
}

func (i *idxGlobal) Store(evt tl.Event) error {
	stmt, names := qb.Insert("events_global").Columns(
		"id",
		"time_block",
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
	q := gocqlx.Query(i.session.Query(stmt), names).BindStruct(cassEvt{
		Event:     evt,
		TimeBlock: evt.Time.Truncate(timeBlockResolution),
	})
	return errors.Wrap(q.ExecRelease(), "cassandra event insert")
}

func (i *idxGlobal) FetchBefore(start time.Time, limit int) ([]tl.Event, error) {
	events := make([]tl.Event, 0)

	blockStart := start.Truncate(timeBlockResolution)

	for j := 0; j < 10; j++ {
		if len(events) >= limit {
			break
		}

		evts, err := i.fetchBlock(blockStart)
		if err != nil {
			return nil, err
		}

		for _, e := range evts {
			if e.Time.After(start) {
				evts = evts[1:]
			} else {
				break
			}
		}

		events = append(events, evts...)
		blockStart = blockStart.Add(-timeBlockResolution)
	}

	return events, nil
}

func (i *idxGlobal) fetchBlock(t time.Time) ([]tl.Event, error) {
	log.Debugf("cassandra: requesting global block %s", t)
	stmt, names := qb.Select("events_global").
		Where(qb.Eq("time_block")).
		OrderBy("time", qb.DESC).
		ToCql()

	q := gocqlx.Query(i.session.Query(stmt), names).BindMap(qb.M{
		"time_block": t,
	})

	var events []cassEvt
	err := q.SelectRelease(&events)
	if err != nil {
		return nil, errors.Wrap(err, "cassandra events fetch")
	}

	rEvents := make([]tl.Event, 0, len(events))

	for _, e := range events {
		rEvents = append(rEvents, e.Event)
	}

	return rEvents, nil
}
