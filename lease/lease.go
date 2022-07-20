package lease

import (
	"errors"
	"time"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/auth"
	"github.com/project-safari/zebra/store"
)

type Lease struct {
	Owner          auth.User           `json:"owner"`
	Request        Request             `json:"request"`
	Resources      *zebra.ResourceMap  `json:"resources"`
	RequestTime    time.Time           `json:"requestTime"`
	ActivationTime time.Time           `json:"activationTime"`
	Status         zebra.ActivityState `json:"status"`
}

type Request struct {
	Duration  time.Duration `json:"duration"`
	Group     string        `json:"group"`
	Resources []ResourceReq `json:"resources"`
}

type ResourceReq struct {
	Type  string `json:"type"`
	Name  string `json:"name"`
	Count uint   `json:"count"`
}

var ErrLeaseActivate = errors.New("tried to activate lease but resources are already active")

// Return a new lease pointer with default values.
func NewLease(owner auth.User, req Request) *Lease {
	// Set default values, don't set activation time yet
	return &Lease{
		Owner:          owner,
		Request:        req,
		Resources:      zebra.NewResourceMap(store.DefaultFactory()),
		RequestTime:    time.Now(),
		ActivationTime: time.Time{},
		Status:         zebra.Inactive,
	}
}

// Activate lease.
func (l *Lease) Activate() error {
	// Check that all resources are indeed free, then change to allocated state
	for _, l := range l.Resources.Resources {
		for _, r := range l.Resources {
			if r.GetStatus().Lease == zebra.Leased {
				return ErrLeaseActivate
			}

			r.GetStatus().State = zebra.Active
		}
	}

	l.ActivationTime = time.Now()
	l.Status = zebra.Active

	return nil
}

// Deactive lease.
func (l *Lease) Deactivate() {
	l.Status = zebra.Inactive
}

/*
// Return estimated wait time for given lease.
func (l *Lease) WaitTime() time.Time {

}
*/
