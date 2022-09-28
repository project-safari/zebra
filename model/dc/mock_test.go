package dc_test

import (
	"context"
	"testing"

	"github.com/project-safari/zebra/model/dc"
	"github.com/stretchr/testify/assert"
)

// Test function that examines the valid creation/generation of a mock dc.
func TestMockDC(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rs := dc.MockDC(10)
	assert.NotEmpty(rs)
	assert.Len(rs, 10)

	ctx := context.Background()
	for _, r := range rs {
		assert.Nil(r.Validate(ctx))
	}
}

// Test function that examines the valid creation/generation of a mock lab.
func TestMockLab(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rs := dc.MockLab(10)
	assert.NotEmpty(rs)
	assert.Len(rs, 10)

	ctx := context.Background()
	for _, r := range rs {
		assert.Nil(r.Validate(ctx))
	}
}

// Test function that examines the valid creation/generation of a mock rack.
func TestMockRack(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rs := dc.MockRack(10)
	assert.NotEmpty(rs)
	assert.Len(rs, 10)

	ctx := context.Background()
	for _, r := range rs {
		assert.Nil(r.Validate(ctx))
	}
}
