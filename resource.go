package zebra

import (
	"context"
	"errors"
	"time"
)

// Errors related to a resource's ability to be leased.
// A resource may be temporarily unavailable, or not leasable.
var (
	ErrNotLeasable  = errors.New(`resource is not leasable`)
	ErrNotAvailable = errors.New(`resource is not currntly available to be leased`)
	ErrTooLong      = errors.New(`lease request took too long`)
)

// Resource interface is implemented by all resources and provides resource,
// validation and label selection methods.
type Resource interface {
	Validate(ctx context.Context) error
	GetMeta() Meta
	GetStatus() Status
}

type BaseResource struct {
	Meta   Meta   `json:"meta"`
	Status Status `json:"status,omitempty"`
}

// //

// Function to change the status.
func (r *BaseResource) ChangeStatus(value string) {
	if r.Leasable() == nil {
		if r.Expired(3) != nil {
			r.Status.State = Inactive

		}

		switch value {
		case "leased":
			r.Status.LeaseStatus = 1
			r.Status.State = State(Leased)

		case "freed":
			r.Status.LeaseStatus = 0
			r.Status.State = Active

		}
	} else {
		r.Status.LeaseStatus = 2
	}
}

// //
// Function to time the duration for lease satisfaction.
func (r *BaseResource) Expired(maxTime int) error {
	now := time.Now()
	timeDifference := now.Sub(r.Meta.ModificationTime)

	if timeDifference > time.Duration(maxTime) {
		return ErrTooLong
	}

	return nil
}

func (r *BaseResource) Validate(ctx context.Context) error {
	if err := r.Meta.Validate(); err != nil {
		return err
	}

	if err := r.Status.Validate(); err != nil {
		return err
	}

	return nil
}

// ///
// Function to see if the resource  leasable.
func (r *BaseResource) Leasable() error {
	if r.Status.LeaseStatus.CanLease() != nil {
		return ErrNotLeasable
	}

	return nil
}

// ///
// Function to see if the reource  is currently available for lease.
func (r *BaseResource) Available() error {
	if r.Status.LeaseStatus.IsFree() != nil {
		return ErrNotAvailable
	}

	return nil
}

func (r *BaseResource) GetMeta() Meta {
	return r.Meta
}

func (r *BaseResource) GetStatus() Status {
	return r.Status
}

func NewBaseResource(rType Type, name, owner, group string) *BaseResource {
	return &BaseResource{
		Meta:   NewMeta(rType, name, group, owner),
		Status: DefaultStatus(),
	}
}
