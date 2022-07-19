package zebra

import (
	"context"
	"errors"
	"strings"
	"time"
)

type Status struct {
	Fault       string    `json:"fault"`
	Lease       string    `json:"lease"`
	UsedBy      string    `json:"usedBy"`
	State       string    `json:"state"`
	CreatedTime time.Time `json:"createdTime"`
}

var (
	ErrFault       = errors.New(`fault is incorrect, must be in ["none", "minor", "major", "critical"]`)
	ErrLease       = errors.New(`lease is incorrect, must be in ["leased", "free", "setup"]`)
	ErrState       = errors.New(`state is incorrect, must be in ["active", "inactive"]`)
	ErrCreatedTime = errors.New(`createdTime is incorrect, must be before current time`)
)

func (s *Status) Validate(ctx context.Context) error {
	if !IsIn(strings.ToLower(s.Fault), []string{"none", "minor", "major", "critical"}) {
		return ErrFault
	}

	if !IsIn(strings.ToLower(s.Lease), []string{"leased", "free", "setup"}) {
		return ErrLease
	}

	if !IsIn(strings.ToLower(s.State), []string{"active", "inactive"}) {
		return ErrState
	}

	if !s.CreatedTime.Before(time.Now()) {
		return ErrCreatedTime
	}

	return nil
}

// DefaultStatus returns a Status object with starting values (i.e. healthy
// resource in a free state, active, no user, and create time as right now).
func DefaultStatus() Status {
	return Status{
		Fault:       "none",
		Lease:       "free",
		UsedBy:      "",
		State:       "active",
		CreatedTime: time.Now(),
	}
}
