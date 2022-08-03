// Package compute_test tests structs and functions pertaining to compute resources
// outlined in the compute package.
package compute_test

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/cmd/herd/pkg"
	"github.com/project-safari/zebra/compute"
	"github.com/stretchr/testify/assert"
)

const Creds string = "Credentials"

func TestServer(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()

	serverType := compute.ServerType()
	assert.NotNil(serverType)

	server, ok := serverType.New().(*compute.Server)
	assert.True(ok)
	assert.NotNil(server)
	assert.NotNil(server.Validate(ctx))

	server.NamedResource = *zebra.NewNamedResource("there", "Server", nil)
	assert.NotNil(server.Validate(ctx))

	server.SerialNumber = "a"
	assert.NotNil(server.Validate(ctx))

	server.BoardIP = net.ParseIP("10.1.0.0")
	assert.NotNil(server.Validate(ctx))

	server.Model = "b"
	assert.NotNil(server.Validate(ctx))

	server.Credentials = *zebra.NewCredentials("e", map[string]string{"password": "e"}, nil)
	assert.NotNil(server.Validate(ctx))

	server.Credentials.Keys["password"] = "actualPassw0rd%9"
	assert.NotNil(server.Validate(ctx))

	server.Labels = pkg.CreateLabels()
	server.Labels = pkg.GroupLabels(server.Labels, "someServer")

	server.Type = "test"
	assert.NotNil(server.Validate(ctx))
}

func TestESX(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()

	esxType := compute.ESXType()
	esx, ok := esxType.New().(*compute.ESX)
	assert.True(ok)
	assert.NotNil(esx.Validate(ctx))

	esx.NamedResource = *zebra.NewNamedResource("adele", "ESX", nil)
	assert.NotNil(esx.Validate(ctx))

	esx.IP = net.ParseIP("10.1.0.0")
	assert.NotNil(esx.Validate(ctx))

	esx.ServerID = "server id"
	assert.NotNil(esx.Validate(ctx))

	esx.Credentials = *zebra.NewCredentials("k", map[string]string{"password": "m"}, nil)
	assert.NotNil(esx.Validate(ctx))

	esx.Credentials.Keys["password"] = "actualPassw0rd%2"
	assert.NotNil(esx.Validate(ctx))

	esx.Type = "notesx"
	assert.NotNil(esx.Validate(ctx))
}

func TestVCenter(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()

	vcType := compute.VCenterType()
	vcenter, ok := vcType.New().(*compute.VCenter)
	assert.True(ok)
	assert.NotNil(vcenter.Validate(ctx))

	vcenter.NamedResource = *zebra.NewNamedResource("blah", "VCenter", nil)
	assert.NotNil(vcenter.Validate(ctx))

	vcenter.IP = net.ParseIP("10.1.0.0")
	assert.NotNil(vcenter.Validate(ctx))

	vcenter.Credentials = *zebra.NewCredentials("n", map[string]string{"password": "p"}, nil)
	assert.NotNil(vcenter.Validate(ctx))

	vcenter.Credentials.Keys["password"] = "actualPassw0rd%4"
	assert.NotNil(vcenter.Validate(ctx))

	vcenter.Type = "test"
	assert.NotNil(vcenter.Validate(ctx))
}

func TestVM(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()

	vmType := compute.VMType()
	machine, ok := vmType.New().(*compute.VM)
	assert.True(ok)
	assert.NotNil(machine.Validate(ctx))

	machine.NamedResource = *zebra.NewNamedResource("test", "VM", nil)
	assert.NotNil(machine.Validate(ctx))

	machine.ESXID = "q"
	assert.NotNil(machine.Validate(ctx))

	machine.ManagementIP = net.ParseIP("10.1.0.0")
	assert.NotNil(machine.Validate(ctx))

	machine.VCenterID = "r"
	assert.NotNil(machine.Validate(ctx))

	machine.Credentials = *zebra.NewCredentials("s", map[string]string{"password": "u"}, nil)
	assert.NotNil(machine.Validate(ctx))

	fmt.Println(machine.Credentials.Keys)
	machine.Credentials.Keys["password"] = "actualPassw0rd%1"
	assert.NotNil(machine.Validate(ctx))

	machine.Labels = pkg.CreateLabels()
	machine.Labels = pkg.GroupLabels(machine.Labels, "someSampleGroup")

	machine.Type = "machine"
	assert.NotNil(machine.Validate(ctx))
}
