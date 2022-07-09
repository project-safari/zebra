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

	root := "teststore"

	t.Cleanup(func() { os.RemoveAll(root) })

	assert.NotNil(store.NewResourceStore(root, nil))
}

func TestInitialize(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore1"

	t.Cleanup(func() { os.RemoveAll(root) })

	rs := store.NewResourceStore(root, nil)
	assert.NotNil(rs)
	assert.Nil(rs.Initialize())
}

func TestWipe(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore2"

	t.Cleanup(func() { os.RemoveAll(root) })

	rs := store.NewResourceStore(root, nil)
	assert.NotNil(rs)
	assert.Nil(rs.Initialize())
	assert.Nil(rs.Wipe())
}

func TestClear(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore3"

	t.Cleanup(func() { os.RemoveAll(root) })

	factory := zebra.Factory()
	factory.Add("VLANPool", func() zebra.Resource { return new(network.VLANPool) })

	rs := store.NewResourceStore(root, factory)
	assert.NotNil(rs)
	assert.Nil(rs.Initialize())

	assert.Nil(rs.Create(getVLAN()))
	assert.Nil(rs.Create(getVLAN()))

	resources, err := rs.Load()
	assert.Nil(err)
	assert.Equal(1, len(resources.Resources))
	assert.Equal(2, len(resources.Resources["VLANPool"].Resources))

	assert.Nil(rs.Clear())

	resources, err = rs.Load()
	assert.Nil(err)
	assert.Empty(len(resources.Resources))
}

func TestLoad(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore4"

	t.Cleanup(func() { os.RemoveAll(root) })

	factory := zebra.Factory()
	factory.Add("VLANPool", func() zebra.Resource { return new(network.VLANPool) })

	rs := store.NewResourceStore(root, factory)
	assert.NotNil(rs)
	assert.Nil(rs.Initialize())

	resources, err := rs.Load()
	assert.Nil(err)
	assert.Empty(len(resources.Resources))

	assert.Nil(rs.Create(getVLAN()))
	resources, err = rs.Load()
	assert.Nil(err)
	assert.Equal(1, len(resources.Resources))
	assert.Equal(1, len(resources.Resources["VLANPool"].Resources))

	assert.Nil(rs.Create(getVLAN()))

	resources, err = rs.Load()
	assert.Nil(err)
	assert.Equal(1, len(resources.Resources))
	assert.Equal(2, len(resources.Resources["VLANPool"].Resources))
}

func TestCreate(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore5"

	t.Cleanup(func() { os.RemoveAll(root) })

	factory := zebra.Factory()
	factory.Add("VLANPool", func() zebra.Resource { return new(network.VLANPool) })

	rs := store.NewResourceStore(root, factory)
	assert.NotNil(rs)
	assert.Nil(rs.Initialize())

	// Invalid resource, should fail
	assert.NotNil(rs.Create(nil))

	// Valid resource, should pass
	vlan := getVLAN()
	assert.Nil(rs.Create(vlan))

	resources, err := rs.Load()
	assert.Nil(err)
	assert.Equal(1, len(resources.Resources))
	assert.Equal(1, len(resources.Resources["VLANPool"].Resources))

	// Duplicate resource, should update
	assert.Nil(rs.Create(vlan))
}

func TestDelete(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore6"

	t.Cleanup(func() { os.RemoveAll(root) })

	factory := zebra.Factory()
	factory.Add("VLANPool", func() zebra.Resource { return new(network.VLANPool) })

	rs := store.NewResourceStore(root, factory)
	assert.NotNil(rs)
	assert.Nil(rs.Initialize())

	// Valid resource, should pass
	vlan := getVLAN()
	assert.Nil(rs.Create(vlan))

	resources, err := rs.Load()
	assert.Nil(err)
	assert.Equal(1, len(resources.Resources))
	assert.Equal(1, len(resources.Resources["VLANPool"].Resources))

	// Delete resource, should pass
	assert.Nil(rs.Delete(vlan))

	// Delete non-existent resource, should fail
	assert.NotNil(rs.Delete(nil))

	// Delete uncreated resource, should pass anyways
	assert.NotNil(rs.Delete(getLab()))
}

func getVLAN() *network.VLANPool {
	return &network.VLANPool{
		BaseResource: *zebra.NewBaseResource("VLANPool", nil),
		RangeStart:   0,
		RangeEnd:     1,
	}
}

func getLab() *dc.Lab {
	br := *zebra.NewBaseResource("Lab", nil)

	return &dc.Lab{
		NamedResource: zebra.NamedResource{
			BaseResource: br,
			Name:         "lab" + br.ID,
		},
	}
}
