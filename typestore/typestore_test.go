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
	assert.Nil(ts.Initialize())

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
	assert.Nil(ts.Initialize())

	resources, err := ts.Load()
	assert.Nil(err)
	assert.Equal(2, len(resources.Resources))

	assert.Nil(ts.Clear())

	resources, err = ts.Load()
	assert.Nil(err)
	assert.Empty(len(resources.Resources))
}

func TestLoad(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resMap := zebra.NewResourceMap(nil)
	resMap.Add(new(network.VLANPool), "VLANPool")
	resMap.Add(new(network.IPAddressPool), "IPAddressPool")

	ts := typestore.NewTypeStore(resMap)
	assert.NotNil(ts)
	assert.Nil(ts.Initialize())

	resources, err := ts.Load()
	assert.Nil(err)
	assert.Equal(2, len(resources.Resources))
	assert.Equal(1, len(resources.Resources["VLANPool"].Resources))
	assert.Equal(1, len(resources.Resources["IPAddressPool"].Resources))
}

func TestCreate(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	vlan1 := getVLAN()
	vlan2 := getVLAN()

	vlan3 := getVLAN()
	vlan3.ID = vlan1.GetID()
	vlan3.Labels = zebra.Labels{
		"test": "true",
	}

	resMap := zebra.NewResourceMap(nil)

	ts := typestore.NewTypeStore(resMap)
	assert.NotNil(ts)
	assert.Nil(ts.Initialize())

	// Create new resource, should pass
	assert.Nil(ts.Create(vlan1))

	resources, err := ts.Load()
	assert.Nil(err)
	assert.Equal(1, len(resources.Resources))
	assert.Equal(1, len(resources.Resources["VLANPool"].Resources))

	// Create duplicate resource, should update
	assert.Nil(ts.Create(vlan3))

	resources, err = ts.Load()
	assert.Nil(err)

	labels := resources.Resources["VLANPool"].Resources[0].GetLabels()
	assert.True(labels.MatchEqual("test", "true"))

	// Create another new resource, should pass
	assert.Nil(ts.Create(vlan2))

	resources, err = ts.Load()
	assert.Nil(err)
	assert.Equal(2, len(resources.Resources["VLANPool"].Resources))
}

func TestDelete(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	vlan1 := getVLAN()
	vlan2 := getVLAN()

	resMap := zebra.NewResourceMap(nil)

	ts := typestore.NewTypeStore(resMap)
	assert.NotNil(ts)
	assert.Nil(ts.Initialize())

	// Create new resource, should pass
	assert.Nil(ts.Create(vlan1))

	// Delete resource, should pass
	assert.Nil(ts.Delete(vlan1))

	// Try to delete non-existent resource, should pass anyways
	assert.Nil(ts.Delete(vlan2))
}

func TestQuery(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resMap := zebra.NewResourceMap(nil)
	resMap.Add(new(network.VLANPool), "VLANPool")
	resMap.Add(new(network.IPAddressPool), "IPAddressPool")

	ts := typestore.NewTypeStore(resMap)
	assert.NotNil(ts)
	assert.Nil(ts.Initialize())

	resources := ts.Query([]string{})
	assert.Empty(len(resources.Resources))

	resources = ts.Query([]string{"VLANPool"})
	assert.Equal(1, len(resources.Resources))
	assert.Equal(1, len(resources.Resources["VLANPool"].Resources))

	resources = ts.Query([]string{"VLANPool", "IPAddressPool"})
	assert.Equal(2, len(resources.Resources))
	assert.Equal(1, len(resources.Resources["VLANPool"].Resources))
	assert.Equal(1, len(resources.Resources["IPAddressPool"].Resources))
}

func getVLAN() *network.VLANPool {
	return &network.VLANPool{
		BaseResource: *zebra.NewBaseResource(network.VLANPoolType(), nil),
		RangeStart:   0,
		RangeEnd:     1,
	}
}
