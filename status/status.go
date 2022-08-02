package status

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"
)

type Status struct {
	lock        sync.RWMutex
	fault       FaultState
	lease       LeaseState
	state       ActivityState
	usedBy      string
	createdTime time.Time
}

type (
	FaultState    uint8
	LeaseState    uint8
	ActivityState uint8
)

const (
	None FaultState = iota
	Minor
	Major
	Critical
)

const (
	Leased LeaseState = iota
	Free
	Setup
)

const (
	Active ActivityState = iota
	Inactive
)

const Unknown = "unknown"

var (
	ErrFault          = errors.New(`fault is incorrect, must be in ["none", "minor", "major", "critical"]`)
	ErrLease          = errors.New(`lease is incorrect, must be in ["leased", "free", "setup"]`)
	ErrState          = errors.New(`state is incorrect, must be in ["active", "inactive"]`)
	ErrCreatedTime    = errors.New(`createdTime is incorrect, must be before current time`)
	ErrLeasedResource = errors.New("tried to set LeaseState to Leased but already leased")
	ErrFreeResource   = errors.New("tried to set LeaseState Free Leased but already free")
)

func (f *FaultState) String() string {
	strs := map[FaultState]string{None: "none", Minor: "minor", Major: "major", Critical: "critical"}
	fstr, ok := strs[*f]

	if !ok {
		return Unknown
	}

	return fstr
}

func (f *FaultState) MarshalText() ([]byte, error) {
	return []byte(f.String()), nil
}

func (f *FaultState) UnmarshalText(data []byte) error {
	fmap := map[string]FaultState{
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

func (l *LeaseState) String() string {
	strs := map[LeaseState]string{Leased: "leased", Free: "free", Setup: "setup"}
	lstr, ok := strs[*l]

	if !ok {
		return Unknown
	}

	return lstr
}

func (l *LeaseState) MarshalText() ([]byte, error) {
	return []byte(l.String()), nil
}

func (l *LeaseState) UnmarshalText(data []byte) error {
	lmap := map[string]LeaseState{
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

func (s *ActivityState) String() string {
	strs := map[ActivityState]string{Active: "active", Inactive: "inactive"}
	sstr, ok := strs[*s]

	if !ok {
		return Unknown
	}

	return sstr
}

func (s *ActivityState) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s *ActivityState) UnmarshalText(data []byte) error {
	smap := map[string]ActivityState{
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
	s.lock.RLock()
	defer s.lock.RUnlock()

	if s.fault > Critical {
		return ErrFault
	}

	if s.lease > Setup {
		return ErrLease
	}

	if s.state > Inactive {
		return ErrState
	}

	if !s.createdTime.Before(time.Now()) {
		return ErrCreatedTime
	}

	return nil
}

// NewStatus takes in values and returns a Status object with those set.
func NewStatus(f FaultState, l LeaseState, s ActivityState, u string, t time.Time) *Status {
	return &Status{
		lock:        sync.RWMutex{},
		fault:       f,
		lease:       l,
		state:       s,
		usedBy:      u,
		createdTime: t,
	}
}

// DefaultStatus returns a Status object with starting values (i.e. healthy
// resource in a free state, active, no user, and create time as right now).
func DefaultStatus() *Status {
	return NewStatus(None, Free, Active, "", time.Now())
}

func (s *Status) Fault() FaultState {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.fault
}

func (s *Status) Lease() LeaseState {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.lease
}

func (s *Status) State() ActivityState {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.state
}

func (s *Status) UsedBy() string {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.usedBy
}

func (s *Status) CreatedTime() time.Time {
	s.lock.RLock()
	defer s.lock.RUnlock()

	return s.createdTime
}

func (s *Status) Activate() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.state = Active
}

func (s *Status) Deactivate() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.state = Inactive
}

func (s *Status) SetLeased() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.lease != Free {
		return ErrLeasedResource
	}

	s.lease = Leased

	return nil
}

func (s *Status) SetFree() error {
	s.lock.Lock()
	defer s.lock.Unlock()

	if s.lease != Leased {
		return ErrLeasedResource
	}

	s.lease = Free

	return nil
}

func (s *Status) SetOwner(owner string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.usedBy = owner
}

func (s *Status) MarshalJSON() ([]byte, error) {
	sp := struct {
		Fault       string    `json:"fault"`
		Lease       string    `json:"lease"`
		State       string    `json:"state"`
		UsedBy      string    `json:"usedBy"`
		CreatedTime time.Time `json:"createdTime"`
	}{
		Fault:       s.fault.String(),
		Lease:       s.lease.String(),
		State:       s.state.String(),
		UsedBy:      s.usedBy,
		CreatedTime: s.createdTime,
	}

	return json.Marshal(sp)
}

func (s *Status) UnmarshalJSON(data []byte) error {
	sp := new(struct {
		Fault       FaultState    `json:"fault"`
		Lease       LeaseState    `json:"lease"`
		State       ActivityState `json:"state"`
		UsedBy      string        `json:"usedBy"`
		CreatedTime time.Time     `json:"createdTime"`
	})

	if err := json.Unmarshal(data, sp); err != nil {
		return err
	}

	s.lock = sync.RWMutex{}
	s.fault = sp.Fault
	s.lease = sp.Lease
	s.state = sp.State
	s.usedBy = sp.UsedBy
	s.createdTime = sp.CreatedTime

	return nil
}
