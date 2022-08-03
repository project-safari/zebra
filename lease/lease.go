package lease

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/auth"
	"github.com/project-safari/zebra/status"
)

const DefaultMaxDuration = 4

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
	ActivationTime time.Time      `json:"activationTime"`
}

var (
	ErrLeaseActivate = errors.New("tried to activate lease but request has not been satisfied entirely")
	ErrLeaseValid    = errors.New("lease is not valid")
)

func (r *ResourceReq) Assign(res zebra.Resource) error {
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
func NewLease(owner auth.User, dur time.Duration, req []*ResourceReq) *Lease {
	// Set default values, don't set activation time yet
	l := &Lease{
		lock:           sync.RWMutex{},
		BaseResource:   *zebra.NewBaseResource("Lease", map[string]string{"system.group": "leases"}),
		Duration:       dur,
		Request:        req,
		ActivationTime: time.Time{},
	}
	l.Status.SetOwner(owner.Email)
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
	// Check that lease has been satisfied and activate only then
	// If it's not, throw error
	if !l.IsSatisfied() {
		return ErrLeaseActivate
	}

	l.lock.Lock()
	defer l.lock.Unlock()

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
