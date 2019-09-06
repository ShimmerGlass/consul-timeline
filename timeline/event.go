package tl

import (
	"time"

	"github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
)

type Status int

const (
	StatusUnknown Status = iota
	StatusMissing
	StatusCritical
	StatusWarning
	StatusPassing
	StatusMaintenance
)

func StatusFromString(s string) Status {
	switch s {
	case api.HealthPassing:
		return StatusPassing
	case api.HealthWarning:
		return StatusWarning
	case api.HealthCritical:
		return StatusCritical
	case api.HealthMaint:
		return StatusMaintenance
	default:
		log.Warnf("unknown status string %s", s)
		return StatusUnknown
	}
}

type Event struct {
	ID int32 `json:"-"`

	Time       time.Time `json:"time,omitempty"`
	Datacenter string    `json:"datacenter,omuitempty"`

	NodeName      string `json:"node_name,omitempty"`
	NodeIP        string `json:"node_ip,omitempty"`
	OldNodeStatus Status `json:"old_node_status"`
	NewNodeStatus Status `json:"new_node_status"`

	ServiceName      string `json:"service_name,omitempty"`
	ServiceID        string `json:"service_id,omitempty"`
	OldServiceStatus Status `json:"old_service_status"`
	NewServiceStatus Status `json:"new_service_status"`
	OldInstanceCount int    `json:"old_instance_count"`
	NewInstanceCount int    `json:"new_instance_count"`

	CheckName      string `json:"check_name,omitempty"`
	OldCheckStatus Status `json:"old_check_status"`
	NewCheckStatus Status `json:"new_check_status"`
	CheckOutput    string `json:"check_output,omitempty"`
}
