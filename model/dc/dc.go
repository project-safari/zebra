// Package dc provides structs and functions pertaining to datacenter resources.
package dc

import (
	"context"
	"errors"

	"github.com/project-safari/zebra"
)

var ErrAddressEmpty = errors.New("address is empty")

var ErrRowEmpty = errors.New("row is empty")

func DatacenterType() zebra.Type {
	return zebra.Type{
		Name:        "dc.datacenter",
		Description: "data center",
	}
}

func EmptyDatacenter() zebra.Resource {
	d := new(Datacenter)
	d.Meta.Type = DatacenterType()

	return d
}

// A Datacenter represents the physical building. It is a named resource also
// with a building address.
type Datacenter struct {
	zebra.BaseResource
	Address string `json:"address"`
}

// create new dc resources.
func NewDatacenter(address, name, owner, group string) *Datacenter {
	return &Datacenter{
		BaseResource: *zebra.NewBaseResource(DatacenterType(), name, owner, group),
		Address:      address,
	}
}

// Validate returns an error if the given Datacenter object has incorrect values.
// Else, it returns nil.
func (dc *Datacenter) Validate(ctx context.Context) error {
	if dc.Address == "" {
		return ErrAddressEmpty
	}

	if dc.Meta.Type.Name != "dc.datacenter" {
		return zebra.ErrWrongType
	}

	return dc.BaseResource.Validate(ctx)
}

func LabType() zebra.Type {
	return zebra.Type{
		Name:        "dc.lab",
		Description: "data center lab",
	}
}

func EmptyLab() zebra.Resource {
	l := new(Lab)
	l.Meta.Type = LabType()

	return l
}

// A Lab represents the lab consisting of a name and an ID.
type Lab struct{ zebra.BaseResource }

// create new dc resources.
func NewLab(name, owner, group string) *Lab {
	return &Lab{
		BaseResource: *zebra.NewBaseResource(LabType(), name, owner, group),
	}
}

func (l *Lab) Validate(ctx context.Context) error {
	if l.Meta.Type.Name != "dc.lab" {
		return zebra.ErrWrongType
	}

	return l.BaseResource.Validate(ctx)
}

func RackType() zebra.Type {
	return zebra.Type{
		Name:        "dc.rack",
		Description: "server rack",
	}
}

func EmptyRack() zebra.Resource {
	r := new(Rack)
	r.Meta.Type = RackType()

	return r
}

// A Rack represents a datacenter rack. It consists of a name, ID, and associated
// row.
type Rack struct {
	zebra.BaseResource
	Row string `json:"row"`
}

// Validate returns an error if the given Rack object has incorrect values.
// Else, it returns nil.
func (r *Rack) Validate(ctx context.Context) error {
	if r.Row == "" {
		return ErrRowEmpty
	}

	if r.Meta.Type.Name != "dc.rack" {
		return zebra.ErrWrongType
	}

	return r.BaseResource.Validate(ctx)
}

func NewRack(row, name, owner, group string) *Rack {
	return &Rack{
		BaseResource: *zebra.NewBaseResource(RackType(), name, owner, group),
		Row:          row,
	}
}
