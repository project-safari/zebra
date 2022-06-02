// Package zebra_test tests structs and functions outlined in the zebra package.
package zebra_test

import (
	"context"
	"testing"

	"github.com/rchamarthy/zebra"
	"github.com/stretchr/testify/assert"
)

// TestBaseResource tests the *BaseResource Validate function with a pass case
// and a fail case.
func TestBaseResource(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()
	r := &zebra.BaseResource{ID: ""}
	assert.NotNil(r.Validate(ctx))

	r.ID = "abracadabra"
	assert.Nil(r.Validate(ctx))
}

// TestBaseResource tests the *NamedResource Validate function with a pass case
// and a fail case.
func TestNamedResource(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()
	res := &zebra.NamedResource{
		BaseResource: zebra.BaseResource{ID: ""},
		Name:         "",
	}
	assert.NotNil(res.Validate(ctx))

	res.ID = "abracadabra"
	assert.NotNil(res.Validate(ctx))

	res.Name = "jasmine"
	assert.Nil(res.Validate(ctx))
}
