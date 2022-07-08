package zebra_test

import (
	"testing"

	"github.com/project-safari/zebra"
	"github.com/stretchr/testify/assert"
)

func TestNewBaseResource(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	br := zebra.NewBaseResource("", nil)
	assert.NotNil(br)
	assert.NotEmpty(br.ID)
	assert.True(br.Type == "BaseResource")
	assert.True(br.Labels == nil)

	labels := zebra.Labels{"owner": "shravya"}

	br = zebra.NewBaseResource("Switch", labels)
	assert.NotNil(br)
	assert.NotEmpty(br.ID)
	assert.True(br.Type == "Switch")
	assert.True(br.Labels != nil)
	assert.True(br.Labels.MatchEqual("owner", "shravya"))
}
