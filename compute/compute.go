// Package compute provides structs and functions pertaining to compute resources.
package compute

import (
	"context"
	"errors"
	"net"

	"github.com/project-safari/zebra"
)

var ErrSerialEmpty = errors.New("serial number is nil")

var ErrIPEmpty = errors.New("ip address is nil")

var ErrModelEmpty = errors.New("model is empty")

var ErrESXEmpty = errors.New("ESX id is empty")

var ErrVCenterEmpty = errors.New("VCenter id is empty")

var ErrServerIDEmtpy = errors.New("server id is empty")

func ServerType() zebra.Type {
	return zebra.Type{
		Name:        "Server",
		Description: "compute server",
		Constructor: func() zebra.Resource { return new(Server) },
	}
}

// A Server represents a server with credentials, a serial number, board IP, and
// model information.
type Server struct {
	zebra.NamedResource
	Credentials  zebra.Credentials `json:"credentials"`
	SerialNumber string            `json:"serialNumber"`
	BoardIP      net.IP            `json:"boardIP"` //nolint:tagliatelle
	Model        string            `json:"model"`
}

func (s *Server) Validate(ctx context.Context) error {
	switch {
	case s.SerialNumber == "":
		return ErrSerialEmpty
	case s.BoardIP == nil:
		return ErrIPEmpty
	case s.Model == "":
		return ErrModelEmpty
	}

	if err := s.Credentials.Validate(ctx); err != nil {
		return err
	}

	return s.NamedResource.Validate(ctx)
}

func ESXType() zebra.Type {
	return zebra.Type{
		Name:        "ESX",
		Description: "VMWare ESX server",
		Constructor: func() zebra.Resource { return new(ESX) },
	}
}

// An ESX represents an ESX server with credentials, an associated server, and IP.
type ESX struct {
	zebra.NamedResource
	Credentials zebra.Credentials `json:"credentials"`
	ServerID    string            `json:"serverID"` //nolint:tagliatelle
	IP          net.IP            `json:"ip"`
}

func (e *ESX) Validate(ctx context.Context) error {
	if e.IP == nil {
		return ErrIPEmpty
	}

	if e.ServerID == "" {
		return ErrServerIDEmtpy
	}

	if credentialsErr := e.Credentials.Validate(ctx); credentialsErr != nil {
		return credentialsErr
	}

	return e.NamedResource.Validate(ctx)
}

func VCenterType() zebra.Type {
	return zebra.Type{
		Name:        "VCenter",
		Description: "VMWare vcenter",
		Constructor: func() zebra.Resource { return new(VCenter) },
	}
}

// A VCenter has credentials and an IP.
type VCenter struct {
	zebra.NamedResource
	Credentials zebra.Credentials `json:"credentials"`
	IP          net.IP            `json:"ip"`
}

func (v *VCenter) Validate(ctx context.Context) error {
	if v.IP == nil {
		return ErrIPEmpty
	}

	if err := v.Credentials.Validate(ctx); err != nil {
		return err
	}

	return v.NamedResource.Validate(ctx)
}

func VMType() zebra.Type {
	return zebra.Type{
		Name:        "VM",
		Description: "virtual machine",
		Constructor: func() zebra.Resource { return new(VM) },
	}
}

// A VM is represented by a set of credentials, associated ESX ID, management IP,
// and VCenterID.
type VM struct {
	zebra.NamedResource
	Credentials  zebra.Credentials `json:"credentials"`
	ESXID        string            `json:"esxID"`        //nolint:tagliatelle
	ManagementIP net.IP            `json:"managementIP"` //nolint:tagliatelle
	VCenterID    string            `json:"vCenterID"`    //nolint:tagliatelle
}

func (v *VM) Validate(ctx context.Context) error {
	switch {
	case v.ESXID == "":
		return ErrESXEmpty
	case v.ManagementIP == nil:
		return ErrIPEmpty
	case v.VCenterID == "":
		return ErrVCenterEmpty
	}

	if err := v.Credentials.Validate(ctx); err != nil {
		return err
	}

	return v.NamedResource.Validate(ctx)
}
