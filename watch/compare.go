package watch

import (
	"encoding/hex"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/hashicorp/consul/types"

	tl "github.com/aestek/consul-timeline/timeline"
	"github.com/hashicorp/consul/agent/structs"
	api "github.com/hashicorp/consul/api"
)

func (w *Watcher) compareServiceStates(at time.Time, old, new []structs.CheckServiceNode) {
	type key struct {
		serviceName string
		serviceID   string
		node        string
	}

	oldInstanceCount := 0
	newInstanceCount := 0

	oldIdx := map[key]structs.CheckServiceNode{}
	newIdx := map[key]structs.CheckServiceNode{}

	for _, v := range old {
		if aggregatedStatus(v.Checks) == tl.StatusPassing {
			oldInstanceCount++
		}
		oldIdx[key{
			serviceName: v.Service.Service,
			serviceID:   v.Service.ID,
			node:        v.Node.Node,
		}] = v
	}

	for _, v := range new {
		if aggregatedStatus(v.Checks) == tl.StatusPassing {
			newInstanceCount++
		}
		newIdx[key{
			serviceName: v.Service.Service,
			serviceID:   v.Service.ID,
			node:        v.Node.Node,
		}] = v
	}

	for key, newState := range newIdx {
		oldChecks := structs.HealthChecks{}
		oldServiceStatus := tl.StatusMissing
		newServiceStatus := aggregatedStatus(newState.Checks)

		oldState, ok := oldIdx[key]
		if ok {
			oldServiceStatus = aggregatedStatus(oldState.Checks)
			oldChecks = oldState.Checks
		}

		base := tl.Event{
			ID:               w.nextEventID(),
			Time:             at,
			ServiceName:      newState.Service.Service,
			ServiceID:        newState.Service.ID,
			NodeName:         newState.Node.Node,
			NodeIP:           newState.Node.Address,
			OldServiceStatus: oldServiceStatus,
			NewServiceStatus: newServiceStatus,
			OldInstanceCount: oldInstanceCount,
			NewInstanceCount: newInstanceCount,
		}

		w.compareChecks(base, oldChecks, newState.Checks)

		if len(oldChecks) == 0 && len(newState.Checks) == 0 && oldServiceStatus != newServiceStatus {
			w.sendEvent(base)
		}
	}

	for key, oldState := range oldIdx {
		_, ok := newIdx[key]
		if ok {
			continue
		}

		base := tl.Event{
			ID:               w.nextEventID(),
			Time:             at,
			ServiceName:      oldState.Service.Service,
			ServiceID:        oldState.Service.ID,
			NodeName:         oldState.Node.Node,
			NodeIP:           oldState.Node.Address,
			OldServiceStatus: aggregatedStatus(oldState.Checks),
			NewServiceStatus: tl.StatusMissing,
			OldInstanceCount: oldInstanceCount,
			NewInstanceCount: newInstanceCount,
		}

		w.compareChecks(base, oldState.Checks, structs.HealthChecks{})

		if len(oldState.Checks) == 0 {
			w.sendEvent(base)
		}
	}
}

func convertToASCII(in string) string {
	var out strings.Builder

	var i int
	for i < len(in) {
		r, n := utf8.DecodeRuneInString(in[i:])

		if r <= utf8.RuneSelf {
			out.WriteRune(r)
		} else {
			out.WriteString("\\u")
			out.WriteString(hex.EncodeToString([]byte(in[i : i+n])))
		}

		i += n
	}

	return out.String()
}

func (w *Watcher) compareChecks(base tl.Event, old structs.HealthChecks, new structs.HealthChecks) {
	oldIdx := map[types.CheckID]*structs.HealthCheck{}
	newIdx := map[types.CheckID]*structs.HealthCheck{}

	for _, s := range old {
		oldIdx[s.CheckID] = s
	}

	for _, s := range new {
		newIdx[s.CheckID] = s
	}

	for _, new := range newIdx {
		oldStatus := tl.StatusMissing
		old, ok := oldIdx[new.CheckID]
		if ok {
			if old.Status == new.Status {
				continue
			}

			oldStatus = tl.StatusFromString(old.Status)
		}

		evt := base
		evt.ID = w.nextEventID()
		evt.CheckName = new.Name
		evt.OldCheckStatus = oldStatus
		evt.NewCheckStatus = tl.StatusFromString(new.Status)
		evt.CheckOutput = convertToASCII(new.Output)
		w.sendEvent(evt)
	}

	for _, old := range oldIdx {
		_, ok := newIdx[old.CheckID]
		if ok {
			continue
		}

		evt := base
		evt.ID = w.nextEventID()
		evt.CheckName = old.Name
		evt.OldCheckStatus = tl.StatusFromString(old.Status)
		evt.NewCheckStatus = tl.StatusMissing
		w.sendEvent(evt)
	}
}

func aggregatedStatus(c structs.HealthChecks) tl.Status {
	if len(c) == 0 {
		return tl.StatusPassing
	}
	var passing, warning, critical, maintenance bool
	for _, check := range c {
		id := string(check.CheckID)
		if id == api.NodeMaint || strings.HasPrefix(id, api.ServiceMaintPrefix) {
			maintenance = true
			continue
		}

		switch check.Status {
		case api.HealthPassing:
			passing = true
		case api.HealthWarning:
			warning = true
		case api.HealthCritical:
			critical = true
		}
	}

	switch {
	case maintenance:
		return tl.StatusMaintenance
	case critical:
		return tl.StatusCritical
	case warning:
		return tl.StatusWarning
	case passing:
		return tl.StatusPassing
	default:
		return tl.StatusUnknown
	}
}
