package zebra

import (
	"context"
)

// Resource interface is implemented by all resources and provides resource
// validation and label selection methods.
type Resource interface {
	Validate(ctx context.Context) error
	GetMeta() Meta
	GetStatus() Status
}

// A BaseResource struct represents a basic resource  with the appropriate meta data and a status.
type BaseResource struct {
	Meta   Meta   `json:"meta"`
	Status Status `json:"status,omitempty"`
}

// Function that calidates a basic resource (BaseResource).
// It returns an error or nil in the absence thereof.
func (r *BaseResource) Validate(ctx context.Context) error {
	if err := r.Meta.Validate(); err != nil {
		return err
	}

	if err := r.Status.Validate(); err != nil {
		return err
	}

	return nil
}

// Function on a pointer to BaseResource to get the meta data from the given resource.
func (r *BaseResource) GetMeta() Meta {
	return r.Meta
}

// Function on a pointer to BaseResource to get the status of the given resource.
func (r *BaseResource) GetStatus() Status {
	return r.Status
}

// Function that returns a pointer to a BaseResource struct
// with meta data (type, name, group, owner) for the resource and its status.
func NewBaseResource(rType Type, name, owner, group string) *BaseResource {
	return &BaseResource{
		Meta:   NewMeta(rType, name, group, owner),
		Status: DefaultStatus(),
	}
}
