package multitenant

import (
	"time"

	"k8s.io/helm/pkg/proto/hapi/chart"
)

// EventAction constants used in action field of Event.
// const (
// 	EventActionPull   = "pull"
// 	EventActionPush   = "push"
// 	EventActionDelete = "delete"
// )

// Event provides the fields required to describe a registry event.
type Event struct {
	// ID provides a unique identifier for the event.
	// ID string `json:"id,omitempty"`

	// Timestamp is the time at which the event occurred.
	Timestamp time.Time `json:"timestamp,omitempty"`

	// Action indicates what action encompasses the provided event.
	Action string `json:"action,omitempty"`

	Repo string `json:"repo,omitempty"`

	Filename string `json:"filename,omitempty"`

	Metadata *chart.Metadata `json:"metadata,omitempty"`
}
