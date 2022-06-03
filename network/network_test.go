// Package network_test tests structs and functions pertaining to network resources
// outlined in the network package.
package network_test

import (
	"context"
	"net"
	"testing"

	"github.com/rchamarthy/zebra/network"
	"github.com/stretchr/testify/assert"
)

// TestSwitch tests the *Switch Validate function with a pass and a fail case.
func TestSwitch(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()
	switch1 := new(network.Switch)
	assert.NotNil(switch1.Validate(ctx))

	switch1.ID = "a"
	assert.NotNil(switch1.Validate(ctx))

	switch1.ManagementIP = net.ParseIP("10.1.0.0")
	assert.NotNil(switch1.Validate(ctx))

	switch1.SerialNumber = "bazar"
	assert.NotNil(switch1.Validate(ctx))

	switch1.Model = "latest and greatest"
	assert.NotNil(switch1.Validate(ctx))

	switch1.NumPorts = 12
	assert.Nil(switch1.Validate(ctx))
}

// TestIPAddressPool tests the *IPAddressPool Validate function with a pass and a fail case.
func TestIPAddressPool(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()
	pool := new(network.IPAddressPool)
	assert.NotNil(pool.Validate(ctx))

	pool.ID = "b"
	assert.NotNil(pool.Validate(ctx))

	pool.IP = net.ParseIP("192.0.2.1")
	assert.NotNil(pool.Validate(ctx))

	pool.Mask = pool.IP.DefaultMask()
	assert.Nil(pool.Validate(ctx))
}

// TestVLANPool tests the *VLANPool Validate function with a pass and a fail case.
func TestVLANPool(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()
	pool := new(network.VLANPool)
	assert.NotNil(pool.Validate(ctx))

	pool.ID = "c"
	pool.RangeStart = 10
	pool.RangeEnd = 1
	assert.NotNil(pool.Validate(ctx))

	pool.RangeEnd = 11
	assert.Nil(pool.Validate(ctx))
}
