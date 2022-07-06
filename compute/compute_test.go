// Package compute_test tests structs and functions pertaining to compute resources
// outlined in the compute package.
package compute_test

import (
	"context"
	"net"
	"testing"

	"github.com/project-safari/zebra/compute"
	"github.com/stretchr/testify/assert"
)

const Creds string = "Credentials"

func TestServer(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()

	server := new(compute.Server)
	assert.NotNil(server.Validate(ctx))

	server.ID = "hello"
	server.Type = "Server"
	server.Name = "there"
	assert.NotNil(server.Validate(ctx))

	server.SerialNumber = "a"
	assert.NotNil(server.Validate(ctx))

	server.BoardIP = net.ParseIP("10.1.0.0")
	assert.NotNil(server.Validate(ctx))

	server.Model = "b"
	assert.NotNil(server.Validate(ctx))

	server.Credentials.Name = "c"
	server.Credentials.Type = Creds
	server.Credentials.ID = "dddd"
	server.Credentials.Keys = make(map[string]string)
	server.Credentials.Keys["password"] = "e"
	assert.NotNil(server.Validate(ctx))

	server.Credentials.Keys["password"] = "actualPassw0rd%9"
	assert.Nil(server.Validate(ctx))
}

func TestESX(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()

	esx := new(compute.ESX)
	assert.NotNil(esx.Validate(ctx))

	esx.ID = "rolling in the deep"
	esx.Type = "ESX"
	esx.Name = "adele"
	assert.NotNil(esx.Validate(ctx))

	esx.IP = net.ParseIP("10.1.0.0")
	assert.NotNil(esx.Validate(ctx))

	esx.ServerID = "server id"
	assert.NotNil(esx.Validate(ctx))

	esx.Credentials.Name = "k"
	esx.Credentials.ID = "lllll"
	esx.Credentials.Type = Creds
	esx.Credentials.Keys = make(map[string]string)
	esx.Credentials.Keys["password"] = "m"
	assert.NotNil(esx.Validate(ctx))

	esx.Credentials.Keys["password"] = "actualPassw0rd%2"
	assert.Nil(esx.Validate(ctx))
}

func TestVCenter(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()

	vcenter := new(compute.VCenter)
	assert.NotNil(vcenter.Validate(ctx))

	vcenter.ID = "blah"
	vcenter.Type = "VCenter"
	vcenter.Name = "blahblah"
	assert.NotNil(vcenter.Validate(ctx))

	vcenter.IP = net.ParseIP("10.1.0.0")
	assert.NotNil(vcenter.Validate(ctx))

	vcenter.Credentials.Name = "n"
	vcenter.Credentials.ID = "oooo"
	vcenter.Credentials.Type = Creds
	vcenter.Credentials.Keys = make(map[string]string)
	vcenter.Credentials.Keys["password"] = "p"
	assert.NotNil(vcenter.Validate(ctx))

	vcenter.Credentials.Keys["password"] = "actualPassw0rd%4"
	assert.Nil(vcenter.Validate(ctx))
}

func TestVM(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()

	machine := new(compute.VM)
	assert.NotNil(machine.Validate(ctx))

	machine.ID = "can you hear me"
	machine.Type = "VM"
	machine.Name = "you'd like to meet"
	assert.NotNil(machine.Validate(ctx))

	machine.ESXID = "q"
	assert.NotNil(machine.Validate(ctx))

	machine.ManagementIP = net.ParseIP("10.1.0.0")
	assert.NotNil(machine.Validate(ctx))

	machine.VCenterID = "r"
	assert.NotNil(machine.Validate(ctx))

	machine.Credentials.Name = "s"
	machine.Credentials.Type = Creds
	machine.Credentials.ID = "tttt"
	machine.Credentials.Keys = make(map[string]string)
	machine.Credentials.Keys["password"] = "u"
	assert.NotNil(machine.Validate(ctx))

	machine.Credentials.Keys["password"] = "actualPassw0rd%1"
	assert.Nil(machine.Validate(ctx))
}
