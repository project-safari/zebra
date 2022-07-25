package zebra_test

import (
	"testing"

	"github.com/project-safari/zebra"
	"github.com/stretchr/testify/assert"
)

func TestNewGroup(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	g := zebra.NewGroup("sj-building-15")
	assert.NotNil(g)
	assert.NotNil(g.Resources)
	assert.NotNil(g.FreePool)
	assert.Empty(g.Count)
}

func TestAddDeleteGroup(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	g := zebra.NewGroup("sj-building-15")
	assert.NotNil(g)

	res1 := zebra.NewBaseResource("", nil)
	res2 := zebra.NewBaseResource("", nil)

	g.Add(res1)
	assert.Equal(1, len(g.Resources.Resources))
	assert.Equal(1, len(g.FreePool.Resources))
	assert.Equal(1, g.Count)

	g.Add(res2)
	assert.Equal(1, len(g.Resources.Resources))
	assert.Equal(1, len(g.FreePool.Resources))
	assert.Equal(2, g.Count)

	g.Delete(res1)
	assert.Equal(1, len(g.Resources.Resources))
	assert.Equal(1, len(g.FreePool.Resources))
	assert.Equal(1, g.Count)

	g.Delete(res2)
	assert.Empty(len(g.Resources.Resources))
	assert.Empty(len(g.FreePool.Resources))
	assert.Empty(g.Count)

	// Try to delete again. Should not panic, but should not update group count.
	g.Delete(res2)
	assert.Empty(g.Count)
}

func TestValidateGroup(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	g := &zebra.Group{
		Name:      "",
		Resources: nil,
		FreePool:  nil,
		Count:     2,
	}
	assert.NotNil(g.Validate())

	g.Name = "not empty"
	assert.NotNil(g.Validate())

	g.Resources = zebra.NewResourceMap(nil)
	assert.NotNil(g.Validate())

	g.FreePool = zebra.NewResourceMap(nil)
	assert.NotNil(g.Validate())

	g.Add(zebra.NewBaseResource("", nil))
	assert.NotNil(g.Validate())

	g.Count = 1
	assert.Nil(g.Validate())
}
