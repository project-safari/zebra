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
	factory.Add(network.VLANPoolType())

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
	factory.Add(network.VLANPoolType())

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
	factory.Add(network.VLANPoolType())

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
	factory.Add(network.VLANPoolType())

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
	assert.Nil(rs.Delete(getLab()))
}

func TestQuery(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore7"

	t.Cleanup(func() { os.RemoveAll(root) })

	factory := zebra.Factory()
	factory.Add(network.VLANPoolType())

	rs := store.NewResourceStore(root, factory)
	assert.NotNil(rs)
	assert.Nil(rs.Initialize())

	// Add 10 resources
	for i := 0; i < 10; i++ {
		assert.Nil(rs.Create(getVLAN()))
	}

	// Query for those 10 resources
	resources := rs.Query()
	assert.Equal(1, len(resources.Resources))
	assert.NotNil(resources.Resources["VLANPool"])
	assert.Equal(10, len(resources.Resources["VLANPool"].Resources))
}

func TestQueryUUID(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore8"

	t.Cleanup(func() { os.RemoveAll(root) })

	factory := zebra.Factory()
	factory.Add(network.VLANPoolType())

	rs := store.NewResourceStore(root, factory)
	assert.NotNil(rs)
	assert.Nil(rs.Initialize())

	ids := make([]string, 0, 5)

	// Add 10 resources
	for i := 0; i < 10; i++ {
		res := getVLAN()
		assert.Nil(rs.Create(res))

		if i%2 == 0 {
			ids = append(ids, res.ID)
		}
	}

	// Query for those 5 resources
	resources := rs.QueryUUID(ids)
	assert.Equal(1, len(resources.Resources))
	assert.NotNil(resources.Resources["VLANPool"])
	assert.Equal(5, len(resources.Resources["VLANPool"].Resources))
}

func TestQueryType(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore9"

	t.Cleanup(func() { os.RemoveAll(root) })

	factory := zebra.Factory()
	factory.Add(network.VLANPoolType())
	factory.Add(dc.LabType())

	rs := store.NewResourceStore(root, factory)
	assert.NotNil(rs)
	assert.Nil(rs.Initialize())

	// Add 10 resources
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			assert.Nil(rs.Create(getLab()))
		} else {
			assert.Nil(rs.Create(getVLAN()))
		}
	}

	// Query for those 5 resources
	resources := rs.QueryType([]string{"Lab"})
	assert.Equal(1, len(resources.Resources))
	assert.NotNil(resources.Resources["Lab"])
	assert.Equal(5, len(resources.Resources["Lab"].Resources))

	resources = rs.QueryType([]string{"VLANPool"})
	assert.Equal(1, len(resources.Resources))
	assert.NotNil(resources.Resources["VLANPool"])
	assert.Equal(5, len(resources.Resources["VLANPool"].Resources))
}

func TestQueryLabel(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore10"

	t.Cleanup(func() { os.RemoveAll(root) })

	factory := zebra.Factory()
	factory.Add(network.VLANPoolType())

	rs := store.NewResourceStore(root, factory)
	assert.NotNil(rs)
	assert.Nil(rs.Initialize())

	// Add 10 resources
	for i := 0; i < 10; i++ {
		res := getVLAN()
		if i%2 == 0 {
			res.Labels = zebra.Labels{"owner": "shravya"}
		}

		assert.Nil(rs.Create(res))
	}

	// Query for those 5 resources
	query := zebra.Query{Op: zebra.MatchEqual, Key: "owner", Values: []string{"shravya"}}
	resources, err := rs.QueryLabel(query)
	assert.Nil(err)
	assert.Equal(1, len(resources.Resources))
	assert.NotNil(resources.Resources["VLANPool"])
	assert.Equal(5, len(resources.Resources["VLANPool"].Resources))

	// Give incorrect query, should return error
	query = zebra.Query{Op: 10, Key: "", Values: []string{""}}
	_, err = rs.QueryLabel(query)
	assert.NotNil(err)
}

func TestQueryProperty(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore11"

	t.Cleanup(func() { os.RemoveAll(root) })

	res1, res2 := getVLAN(), getLab()

	f := zebra.Factory()
	f.Add(network.VLANPoolType())

	rs := store.NewResourceStore(root, f)
	assert.Nil(rs.Initialize())
	assert.Nil(rs.Create(res1))
	assert.Nil(rs.Create(res2))

	query1 := zebra.Query{Op: zebra.MatchEqual, Key: "Type", Values: []string{"VLANPool", "Lab"}}
	query2 := zebra.Query{Op: zebra.MatchIn, Key: "Type", Values: []string{"VLANPool"}}
	query3 := zebra.Query{Op: zebra.MatchNotEqual, Key: "Type", Values: []string{"VLANPool", "Lab"}}
	query4 := zebra.Query{Op: zebra.MatchNotIn, Key: "Type", Values: []string{"VLANPool", "Lab"}}

	// Should fail on query 1 and query 3.
	_, err := rs.QueryProperty(query1)
	assert.NotNil(err)

	_, err = rs.QueryProperty(query3)
	assert.NotNil(err)

	// Update query 1, should succeed.
	query1.Values = []string{"Lab"}
	resMap, err := rs.QueryProperty(query1)
	assert.Nil(err)
	assert.Equal(1, len(resMap.Resources))
	assert.Equal(1, len(resMap.Resources["Lab"].Resources))
	assert.Equal(res2.ID, resMap.Resources["Lab"].Resources[0].GetID())

	// Should succeed on query 2, return first resource.
	resMap, err = rs.QueryProperty(query2)
	assert.Nil(err)
	assert.Equal(1, len(resMap.Resources))
	assert.Equal(res1.ID, resMap.Resources["VLANPool"].Resources[0].GetID())

	// Should succeed on query 4, return no resources.
	resMap, err = rs.QueryProperty(query4)
	assert.Nil(err)
	assert.Empty(len(resMap.Resources))

	// Update query 3 to be valid, return 1 resource.
	query3.Values = []string{"Lab"}
	resMap, err = rs.QueryProperty(query3)
	assert.Nil(err)
	assert.Equal(1, len(resMap.Resources))
	assert.Equal(res1.ID, resMap.Resources["VLANPool"].Resources[0].GetID())

	resMap, err = rs.QueryProperty(zebra.Query{Op: 0x7, Key: "", Values: []string{""}})
	assert.Nil(resMap)
	assert.NotNil(err)
}

