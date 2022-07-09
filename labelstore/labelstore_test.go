package labelstore_test

import (
	"testing"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/labelstore"
	"github.com/project-safari/zebra/network"
	"github.com/stretchr/testify/assert"
)

func TestNewTestStore(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resMap := zebra.NewResourceMap(nil)
	assert.NotNil(labelstore.NewLabelStore(resMap))
}

func TestInitialize(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resMap := zebra.NewResourceMap(nil)

	ls := labelstore.NewLabelStore(resMap)
	assert.NotNil(ls)

	assert.Nil(ls.Initialize())
}

func TestWipe(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resMap := zebra.NewResourceMap(nil)

	ls := labelstore.NewLabelStore(resMap)
	assert.NotNil(ls)

	assert.Nil(ls.Wipe())
}

func TestClear(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	vlan1 := &network.VLANPool{
		BaseResource: zebra.BaseResource{
			ID:   "0100001",
			Type: "VLANPool",
			Labels: zebra.Labels{
				"owner": "Shravya Nandyala",
			},
		},
		RangeStart: 43,
		RangeEnd:   47,
	}

	resMap := zebra.NewResourceMap(nil)
	resMap.Add(vlan1, "VLANPool")

	ls := labelstore.NewLabelStore(resMap)
	assert.NotNil(ls)

	resources, err := ls.Load()
	assert.Nil(err)
	assert.True(len(resources.Resources) == 1)
	assert.True(len(resources.Resources["owner = Shravya Nandyala"].Resources) == 1)

	assert.Nil(ls.Clear())

	resources, err = ls.Load()
	assert.Nil(err)
	assert.True(len(resources.Resources) == 0)
}

func TestLoad(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	vlan1 := &network.VLANPool{
		BaseResource: zebra.BaseResource{
			ID:   "0100001",
			Type: "VLANPool",
			Labels: zebra.Labels{
				"owner":  "Shravya Nandyala",
				"player": "1",
			},
		},
		RangeStart: 43,
		RangeEnd:   47,
	}

	vlan2 := &network.VLANPool{
		BaseResource: zebra.BaseResource{
			ID:   "0100002",
			Type: "VLANPool",
			Labels: zebra.Labels{
				"owner":  "Shravya Nandyala",
				"player": "2",
			},
		},
		RangeStart: 1,
		RangeEnd:   2,
	}

	resMap := zebra.NewResourceMap(nil)
	resMap.Add(vlan1, "VLANPool")
	resMap.Add(vlan2, "VLANPool")

	ls := labelstore.NewLabelStore(resMap)
	assert.NotNil(ls)

	resources, err := ls.Load()
	assert.Nil(err)
	assert.True(len(resources.Resources) == 3)
	assert.True(len(resources.Resources["owner = Shravya Nandyala"].Resources) == 2)
	assert.True(len(resources.Resources["player = 1"].Resources) == 1)
	assert.True(len(resources.Resources["player = 2"].Resources) == 1)
}

func TestCreate(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	vlan1 := &network.VLANPool{
		BaseResource: zebra.BaseResource{
			ID:   "0100001",
			Type: "VLANPool",
			Labels: zebra.Labels{
				"a": "i",
				"b": "j",
				"c": "k",
			},
		},
		RangeStart: 43,
		RangeEnd:   47,
	}

	vlan2 := &network.VLANPool{
		BaseResource: zebra.BaseResource{
			ID:   "0100002",
			Type: "VLANPool",
			Labels: zebra.Labels{
				"a": "1",
				"b": "2",
				"c": "3",
			},
		},
		RangeStart: 18,
		RangeEnd:   30,
	}

	resMap := zebra.NewResourceMap(nil)

	ls := labelstore.NewLabelStore(resMap)
	assert.NotNil(ls)

	// Create new resource, should pass
	assert.Nil(ls.Create(vlan1))

	resources, err := ls.Load()
	assert.Nil(err)
	assert.True(len(resources.Resources) == 3)
	assert.True(len(resources.Resources["a = i"].Resources) == 1)

	// Create another new resource, should pass
	assert.Nil(ls.Create(vlan2))

	// Create duplicate resource, should fail
	assert.NotNil(ls.Create(vlan1))
}

