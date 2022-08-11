// Package compute_test tests structs and functions pertaining to compute resources
// outlined in the compute package.
package compute_test

import (
	"context"
	"net"
	"testing"

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
	assert.NotNil(machine.Validate(ctx))

	machine.Labels = pkg.CreateLabels()
	machine.Labels = pkg.GroupLabels(machine.Labels, "someSampleGroup")

	machine.Type = "machine"
	assert.NotNil(machine.Validate(ctx))
}

// test for vcenter generator.
func TestNewVCenter(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	labels := pkg.CreateLabels()

	labels = pkg.GroupLabels(labels, "group")

	vc := compute.NewVCenter(pkg.Name(), net.IP("123.111.001"), labels)

	assert.NotNil(vc)
}

func TestNewServer(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	labels := pkg.CreateLabels()

	labels = pkg.GroupLabels(labels, "group")

	arr := []string{pkg.Serials(), pkg.Models(), pkg.Name()}
	srv := compute.NewServer(arr, net.IP("123.111.001"), labels)

	assert.NotNil(srv)
}

func TestNewESX(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	labels := pkg.CreateLabels()

	labels = pkg.GroupLabels(labels, "group")

	esx := compute.NewESX(pkg.Name(), pkg.SelectServer(), net.IP("123.111.001"), labels)

	assert.NotNil(esx)
}

func TestNewVM(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	labels := pkg.CreateLabels()

	labels = pkg.GroupLabels(labels, "group")

	arr := []string{pkg.Name(), pkg.SelectESX(), pkg.SelectVcenter()}
	vm := compute.NewVM(arr, net.IP("123.111.001"), labels)

	assert.NotNil(vm)
}
