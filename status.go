package zebra

import (
	"context"
	"errors"
	"strings"
	"time"
)

// Status is a struct that sets the status of a resource.
//
// It should contain a fault, a lease,
// a string that represents a user's name, a state, and the time when it was created.
type Status struct {
	Fault       Fault     `json:"fault"`
	Lease       Lease     `json:"lease"`
	UsedBy      string    `json:"usedBy"`
	State       State     `json:"state"`
	CreatedTime time.Time `json:"createdTime"`
}

// Set types for Fault, Lease, State.
type (
	Fault uint8
	Lease uint8
	State uint8
)

// Default values to be used for status.
const (
	None Fault = iota
	Minor
	Major
	Critical
)

// Default values to be used for lease.
const (
	Leased Lease = iota
	Free
	Setup
)

// Default values to be used for state.
const (
	Active State = iota
	Inactive
)

const Unknown = "unknown"

// Errors that can occur in a lease or resource state.
var (
	ErrFault       = errors.New(`fault is incorrect, must be in ["none", "minor", "major", "critical"]`)
	ErrLease       = errors.New(`lease is incorrect, must be in ["leased", "free", "setup"]`)
	ErrState       = errors.New(`state is incorrect, must be in ["active", "inactive"]`)
	ErrCreatedTime = errors.New(`createdTime is incorrect, must be before current time`)
)

func (f *Fault) String() string {
	strs := map[Fault]string{None: "none", Minor: "minor", Major: "major", Critical: "critical"}
	fstr, ok := strs[*f]

	if !ok {
		return Unknown
	}

	return fstr
}

func (f *Fault) MarshalText() ([]byte, error) {
	return []byte(f.String()), nil
}

func (f *Fault) UnmarshalText(data []byte) error {
	fmap := map[string]Fault{
		"none":     None,
		"minor":    Minor,
		"major":    Major,
		"critical": Critical,
	}

	fval, ok := fmap[strings.ToLower(string(data))]
	if !ok {
		return ErrFault
	}

	*f = fval

	return nil
}

// Function that deals with lease status.
func (l Lease) String() string {
	strs := map[Lease]string{Leased: "leased", Free: "free", Setup: "setup"}
	lstr, ok := strs[l]

	if !ok {
		return Unknown
	}

	return lstr
}

func (l *Lease) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

func (l *Lease) UnmarshalText(data []byte) error {
	lmap := map[string]Lease{
		"leased": Leased,
		"free":   Free,
		"setup":  Setup,
	}

	lval, ok := lmap[strings.ToLower(string(data))]
	if !ok {
		return ErrFault
	}

	*l = lval

	return nil
}

// Function that deals with state status.
func (s State) String() string {
	strs := map[State]string{Active: "active", Inactive: "inactive"}
	sstr, ok := strs[s]

	if !ok {
		return Unknown
	}

	return sstr
}

func (s *State) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s *State) UnmarshalText(data []byte) error {
	smap := map[string]State{
		"active":   Active,
		"inactive": Inactive,
	}

	sval, ok := smap[strings.ToLower(string(data))]
	if !ok {
		return ErrFault
	}

	*s = sval

	return nil
}

func (s *Status) Validate(ctx context.Context) error {
	if s.Fault > Critical {
		return ErrFault
	}

	if s.Lease > Setup {
		return ErrLease
	}

	if s.State > Inactive {
		return ErrState
	}

	if !s.CreatedTime.Before(time.Now()) {
		return ErrCreatedTime
	}

	return nil
}

// DefaultStatus returns a Status object with starting values (i.e. healthy
// resource in a free state, active, no user, and create time as right now).
func DefaultStatus() *Status {
	return &Status{
		Fault:       None,
		Lease:       Free,
		UsedBy:      "",
		State:       Inactive,
		CreatedTime: time.Now(),
	}
}
