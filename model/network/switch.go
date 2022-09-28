// Package network provides structs and functions pertaining to network resources.
package network

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/project-safari/zebra"
)

var ErrIPEmpty = errors.New("ip address is nil")

var ErrSerialNumberEmpty = errors.New("serial number is empty")

var ErrModelEmpty = errors.New("model is empty")

var ErrNumPortsEmpty = errors.New("number of ports is 0")

// Function that returns a zabra type of name switch and network category.
func SwitchType() zebra.Type {
	return zebra.Type{
		Name:        "network.switch",
		Description: "network server",
	}
}

func EmptySwitch() zebra.Resource {
	s := new(Switch)
	s.Meta.Type = SwitchType()

	return s
}

// A Switch represents a switching device which has an ID, an associated IP
// address, a serial number, model, and ports.
type Switch struct {
	zebra.BaseResource
	Credentials  zebra.Credentials `json:"credentials"`
	ManagementIP net.IP            `json:"managementIp"`
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

	if s.Meta.Type.Name != "network.switch" {
		return zebra.ErrWrongType
	}

	if err := s.Credentials.Validate(); err != nil {
		return err
	}

	return s.BaseResource.Validate(ctx)
}

// Function that creates a new resource of type switch.
//
// It takes in a name, an owner, and a group,
// and returns a pointer to Switch.
func NewSwitch(name, owner, group string) *Switch {
	r := zebra.NewBaseResource(SwitchType(), name, owner, group)

	return &Switch{
		BaseResource: *r,
	}
}

// Function that generates "mock" switches as sample data.
//
// It takes in the number of resources to generate and returns a list of zebra resources.
func MockSwitch(num int) []zebra.Resource {
	models := []string{"SWITCH-MODEL-1", "SWITCH-MODEL-2", "SWITCH-MODEL-3"}
	rs := make([]zebra.Resource, 0, num)

	for i := 1; i <= num; i++ {
		r := NewSwitch(
			fmt.Sprintf("mock-switch-%d", i),
			"mocker",
			"switch",
		)

		r.Model = models[i%3]
		r.SerialNumber = fmt.Sprintf("SWITCH-SERIAL-%d", i)
		r.ManagementIP = net.IP{14, 14, 14, byte(i)}
		r.Credentials = zebra.NewCredentials("admin")
		_ = r.Credentials.Add("password", fmt.Sprintf("SwitchSecret%d!!!", i))
		r.NumPorts = 100

		rs = append(rs, r)
	}

	return rs
}
