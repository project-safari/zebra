package zebra

import (
	"context"
	"errors"
)

// Resource interface is implemented by all resources and provides resource
// validation and label selection methods.
type Resource interface {
	Validate(ctx context.Context) error
}

var ErrNameEmpty = errors.New("name is empty")

var ErrIDEmpty = errors.New("id is empty")

type BaseResource struct {
	ID string `json:"id"`
}

func (r *BaseResource) Validate(ctx context.Context) error {
	if r.ID == "" {
		return ErrIDEmpty
	}

	return nil
}

type NamedResource struct {
	BaseResource
	Name string `json:"name"`
}

func (r *NamedResource) Validate(ctx context.Context) error {
	if r.Name == "" {
		return ErrNameEmpty
	}

	return r.BaseResource.Validate(ctx)
}
