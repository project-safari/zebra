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

	assert.Nil(rs.Wipe())
}

func TestClear(t *testing.T) {
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
	resMap.Add(vlan1, "VLANPool")
	resMap.Add(vlan2, "VLANPool")

	rs := idstore.NewIDStore(resMap)
	assert.NotNil(rs)

	resources, err := rs.Load()
	assert.Nil(err)
	assert.True(len(resources.Resources) == 2)

	assert.Nil(rs.Clear())

	resources, err = rs.Load()
	assert.Nil(err)
	assert.True(len(resources.Resources) == 0)
}

func TestLoad(t *testing.T) {
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
	resMap.Add(vlan1, "VLANPool")
	resMap.Add(vlan2, "VLANPool")

	rs := idstore.NewIDStore(resMap)
	assert.NotNil(rs)

	resources, err := rs.Load()
	assert.Nil(err)
	assert.True(len(resources.Resources) == 2)
	assert.True(len(resources.Resources["0100001"].Resources) == 1)
	assert.True(len(resources.Resources["0100002"].Resources) == 1)
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

	rs := idstore.NewIDStore(resMap)
	assert.NotNil(rs)

	// Create new resource, should pass
	assert.Nil(rs.Create(vlan1))

	resources, err := rs.Load()
	assert.Nil(err)
	assert.True(len(resources.Resources) == 1)
	assert.True(len(resources.Resources["0100001"].Resources) == 1)

	// Create another new resource, should pass
	assert.Nil(rs.Create(vlan2))

	// Create duplicate resource, should fail
	assert.NotNil(rs.Create(vlan1))
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

	rs := idstore.NewIDStore(resMap)
	assert.NotNil(rs)

	// Create new resource, should pass
	assert.Nil(rs.Create(vlan1))

	// Try to update, should pass
	assert.Nil(rs.Update(vlan1))

	// Try to update non-existent resource, should fail
	assert.NotNil(rs.Update(vlan2))

	// Try to update an invalid resource, should fail
	assert.NotNil(rs.Update(new(network.IPAddressPool)))
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

	rs := idstore.NewIDStore(resMap)
	assert.NotNil(rs)

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
	resMap.Add(vlan1, "VLANPool")
	resMap.Add(vlan2, "VLANPool")

	rs := idstore.NewIDStore(resMap)
	assert.NotNil(rs)

	resources := rs.Query([]string{})
	assert.True(len(resources.Resources) == 0)

	resources = rs.Query([]string{"0100001"})
	assert.True(len(resources.Resources) == 1)
	assert.True(len(resources.Resources["VLANPool"].Resources) == 1)

	resources = rs.Query([]string{"0100001", "0100002"})
	assert.True(len(resources.Resources) == 1)
	assert.True(len(resources.Resources["VLANPool"].Resources) == 2)
}
