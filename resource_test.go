// Package zebra_test tests structs and functions outlined in the zebra package.
package zebra_test

import (
	"context"
	"testing"

	"github.com/project-safari/zebra"
	"github.com/stretchr/testify/assert"
)

// TestBaseResource tests the *BaseResource Validate function with a pass case
// and a fail case.
func TestBaseResource(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	d, _ := dummyType()
	ctx := context.Background()
	res := &zebra.BaseResource{
		Meta:   zebra.NewMeta(d, "", "", ""),
		Status: zebra.DefaultStatus(),
	}
	res.Meta.Name = ""
	assert.NotNil(res.Validate(ctx))

	res.Meta.ID = "abracadabra"
	assert.NotNil(res.Validate(ctx))
	assert.Equal("abracadabra", res.GetMeta().ID)
}

func TestGettingStatus(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	d, _ := dummyType()
	ctx := context.Background()

	res := zebra.NewBaseResource(d, "dummy", "dummy", "dummy")
	assert.Nil(res.Validate(ctx))
	assert.NotNil(res.GetStatus())

	res.Status.Fault = zebra.Fault(100)
	assert.NotNil(res.Validate(ctx))
}
