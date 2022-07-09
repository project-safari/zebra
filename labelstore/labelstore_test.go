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

	vlan1 := getVLAN()
	vlan1.Labels = zebra.Labels{"owner": "Shravya Nandyala"}

	resMap := zebra.NewResourceMap(nil)
	resMap.Add(vlan1, "VLANPool")

	ls := labelstore.NewLabelStore(resMap)
	assert.NotNil(ls)

	resources, err := ls.Load()
	assert.Nil(err)
	assert.Equal(1, len(resources.Resources))
	assert.Equal(1, len(resources.Resources["owner = Shravya Nandyala"].Resources))

	assert.Nil(ls.Clear())

	resources, err = ls.Load()
	assert.Nil(err)
	assert.Empty(len(resources.Resources))
}

func TestLoad(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	vlan1 := getVLAN()
	vlan1.Labels = zebra.Labels{
		"owner":  "Shravya Nandyala",
		"player": "1",
	}

	vlan2 := getVLAN()
	vlan2.Labels = zebra.Labels{
		"owner":  "Shravya Nandyala",
		"player": "2",
	}

	resMap := zebra.NewResourceMap(nil)
	resMap.Add(vlan1, "VLANPool")
	resMap.Add(vlan2, "VLANPool")

	ls := labelstore.NewLabelStore(resMap)
	assert.NotNil(ls)

	resources, err := ls.Load()
	assert.Nil(err)
	assert.Equal(3, len(resources.Resources))
	assert.Equal(2, len(resources.Resources["owner = Shravya Nandyala"].Resources))
	assert.Equal(1, len(resources.Resources["player = 1"].Resources))
	assert.Equal(1, len(resources.Resources["player = 2"].Resources))
}

func TestCreate(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	vlan1 := getVLAN()
	vlan1.Labels = zebra.Labels{
		"a": "i",
		"b": "j",
		"c": "k",
	}

	vlan2 := getVLAN()
	vlan2.Labels = zebra.Labels{
		"a": "1",
		"b": "2",
		"c": "3",
	}

	resMap := zebra.NewResourceMap(nil)

	ls := labelstore.NewLabelStore(resMap)
	assert.NotNil(ls)

	// Create new resource, should pass
	assert.Nil(ls.Create(vlan1))

	resources, err := ls.Load()
	assert.Nil(err)
	assert.Equal(3, len(resources.Resources))
	assert.Equal(1, len(resources.Resources["a = i"].Resources))

	// Create another new resource, should pass
	assert.Nil(ls.Create(vlan2))

	// Create duplicate resource, should update
	assert.Nil(ls.Create(vlan1))
}

func TestDelete(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	vlan1 := getVLAN()
	vlan1.Labels = zebra.Labels{
		"a": "i",
		"b": "j",
		"c": "k",
	}

	vlan2 := getVLAN()
	vlan2.Labels = zebra.Labels{
		"a": "1",
		"b": "2",
		"c": "3",
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

	vlan1 := getVLAN()
	vlan1.Labels = zebra.Labels{
		"a": "i",
	}

	vlan2 := getVLAN()
	vlan2.Labels = zebra.Labels{
		"a": "1",
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
	assert.Equal(1, len(resources.Resources))
	assert.Equal(1, len(resources.Resources["VLANPool"].Resources))

	resources = ls.Query(query2)
	assert.Equal(1, len(resources.Resources))
	assert.Equal(2, len(resources.Resources["VLANPool"].Resources))

	assert.Nil(ls.Query(query3))

	query3.Values = []string{"i"}
	assert.Equal(1, len(ls.Query(query3).Resources))
	assert.Equal(1, len(ls.Query(query4).Resources))
	assert.Nil(ls.Query(invalid))
}

func getVLAN() *network.VLANPool {
	return &network.VLANPool{
		BaseResource: *zebra.NewBaseResource("VLANPool", nil),
		RangeStart:   0,
		RangeEnd:     1,
	}
}
