package zebra

import (
	"context"
	"errors"
)

var ErrorNotLeasable = errors.New(`resource is not leasable`)

// Resource interface is implemented by all resources and provides resource
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

// function to see if the reource  leasable.
func (r *BaseResource) Leasable() error {
	if r.Status.LeaseStatus.CanLease() != nil {
		return ErrorNotLeasable
	}

	return nil
}

// function to see if the reource  is currently available for lease.
func (r *BaseResource) Available() error {
	if r.Status.LeaseStatus.IsFree() != nil {
		return ErrorNotLeasable
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
