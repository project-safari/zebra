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
	res := &zebra.BaseResource{
		ID:     "",
		Type:   "",
		Labels: zebra.Labels{"key": "value"},
	}
	assert.NotNil(res.Validate(ctx))

	res.ID = "abracadabra"
	assert.NotNil(res.Validate(ctx))

	res.Type = "BaseResource"
	assert.Nil(res.Validate(ctx))

	assert.True(res.GetID() == res.ID)
	assert.True(res.GetType() == res.Type)
	assert.True(res.GetLabels().HasKey("key"))
}

// TestBaseResource tests the *NamedResource Validate function with a pass case
// and a fail case.
func TestNamedResource(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()
	res := &zebra.NamedResource{
		BaseResource: zebra.BaseResource{
			ID:     "",
			Type:   "",
			Labels: zebra.Labels{"key": "value"},
		},
		Name: "",
	}
	assert.NotNil(res.Validate(ctx))

	res.ID = "abracadabra"
	assert.NotNil(res.Validate(ctx))
	assert.True(res.GetID() == res.ID)

	res.Type = "NamedResource"
	assert.NotNil(res.Validate(ctx))
	assert.True(res.GetType() == res.Type)

	res.Name = "jasmine"
	assert.Nil(res.Validate(ctx))

	assert.True(res.GetLabels().HasKey("key"))
}

func TestCredentials(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()

	credentials := zebra.Credentials{
		NamedResource: zebra.NamedResource{
			BaseResource: zebra.BaseResource{
				ID:     "",
				Type:   "Credentials",
				Labels: zebra.Labels{},
			},
			Name: "",
		},
		Keys: nil,
	}
	assert.NotNil(credentials.Validate(ctx))

	credentials.ID = "id123"
	assert.NotNil(credentials.Validate(ctx))

	credentials.Type = "Credentials"
	assert.NotNil(credentials.Validate(ctx))

	credentials.Name = "name"
	assert.NotNil(credentials.Validate(ctx))

	credentials.Keys = make(map[string]string)
	assert.Nil(credentials.Validate(ctx))

	credentials.Keys["password"] = "a"
	credentials.Keys["ssh-key"] = "test"
	assert.NotNil(credentials.Validate(ctx))

	credentials.Keys["password"] = "abcdefghijklm"
	assert.NotNil(credentials.Validate(ctx))

	credentials.Keys["password"] = "ABCDEFGHIJKLM"
	assert.NotNil(credentials.Validate(ctx))

	credentials.Keys["password"] = "ABCDEFghijklm"
	assert.NotNil(credentials.Validate(ctx))

	credentials.Keys["password"] = "ABCDEFghijklm1"
	assert.NotNil(credentials.Validate(ctx))

	credentials.Keys["password"] = "properPass123$"
	assert.Nil(credentials.Validate(ctx))
}
