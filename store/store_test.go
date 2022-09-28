package store_test

import (
	"os"
	"testing"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/store"
	"github.com/stretchr/testify/assert"
)

func getResMap() *zebra.ResourceMap {
	f := factory()

	// make 10 resources and add them to list
	resMap := zebra.NewResourceMap(f)

	for i := 0; i < 10; i++ {
		if e := resMap.Add(f.New("dummy-1")); e != nil {
			return nil
		}
	}

	return resMap
}

func TestNewResourceStore(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore_new_resource_Store"

	defer func() { os.RemoveAll(root) }()

	assert.NotNil(store.NewResourceStore(root, factory()))
}

func TestStoreInitialize(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore_initialize"

	defer func() { os.RemoveAll(root) }()

	rs := store.NewResourceStore(root, factory())
	assert.NotNil(rs)
	assert.Nil(rs.Initialize())
}

func TestStoreWipe(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore_Wipe"

	defer func() { os.RemoveAll(root) }()

	rs := store.NewResourceStore(root, factory())
	assert.NotNil(rs)
	assert.Nil(rs.Initialize())
	assert.Nil(rs.Wipe())
}

func TestStoreDelete(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore_delete"

	defer func() { os.RemoveAll(root) }()

	f := factory()

	rs := store.NewResourceStore(root, f)
	assert.NotNil(rs)
	assert.Nil(rs.Initialize())

	// Valid resource, should pass
	r := f.New("dummy-1")
	assert.Nil(rs.Create(r))

	resources, err := rs.Load()
	assert.Nil(err)
	assert.Equal(1, len(resources.Resources))

	// Delete resource, should pass
	assert.Nil(rs.Delete(r))

	// Delete non-existent resource, should fail
	assert.NotNil(rs.Delete(nil))

	// Delete uncreated resource, should pass anyways
	assert.NotNil(rs.Delete(f.New("dummy-2")))
}

func TestStoreLoad(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore_load"

	defer func() { os.RemoveAll(root) }()

	f := factory()

	rs := store.NewResourceStore(root, f)
	assert.NotNil(rs)
	assert.Nil(rs.Initialize())

	resources, err := rs.Load()
	assert.Nil(err)
	assert.Empty(len(resources.Resources))

	assert.Nil(rs.Create(f.New("dummy-1")))
	resources, err = rs.Load()
	assert.Nil(err)
	assert.Len(resources.Resources, 1)
	assert.Len(resources.Resources["dummy-1"].Resources, 1)

	assert.Nil(rs.Create(f.New("dummy-1")))

	resources, err = rs.Load()
	assert.Nil(err)
	assert.Len(resources.Resources, 1)
	assert.Len(resources.Resources["dummy-1"].Resources, 2)
}

func TestStoreCreate(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore_create"

	defer func() { os.RemoveAll(root) }()

	f := factory()

	rs := store.NewResourceStore(root, f)
	assert.NotNil(rs)
	assert.Nil(rs.Initialize())

	// Invalid resource, should fail
	assert.NotNil(rs.Create(nil))

	// Valid resource, should pass
	r := f.New("dummy-1")
	assert.Nil(rs.Create(r))

	resources, err := rs.Load()
	assert.Nil(err)
	assert.Len(resources.Resources, 1)
	assert.Len(resources.Resources["dummy-1"].Resources, 1)

	// Duplicate resource, should update
	assert.Nil(rs.Create(r))
	resources, err = rs.Load()
	assert.Nil(err)
	assert.Len(resources.Resources, 1)
	assert.Len(resources.Resources["dummy-1"].Resources, 1)
}

