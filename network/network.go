// Package network provides structs and functions pertaining to network resources.
package network

import (
	"context"
	"errors"
	"net"

	"github.com/rchamarthy/zebra"
)

var ErrIPEmpty = errors.New("ip address is nil")

var ErrSerialNumberEmpty = errors.New("serial number is empty")

var ErrModelEmpty = errors.New("model is empty")

var ErrNumPortsEmpty = errors.New("number of ports is 0")

var ErrMaskEmpty = errors.New("mask is nil")

var ErrInvalidRange = errors.New("range bounds are invalid, start is greater than end")

// A Switch represents a switching device which has an ID, an associated IP
// address, a serial number, model, and ports.
type Switch struct {
	zebra.BaseResource
	ManagementIP net.IP `json:"managementIp"`
	SerialNumber string `json:"serialNumber"`
	Model        string `json:"model"`
	NumPorts     uint32 `json:"numPorts"`
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

	return s.BaseResource.Validate(ctx)
}

// An IPAddressPool represents a range of consecutive IP addresses belonging
// to the same network.
type IPAddressPool struct {
	zebra.BaseResource
	net.IPNet
}

// Validate returns an error if the given IPAddressPool object has incorrect values.
// Else, it returns nil.
func (p *IPAddressPool) Validate(ctx context.Context) error {
	if p.IP == nil {
		return ErrIPEmpty
	} else if p.Mask == nil {
		return ErrMaskEmpty
	}

	return p.BaseResource.Validate(ctx)
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

	return v.BaseResource.Validate(ctx)
}
