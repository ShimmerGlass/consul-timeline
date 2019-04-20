package server

import (
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/aestek/consul-timeline/timeline"
	"github.com/pkg/errors"
)

type filter struct {
	Service string
	Start   time.Time
	Limit   int
}

func filterFromQuery(q url.Values) (filter, error) {
	f := filter{
		Service: strings.TrimSpace(q.Get("service")),
	}
	if q.Get("start") != "" {
		ts, err := strconv.ParseInt(q.Get("start"), 10, 64)
		if err != nil {
			return f, errors.Wrap(err, "error parsing start filter")
		}
		f.Start = time.Unix(ts, 0)
	}
	if q.Get("limit") != "" {
		l, err := strconv.Atoi(q.Get("limit"))
		if err != nil {
			return f, errors.Wrap(err, "error parsing limit filter")
		}
		f.Limit = l
	}
	return f, nil
}

func (f filter) Match(e tl.Event) bool {
	if f.Service != "" && e.ServiceName != f.Service {
		return false
	}

	return true
}