func TestStoreQueryLabel(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore_query_label"

	defer func() { os.RemoveAll(root) }()

	f := factory()

	rs := store.NewResourceStore(root, f)
	assert.NotNil(rs)
	assert.Nil(rs.Initialize())

	// Add 10 resources
	for i := 0; i < 10; i++ {
		res := f.New("dummy-1")

		if i%2 == 0 {
			res.GetMeta().Labels.Add("owner", "test_owner")
		}

		assert.Nil(rs.Create(res))
	}

	// Query for incorrect values
	query := zebra.Query{Op: zebra.MatchEqual, Key: "owner", Values: []string{"unknown"}}
	resources, err := rs.QueryLabel(query)
	assert.Nil(err)
	assert.Len(resources.Resources, 0)

	// Query for incorrect values
	query = zebra.Query{Op: zebra.MatchEqual, Key: "owner", Values: []string{"test_owner"}}
	resources, err = rs.QueryLabel(query)
	assert.Nil(err)
	assert.Len(resources.Resources["dummy-1"].Resources, 5)

	// Give incorrect query, should return error
	query = zebra.Query{Op: 10, Key: "", Values: []string{""}}
	_, err = rs.QueryLabel(query)
	assert.NotNil(err)
}

func TestStoreQuery(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore_query"

	defer func() { os.RemoveAll(root) }()

	f := factory()

	rs := store.NewResourceStore(root, f)
	assert.NotNil(rs)
	assert.Nil(rs.Initialize())

	// Add 10 resources
	for i := 0; i < 10; i++ {
		assert.Nil(rs.Create(f.New("dummy-1")))
	}

	// Query for those 10 resources
	resources := rs.Query()
	assert.Len(resources.Resources, 1)
	assert.Len(resources.Resources["dummy-1"].Resources, 10)
}

func TestStoreQueryUUID(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore_query_uuid"

	defer func() { os.RemoveAll(root) }()

	ids := make([]string, 0, 5)
	f := factory()

	rs := store.NewResourceStore(root, f)
	assert.NotNil(rs)
	assert.Nil(rs.Initialize())

	// Add 10 resources
	for i := 0; i < 10; i++ {
		r := f.New("dummy-1")
		assert.Nil(rs.Create(r))

		if i%2 == 0 {
			ids = append(ids, r.GetMeta().ID)
		}
	}

	// Query for those 5 resources
	resources := rs.QueryUUID(ids)
	assert.Len(resources.Resources, 1)
	assert.Len(resources.Resources["dummy-1"].Resources, 5)
}

func TestStoreQueryType(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore_query_type"

	defer func() { os.RemoveAll(root) }()

	f := factory()

	rs := store.NewResourceStore(root, f)
	assert.NotNil(rs)
	assert.Nil(rs.Initialize())

	// Add 10 resources
	for i := 0; i < 10; i++ {
		if i%2 == 0 {
			assert.Nil(rs.Create(f.New("dummy-1")))
		} else {
			assert.Nil(rs.Create(f.New("dummy-2")))
		}
	}

	// Query for those 5 resources
	resources := rs.QueryType([]string{"dummy-1"})
	assert.Len(resources.Resources, 1)
	assert.Len(resources.Resources["dummy-1"].Resources, 5)

	resources = rs.QueryType([]string{"dummy-2"})
	assert.Len(resources.Resources, 1)
	assert.Len(resources.Resources["dummy-2"].Resources, 5)
}

func TestStoreFilterUUID(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resMap := getResMap()
	f := factory()
	r1 := f.New("dummy-1")
	id := r1.GetMeta().ID

	assert.Nil(resMap.Add(r1))

	resMap, err := store.FilterUUID([]string{id}, resMap)
	assert.Nil(err)
	assert.Len(resMap.Resources, 1)
	assert.Len(resMap.Resources["dummy-1"].Resources, 1)
	assert.Equal(id, resMap.Resources["dummy-1"].Resources[0].GetMeta().ID)
}

