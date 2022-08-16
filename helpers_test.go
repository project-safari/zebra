package zebra_test

import (
	"testing"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/network"
	"github.com/stretchr/testify/assert"
)

func TestNewBaseResource(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	br := zebra.NewBaseResource(zebra.DefaultType(), nil)
	assert.NotNil(br)
	assert.NotEmpty(br.ID)
	assert.True(br.Type.Name == "BaseResource")

	labels := zebra.Labels{"owner": "shravya"}

	br = zebra.NewBaseResource(network.SwitchType(), labels)
	assert.NotNil(br)
	assert.NotEmpty(br.ID)
	assert.True(br.Type.Name == "Switch")
	assert.True(br.Labels != nil)
	assert.True(br.Labels.MatchEqual("owner", "shravya"))
}

func TestIsIn(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	list := []string{"hi", "hello", "goodbye"}

	assert.True(zebra.IsIn("hello", list))
	assert.False(zebra.IsIn("hey", list))
}
