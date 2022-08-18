// Package zebra_test tests structs and functions outlined in the zebra package.
package zebra_test

import (
	"context"
	"testing"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/cmd/herd/pkg"
	"github.com/project-safari/zebra/status"
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
		Status: status.DefaultStatus(),
	}
	assert.NotNil(res.Validate(ctx))

	res.ID = "ab"
	assert.NotNil(res.Validate(ctx))
	assert.Equal("ab", res.GetName())

	res.ID = "abracadabra"
	assert.NotNil(res.Validate(ctx))

	res.Type = "BaseResource"
	assert.NotNil(res.Validate(ctx))

	res.Labels.Add("system.group", "test")

	assert.Nil(res.Validate(ctx))
	assert.Equal(res.ID, res.GetID())
	assert.Equal(res.Type, res.GetType())
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
			Type:   "",
			Labels: zebra.Labels{"key": "value"},
			Status: status.DefaultStatus(),
		},
		Name: "",
	}
	assert.NotNil(res.Validate(ctx))

	res.ID = "abracadabra"
	assert.NotNil(res.Validate(ctx))
	assert.Equal(res.ID, res.GetID())

	res.Type = "NamedResource"
	assert.NotNil(res.Validate(ctx))
	assert.Equal(res.Type, res.GetType())

	res.Name = "jasmine"
	assert.NotNil(res.Validate(ctx))
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
				Type:   "Credentials",
				Labels: zebra.Labels{},
				Status: status.DefaultStatus(),
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
	assert.NotNil(credentials.Validate(ctx))

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
	assert.NotNil(credentials.Validate(ctx))
}

func TestLabelsValidation(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// first test - with a correct default label
	mapOne := map[string]string{
		"system.group": "Americas",
		"color":        "red",
	}

	resOne := zebra.NewBaseResource("", mapOne)
	assert.Nil(resOne.Validate(context.Background()))

	// second test - with an incorrect default label
	mapTwo := map[string]string{
		"letter": "alpha",
		"color":  "blue",
	}

	resTwo := zebra.NewBaseResource("", mapTwo)

	assert.Nil(resTwo.Validate(context.Background()))
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
			Type:   "Switch",
			Labels: zebra.Labels{"key": "value"},
			Status: status.DefaultStatus(),
		},
		Name: "",
	}
	assert.NotNil(res.Validate(ctx))
	assert.NotNil(res.GetStatus())
}
