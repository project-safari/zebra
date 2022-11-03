package compute_test

import (
	"context"
	"testing"

	"github.com/project-safari/zebra/model/compute"
	"github.com/stretchr/testify/assert"
)

// Test function that examines the valid creation/generation of a mock server.
func TestMockServer(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rs := compute.MockServer(10)
	assert.NotEmpty(rs)
	assert.Len(rs, 10)

	ctx := context.Background()
	for _, r := range rs {
		assert.Nil(r.Validate(ctx))
	}
}

// Test function that examines the valid creation/generation of a mock esx server.
func TestMockESX(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rs := compute.MockESX(10)
	assert.NotEmpty(rs)
	assert.Len(rs, 10)

	ctx := context.Background()
	for _, r := range rs {
		assert.Nil(r.Validate(ctx))
	}
}

// Test function that examines the valid creation/generation of a mock vcenter.
func TestMockVCenter(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rs := compute.MockVCenter(10)
	assert.NotEmpty(rs)
	assert.Len(rs, 10)

	ctx := context.Background()
	for _, r := range rs {
		assert.Nil(r.Validate(ctx))
	}
}

// Test function that examines the valid creation/generation of a mock vm.
func TestMockVM(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rs := compute.MockVM(10)
	assert.NotEmpty(rs)
	assert.Len(rs, 10)

	ctx := context.Background()
	for _, r := range rs {
		assert.Nil(r.Validate(ctx))
	}
}
