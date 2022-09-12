// Package network provides structs and functions pertaining to network resources.
package network

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/project-safari/zebra"
)

var ErrMaskEmpty = errors.New("mask is nil")

// An IPAddressPool represents a range of consecutive IP addresses belonging
// to the same network.
type IPAddressPool struct {
	zebra.BaseResource
	Subnets []net.IPNet `json:"subnets"`
}

func IPAddressPoolType() zebra.Type {
	return zebra.Type{
		Name:        "network.ipAddressPool",
		Description: "a network ip address pool",
	}
}

func EmptyIPAddressPool() zebra.Resource {
	r := new(IPAddressPool)
	r.Meta.Type = IPAddressPoolType()

	return r
}

// Validate returns an error if the given IPAddressPool object has incorrect values.
// Else, it returns nil.
func (p *IPAddressPool) Validate(ctx context.Context) error {
	if len(p.Subnets) == 0 {
		return ErrIPEmpty
	}

	for _, ip := range p.Subnets {
		if ip.IP == nil {
			return ErrIPEmpty
		} else if ip.Mask == nil {
			return ErrMaskEmpty
		}
	}

	if p.Meta.Type.Name != "network.ipAddressPool" {
		return zebra.ErrWrongType
	}

	return p.BaseResource.Validate(ctx)
}

func NewIPAddressPool(name, owner, group string) *IPAddressPool {
	r := zebra.NewBaseResource(IPAddressPoolType(), name, owner, group)

	return &IPAddressPool{
		BaseResource: *r,
	}
}

func MockIPAddressPool(num int) []zebra.Resource {
	maxByte := 254
	rs := make([]zebra.Resource, 0, num)

	for i := 1; i <= num; i++ {
		r := NewIPAddressPool(
			fmt.Sprintf("mock-ip-pool-%d", i),
			"mocker",
			"ip",
		)

		b := byte(i % maxByte)
		r.Subnets = []net.IPNet{{
			IP:   net.IP{b, b, b, 0},
			Mask: net.IPMask{255, 255, 255, 0},
		}}

		rs = append(rs, r)
	}

	return rs
}
