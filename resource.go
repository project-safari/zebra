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

// BaseResource must be embedded in all resource structs, ensuring each resource is
// assigned an ID string.
type BaseResource struct {
	ID string `json:"id"`
}

// Validate returns an error if the given BaseResource object has incorrect values.
// Else, it returns nil.
func (r *BaseResource) Validate(ctx context.Context) error {
	if r.ID == "" {
		return ErrIDEmpty
	}

	return nil
}

// NamedResource is implemented by all resources assigned both a string ID and
// a name.
type NamedResource struct {
	BaseResource
	Name string `json:"name"`
}

// Validate returns an error if the given NamedResource object has incorrect values.
// Else, it returns nil.
func (r *NamedResource) Validate(ctx context.Context) error {
	if r.Name == "" {
		return ErrNameEmpty
	}

	return r.BaseResource.Validate(ctx)
}
