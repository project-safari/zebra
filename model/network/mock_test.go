package network_test

import (
	"context"
	"testing"

	"github.com/project-safari/zebra/model/network"
	"github.com/stretchr/testify/assert"
)

func TestMockSwitch(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rs := network.MockSwitch(10)
	assert.NotEmpty(rs)
	assert.Len(rs, 10)

	ctx := context.Background()
	for _, r := range rs {
		assert.Nil(r.Validate(ctx))
	}
}

func TestMockVLANPool(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rs := network.MockVLANPool(10)
	assert.NotEmpty(rs)
	assert.Len(rs, 10)

	ctx := context.Background()
	for _, r := range rs {
		assert.Nil(r.Validate(ctx))
	}
}

func TestMockIPAddressPool(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rs := network.MockIPAddressPool(10)
	assert.NotEmpty(rs)
	assert.Len(rs, 10)

	ctx := context.Background()
	for _, r := range rs {
		assert.Nil(r.Validate(ctx))
	}
}
