package compute_test

import (
	"context"
	"testing"

	"github.com/project-safari/zebra/model/compute"
	"github.com/stretchr/testify/assert"
)

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
