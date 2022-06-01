package zebra

import (
	"context"
	"fmt"
)

// Resource interface is implemented by all resources and provides resource
// validation and label selection methods.
type Resource interface {
	Validate(ctx context.Context) error
}

type BaseResource struct {
	ID string `json:"id"`
}

func (r *BaseResource) Validate(ctx context.Context) error {
	if r.ID == "" {
		return fmt.Errorf("Resource ID is empty")
	}

	return nil
}

type NamedResource struct {
	BaseResource
	Name string `json:"name"`
}

func (r *NamedResource) Validate(ctx context.Context) error {
	if r.Name == "" {
		return fmt.Errorf("Resource Name is empty")
	}

	return r.BaseResource.Validate(ctx)
}
