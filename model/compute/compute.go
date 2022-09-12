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
		Name:        "compute.server",
		Description: "compute server",
	}
}

func EmptyServer() zebra.Resource {
	s := new(Server)
	s.Meta.Type = ServerType()

	return s
}

// A Server represents a server with credentials, a serial number, board IP, and
// model information.
type Server struct {
	zebra.BaseResource
	Credentials  zebra.Credentials `json:"credentials"`
	SerialNumber string            `json:"serialNumber"`
	BoardIP      net.IP            `json:"boardIp"`
	Model        string            `json:"model"`
}

func NewServer(serial, model, name, owner, group string) *Server {
	return &Server{
		BaseResource: *zebra.NewBaseResource(ServerType(), name, owner, group),
		SerialNumber: serial,
		Model:        model,
	}
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

	if s.Meta.Type.Name != "compute.server" {
		return zebra.ErrWrongType
	}

	if err := s.Credentials.Validate(); err != nil {
		return err
	}

	return s.BaseResource.Validate(ctx)
}

func ESXType() zebra.Type {
	return zebra.Type{
		Name:        "compute.esx",
		Description: "VMWare ESX server",
	}
}

func EmptyESX() zebra.Resource {
	e := new(ESX)
	e.Meta.Type = ESXType()

	return e
}

// An ESX represents an ESX server with credentials, an associated server, and IP.
type ESX struct {
	zebra.BaseResource
	Credentials zebra.Credentials `json:"credentials"`
	ServerID    string            `json:"serverId"`
	IP          net.IP            `json:"ip"`
}

func NewESX(serverID, name, owner, group string) *ESX {
	return &ESX{
		BaseResource: *zebra.NewBaseResource(ESXType(), name, owner, group),
		ServerID:     serverID,
	}
}

func (e *ESX) Validate(ctx context.Context) error {
	if e.IP == nil {
		return ErrIPEmpty
	}

	if e.ServerID == "" {
		return ErrServerIDEmtpy
	}

	if e.Meta.Type.Name != "compute.esx" {
		return zebra.ErrWrongType
	}

	if err := e.Credentials.Validate(); err != nil {
		return err
	}

	return e.BaseResource.Validate(ctx)
}

func VCenterType() zebra.Type {
	return zebra.Type{
		Name:        "compute.vcenter",
		Description: "VMWare vcenter",
	}
}

func EmptyVCenter() zebra.Resource {
	e := new(VCenter)
	e.Meta.Type = VCenterType()

	return e
}

// A VCenter has credentials and an IP.
type VCenter struct {
	zebra.BaseResource
	Credentials zebra.Credentials `json:"credentials"`
	IP          net.IP            `json:"ip"`
}

func NewVCenter(name, owner, group string) *VCenter {
	return &VCenter{
		BaseResource: *zebra.NewBaseResource(VCenterType(), name, owner, group),
	}
}

func (v *VCenter) Validate(ctx context.Context) error {
	if v.IP == nil {
		return ErrIPEmpty
	}

	if v.Meta.Type.Name != "compute.vcenter" {
		return zebra.ErrWrongType
	}

	if err := v.Credentials.Validate(); err != nil {
		return err
	}

	return v.BaseResource.Validate(ctx)
}

func VMType() zebra.Type {
	return zebra.Type{
		Name:        "compute.vm",
		Description: "virtual machine",
	}
}

func EmptyVM() zebra.Resource {
	vm := new(VM)
	vm.Meta.Type = VMType()

	return vm
}

// A VM is represented by a set of credentials, associated ESX ID, management IP,
// and VCenterID.
type VM struct {
	zebra.BaseResource
	Credentials  zebra.Credentials `json:"credentials"`
	ESXID        string            `json:"esxId"`
	ManagementIP net.IP            `json:"managementIp"`
	VCenterID    string            `json:"vCenterId"`
}

func NewVM(esx, name, owner, group string) *VM {
	return &VM{
		BaseResource: *zebra.NewBaseResource(VMType(), name, owner, group),
		ESXID:        esx,
	}
}

func (v *VM) Validate(ctx context.Context) error {
	switch {
	case v.ESXID == "":
		return ErrESXEmpty
	case v.ManagementIP == nil:
		return ErrIPEmpty
	}

	if v.Meta.Type.Name != "compute.vm" {
		return zebra.ErrWrongType
	}

	if err := v.Credentials.Validate(); err != nil {
		return err
	}

	return v.BaseResource.Validate(ctx)
}
