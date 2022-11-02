package zebra

import (
	"errors"
	"strings"
)

// The Status struct is a nested struct that contains substructs with info about a lease.
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
	// ErrFault  occurs if something is wrong to the fault.
	ErrFault = errors.New(`fault is incorrect, must be in ["none", "minor", "major", "critical"]`)
	// ErrLeaseStatus occurs if something is wrong with the status of the lease.
	ErrLeaseStatus = errors.New(`lease is incorrect, must be in ["leased", "free", "setup"]`)
	// ErrState occurs if something is wrong with the state.
	ErrState = errors.New(`state is incorrect, must be in ["active", "inactive"]`)
	// ErrCreatedTime occurs if the creation time is greater than the current time.
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

// Function on a pointer to the fault struct to marshal a Fault.
//
// It returns a byte array and an error or nil, in the absence thereof.
func (f *Fault) MarshalText() ([]byte, error) {
	return []byte(f.String()), nil
}

// Function on a pointer to the fault struct to unMarshal a Fault.
//
// It takes in a byte array.
// It returns an error or nil, in the absence thereof.
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

// Funciton that returns a lease's status as a string.
func (l LeaseStatus) String() string {
	strs := map[LeaseStatus]string{Leased: "leased", Free: "free", Setup: "setup"}
	lstr, ok := strs[l]

	if !ok {
		return Unknown
	}

	return lstr
}

// Function on a pointer to the leaseStatus struct to marshal a leaseStatus.
//
// It returns a byte array and an error or nil, in the absence thereof.
func (l *LeaseStatus) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

// Function on a pointer to the leaseStatus struct to unMarshal a LeaseStatus.
//
// It takes in a byte array.
// It returns an error or nil, in the absence thereof.
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

// Funciton that returns a lease's state as a string.
func (s State) String() string {
	strs := map[State]string{Active: "active", Inactive: "inactive"}
	sstr, ok := strs[s]

	if !ok {
		return Unknown
	}

	return sstr
}

// Function on a pointer to the status struct to marshal a status.
//
// It returns a byte array and an error or nil, in the absence thereof.
func (s *State) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

// Function on a pointer to the status struct to unMarshal a status.
//
// It takes in a byte array.
// It returns an error or nil, in the absence thereof.
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

// Function on an instance of the status struct.
//
// It validates a status and returns an error or nil, in the absence thereof.
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

// DefaultStatus returns a Status object with starting values (i.e. healthy
// resource in a free state, active, no user, and create time as right now).
func DefaultStatus() Status {
	return Status{
		Fault:       None,
		LeaseStatus: Free,
		UsedBy:      "",
		State:       Inactive,
	}
}
