package network_test

import (
	"context"
	"net"
	"testing"

	"github.com/project-safari/zebra/model/network"
	"github.com/stretchr/testify/assert"
)

func TestIPAddressPool(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()

	ip := network.EmptyIPAddressPool()
	assert.NotNil(ip.Validate(ctx))

	ip1 := network.NewIPAddressPool("test_ip", "test_owner", "test_group")
	assert.NotNil(ip1.Validate(ctx))

	ip1.Subnets = []net.IPNet{{}}
	assert.NotNil(ip1.Validate(ctx))
	ip1.Subnets = []net.IPNet{{IP: net.IP{1, 1, 1, 1}}}
	assert.NotNil(ip1.Validate(ctx))

	ip1.Subnets = []net.IPNet{{IP: net.IP{1, 1, 1, 1}, Mask: net.IPMask{255, 255, 255, 0}}}
	assert.Nil(ip1.Validate(ctx))

	ip1.Meta.Type.Name = "junk"
	assert.NotNil(ip1.Validate(ctx))
}

func TestNewIP(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ip := network.NewIPAddressPool("test_ip", "test_owner", "test_group")
	assert.NotNil(ip)

	ip.Subnets = []net.IPNet{{IP: net.IP{1, 1, 1, 1}, Mask: net.IPMask{255, 255, 255, 0}}}
	assert.Nil(ip.Validate(context.Background()))
}
