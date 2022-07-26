// Package dc_test tests structs and functions pertaining to datacenter resources
// outlined in the dc package.
package dc_test

import (
	"context"
	"testing"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/cmd/herd/pkg"
	"github.com/project-safari/zebra/dc"
	"github.com/stretchr/testify/assert"
)

func EmptyType() zebra.Type {
	return zebra.Type{
		Name:        "Empty",
		Description: "Empty Type",
		Constructor: func() zebra.Resource { return nil },
	}
}

// TestDatacenter tests the *Datacenter Validate function with a pass and a fail
// case.
func TestDatacenter(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()
	dcType := dc.DataCenterType()
	assert.NotNil(dcType)

	datacenter, ok := dcType.New().(*dc.Datacenter)
	assert.True(ok)
	assert.NotNil(datacenter.Validate(ctx))

	datacenter.ID = "bahbah"
	datacenter.Type = dc.DataCenterType()
	datacenter.Name = "jasmine"
	datacenter.Address = "1 palace st, agrabah"
	assert.NotNil(datacenter.Validate(ctx))

	datacenter.Labels = pkg.CreateLabels()
	datacenter.Labels = pkg.GroupLabels(datacenter.Labels, "someGroup")
	assert.Nil(datacenter.Validate(ctx))

	datacenter.Type = EmptyType()
	assert.NotNil(datacenter.Validate(ctx))

	labType := dc.LabType()
	assert.NotNil(labType)
}

func TestLab(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()
	labType := dc.LabType()
	lab, ok := labType.New().(*dc.Lab)
	assert.True(ok)
	assert.NotNil(lab.Validate(ctx))

	lab.ID = "abracadabra"
	lab.Type = dc.LabType()
	lab.Name = "sher"

	lab.Labels = pkg.CreateLabels()
	assert.NotNil(lab.Validate(ctx))

	lab.Labels = pkg.GroupLabels(lab.Labels, "oneGroup")
	assert.Nil(lab.Validate(ctx))

	lab.Type = EmptyType()
	assert.NotNil(lab.Validate(ctx))
}

// TestRack tests the *Rack Validate function with a pass and a fail case.
func TestRack(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()
	rackType := dc.RackType()
	rack, ok := rackType.New().(*dc.Rack)
	assert.True(ok)
	assert.NotNil(rack.Validate(ctx))

	rack.ID = "abracadabra"
	rack.Type = dc.RackType()
	rack.Name = "sher"
	rack.Row = "bazar"
	assert.NotNil(rack.Validate(ctx))

	rack.Type = EmptyType()
	assert.NotNil(rack.Validate(ctx))
}
