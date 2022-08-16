/// Package network_test tests structs and functions pertaining to network resources
// outlined in the network package.
package network_test

import (
	"context"
	"net"
	"testing"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/cmd/herd/pkg"
	"github.com/project-safari/zebra/network"
	"github.com/stretchr/testify/assert"
)

func Type() zebra.Type {
	return zebra.Type{
		Name:        "TestType",
		Description: "Test Type",
		Constructor: func() zebra.Resource { return nil },
	}
}

// TestSwitch tests the *Switch Validate function with a pass and a fail case.
func TestSwitch(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()
	switchType := network.SwitchType()
	switch1, ok := switchType.New().(*network.Switch)
	assert.True(ok)
	assert.NotNil(switch1.Validate(ctx, switch1.Type.Name))

	switch1.ID = "aaaa"
	assert.NotNil(switch1.Validate(ctx, switch1.Type.Name))

	switch1.Type = network.SwitchType()
	assert.NotNil(switch1.Validate(ctx, switch1.Type.Name))

	switch1.ManagementIP = net.ParseIP("10.1.0.0")
	assert.NotNil(switch1.Validate(ctx, switch1.Type.Name))

	switch1.SerialNumber = "bazar"
	assert.NotNil(switch1.Validate(ctx, switch1.Type.Name))

	switch1.Model = "latest and greatest"
	assert.NotNil(switch1.Validate(ctx, switch1.Type.Name))

	switch1.NumPorts = 12
	assert.NotNil(switch1.Validate(ctx, switch1.Type.Name))

	switch1.Credentials = zebra.Credentials{
		NamedResource: zebra.NamedResource{
			BaseResource: zebra.BaseResource{
				ID:     "blahblah",
				Type:   zebra.CredentialsType(),
				Labels: nil,
				Status: zebra.DefaultStatus(),
			},
			Name: "blah",
		},
		Keys: nil,
	}
	assert.NotNil(switch1.Validate(ctx, switch1.Type.Name))

	switch1.Credentials.Keys = make(map[string]string)
	assert.NotNil(switch1.Validate(ctx, switch1.Type.Name))

	switch1.Type = zebra.DefaultType()
	assert.NotNil(switch1.Validate(ctx, switch1.Type.Name))
}

// TestIPAddressPool tests the *IPAddressPool Validate function with a pass and a fail case.
func TestIPAddressPool(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()
	ipPoolType := network.IPAddressPoolType()
	pool, ok := ipPoolType.New().(*network.IPAddressPool)
	assert.True(ok)
	assert.NotNil(pool.Validate(ctx, pool.Type.Name))

	pool.ID = "aaaa"
	pool.Type = network.IPAddressPoolType()

	pool.Labels = make(map[string]string)
	pool.Labels = pkg.GroupLabels(pool.Labels, "groupSample")

	assert.Nil(pool.Validate(ctx, pool.Type.Name))

	ipnet := net.IPNet{IP: net.ParseIP("192.0.2.1"), Mask: nil}
	ipnet.Mask = ipnet.IP.DefaultMask()
	pool.Subnets = append(pool.Subnets, ipnet)
	assert.Nil(pool.Validate(ctx, pool.Type.Name))

	pool = new(network.IPAddressPool)
	pool.ID = "bbbb"
	ipnet1 := net.IPNet{IP: net.ParseIP("192.0.2.1"), Mask: nil}
	pool.Subnets = append(pool.Subnets, ipnet1)
	assert.NotNil(pool.Validate(ctx, pool.Type.Name))

	pool = new(network.IPAddressPool)
	pool.ID = "cccc"
	ipnet2 := net.IPNet{IP: nil, Mask: nil}
	pool.Subnets = append(pool.Subnets, ipnet2)
	assert.NotNil(pool.Validate(ctx, pool.Type.Name))

	pool.Type = Type()
	assert.NotNil(pool.Validate(ctx, pool.Type.Name))
}

// TestVLANPool tests the *VLANPool Validate function with a pass and a fail case.
func TestVLANPool(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()
	vlanPoolType := network.VLANPoolType()
	pool, ok := vlanPoolType.New().(*network.VLANPool)
	assert.True(ok)
	assert.NotNil(pool.Validate(ctx, pool.Type.Name))

	pool.ID = "cccc"
	pool.Type = network.VLANPoolType()
	pool.RangeStart = 10
	pool.RangeEnd = 1
	assert.NotNil(pool.Validate(ctx, pool.Type.Name))

	pool.RangeEnd = 11
	assert.NotNil(pool.Validate(ctx, pool.Type.Name))

	pool.Type = Type()
	assert.NotNil(pool.Validate(ctx, pool.Type.Name))
}

func TestNewVlan(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	labels := pkg.CreateLabels()

	labels = pkg.GroupLabels(labels, "group")

	newV := network.NewVlanPool(1, 100, labels)
	assert.NotNil(newV)
}

func TestNewSwitch(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	labels := pkg.CreateLabels()

	labels = pkg.GroupLabels(labels, "group")

	arr := []string{pkg.Serials(), pkg.Models(), pkg.Name()}
	newV := network.NewSwitch(arr, pkg.Ports(), net.IP("123.111.001"), labels)
	assert.NotNil(newV)
}

func TestNewIP(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	labels := pkg.CreateLabels()

	labels = pkg.GroupLabels(labels, "group")

	newV := network.NewIPAddressPool(pkg.CreateIPArr(3), labels)
	assert.NotNil(newV)
}
