package dc_test

import (
	"context"
	"testing"

	"github.com/rchamarthy/zebra/dc"
	"github.com/stretchr/testify/assert"
)

func TestDatacenter(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()
	datacenter := new(dc.Datacenter)
	assert.NotNil(datacenter.Validate(ctx))

	datacenter.ID = "abracadabra"
	datacenter.Name = "jasmine"
	datacenter.Address = "1 palace st, agrabah"
	assert.Nil(datacenter.Validate(ctx))
}

func TestRack(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()
	rack := new(dc.Rack)
	assert.NotNil(rack.Validate(ctx))

	rack.ID = "abracadabra"
	rack.Name = "sher"
	rack.Row = "bazar"
	assert.Nil(rack.Validate(ctx))
}
