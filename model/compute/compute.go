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

// Function that returns a zabra type of name server and compute category.
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

// Function that creates a new resource of type server.
//
// It takes in a serial number, a name, an owner, and a group,
// and returns a pointer to Server.
func NewServer(serial, model, name, owner, group string) *Server {
	return &Server{
		BaseResource: *zebra.NewBaseResource(ServerType(), name, owner, group),
		SerialNumber: serial,
		Model:        model,
	}
}

// Function to validate a server resource, given a pointer to the server struct.
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

// Function that returns a zabra type of name esx server and compute category.
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

// Function that creates a new resource of type esx server.
//
// It takes in a server ID, a name, an owner, and a group,
// and returns a pointer to ESX.
func NewESX(serverID, name, owner, group string) *ESX {
	return &ESX{
		BaseResource: *zebra.NewBaseResource(ESXType(), name, owner, group),
		ServerID:     serverID,
	}
}

// Function to validate an esx resource, given a pointer to the esx struct.
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

// Function that returns a zebra type of name vcenter and compute category.
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

// Function that creates a new resource of type vcenter.
//
// It takes in a name, an owner, and a group,
// and returns a pointer to VCenter.
func NewVCenter(name, owner, group string) *VCenter {
	return &VCenter{
		BaseResource: *zebra.NewBaseResource(VCenterType(), name, owner, group),
	}
}

// Function to validate a vcenter resource, given a pointer to the vcenter struct.
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

// Function that returns a zabra type of name vm and compute category.
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

// Function that creates a new resource of type vm.
//
// It takes in an esx ID, a name, an owner, and a group,
// and returns a pointer to VM.
func NewVM(esx, name, owner, group string) *VM {
	return &VM{
		BaseResource: *zebra.NewBaseResource(VMType(), name, owner, group),
		ESXID:        esx,
	}
}

// Function to validate a vm resource, given a pointer to the vm struct.
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
