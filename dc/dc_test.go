// Package dc_test tests structs and functions pertaining to datacenter resources
// outlined in the dc package.
package dc_test

import (
	"context"
	"testing"

	"github.com/project-safari/zebra/dc"
	"github.com/stretchr/testify/assert"
)

// TestDatacenter tests the *Datacenter Validate function with a pass and a fail
// case.
func TestDatacenter(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()
	datacenter := new(dc.Datacenter)
	assert.NotNil(datacenter.Validate(ctx))

	datacenter.ID = "abracadabra"
	datacenter.Type = "Datacenter"
	datacenter.Name = "jasmine"
	datacenter.Address = "1 palace st, agrabah"
	assert.Nil(datacenter.Validate(ctx))
}

// TestRack tests the *Rack Validate function with a pass and a fail case.
func TestRack(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()
	rack := new(dc.Rack)
	assert.NotNil(rack.Validate(ctx))

	rack.ID = "abracadabra"
	rack.Type = "Rack"
	rack.Name = "sher"
	rack.Row = "bazar"
	assert.Nil(rack.Validate(ctx))
}
