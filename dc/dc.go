// Package dc provides structs and functions pertaining to datacenter resources.
package dc

import (
	"context"
	"errors"

	"github.com/project-safari/zebra"
)

var ErrAddressEmpty = errors.New("address is empty")

var ErrRowEmpty = errors.New("row is empty")

func DataCenterType() zebra.Type {
	return zebra.Type{
		Name:        "Datacenter",
		Description: "data center",
		Constructor: func() zebra.Resource { return new(Datacenter) },
	}
}

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

	if dc.Type != "Datacenter" {
		return zebra.ErrWrongType
	}

	return dc.NamedResource.Validate(ctx)
}

func LabType() zebra.Type {
	return zebra.Type{
		Name:        "Lab",
		Description: "data center lab",
		Constructor: func() zebra.Resource { return new(Lab) },
	}
}

// A Lab represents the lab consisting of a name and an ID.
type Lab struct {
	zebra.NamedResource
}

func (l *Lab) Validate(ctx context.Context) error {
	if l.Type != "Lab" {
		return zebra.ErrWrongType
	}

	return l.NamedResource.Validate(ctx)
}

func RackType() zebra.Type {
	return zebra.Type{
		Name:        "Rack",
		Description: "server rack",
		Constructor: func() zebra.Resource { return new(Rack) },
	}
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

	if r.Type != "Rack" {
		return zebra.ErrWrongType
	}

	return r.NamedResource.Validate(ctx)
}
