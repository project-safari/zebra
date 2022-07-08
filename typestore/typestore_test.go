package typestore_test

import (
	"testing"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/network"
	"github.com/project-safari/zebra/typestore"
	"github.com/stretchr/testify/assert"
)

func TestNewTestStore(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resMap := zebra.NewResourceMap(nil)
	assert.NotNil(typestore.NewTypeStore(resMap))
}

func TestInitialize(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resMap := zebra.NewResourceMap(nil)

	ts := typestore.NewTypeStore(resMap)
	assert.NotNil(ts)

	assert.Nil(ts.Initialize())
}

func TestWipe(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resMap := zebra.NewResourceMap(nil)

	ts := typestore.NewTypeStore(resMap)
	assert.NotNil(ts)

	assert.Nil(ts.Wipe())
}

func TestClear(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resMap := zebra.NewResourceMap(nil)
	resMap.Add(new(network.VLANPool), "VLANPool")
	resMap.Add(new(network.IPAddressPool), "IPAddressPool")

	ts := typestore.NewTypeStore(resMap)
	assert.NotNil(ts)

	resources, err := ts.Load()
	assert.Nil(err)
	assert.True(len(resources.Resources) == 2)

	assert.Nil(ts.Clear())

	resources, err = ts.Load()
	assert.Nil(err)
	assert.True(len(resources.Resources) == 0)
}

func TestLoad(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resMap := zebra.NewResourceMap(nil)
	resMap.Add(new(network.VLANPool), "VLANPool")
	resMap.Add(new(network.IPAddressPool), "IPAddressPool")

	ts := typestore.NewTypeStore(resMap)
	assert.NotNil(ts)

	resources, err := ts.Load()
	assert.Nil(err)
	assert.True(len(resources.Resources) == 2)
	assert.True(len(resources.Resources["VLANPool"].Resources) == 1)
	assert.True(len(resources.Resources["IPAddressPool"].Resources) == 1)
}

func TestCreate(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	vlan1 := &network.VLANPool{
		BaseResource: zebra.BaseResource{
			ID:     "0100001",
			Type:   "VLANPool",
			Labels: nil,
		},
		RangeStart: 43,
		RangeEnd:   47,
	}

	vlan2 := &network.VLANPool{
		BaseResource: zebra.BaseResource{
			ID:     "0100002",
			Type:   "VLANPool",
			Labels: nil,
		},
		RangeStart: 18,
		RangeEnd:   30,
	}

	resMap := zebra.NewResourceMap(nil)

	ts := typestore.NewTypeStore(resMap)
	assert.NotNil(ts)

	// Create new resource, should pass
	assert.Nil(ts.Create(vlan1))

	resources, err := ts.Load()
	assert.Nil(err)
	assert.True(len(resources.Resources) == 1)
	assert.True(len(resources.Resources["VLANPool"].Resources) == 1)

	// Create another new resource, should pass
	assert.Nil(ts.Create(vlan2))

	// Create duplicate resource, should fail
	assert.NotNil(ts.Create(vlan1))

	// Create invalid resource, should fail
	assert.NotNil(ts.Create(new(network.VLANPool)))
}

func TestUpdate(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	vlan1 := &network.VLANPool{
		BaseResource: zebra.BaseResource{
			ID:     "0100001",
			Type:   "VLANPool",
			Labels: nil,
		},
		RangeStart: 43,
		RangeEnd:   47,
	}

	vlan2 := &network.VLANPool{
		BaseResource: zebra.BaseResource{
			ID:     "0100002",
			Type:   "VLANPool",
			Labels: nil,
		},
		RangeStart: 18,
		RangeEnd:   30,
	}

	resMap := zebra.NewResourceMap(nil)

	ts := typestore.NewTypeStore(resMap)
	assert.NotNil(ts)

	// Create new resource, should pass
	assert.Nil(ts.Create(vlan1))

	// Try to update, should pass
	assert.Nil(ts.Update(vlan1))

	// Try to update non-existent resource, should fail
	assert.NotNil(ts.Update(vlan2))

	// Try to update an invalid resource, should fail
	assert.NotNil(ts.Update(new(network.IPAddressPool)))
}

func TestDelete(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	vlan1 := &network.VLANPool{
		BaseResource: zebra.BaseResource{
			ID:     "0100001",
			Type:   "VLANPool",
			Labels: nil,
		},
		RangeStart: 43,
		RangeEnd:   47,
	}

	vlan2 := &network.VLANPool{
		BaseResource: zebra.BaseResource{
			ID:     "0100002",
			Type:   "VLANPool",
			Labels: nil,
		},
		RangeStart: 18,
		RangeEnd:   30,
	}

	resMap := zebra.NewResourceMap(nil)

	ts := typestore.NewTypeStore(resMap)
	assert.NotNil(ts)

	// Create new resource, should pass
	assert.Nil(ts.Create(vlan1))

	// Delete resource, should pass
	assert.Nil(ts.Delete(vlan1))

	// Try to delete non-existent resource, should pass anyways
	assert.Nil(ts.Delete(vlan2))

	// Try to delete invalid resource, should fail
	assert.NotNil(ts.Delete(new(network.Switch)))
}

func TestQuery(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resMap := zebra.NewResourceMap(nil)
	resMap.Add(new(network.VLANPool), "VLANPool")
	resMap.Add(new(network.IPAddressPool), "IPAddressPool")

	ts := typestore.NewTypeStore(resMap)
	assert.NotNil(ts)

	resources := ts.Query([]string{})
	assert.True(len(resources.Resources) == 0)

	resources = ts.Query([]string{"VLANPool"})
	assert.True(len(resources.Resources) == 1)
	assert.True(len(resources.Resources["VLANPool"].Resources) == 1)

	resources = ts.Query([]string{"VLANPool", "IPAddressPool"})
	assert.True(len(resources.Resources) == 2)
	assert.True(len(resources.Resources["VLANPool"].Resources) == 1)
	assert.True(len(resources.Resources["IPAddressPool"].Resources) == 1)
}
