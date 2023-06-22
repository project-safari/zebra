package zebra

import (
	"context"
	"errors"
)

// Errors related to a resource's ability to be leased.
// A resource may be temporarily unavailable, or not leasable.
var (
	ErrNotLeasable  = errors.New(`resource is not leasable`)
	ErrNotAvailable = errors.New(`resource is not currntly available to be leased`)
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

func (r *BaseResource) Validate(ctx context.Context) error {
	if err := r.Meta.Validate(); err != nil {
		return err
	}

	if err := r.Status.Validate(); err != nil {
		return err
	}

	return nil
}

// Function to see if the reource  leasable.
func (r *BaseResource) Leasable() error {
	if r.Status.LeaseStatus.CanLease() != nil {
		return ErrNotLeasable
	}

	return nil
}

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