func TestUpdate(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	vlan1 := &network.VLANPool{
		BaseResource: zebra.BaseResource{
			ID:   "0100001",
			Type: "VLANPool",
			Labels: zebra.Labels{
				"a": "i",
				"b": "j",
				"c": "k",
			},
		},
		RangeStart: 43,
		RangeEnd:   47,
	}

	vlan2 := &network.VLANPool{
		BaseResource: zebra.BaseResource{
			ID:   "0100002",
			Type: "VLANPool",
			Labels: zebra.Labels{
				"a": "1",
				"b": "2",
				"c": "3",
			},
		},
		RangeStart: 18,
		RangeEnd:   30,
	}

	resMap := zebra.NewResourceMap(nil)

	ls := labelstore.NewLabelStore(resMap)
	assert.NotNil(ls)

	// Create new resource, should pass
	assert.Nil(ls.Create(vlan1))

	// Try to update, should pass
	assert.Nil(ls.Update(vlan1))

	// Try to update non-existent resource, should fail
	assert.NotNil(ls.Update(vlan2))

	// Try to update an invalid resource, should fail
	assert.NotNil(ls.Update(new(network.IPAddressPool)))
}

func TestDelete(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	vlan1 := &network.VLANPool{
		BaseResource: zebra.BaseResource{
			ID:   "0100001",
			Type: "VLANPool",
			Labels: zebra.Labels{
				"a": "i",
				"b": "j",
				"c": "k",
			},
		},
		RangeStart: 43,
		RangeEnd:   47,
	}

	vlan2 := &network.VLANPool{
		BaseResource: zebra.BaseResource{
			ID:   "0100002",
			Type: "VLANPool",
			Labels: zebra.Labels{
				"a": "1",
				"b": "2",
				"c": "3",
			},
		},
		RangeStart: 18,
		RangeEnd:   30,
	}

	resMap := zebra.NewResourceMap(nil)

	ls := labelstore.NewLabelStore(resMap)
	assert.NotNil(ls)

	// Create new resource, should pass
	assert.Nil(ls.Create(vlan1))

	// Delete resource, should pass
	assert.Nil(ls.Delete(vlan1))

	// Try to delete non-existent resource, should pass anyways
	assert.Nil(ls.Delete(vlan2))
}

func TestQuery(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	vlan1 := &network.VLANPool{
		BaseResource: zebra.BaseResource{
			ID:   "0100001",
			Type: "VLANPool",
			Labels: zebra.Labels{
				"a": "i",
			},
		},
		RangeStart: 43,
		RangeEnd:   47,
	}

	vlan2 := &network.VLANPool{
		BaseResource: zebra.BaseResource{
			ID:   "0100002",
			Type: "VLANPool",
			Labels: zebra.Labels{
				"a": "1",
			},
		},
		RangeStart: 18,
		RangeEnd:   30,
	}

	resMap := zebra.NewResourceMap(nil)
	resMap.Add(vlan1, "VLANPool")
	resMap.Add(vlan2, "VLANPool")

	ls := labelstore.NewLabelStore(resMap)
	assert.NotNil(ls)

	query1 := labelstore.Query{Op: labelstore.MatchEqual, Key: "a", Values: []string{"i", "1"}}
	query2 := labelstore.Query{Op: labelstore.MatchIn, Key: "a", Values: []string{"i", "1"}}
	query3 := labelstore.Query{Op: labelstore.MatchNotEqual, Key: "a", Values: []string{"i", "1"}}
	query4 := labelstore.Query{Op: labelstore.MatchNotIn, Key: "a", Values: []string{"i"}}
	invalid := labelstore.Query{Op: 11, Key: "a", Values: []string{"i"}}

	assert.Nil(ls.Query(query1))

	query1.Values = []string{"i"}
	resources := ls.Query(query1)
	assert.True(len(resources.Resources) == 1)
	assert.True(len(resources.Resources["VLANPool"].Resources) == 1)

	resources = ls.Query(query2)
	assert.True(len(resources.Resources) == 1)
	assert.True(len(resources.Resources["VLANPool"].Resources) == 2)

	assert.Nil(ls.Query(query3))

	query3.Values = []string{"i"}
	assert.True(len(ls.Query(query3).Resources) == 1)

	assert.True(len(ls.Query(query4).Resources) == 1)

	assert.Nil(ls.Query(invalid))
}
