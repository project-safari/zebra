package lease

import (
	"errors"
	"sync"
	"time"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/auth"
	"github.com/project-safari/zebra/store"
)

type Lease struct {
	lock sync.RWMutex
	zebra.BaseResource
	duration       time.Duration
	request        map[*zebra.Group]*ResourceReq
	resources      *zebra.ResourceMap
	activationTime time.Time
	progress       map[*zebra.Group]bool
}
type ResourceReq struct {
	Type  string `json:"type"`
	Group string `json:"group"`
	Name  string `json:"name"`
	Count uint   `json:"count"`
}

var ErrLeaseActivate = errors.New("tried to activate lease but request has not been satisfied entirely")

// Return a new lease pointer with default values.
func NewLease(owner auth.User, dur time.Duration, req map[*zebra.Group]*ResourceReq) *Lease {
	// Set default values, don't set activation time yet
	l := &Lease{
		lock:           sync.RWMutex{},
		BaseResource:   *zebra.NewBaseResource("Lease", nil),
		duration:       dur,
		request:        req,
		resources:      zebra.NewResourceMap(store.DefaultFactory()),
		activationTime: time.Time{},
		progress: func() map[*zebra.Group]bool {
			m := make(map[*zebra.Group]bool, len(req))
			for k := range req {
				m[k] = false
			}

			return m
		}(),
	}
	l.Status.CreatedTime = time.Now()
	l.Status.UsedBy = owner.Email
	l.Status.State = zebra.Inactive

	return l
}

// Returns email of user associated with lease.
func (l *Lease) Owner() string {
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

	l.activationTime = time.Now()
	l.Status.State = zebra.Active

	return nil
}

// Deactive lease.
func (l *Lease) Deactivate() {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.Status.State = zebra.Inactive
}

// Must be called with the lock.
func (l *Lease) IsSatisfied() bool {
	l.lock.RLock()
	defer l.lock.RUnlock()

	for _, b := range l.progress {
		if !b {
			return false
		}
	}

	return true
}

func (l *Lease) IsValid() bool {
	l.lock.RLock()
	defer l.lock.RUnlock()

	// Return if lease has not expired yet
	return time.Now().Before(l.activationTime.Add(l.duration)) && l.Status.State == zebra.Active
}

func (l *Lease) IsExpired() bool {
	l.lock.RLock()
	defer l.lock.RUnlock()

	// Return if lease is expired
	return time.Now().After(l.activationTime.Add(l.duration)) || l.Status.State == zebra.Inactive
}

func (l *Lease) Assign(res zebra.Resource) {
	l.lock.RLock()
	defer l.lock.RUnlock()

	res.GetStatus().Lease = zebra.Leased

	l.resources.Add(res, res.GetType())
}

func (l *Lease) Finish(g *zebra.Group) {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.progress[g] = true
}

func (l *Lease) Request() map[*zebra.Group]*ResourceReq {
	l.lock.RLock()
	defer l.lock.RUnlock()

	return l.request
}

/*
// Return estimated wait time for given lease.
func (l *Lease) WaitTime() time.Time {
	// Return longest wait time for any given queue
}
*/
