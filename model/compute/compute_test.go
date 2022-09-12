// Package compute_test tests structs and functions pertaining to compute resources
// outlined in the compute package.
package compute_test

import (
	"context"
	"net"
	"testing"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/model/compute"
	"github.com/stretchr/testify/assert"
)

const Creds string = "Credentials"

func TestServer(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()

	assert.NotNil(compute.EmptyServer().Validate(ctx))

	s := compute.NewServer("", "", "test_server", "test_owner", "test_group")

	s.SerialNumber = "some_serial"
	assert.NotNil(s.Validate(ctx))

	s.BoardIP = net.IP{1, 1, 1, 1}
	assert.NotNil(s.Validate(ctx))

	s.Model = "latest"

	assert.NotNil(ctx)

	n := s.Meta.Type.Name
	s.Meta.Type.Name = "junk"
	assert.NotNil(s.Validate(ctx))
	s.Meta.Type.Name = n
	assert.NotNil(s.Validate(ctx))

	s.Credentials = zebra.NewCredentials("tester")
	assert.Nil(s.Credentials.Add("password", "thisIsAGoodPassword!123"))
	assert.Nil(s.Validate(ctx))
}

func TestESX(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()
	assert.NotNil(compute.EmptyESX().Validate(ctx))

	e := compute.NewESX("", "test_esx", "test_owner", "test_group")
	assert.NotNil(e.Validate(ctx))

	e.IP = net.IP{1, 1, 1, 1}
	assert.NotNil(e.Validate(ctx))

	e.ServerID = "some_server"
	assert.NotNil(e.Validate(ctx))

	n := e.Meta.Type.Name
	e.Meta.Type.Name = ""
	assert.NotNil(e.Validate(ctx))
	e.Meta.Type.Name = n

	e.Credentials = zebra.NewCredentials("tester")
	assert.Nil(e.Credentials.Add("password", "thisIsAGoodPassword!123"))
	assert.Nil(e.Validate(ctx))
}

func TestVCenter(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()
	assert.NotNil(compute.EmptyVCenter().Validate(ctx))

	v := compute.NewVCenter("test_vcenter", "test_owner", "test_group")
	assert.NotNil(v.Validate(ctx))

	v.IP = net.IP{1, 1, 1, 1}
	assert.NotNil(v.Validate(ctx))

	n := v.Meta.Type.Name
	v.Meta.Type.Name = ""
	assert.NotNil(v.Validate(ctx))
	v.Meta.Type.Name = n

	v.Credentials = zebra.NewCredentials("tester")
	assert.Nil(v.Credentials.Add("password", "thisIsAGoodPassword!123"))
	assert.Nil(v.Validate(ctx))
}

func TestVM(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()
	assert.NotNil(compute.EmptyVM().Validate(ctx))

	v := compute.NewVM("", "test_vm", "test_owner", "test_group")
	assert.NotNil(v.Validate(ctx))

	v.ESXID = "some_esx"
	assert.NotNil(v.Validate(ctx))

	v.ManagementIP = net.IP{1, 1, 1, 1}
	assert.NotNil(v.Validate(ctx))

	n := v.Meta.Type.Name
	v.Meta.Type.Name = ""
	assert.NotNil(v.Validate(ctx))
	v.Meta.Type.Name = n

	v.Credentials = zebra.NewCredentials("tester")
	assert.Nil(v.Credentials.Add("password", "thisIsAGoodPassword!123"))
	assert.Nil(v.Validate(ctx))
}
