package zebra

import (
	"errors"
	"strings"
)

type Status struct {
	Fault       Fault       `json:"fault,omitempty"`
	LeaseStatus LeaseStatus `json:"lease,omitempty"`
	UsedBy      string      `json:"usedBy,omitempty"`
	State       State       `json:"state,omitempty"`
}

type (
	Fault       uint8
	LeaseStatus uint8
	State       uint8
)

const (
	None Fault = iota
	Minor
	Major
	Critical
)

const (
	Free LeaseStatus = iota
	Leased
	Setup
)

const (
	Inactive State = iota
	Active
)

const Unknown = "unknown"

var (
	ErrFault       = errors.New(`fault is incorrect, must be in ["none", "minor", "major", "critical"]`)
	ErrLeaseStatus = errors.New(`lease is incorrect, must be in ["leased", "free", "setup"]`)
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

func (l LeaseStatus) String() string {
	strs := map[LeaseStatus]string{Leased: "leased", Free: "free", Setup: "setup"}
	lstr, ok := strs[l]

	if !ok {
		return Unknown
	}

	return lstr
}

func (l *LeaseStatus) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

func (l *LeaseStatus) UnmarshalText(data []byte) error {
	lmap := map[string]LeaseStatus{
		"leased": Leased,
		"free":   Free,
		"setup":  Setup,
	}

	lval, ok := lmap[strings.ToLower(string(data))]
	if !ok {
		return ErrLeaseStatus
	}

	*l = lval

	return nil
}

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
		return ErrState
	}

	*s = sval

	return nil
}

func (s Status) Validate() error {
	if s.Fault > Critical {
		return ErrFault
	}

	if s.LeaseStatus > Setup {
		return ErrLeaseStatus
	}

	if s.State > Active {
		return ErrState
	}

	return nil
}

func (s *Status) UpdateLeaseState(num int) {
	switch {
	case num == 0:
		s.LeaseStatus = 0
	case num == 1:
		s.LeaseStatus = 1
	default:
		s.LeaseStatus = 2
	}
}

// DefaultStatus returns a Status object with starting values (i.e. healthy
// resource in a free state, active, no user, and create time as right now).
func DefaultStatus() Status {
	return Status{
		Fault:       None,
		LeaseStatus: Free,
		UsedBy:      "",
		State:       Active,
	}
}