func TestStoreFilterType(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resMap := getResMap()
	f := factory()
	r1 := f.New("dummy-2")
	id := r1.GetMeta().ID

	assert.Nil(resMap.Add(r1))

	resMap, err := store.FilterType([]string{"dummy-2"}, resMap)
	assert.Nil(err)
	assert.Len(resMap.Resources, 1)
	assert.Len(resMap.Resources["dummy-2"].Resources, 1)
	assert.Equal(id, resMap.Resources["dummy-2"].Resources[0].GetMeta().ID)

	resMap, err = store.FilterType([]string{"blah"}, resMap)
	assert.Nil(err)
	assert.Equal(0, len(resMap.Resources))
}

func TestStoreFilterLabel(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resMap := getResMap()

	f := factory()
	r := f.New("dummy-1")
	id := r.GetMeta().ID

	r.GetMeta().Labels.Add("owner", "test_owner")
	assert.Nil(resMap.Add(r))

	query := zebra.Query{Op: 10, Key: "owner", Values: []string{"test_owner"}}
	resMap, err := store.FilterLabel(query, resMap)
	assert.NotNil(err)

	query = zebra.Query{Op: zebra.MatchEqual, Key: "owner", Values: []string{"test_owner", "someone"}}
	resMap, err = store.FilterLabel(query, resMap)
	assert.NotNil(err)

	query = zebra.Query{Op: zebra.MatchNotEqual, Key: "owner", Values: []string{"test_owner", "someone"}}
	resMap, err = store.FilterLabel(query, resMap)
	assert.NotNil(err)

	query = zebra.Query{Op: zebra.MatchEqual, Key: "owner", Values: []string{"test_owner"}}
	resMap, err = store.FilterLabel(query, resMap)
	assert.Nil(err)
	assert.Len(resMap.Resources, 1)
	assert.Len(resMap.Resources["dummy-1"].Resources, 1)
	assert.Equal(id, resMap.Resources["dummy-1"].Resources[0].GetMeta().ID)

	query = zebra.Query{Op: zebra.MatchNotEqual, Key: "owner", Values: []string{"test_owner"}}
	resMap, err = store.FilterLabel(query, resMap)
	assert.Nil(err)
	assert.Empty(resMap.Resources)
}

type propRes struct {
	zebra.BaseResource
	Prop1 string `json:"prop1"`
	Prop2 string `json:"prop2"`
}

func propType() (zebra.Type, zebra.TypeConstructor) {
	t := zebra.Type{Name: "propType", Description: "dummy prop resource"}

	return t, func() zebra.Resource {
		r := new(propRes)
		r.BaseResource = *zebra.NewBaseResource(t, "prop_test", "prop_test_owner", "prop_test_group")
		r.Meta.Type = t

		return r
	}
}

func TestStoreFilterProperty(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	f := factory()
	resMap := zebra.NewResourceMap(f)
	ty, c := propType()
	f.Add(ty, c)

	r1, ok := f.New("propType").(*propRes)
	assert.True(ok)

	r1.Meta.Type = ty
	r1.Prop1 = "test_prop_value"

	assert.Nil(resMap.Add(r1))

	query := zebra.Query{Op: 10, Key: "Type", Values: []string{"dummy-2"}}
	resMap, err := store.FilterProperty(query, resMap)
	assert.NotNil(err)

	query = zebra.Query{Op: zebra.MatchEqual, Key: "Type", Values: []string{"dummy-1", "dummy-2"}}
	resMap, err = store.FilterProperty(query, resMap)
	assert.NotNil(err)

	query = zebra.Query{Op: zebra.MatchNotEqual, Key: "type", Values: []string{"dummy-1", "dummy-2"}}
	resMap, err = store.FilterProperty(query, resMap)
	assert.NotNil(err)

	query = zebra.Query{Op: zebra.MatchEqual, Key: "Prop1", Values: []string{"test_prop_value"}}
	resMap, err = store.FilterProperty(query, resMap)
	assert.Nil(err)
	assert.Len(resMap.Resources, 1)
	assert.Len(resMap.Resources["propType"].Resources, 1)
	assert.Equal(r1.GetMeta().ID, resMap.Resources["propType"].Resources[0].GetMeta().ID)

	query = zebra.Query{Op: zebra.MatchNotEqual, Key: "type", Values: []string{"Lab"}}
	resMap, err = store.FilterProperty(query, resMap)
	assert.Nil(err)
	assert.NotNil(resMap)
}

