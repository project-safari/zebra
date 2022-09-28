package lease

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/project-safari/zebra"
)

// Function that returns a zabra type of a lease request.
func Type() zebra.Type {
	return zebra.Type{
		Name:        "system.lease",
		Description: "lease request from user",
	}
}

func Empty() zebra.Resource {
	l := new(Lease)
	l.Meta.Type = Type()

	return l
}

// A ResourceReq struct represents a Request for a lease with type, an associated group, a name,
// a count - number - of requested resources, and two arrays - for filters and for resources.
type ResourceReq struct {
	Type      string           `json:"type"`
	Group     string           `json:"group"`
	Name      string           `json:"name"`
	Count     int              `json:"count"`
	Filters   []zebra.Query    `json:"filters,omitempty"`
	Resources []zebra.Resource `json:"resources,omitempty"`
}

// An ESX struct represents a Lease for a certain resource,
// with its zebra.BaseResource data, a mutex lock, a duration for the lease,
// a pointer to the request and an associated activation time.
type Lease struct {
	zebra.BaseResource
	lock           sync.RWMutex
	Duration       time.Duration  `json:"duration"`
	Request        []*ResourceReq `json:"request"`
	ActivationTime time.Time      `json:"activationTime"`
}

// Errors that can occur with leases.
var (
	// ErrLeaseActivate occurs when the activation of a lease failed.
	ErrLeaseActivate = errors.New("tried to activate lease but request has not been satisfied entirely")
	// ErrLeaseValid occurs when a lease is not valid.
	ErrLeaseValid = errors.New("lease is not valid")
)

// Function on a pointer to ResourceReq.
// The function appends resources to a lease request.
//
// It takes in a zebra.Resource and returns an error or nil in the absence thereof.
func (r *ResourceReq) Assign(res zebra.Resource) error {
	if r.Resources == nil {
		r.Resources = make([]zebra.Resource, 0)
	}

	r.Resources = append(r.Resources, res)

	return nil
}

/// Function on a pointer to ResourceReq.
// The function verifies if a request was satisfied.
//
// It returns a boolean value.
func (r *ResourceReq) IsSatisfied() bool {
	return len(r.Resources) == r.Count
}

// Return a new lease pointer with default values.
func NewLease(userEmail string, dur time.Duration, req []*ResourceReq) *Lease {
	// Set default values, don't set activation time yet
	l := &Lease{
		lock:           sync.RWMutex{},
		BaseResource:   *zebra.NewBaseResource(Type(), "", userEmail, "system.leases"),
		Duration:       dur,
		Request:        req,
		ActivationTime: time.Time{},
	}
	l.Status.UsedBy = userEmail
	l.Status.State = zebra.Inactive

	return l
}

// Returns email of user associated with lease.
func (l *Lease) Owner() string {
	l.lock.RLock()
	defer l.lock.RUnlock()

	return l.Status.UsedBy
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
	l.Status.State = zebra.Active

	return nil
}

// Deactive lease.
func (l *Lease) Deactivate() {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.Status.State = zebra.Inactive
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

// Function on a pointer to Lease.
// The function checks if a lease is valid.
//
// It returns a boolean value.
func (l *Lease) IsValid() bool {
	l.lock.RLock()
	defer l.lock.RUnlock()

	// Return if lease has not expired yet
	return time.Now().Before(l.ActivationTime.Add(l.Duration)) && l.Status.State == zebra.Active
}

// Function on a pointer to Lease.
// The function checks if a lease expired.
//
// It returns a boolean value.
func (l *Lease) IsExpired() bool {
	l.lock.RLock()
	defer l.lock.RUnlock()

	// Return if lease is expired
	return time.Now().After(l.ActivationTime.Add(l.Duration)) || l.Status.State == zebra.Inactive
}

// Function on a pointer to ResourceReq.
// It returns a pointer to the lease's respective ResourceReq.
func (l *Lease) RequestList() []*ResourceReq {
	l.lock.RLock()
	defer l.lock.RUnlock()

	return l.Request
}

// Validate funtcion for Lease.
func (l *Lease) Validate(ctx context.Context) error {
	if l.Duration.Hours() > zebra.DefaultMaxDuration {
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
