// Package dc provides structs and functions pertaining to datacenter resources.
package dc

import (
	"context"
	"errors"

	"github.com/rchamarthy/zebra"
)

var ErrAddressEmpty = errors.New("address is empty")

var ErrRowEmpty = errors.New("row is empty")

// A Datacenter represents the physical building. It is a named resource also
// with a building address.
type Datacenter struct {
	zebra.NamedResource
	Address string `json:"address"`
}

// Validate returns an error if the given Datacenter object has incorrect values.
// Else, it returns nil.
func (dc *Datacenter) Validate(ctx context.Context) error {
	if dc.Address == "" {
		return ErrAddressEmpty
	}

	return dc.NamedResource.Validate(ctx)
}

// A Lab represents the lab consisting of a name and an ID.
type Lab struct {
	zebra.NamedResource
}

// A Rack represents a datacenter rack. It consists of a name, ID, and associated
// row.
type Rack struct {
	zebra.NamedResource
	Row string `json:"row"`
}

// Validate returns an error if the given Rack object has incorrect values.
// Else, it returns nil.
func (r *Rack) Validate(ctx context.Context) error {
	if r.Row == "" {
		return ErrRowEmpty
	}

	return r.NamedResource.Validate(ctx)
}
