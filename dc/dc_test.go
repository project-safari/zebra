// Package dc_test tests structs and functions pertaining to datacenter resources
// outlined in the dc package.
package dc_test

import (
	"context"
	"testing"

	"github.com/project-safari/zebra/cmd/herd/pkg"
	"github.com/project-safari/zebra/dc"
	"github.com/stretchr/testify/assert"
)

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
	datacenter.Type = "Datacenter"
	datacenter.Name = "jasmine"
	datacenter.Address = "1 palace st, agrabah"
	assert.NotNil(datacenter.Validate(ctx))

	datacenter.Labels = pkg.CreateLabels()
	datacenter.Labels = pkg.GroupLabels(datacenter.Labels, "someGroup")
	assert.Nil(datacenter.Validate(ctx))

	datacenter.Type = "test1"
	assert.NotNil(datacenter.Validate(ctx))

	labType := dc.LabType()
	assert.NotNil(labType)
}

// Test lab data.
func TestLab(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()
	labType := dc.LabType()
	lab, ok := labType.New().(*dc.Lab)
	assert.True(ok)
	assert.NotNil(lab.Validate(ctx))

	lab.ID = "abracadabra"
	lab.Type = "Lab"
	lab.Name = "sher"

	lab.Labels = pkg.CreateLabels()
	assert.Nil(lab.Validate(ctx))

	lab.Labels = pkg.GroupLabels(lab.Labels, "oneGroup")
	assert.Nil(lab.Validate(ctx))

	lab.Type = "test2"
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
	rack.Type = "Rack"
	rack.Name = "sher"
	rack.Row = "bazar"
	assert.NotNil(rack.Validate(ctx))

	rack.Type = "test3"
	assert.NotNil(rack.Validate(ctx))
}

// Test for datacenter generation.
func TestNewDc(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	labels := pkg.CreateLabels()

	labels = pkg.GroupLabels(labels, "group")

	dataC := dc.NewDatacenter(pkg.Addresses(), pkg.Name(), labels)
	assert.NotNil(dataC)
}

// Test for lab generation.
func TestNewLab(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	labels := pkg.CreateLabels()

	labels = pkg.GroupLabels(labels, "group")

	dataC := dc.NewLab(pkg.Name(), labels)
	assert.NotNil(dataC)
}

// Test for rack generation.
func TestNewRack(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	labels := pkg.CreateLabels()

	labels = pkg.GroupLabels(labels, "group")

	dataC := dc.NewRack(pkg.Name(), pkg.Rows(), labels)
	assert.NotNil(dataC)
}
