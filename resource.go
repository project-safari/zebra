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
// should a status change be needed, call this function and
// pass a string value for the event that triggers the chenge.
func (r *BaseResource) ChangeStatus(value string) {
	if r.Leasable() == nil {
		if r.Expired(3) != nil {
			r.Status.State = Inactive

		}

		switch value {
		case "leased":
			r.Status.LeaseStatus = 1
			r.Status.State = Active

		case "freed":
			r.Status.LeaseStatus = 0
			r.Status.State = Inactive

		}
	} else {
		r.Status.LeaseStatus = 2
	}

	// May want to add more cases:
	//
	// a resource is in the process of being reserved --> can add a pending status/state.
	//
	// a resource is defficient --> can add a state to represent this (could also split into multiple).
}

/*
func (r *BaseResource) StatusEvent() string {

}
*/

// //
// Function to time the duration for lease satisfaction.
// This function should be used if the lease has/has not been satisfied
// within an alloted time frame.
func (r *BaseResource) Expired(maxTime int) error {
	now := time.Now()
	timeDifference := now.Sub(r.Meta.ModificationTime)

	if time.Duration(timeDifference) > time.Duration(maxTime) {
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
// Function to see if the resource is leasable.
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
