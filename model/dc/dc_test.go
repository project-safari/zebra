// Package dc_test tests structs and functions pertaining to datacenter resources
// outlined in the dc package.
package dc_test

import (
	"context"
	"testing"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/model/dc"
	"github.com/stretchr/testify/assert"
)

const name string = "junk"

// TestDatacenter tests the *Datacenter Validate function with a pass and a fail
// case.
func TestDatacenter(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()
	d, ok := dc.EmptyDataCenter().(*dc.Datacenter)
	assert.True(ok)
	assert.NotNil(d.Validate(ctx))

	d.Address = "on earth"
	assert.NotNil(d.Validate(ctx))

	d.Meta.Type.Name = name
	assert.Equal(zebra.ErrWrongType, d.Validate(ctx))

	d = dc.NewDatacenter(d.Address, "test", "test_owner", "test_group")
	assert.Nil(d.Validate(ctx))
}

// TestLab is a test function for the validate function on a lab resource.
func TestLab(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()
	l, ok := dc.EmptyLab().(*dc.Lab)
	assert.True(ok)
	assert.NotNil(l.Validate(ctx))

	l.Meta.Type.Name = name
	assert.Equal(zebra.ErrWrongType, l.Validate(ctx))

	l = dc.NewLab("test_lab", "test_owner", "test_group")
	assert.Nil(l.Validate(ctx))
}

// TestRack tests the *Rack Validate function with a pass and a fail case.
func TestRack(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()
	r, ok := dc.EmptyRack().(*dc.Rack)
	assert.True(ok)
	assert.NotNil(r.Validate(ctx))

	r.Row = "r"
	r.Meta.Type.Name = name
	assert.Equal(zebra.ErrWrongType, r.Validate(ctx))

	r = dc.NewRack("test_row", "test_rack", "test_owner", "test_group")
	assert.Nil(r.Validate(ctx))
}
