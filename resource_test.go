// Package zebra_test tests structs and functions outlined in the zebra package.
package zebra_test

import (
	"context"
	"testing"

	"github.com/project-safari/zebra"
	"github.com/stretchr/testify/assert"
)

func EmptyType() zebra.Type {
	return zebra.Type{
		Name:        "",
		Description: "Empty Type",
		Constructor: func() zebra.Resource { return nil },
	}
}

func BaseResourceType() zebra.Type {
	return zebra.Type{
		Name:        "BaseResource",
		Description: "Base Resource",
		Constructor: func() zebra.Resource { return nil },
	}
}

func NamedResourceType() zebra.Type {
	return zebra.Type{
		Name:        "NamedResource",
		Description: "Named Resource",
		Constructor: func() zebra.Resource { return nil },
	}
}

// TestBaseResource tests the *BaseResource Validate function with a pass case
// and a fail case.
func TestBaseResource(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()
	res := &zebra.BaseResource{
		ID:     "",
		Type:   EmptyType(),
		Labels: zebra.Labels{"key": "value"},
		Status: zebra.DefaultStatus(),
	}
	assert.NotNil(res.Validate(ctx))

	res.ID = "ab"
	assert.NotNil(res.Validate(ctx))

	res.ID = "abracadabra"
	assert.NotNil(res.Validate(ctx))

	res.Type = BaseResourceType()
	assert.NotNil(res.Validate(ctx))

	assert.Equal(res.ID, res.GetID())
	assert.Equal(res.Type.Name, res.GetType().Name)
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
			Type:   EmptyType(),
			Labels: zebra.Labels{"key": "value"},
			Status: zebra.DefaultStatus(),
		},
		Name: "",
	}
	assert.NotNil(res.Validate(ctx))

	res.ID = "abracadabra"
	assert.NotNil(res.Validate(ctx))
	assert.Equal(res.ID, res.GetID())

	res.Type = NamedResourceType()
	assert.NotNil(res.Validate(ctx))
	assert.Equal(res.Type.Name, res.GetType().Name)

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
				Type:   zebra.CredentialsType(),
				Labels: zebra.Labels{},
				Status: zebra.DefaultStatus(),
			},
			Name: "",
		},
		Keys: nil,
	}
	assert.NotNil(credentials.Validate(ctx))

	credentials.ID = "id123"
	assert.NotNil(credentials.Validate(ctx))

	credentials.Type = zebra.CredentialsType()
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
