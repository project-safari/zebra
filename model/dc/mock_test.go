package dc_test

import (
	"context"
	"testing"

	"github.com/project-safari/zebra/model/dc"
	"github.com/stretchr/testify/assert"
)

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
