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
	assert.Equal(1, len(g.FreePool().Resources))

	g.Add(res2)
	assert.Equal(1, len(g.Resources.Resources))
	assert.Equal(1, len(g.FreePool().Resources))

	g.Delete(res1)
	assert.Equal(1, len(g.Resources.Resources))
	assert.Equal(1, len(g.FreePool().Resources))

	g.Delete(res2)
	assert.Empty(len(g.Resources.Resources))
	assert.Empty(len(g.FreePool().Resources))

	// Try to delete again. Should not panic.
	g.Delete(res2)
}

func TestFreeLease(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	res := zebra.NewBaseResource("", map[string]string{"system.group": "test-group"})
	resMap := zebra.NewResourceMap(nil)
	resMap.Add(res, res.Type)

	free := zebra.NewResourceMap(nil)
	zebra.CopyResourceMap(free, resMap)

	g := &zebra.Group{
		Name:      "test-group",
		Resources: resMap,
	}

	assert.Nil(g.Lease(res))
	assert.NotNil(g.Lease(res))
	assert.Nil(g.Free(res))
	assert.NotNil(g.Free(res))

	assert.NotNil(g.Lease(zebra.NewBaseResource("", nil)))
	assert.NotNil(g.Free(zebra.NewBaseResource("", nil)))
}

func TestValidateGroup(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	g := &zebra.Group{
		Name:      "",
		Resources: nil,
	}
	assert.NotNil(g.Validate())

	g.Name = "not empty"
	assert.NotNil(g.Validate())

	g.Resources = zebra.NewResourceMap(nil)
	assert.Nil(g.Validate())
}
