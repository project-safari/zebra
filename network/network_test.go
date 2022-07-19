// Package network_test tests structs and functions pertaining to network resources
// outlined in the network package.
package network_test

import (
	"context"
	"net"
	"testing"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/network"
	"github.com/stretchr/testify/assert"
)

// TestSwitch tests the *Switch Validate function with a pass and a fail case.
func TestSwitch(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()
	switchType := network.SwitchType()
	switch1, ok := switchType.New().(*network.Switch)
	assert.True(ok)
	assert.NotNil(switch1.Validate(ctx))

	switch1.ID = "aaaa"
	assert.NotNil(switch1.Validate(ctx))

	switch1.Type = "Switch"
	assert.NotNil(switch1.Validate(ctx))

	switch1.ManagementIP = net.ParseIP("10.1.0.0")
	assert.NotNil(switch1.Validate(ctx))

	switch1.SerialNumber = "bazar"
	assert.NotNil(switch1.Validate(ctx))

	switch1.Model = "latest and greatest"
	assert.NotNil(switch1.Validate(ctx))

	switch1.NumPorts = 12
	assert.NotNil(switch1.Validate(ctx))

	switch1.Credentials = zebra.Credentials{
		NamedResource: zebra.NamedResource{
			BaseResource: zebra.BaseResource{
				ID:     "blahblah",
				Type:   "Credentials",
				Labels: nil,
			},
			Name: "blah",
		},
		Keys: nil,
	}
	assert.NotNil(switch1.Validate(ctx))

	switch1.Credentials.Keys = make(map[string]string)
	assert.Nil(switch1.Validate(ctx))

	switch1.Type = "test"
	assert.NotNil(switch1.Validate(ctx))
}

// TestIPAddressPool tests the *IPAddressPool Validate function with a pass and a fail case.
func TestIPAddressPool(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()
	ipPoolType := network.IPAddressPoolType()
	pool, ok := ipPoolType.New().(*network.IPAddressPool)
	assert.True(ok)
	assert.NotNil(pool.Validate(ctx))

	pool.ID = "aaaa"
	pool.Type = "IPAddressPool"
	assert.Nil(pool.Validate(ctx))

	ipnet := net.IPNet{IP: net.ParseIP("192.0.2.1"), Mask: nil}
	ipnet.Mask = ipnet.IP.DefaultMask()
	pool.Subnets = append(pool.Subnets, ipnet)
	assert.Nil(pool.Validate(ctx))

	pool = new(network.IPAddressPool)
	pool.ID = "bbbb"
	ipnet1 := net.IPNet{IP: net.ParseIP("192.0.2.1"), Mask: nil}
	pool.Subnets = append(pool.Subnets, ipnet1)
	assert.NotNil(pool.Validate(ctx))

	pool = new(network.IPAddressPool)
	pool.ID = "cccc"
	ipnet2 := net.IPNet{IP: nil, Mask: nil}
	pool.Subnets = append(pool.Subnets, ipnet2)
	assert.NotNil(pool.Validate(ctx))

	pool.Type = "test"
	assert.NotNil(pool.Validate(ctx))
}

// TestVLANPool tests the *VLANPool Validate function with a pass and a fail case.
func TestVLANPool(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()
	vlanPoolType := network.VLANPoolType()
	pool, ok := vlanPoolType.New().(*network.VLANPool)
	assert.True(ok)
	assert.NotNil(pool.Validate(ctx))

	pool.ID = "cccc"
	pool.Type = "VLANPool"
	pool.RangeStart = 10
	pool.RangeEnd = 1
	assert.NotNil(pool.Validate(ctx))

	pool.RangeEnd = 11
	assert.Nil(pool.Validate(ctx))

	pool.Type = "test123"
	assert.NotNil(pool.Validate(ctx))
}
