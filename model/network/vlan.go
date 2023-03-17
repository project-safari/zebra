// Package network provides structs and functions pertaining to network resources.
package network

import (
	"context"
	"errors"
	"fmt"

	"github.com/project-safari/zebra"
)

var ErrInvalidRange = errors.New("range bounds are invalid, start is greater than end")

// Function that returns a zabra type of name VlanPool and network category.
func VLANPoolType() zebra.Type {
	return zebra.Type{
		Name:        "network.vlanPool",
		Description: "network vlan pool",
	}
}

func EmptyVLANPool() zebra.Resource {
	r := new(VLANPool)
	r.Meta.Type = VLANPoolType()

	return r
}

// A VLANPool represents a pool of VLANs belonging to the same network.
type VLANPool struct {
	zebra.BaseResource
	RangeStart uint16 `json:"rangeStart"`
	RangeEnd   uint16 `json:"rangeEnd"`
}

// Validate returns an error if the given VLANPool object has incorrect values.
// Else, it returns nil.
func (v *VLANPool) Validate(ctx context.Context) error {
	if v.RangeStart > v.RangeEnd {
		return ErrInvalidRange
	}

	if v.Meta.Type.Name != "network.vlanPool" {
		return zebra.ErrWrongType
	}

	return v.BaseResource.Validate(ctx)
}

func (v *VLANPool) String() string {
	return fmt.Sprintf("%d-%d", v.RangeStart, v.RangeEnd)
}

// Function that creates a new resource of type vlan.
//
// It takes in a name, an owner, and a group,
// and returns a pointer to VLANPool.
func NewVLANPool(name, owner, group string) *VLANPool {
	r := zebra.NewBaseResource(VLANPoolType(), name, owner, group)

	return &VLANPool{
		BaseResource: *r,
	}
}

// Function that generates "mock" VLANPools as sample data.
//
// It takes in the number of resources to generate and returns a list of zebra resources.
func MockVLANPool(num int) []zebra.Resource {
	rs := make([]zebra.Resource, 0, num)
	prevStart := uint16(1)
	size := 9

	for i := 1; i <= num; i++ {
		r := NewVLANPool(
			fmt.Sprintf("mock-vlan-pool-%d", i),
			"mocker",
			"ip",
		)

		r.RangeStart = prevStart
		r.RangeEnd = prevStart + uint16(size)
		prevStart += 10

		rs = append(rs, r)
	}

	return rs
}
