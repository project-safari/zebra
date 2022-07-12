package idstore_test

import (
	"testing"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/idstore"
	"github.com/project-safari/zebra/network"
	"github.com/stretchr/testify/assert"
)

func TestNewIDStore(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resMap := zebra.NewResourceMap(nil)
	assert.NotNil(idstore.NewIDStore(resMap))
}

func TestInitialize(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resMap := zebra.NewResourceMap(nil)

	rs := idstore.NewIDStore(resMap)
	assert.NotNil(rs)

	assert.Nil(rs.Initialize())
}

func TestWipe(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resMap := zebra.NewResourceMap(nil)

	rs := idstore.NewIDStore(resMap)
	assert.NotNil(rs)
	assert.Nil(rs.Initialize())

	assert.Nil(rs.Wipe())
}

func TestClear(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	vlan1 := getVLAN()
	vlan2 := getVLAN()

	resMap := zebra.NewResourceMap(nil)
	resMap.Add(vlan1, "VLANPool")
	resMap.Add(vlan2, "VLANPool")

	rs := idstore.NewIDStore(resMap)
	assert.NotNil(rs)
	assert.Nil(rs.Initialize())

	resources, err := rs.Load()
	assert.Nil(err)
	assert.Equal(2, len(resources.Resources))

	assert.Nil(rs.Clear())

	resources, err = rs.Load()
	assert.Nil(err)
	assert.Empty(len(resources.Resources))
}

func TestLoad(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	vlan1 := getVLAN()
	vlan2 := getVLAN()

	resMap := zebra.NewResourceMap(nil)
	resMap.Add(vlan1, "VLANPool")
	resMap.Add(vlan2, "VLANPool")

	rs := idstore.NewIDStore(resMap)
	assert.NotNil(rs)
	assert.Nil(rs.Initialize())

	resources, err := rs.Load()
	assert.Nil(err)
	assert.Equal(2, len(resources.Resources))
	assert.Equal(1, len(resources.Resources[vlan1.ID].Resources))
	assert.Equal(1, len(resources.Resources[vlan2.ID].Resources))
}

func TestCreate(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	vlan1 := getVLAN()
	vlan2 := getVLAN()

	resMap := zebra.NewResourceMap(nil)

	rs := idstore.NewIDStore(resMap)
	assert.NotNil(rs)
	assert.Nil(rs.Initialize())

	// Create new resource, should pass
	assert.Nil(rs.Create(vlan1))

	resources, err := rs.Load()
	assert.Nil(err)
	assert.Equal(1, len(resources.Resources))
	assert.Equal(1, len(resources.Resources[vlan1.ID].Resources))

	// Create another new resource, should pass
	assert.Nil(rs.Create(vlan2))

	// Create duplicate resource, should update
	assert.Nil(rs.Create(vlan1))
}

func TestDelete(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	vlan1 := getVLAN()
	vlan2 := getVLAN()

	resMap := zebra.NewResourceMap(nil)

	rs := idstore.NewIDStore(resMap)
	assert.NotNil(rs)
	assert.Nil(rs.Initialize())

	// Create new resource, should pass
	assert.Nil(rs.Create(vlan1))

	// Delete resource, should pass
	assert.Nil(rs.Delete(vlan1))

	// Try to delete non-existent resource, should pass anyways
	assert.Nil(rs.Delete(vlan2))
}

func TestQuery(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	vlan1 := getVLAN()

	vlan2 := getVLAN()

	resMap := zebra.NewResourceMap(nil)
	resMap.Add(vlan1, "VLANPool")
	resMap.Add(vlan2, "VLANPool")

	rs := idstore.NewIDStore(resMap)
	assert.NotNil(rs)
	assert.Nil(rs.Initialize())

	resources := rs.Query([]string{})
	assert.Empty(len(resources.Resources))

	resources = rs.Query([]string{vlan1.ID})
	assert.Equal(1, len(resources.Resources))
	assert.Equal(1, len(resources.Resources["VLANPool"].Resources))

	resources = rs.Query([]string{vlan1.ID, vlan2.ID})
	assert.Equal(1, len(resources.Resources))
	assert.Equal(2, len(resources.Resources["VLANPool"].Resources))

	resources = rs.Query([]string{"random id"})
	assert.Empty(resources.Resources)
}

func getVLAN() *network.VLANPool {
	return &network.VLANPool{
		BaseResource: *zebra.NewBaseResource("VLANPool", nil),
		RangeStart:   0,
		RangeEnd:     1,
	}
}
