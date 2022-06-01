package zebra_test

import (
	"context"
	"testing"

	"github.com/rchamarthy/zebra"
	"github.com/stretchr/testify/assert"
)

func TestBaseResource(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	r := &zebra.BaseResource{}
	assert.NotNil(r.Validate(ctx))

	r.ID = "abracadabra"
	assert.Nil(r.Validate(ctx))
}

func TestNamedResource(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	r := &zebra.NamedResource{}
	assert.NotNil(r.Validate(ctx))

	r.ID = "abracadabra"
	assert.NotNil(r.Validate(ctx))

	r.Name = "jasmine"
	assert.Nil(r.Validate(ctx))
}
