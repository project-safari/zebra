package network_test

import (
	"context"
	"testing"

	"github.com/project-safari/zebra/model/network"
	"github.com/stretchr/testify/assert"
)

func TestVLANPool(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()

	v := network.EmptyVLANPool()
	assert.NotNil(v.Validate(ctx))

	v1 := network.NewVLANPool("test_vlan", "test_owner", "test_group")
	assert.Nil(v1.Validate(ctx))

	v1.RangeStart = 1
	assert.Equal(network.ErrInvalidRange, v1.Validate(ctx))

	v1.RangeStart = 0
	v1.RangeEnd = 10
	assert.Nil(v1.Validate(ctx))
	assert.Equal("0-10", v1.String())

	v1.Meta.Type.Name = "duck"
	assert.NotNil(v1.Validate(ctx))
}

func TestNewVlan(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	newV := network.NewVLANPool("test_vlan", "test_owner", "test_group")
	assert.NotNil(newV)
}
