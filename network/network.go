// Package network provides structs and functions pertaining to network resources.
package network

import (
	"context"
	"errors"
	"net"

	"github.com/project-safari/zebra"
)

var ErrIPEmpty = errors.New("ip address is nil")

var ErrSerialNumberEmpty = errors.New("serial number is empty")

var ErrModelEmpty = errors.New("model is empty")

var ErrNumPortsEmpty = errors.New("number of ports is 0")

var ErrMaskEmpty = errors.New("mask is nil")

var ErrInvalidRange = errors.New("range bounds are invalid, start is greater than end")

func SwitchType() zebra.Type {
	return zebra.Type{
		Name:        "Switch",
		Description: "network server",
		Constructor: func() zebra.Resource { return new(Switch) },
	}
}

// A Switch represents a switching device which has an ID, an associated IP
// address, a serial number, model, and ports.
type Switch struct {
	zebra.BaseResource
	Credentials  zebra.Credentials `json:"credentials"`
	ManagementIP net.IP            `json:"managementIP"` //nolint:tagliatelle
	SerialNumber string            `json:"serialNumber"`
	Model        string            `json:"model"`
	NumPorts     uint32            `json:"numPorts"`
}

// Validate returns an error if the given Switch object has incorrect values.
// Else, it returns nil.
func (s *Switch) Validate(ctx context.Context) error {
	switch {
	case s.ManagementIP == nil:
		return ErrIPEmpty
	case s.SerialNumber == "":
		return ErrSerialNumberEmpty
	case s.Model == "":
		return ErrModelEmpty
	case s.NumPorts == 0:
		return ErrNumPortsEmpty
	}

	if s.Type != "Switch" {
		return zebra.ErrWrongType
	}

	if err := s.Credentials.Validate(ctx); err != nil {
		return err
	}

	return s.BaseResource.Validate(ctx)
}

func IPAddressPoolType() zebra.Type {
	return zebra.Type{
		Name:        "IPAddressPool",
		Description: "ip address pool",
		Constructor: func() zebra.Resource { return new(IPAddressPool) },
	}
}

// An IPAddressPool represents a range of consecutive IP addresses belonging
// to the same network.
type IPAddressPool struct {
	zebra.BaseResource
	Subnets []net.IPNet `json:"subnets"`
}

// Validate returns an error if the given IPAddressPool object has incorrect values.
// Else, it returns nil.
func (p *IPAddressPool) Validate(ctx context.Context) error {
	for _, ip := range p.Subnets {
		if ip.IP == nil {
			return ErrIPEmpty
		} else if ip.Mask == nil {
			return ErrMaskEmpty
		}
	}

	if p.Type != "IPAddressPool" {
		return zebra.ErrWrongType
	}

	return p.BaseResource.Validate(ctx)
}

func VLANPoolType() zebra.Type {
	return zebra.Type{
		Name:        "VLANPool",
		Description: "vlan pool",
		Constructor: func() zebra.Resource { return new(VLANPool) },
	}
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

	if v.Type != "VLANPool" {
		return zebra.ErrWrongType
	}

	return v.BaseResource.Validate(ctx)
}

func NewVlanPool(start uint16, end uint16, labels zebra.Labels) *VLANPool {
	theRes := zebra.NewBaseResource("VlanPool", labels)
	ret := &VLANPool{
		BaseResource: *theRes,
		RangeStart:   start,
		RangeEnd:     end,
	}

	return ret
}
