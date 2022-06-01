package dc_test

import (
	"context"
	"testing"

	"github.com/rchamarthy/zebra/dc"
	"github.com/stretchr/testify/assert"
)

func TestDatacenter(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	d := &dc.Datacenter{}
	assert.NotNil(d.Validate(ctx))

	d.ID = "abracadabra"
	d.Name = "jasmine"
	d.Address = "1 palace st, agrabah"
	assert.Nil(d.Validate(ctx))
}

func TestRack(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	r := &dc.Rack{}
	assert.NotNil(r.Validate(ctx))

	r.ID = "abracadabra"
	r.Name = "sher"
	r.Row = "bazar"
	assert.Nil(r.Validate(ctx))
}
