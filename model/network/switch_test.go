package network_test

import (
	"context"
	"net"
	"testing"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/model/network"
	"github.com/stretchr/testify/assert"
)

func TestSwitch(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()
	s := network.EmptySwitch()

	assert.NotNil(s.Validate(ctx))

	s1, ok := s.(*network.Switch)
	assert.True(ok)
	assert.NotNil(s1.Validate(ctx))

	s1.ManagementIP = net.IP{0, 0, 0, 0}
	assert.NotNil(s1.Validate(ctx))

	s1.SerialNumber = "fake-serial"
	assert.NotNil(s1.Validate(ctx))

	s1.Model = "fake-model"
	assert.NotNil(s1.Validate(ctx))

	s1.NumPorts = 96
	assert.NotNil(s1.Validate(ctx))

	s1.Meta.Type.Name = "blah"
	assert.NotNil(s1.Validate(ctx))
	s1.Meta.Type.Name = "network.switch"
	assert.NotNil(s1.Validate(ctx))

	s1.Credentials = zebra.NewCredentials("test_user")
	assert.Nil(s1.Credentials.Add("password", "aNewPassword123!"))
	assert.NotNil(s1.Validate(ctx))
}

func TestNewSwitch(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()

	s := network.NewSwitch("test_switch", "test_owner", "test_group")
	assert.NotNil(s)
	assert.NotNil(s.Validate(ctx))

	s.ManagementIP = net.IP{1, 1, 1, 1}
	s.SerialNumber = "fake-serial"
	s.Model = "fake-model"
	s.NumPorts = 96
	s.Meta.Type.Name = "blah"
	s.Meta.Type.Name = "network.switch"
	s.Credentials = zebra.NewCredentials("test_user")
	assert.Nil(s.Credentials.Add("password", "aNewPassword123!"))
	assert.Nil(s.Validate(ctx))
}
