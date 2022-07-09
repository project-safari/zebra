package store_test

import (
	"os"
	"testing"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/dc"
	"github.com/project-safari/zebra/network"
	"github.com/project-safari/zebra/store"
	"github.com/stretchr/testify/assert"
)

func TestNewResourceStore(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	t.Cleanup(func() { os.RemoveAll("teststore") })

	assert.NotNil(store.NewResourceStore("teststore", nil))
}

func TestInitialize(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	t.Cleanup(func() { os.RemoveAll("teststore1") })

	rs := store.NewResourceStore("teststore1", nil)
	assert.NotNil(rs)
	assert.Nil(rs.Initialize())
}

func TestWipe(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	t.Cleanup(func() { os.RemoveAll("teststore2") })

	rs := store.NewResourceStore("teststore2", nil)
	assert.NotNil(rs)
	assert.Nil(rs.Initialize())
	assert.Nil(rs.Wipe())
}

func TestClear(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	t.Cleanup(func() { os.RemoveAll("teststore3") })

	factory := zebra.Factory()
	factory.Add("VLANPool", func() zebra.Resource { return new(network.VLANPool) })

	rs := store.NewResourceStore("teststore3", factory)
	assert.NotNil(rs)
	assert.Nil(rs.Initialize())

	assert.Nil(rs.Create(getVLAN()))
	assert.Nil(rs.Create(getVLAN()))

	resources, err := rs.Load()
	assert.Nil(err)
	assert.True(len(resources.Resources) == 1)
	assert.True(len(resources.Resources["VLANPool"].Resources) == 2)

	assert.Nil(rs.Clear())

	resources, err = rs.Load()
	assert.Nil(err)
	assert.True(len(resources.Resources) == 0)
}

func TestLoad(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	t.Cleanup(func() { os.RemoveAll("teststore4") })

	factory := zebra.Factory()
	factory.Add("VLANPool", func() zebra.Resource { return new(network.VLANPool) })

	rs := store.NewResourceStore("teststore4", factory)
	assert.NotNil(rs)
	assert.Nil(rs.Initialize())

	resources, err := rs.Load()
	assert.Nil(err)
	assert.True(len(resources.Resources) == 0)

	assert.Nil(rs.Create(getVLAN()))
	resources, err = rs.Load()
	assert.Nil(err)
	assert.True(len(resources.Resources) == 1)
	assert.True(len(resources.Resources["VLANPool"].Resources) == 1)

	assert.Nil(rs.Create(getVLAN()))

	resources, err = rs.Load()
	assert.Nil(err)
	assert.True(len(resources.Resources) == 1)
	assert.True(len(resources.Resources["VLANPool"].Resources) == 2)
}

func TestCreate(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	t.Cleanup(func() { os.RemoveAll("teststore5") })

	factory := zebra.Factory()
	factory.Add("VLANPool", func() zebra.Resource { return new(network.VLANPool) })

	rs := store.NewResourceStore("teststore5", factory)
	assert.NotNil(rs)
	assert.Nil(rs.Initialize())

	// Invalid resource, should fail
	assert.NotNil(rs.Create(nil))

	// Valid resource, should pass
	vlan := getVLAN()
	assert.Nil(rs.Create(vlan))

	resources, err := rs.Load()
	assert.Nil(err)
	assert.True(len(resources.Resources) == 1)
	assert.True(len(resources.Resources["VLANPool"].Resources) == 1)

	// Duplicate resource, should fail
	assert.NotNil(rs.Create(vlan))
}

func TestUpdate(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	t.Cleanup(func() { os.RemoveAll("teststore6") })

	factory := zebra.Factory()
	factory.Add("VLANPool", func() zebra.Resource { return new(network.VLANPool) })

	rs := store.NewResourceStore("teststore6", factory)
	assert.NotNil(rs)
	assert.Nil(rs.Initialize())

	// Valid resource, should pass
	vlan := getVLAN()
	assert.Nil(rs.Create(vlan))

	resources, err := rs.Load()
	assert.Nil(err)
	assert.True(len(resources.Resources) == 1)
	assert.True(len(resources.Resources["VLANPool"].Resources) == 1)

	// Update resource, should pass
	assert.Nil(rs.Update(vlan))

	// Update non-existent resource, should fail
	assert.NotNil(rs.Update(nil))

	// Update uncreated resource, should fail
	assert.NotNil(rs.Update(getLab()))
}

func TestDelete(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	t.Cleanup(func() { os.RemoveAll("teststore7") })

	factory := zebra.Factory()
	factory.Add("VLANPool", func() zebra.Resource { return new(network.VLANPool) })

	rs := store.NewResourceStore("teststore7", factory)
	assert.NotNil(rs)
	assert.Nil(rs.Initialize())

	// Valid resource, should pass
	vlan := getVLAN()
	assert.Nil(rs.Create(vlan))

	resources, err := rs.Load()
	assert.Nil(err)
	assert.True(len(resources.Resources) == 1)
	assert.True(len(resources.Resources["VLANPool"].Resources) == 1)

	// Delete resource, should pass
	assert.Nil(rs.Delete(vlan))

	// Delete non-existent resource, should fail
	assert.NotNil(rs.Update(nil))

	// Delete uncreated resource, should pass anyways
	assert.NotNil(rs.Update(getLab()))
}

func getVLAN() zebra.Resource {
	return &network.VLANPool{
		BaseResource: *zebra.NewBaseResource("VLANPool", nil),
		RangeStart:   0,
		RangeEnd:     1,
	}
}

func getLab() zebra.Resource {
	br := *zebra.NewBaseResource("Lab", nil)

	return &dc.Lab{
		NamedResource: zebra.NamedResource{
			BaseResource: br,
			Name:         "lab" + br.ID,
		},
	}
}
