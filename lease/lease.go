package lease

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/status"
)

const DefaultMaxDuration = 4

func Type() zebra.Type {
	return zebra.Type{
		Name:        "Lease",
		Description: "lease request from user",
		Constructor: func() zebra.Resource { return new(Lease) },
	}
}

type ResourceReq struct {
	Type      string           `json:"type"`
	Group     string           `json:"group"`
	Name      string           `json:"name"`
	Count     int              `json:"count"`
	Filters   []zebra.Query    `json:"filters,omitempty"`
	Resources []zebra.Resource `json:"resources,omitempty"`
}

type Lease struct {
	zebra.BaseResource
	lock           sync.RWMutex
	Duration       time.Duration  `json:"duration"`
	Request        []*ResourceReq `json:"request"`
	ActivationTime time.Time      `json:"activationTime,omitempty"`
}

var (
	ErrLeaseActivate = errors.New("tried to activate lease but request has not been satisfied entirely")
	ErrLeaseValid    = errors.New("lease is not valid")
)

func (r *ResourceReq) Assign(res zebra.Resource) error {
	if err := res.GetStatus().SetLeased(); err != nil {
		return err
	}

	if r.Resources == nil {
		r.Resources = make([]zebra.Resource, 0)
	}

	r.Resources = append(r.Resources, res)

	return nil
}

func (r *ResourceReq) IsSatisfied() bool {
	return len(r.Resources) == r.Count
}

// Return a new lease pointer with default values.
func NewLease(email string, dur time.Duration, req []*ResourceReq) *Lease {
	// Set default values, don't set activation time yet
	l := &Lease{
		lock:           sync.RWMutex{},
		BaseResource:   *zebra.NewBaseResource("Lease", map[string]string{"system.group": "leases"}),
		Duration:       dur,
		Request:        req,
		ActivationTime: time.Time{},
	}
	l.Status.SetUser(email)
	l.Status.Deactivate()

	return l
}

// Returns email of user associated with lease.
func (l *Lease) Owner() string {
	l.lock.RLock()
	defer l.lock.RUnlock()

	return l.Status.UsedBy()
}

// Activate lease.
func (l *Lease) Activate() error {
	l.lock.Lock()
	defer l.lock.Unlock()

	if !l.isSatisfied() {
		return ErrLeaseActivate
	}

	l.ActivationTime = time.Now()
	l.Status.Activate()

	return nil
}

// Deactive lease.
func (l *Lease) Deactivate() {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.Status.Deactivate()
}

func (l *Lease) IsSatisfied() bool {
	l.lock.RLock()
	defer l.lock.RUnlock()

	return l.isSatisfied()
}

// Must hold lock when calling this function.
func (l *Lease) isSatisfied() bool {
	for _, r := range l.Request {
		if !r.IsSatisfied() {
			return false
		}
	}

	return true
}

func (l *Lease) IsValid() bool {
	l.lock.RLock()
	defer l.lock.RUnlock()

	// Return if lease has not expired yet
	return time.Now().Before(l.ActivationTime.Add(l.Duration)) && l.Status.State() == status.Active
}

func (l *Lease) IsExpired() bool {
	l.lock.RLock()
	defer l.lock.RUnlock()

	// Return if lease is expired
	return time.Now().After(l.ActivationTime.Add(l.Duration)) || l.Status.State() == status.Inactive
}

func (l *Lease) RequestList() []*ResourceReq {
	l.lock.RLock()
	defer l.lock.RUnlock()

	return l.Request
}

func (l *Lease) Validate(ctx context.Context) error {
	l.lock.RLock()
	defer l.lock.RUnlock()

	if l.Duration.Hours() > DefaultMaxDuration {
		return ErrLeaseValid
	}

	if l.Request == nil {
		return ErrLeaseValid
	}

	if l.ActivationTime.After(time.Now()) {
		return ErrLeaseValid
	}

	return l.BaseResource.Validate(ctx)
}

func (l *Lease) MarshalJSON() ([]byte, error) {
	return json.Marshal(l)
}

func (l *Lease) UnmarshalJSON(data []byte) error {
	lease := &struct {
		zebra.BaseResource
		Duration       time.Duration  `json:"duration"`
		Request        []*ResourceReq `json:"request"`
		ActivationTime time.Time      `json:"activationTime,omitempty"`
	}{}

	if err := json.Unmarshal(data, lease); err != nil {
		return err
	}

	l.BaseResource = lease.BaseResource
	l.Duration = lease.Duration
	l.Request = lease.Request
	l.ActivationTime = lease.ActivationTime
	l.lock = sync.RWMutex{}

	return nil
}
