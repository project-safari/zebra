// Package zebra_test tests structs and functions outlined in the zebra package.
package zebra_test

import (
	"context"
	"testing"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/cmd/herd/pkg"
	"github.com/project-safari/zebra/network"
	"github.com/stretchr/testify/assert"
)

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
		Type:   zebra.DefaultType(),
		Labels: zebra.Labels{"key": "value"},
		Status: zebra.DefaultStatus(),
	}
	assert.NotNil(res.Validate(ctx, res.Type.Name))

	res.ID = "ab"
	assert.NotNil(res.Validate(ctx, res.Type.Name))
	assert.Equal("ab", res.GetName())

	res.ID = "abracadabra"
	assert.NotNil(res.Validate(ctx, res.Type.Name))

	res.Type = BaseResourceType()
	assert.NotNil(res.Validate(ctx, res.Type.Name))

	res.Labels.Add("system.group", "test")

	assert.Nil(res.Validate(ctx, res.Type.Name))
	assert.Equal(res.ID, res.GetID())
	assert.Equal(res.Type.Name, res.GetType().Name)
	assert.True(res.GetLabels().HasKey("key"))
	assert.Equal(res.ID[:7], res.GetName())
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
			Type:   zebra.DefaultType(),
			Labels: zebra.Labels{"key": "value"},
			Status: zebra.DefaultStatus(),
		},
		Name: "",
	}
	assert.NotNil(res.Validate(ctx, res.Type.Name))

	res.ID = "abracadabra"
	assert.NotNil(res.Validate(ctx, res.Type.Name))
	assert.Equal(res.ID, res.GetID())

	res.Type = NamedResourceType()
	assert.NotNil(res.Validate(ctx, res.Type.Name))
	assert.Equal(res.Type.Name, res.GetType().Name)

	res.Name = "jasmine"
	assert.NotNil(res.Validate(ctx, res.Type.Name))
	assert.Equal("jasmine", res.GetName())

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
				Type:   zebra.DefaultType(),
				Labels: zebra.Labels{},
				Status: zebra.DefaultStatus(),
			},
			Name: "",
		},
		Keys: nil,
	}
	assert.NotNil(credentials.Validate(ctx, credentials.Type.Name))

	credentials.ID = "id123"
	assert.NotNil(credentials.Validate(ctx, credentials.Type.Name))

	credentials.Type = zebra.CredentialsType()
	assert.NotNil(credentials.Validate(ctx, credentials.Type.Name))

	credentials.Name = "name"
	assert.NotNil(credentials.Validate(ctx, credentials.Type.Name))

	credentials.Keys = make(map[string]string)
	assert.NotNil(credentials.Validate(ctx, credentials.Type.Name))

	credentials.Keys["password"] = "a"
	credentials.Keys["ssh-key"] = "test"
	assert.NotNil(credentials.Validate(ctx, credentials.Type.Name))

	credentials.Keys["password"] = "abcdefghijklm"
	assert.NotNil(credentials.Validate(ctx, credentials.Type.Name))

	credentials.Keys["password"] = "ABCDEFGHIJKLM"
	assert.NotNil(credentials.Validate(ctx, credentials.Type.Name))

	credentials.Keys["password"] = "ABCDEFghijklm"
	assert.NotNil(credentials.Validate(ctx, credentials.Type.Name))

	credentials.Keys["password"] = "ABCDEFghijklm1"
	assert.NotNil(credentials.Validate(ctx, credentials.Type.Name))

	credentials.Keys["password"] = "properPass123$"
	assert.NotNil(credentials.Validate(ctx, credentials.Type.Name))
}

func TestLabelsValidation(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// first test - with a correct default label
	mapOne := map[string]string{
		"system.group": "Americas",
		"color":        "red",
	}

	resOne := zebra.NewBaseResource(zebra.DefaultType(), mapOne)
	assert.Nil(resOne.Validate(context.Background(), resOne.Type.Name))

	// second test - with an incorrect default label
	mapTwo := map[string]string{
		"letter": "alpha",
		"color":  "blue",
	}

	resTwo := zebra.NewBaseResource(zebra.DefaultType(), mapTwo)

	assert.Nil(resTwo.Validate(context.Background(), resTwo.Type.Name))
}

func TestNewCred(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	labels := pkg.CreateLabels()

	labels = pkg.GroupLabels(labels, "group")

	creds := zebra.NewCredential(pkg.Name(), labels)

	assert.NotNil(creds)
}

func TestGettingStatus(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ctx := context.Background()
	res := &zebra.NamedResource{
		BaseResource: zebra.BaseResource{
			ID:     "",
			Type:   network.SwitchType(),
			Labels: zebra.Labels{"key": "value"},
			Status: zebra.DefaultStatus(),
		},
		Name: "",
	}
	assert.NotNil(res.Validate(ctx, res.Type.Name))
	assert.NotNil(res.GetStatus())
}