func TestFilterUUID(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resMap := getResMap()

	vlan := getVLAN()
	id := vlan.ID

	resMap.Add(vlan, "VLANPool")

	resMap, err := store.FilterUUID([]string{id}, resMap)
	assert.Nil(err)
	assert.Equal(1, len(resMap.Resources))
	assert.Equal(1, len(resMap.Resources["VLANPool"].Resources))
	assert.Equal(id, resMap.Resources["VLANPool"].Resources[0].GetID())
}

func TestFilterType(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resMap := getResMap()

	lab := getLab()
	id := lab.ID

	resMap.Add(lab, "Lab")

	resMap, err := store.FilterType([]string{"Lab"}, resMap)
	assert.Nil(err)
	assert.Equal(1, len(resMap.Resources))
	assert.Equal(1, len(resMap.Resources["Lab"].Resources))
	assert.Equal(id, resMap.Resources["Lab"].Resources[0].GetID())

	resMap, err = store.FilterType([]string{"blah"}, resMap)
	assert.Nil(err)
	assert.Equal(0, len(resMap.Resources))
}

func TestFilterLabel(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resMap := getResMap()

	lab := getLab()
	id := lab.ID
	lab.Labels = zebra.Labels{"owner": "shravya"}

	resMap.Add(lab, "Lab")

	query := zebra.Query{Op: 10, Key: "owner", Values: []string{"shravya"}}
	resMap, err := store.FilterLabel(query, resMap)
	assert.NotNil(err)

	query = zebra.Query{Op: zebra.MatchEqual, Key: "owner", Values: []string{"shravya", "nandyala"}}
	resMap, err = store.FilterLabel(query, resMap)
	assert.NotNil(err)

	query = zebra.Query{Op: zebra.MatchNotEqual, Key: "owner", Values: []string{"shravya", "nandyala"}}
	resMap, err = store.FilterLabel(query, resMap)
	assert.NotNil(err)

	query = zebra.Query{Op: zebra.MatchEqual, Key: "owner", Values: []string{"shravya"}}
	resMap, err = store.FilterLabel(query, resMap)
	assert.Nil(err)
	assert.Equal(1, len(resMap.Resources))
	assert.Equal(1, len(resMap.Resources["Lab"].Resources))
	assert.Equal(id, resMap.Resources["Lab"].Resources[0].GetID())

	query = zebra.Query{Op: zebra.MatchNotEqual, Key: "owner", Values: []string{"shravya"}}
	resMap, err = store.FilterLabel(query, resMap)
	assert.Nil(err)
	assert.Equal(0, len(resMap.Resources))
}

func TestFilterProperty(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resMap := getResMap()

	lab := getLab()

	resMap.Add(lab, "Lab")
	resMap.Add(getVLAN(), "VLANPool")

	query := zebra.Query{Op: 10, Key: "Type", Values: []string{"Lab"}}
	resMap, err := store.FilterProperty(query, resMap)
	assert.NotNil(err)

	query = zebra.Query{Op: zebra.MatchEqual, Key: "Type", Values: []string{"Lab", "VLANPool"}}
	resMap, err = store.FilterProperty(query, resMap)
	assert.NotNil(err)

	query = zebra.Query{Op: zebra.MatchNotEqual, Key: "type", Values: []string{"Lab", "VLANPool"}}
	resMap, err = store.FilterProperty(query, resMap)
	assert.NotNil(err)

	query = zebra.Query{Op: zebra.MatchEqual, Key: "Type", Values: []string{"Lab"}}
	resMap, err = store.FilterProperty(query, resMap)
	assert.Nil(err)
	assert.Equal(1, len(resMap.Resources))
	assert.Equal(1, len(resMap.Resources["Lab"].Resources))
	assert.Equal(lab.ID, resMap.Resources["Lab"].Resources[0].GetID())

	query = zebra.Query{Op: zebra.MatchNotEqual, Key: "type", Values: []string{"Lab"}}
	resMap, err = store.FilterProperty(query, resMap)
	assert.Nil(err)
	assert.Equal(0, len(resMap.Resources))
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

func getResMap() *zebra.ResourceMap {
	// make 10 resources and add them to list
	resMap := zebra.NewResourceMap(nil)

	for i := 0; i < 10; i++ {
		resMap.Add(getVLAN(), "VLANPool")
	}

	return resMap
}