func TestStoreClear(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore_clear"

	defer func() { os.RemoveAll(root) }()

	f := factory()
	rs := store.NewResourceStore(root, f)

	assert.NotNil(rs)
	assert.Nil(rs.Initialize())

	assert.Nil(rs.Create(f.New("dummy-1")))
	assert.Nil(rs.Create(f.New("dummy-1")))

	resources, err := rs.Load()
	assert.Nil(err)
	assert.Len(resources.Resources, 1)

	assert.Nil(rs.Clear())

	resources, err = rs.Load()
	assert.Nil(err)
	assert.Empty(resources.Resources)
}

func TestStoreQueryProperty(t *testing.T) { //nolint:funlen
	t.Parallel()
	assert := assert.New(t)

	root := "teststore_query_property"

	defer func() { os.RemoveAll(root) }()

	f := factory()
	ty, c := propType()
	f.Add(ty, c)

	res1, ok := f.New("propType").(*propRes)
	assert.True(ok)

	res2, ok := f.New("propType").(*propRes)
	assert.True(ok)

	res1.Prop1 = "r1p1"
	res1.Prop2 = "r1p2"
	res2.Prop1 = "r2p1"
	res2.Prop2 = "r2p2"

	rs := store.NewResourceStore(root, f)
	assert.Nil(rs.Initialize())
	assert.Nil(rs.Create(res1))
	assert.Nil(rs.Create(res2))

	query1 := zebra.Query{Op: zebra.MatchEqual, Key: "Prop1", Values: []string{"r1p1", "r2p1"}}
	query2 := zebra.Query{Op: zebra.MatchIn, Key: "Prop1", Values: []string{"r1p1"}}
	query3 := zebra.Query{Op: zebra.MatchNotEqual, Key: "Prop1", Values: []string{"VLANPool", "Lab"}}
	query4 := zebra.Query{Op: zebra.MatchNotIn, Key: "Prop1", Values: []string{"r1p1", "r2p1"}}

	// Should fail on query 1 and query 3.
	_, err := rs.QueryProperty(query1)
	assert.NotNil(err)

	_, err = rs.QueryProperty(query3)
	assert.NotNil(err)

	// Update query 1, should succeed.
	query1.Values = []string{"r1p1"}
	resMap, err := rs.QueryProperty(query1)
	assert.Nil(err)
	assert.Len(resMap.Resources, 1)
	assert.Equal(res1.GetMeta().ID, resMap.Resources["propType"].Resources[0].GetMeta().ID)

	// Should succeed on query 2, return first resource.
	resMap, err = rs.QueryProperty(query2)
	assert.Nil(err)
	assert.Len(resMap.Resources, 1)
	assert.Equal(res1.GetMeta().ID, resMap.Resources["propType"].Resources[0].GetMeta().ID)

	// Should succeed on query 4, return no resources.
	resMap, err = rs.QueryProperty(query4)
	assert.Nil(err)
	assert.Empty(resMap.Resources)

	// Update query 3 to be valid, return 1 resource.
	query3.Values = []string{"r1p1"}
	resMap, err = rs.QueryProperty(query3)
	assert.Nil(err)
	assert.Len(resMap.Resources, 1)
	assert.Equal(res2.GetMeta().ID, resMap.Resources["propType"].Resources[0].GetMeta().ID)

	resMap, err = rs.QueryProperty(zebra.Query{Op: 0x7, Key: "", Values: []string{""}})
	assert.Nil(resMap)
	assert.NotNil(err)
}
