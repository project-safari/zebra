package store_test

import (
	"testing"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/store"
	"github.com/stretchr/testify/assert"
)

// Test function for  a new  label store.
func TestNewTestStore(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resMap := zebra.NewResourceMap(factory())
	assert.NotNil(store.NewLabelStore(resMap))
}

// Test function for initializing a label store.
func TestLSInitialize(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resMap := zebra.NewResourceMap(factory())

	ls := store.NewLabelStore(resMap)
	assert.NotNil(ls)
	assert.Nil(ls.Initialize())
}

// Test function for label store wipe.
func TestLSWipe(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resMap := zebra.NewResourceMap(factory())

	ls := store.NewLabelStore(resMap)
	assert.NotNil(ls)
	assert.Nil(ls.Initialize())

	assert.Nil(ls.Wipe())
}

// Test function for label store clear.
func TestLSClear(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	f := factory()
	r1 := f.New("dummy-1")
	r1.GetMeta().Labels.Add("a", "b")

	resMap := zebra.NewResourceMap(f)
	assert.Nil(resMap.Add(r1))

	ls := store.NewLabelStore(resMap)
	assert.NotNil(ls)
	assert.Nil(ls.Initialize())

	resources, err := ls.Load()
	assert.Nil(err)
	assert.Equal(1, len(resources.Resources))
	assert.Equal(1, len(resources.Resources["dummy-1"].Resources))

	assert.Nil(ls.Clear())

	resources, err = ls.Load()
	assert.Nil(err)
	assert.Empty(len(resources.Resources))
}

// Test function for loading a label store.
func TestLSLoad(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	f := factory()
	r1 := f.New("dummy-1")
	r1.GetMeta().Labels.Add("owner", "test_owner")
	r1.GetMeta().Labels.Add("player", "1")

	r2 := f.New("dummy-1")
	r2.GetMeta().Labels.Add("owner", "test_owner")
	r2.GetMeta().Labels.Add("player", "2")

	resMap := zebra.NewResourceMap(f)
	assert.Nil(resMap.Add(r1))
	assert.Nil(resMap.Add(r2))

	ls := store.NewLabelStore(resMap)
	assert.NotNil(ls)
	assert.Nil(ls.Initialize())

	resources, err := ls.Load()
	assert.Nil(err)

	for t, l := range resources.Resources {
		assert.Equal("dummy-1", t)

		for _, r := range l.Resources {
			assert.GreaterOrEqual(len(r.GetMeta().Labels), 2)
			assert.True(r.GetMeta().Labels.HasKey("owner"))
			assert.True(r.GetMeta().Labels.HasKey("player"))
			assert.Equal("test_owner", r.GetMeta().Labels["owner"])
			assert.True(r.GetMeta().Labels.MatchIn("player", "1", "2"))
		}
	}
}

// Test function for creating a new resource for a label store.
func TestLSCreate(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	f := factory()
	r1 := f.New("dummy-1")
	r1.GetMeta().Labels.Add("a", "i")
	r1.GetMeta().Labels.Add("b", "j")
	r1.GetMeta().Labels.Add("c", "k")

	r2 := f.New("dummy-1")
	r2.GetMeta().Labels.Add("i", "1")
	r2.GetMeta().Labels.Add("j", "2")
	r2.GetMeta().Labels.Add("k", "3")

	resMap := zebra.NewResourceMap(f)

	ls := store.NewLabelStore(resMap)
	assert.NotNil(ls)
	assert.Nil(ls.Initialize())

	// Create new resource, should pass
	assert.Nil(ls.Create(r1))

	resources, err := ls.Load()
	assert.Nil(err)
	assert.Len(resources.Resources, 1)
	assert.Len(resources.Resources["dummy-1"].Resources, 1)
	assert.Equal(r1.GetMeta().ID, resources.Resources["dummy-1"].Resources[0].GetMeta().ID)

	// Create another new resource, should pass
	assert.Nil(ls.Create(r2))

	// Create duplicate resource, should update
	assert.Nil(ls.Create(r1))
}

// Test function for a label store delete.
func TestLSDelete(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	f := factory()
	r1 := f.New("dummy-1")
	r1.GetMeta().Labels.Add("a", "i")
	r1.GetMeta().Labels.Add("b", "j")
	r1.GetMeta().Labels.Add("c", "k")

	r2 := f.New("dummy-1")
	r2.GetMeta().Labels.Add("i", "1")
	r2.GetMeta().Labels.Add("j", "2")
	r2.GetMeta().Labels.Add("k", "3")

	resMap := zebra.NewResourceMap(f)

	ls := store.NewLabelStore(resMap)
	assert.NotNil(ls)
	assert.Nil(ls.Initialize())

	// Create new resource, should pass
	assert.Nil(ls.Create(r1))

	// Delete resource, should pass
	assert.Nil(ls.Delete(r1))

	// Try to delete non-existent resource, should fail
	assert.NotNil(ls.Delete(r2))
}

// Test function for a label store query.
func TestLSQuery(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	f := factory()
	r1 := f.New("dummy-1")
	r1.GetMeta().Labels.Add("a", "i")
	r1.GetMeta().Labels.Add("b", "j")
	r1.GetMeta().Labels.Add("c", "k")

	r2 := f.New("dummy-1")
	r2.GetMeta().Labels.Add("a", "1")
	r2.GetMeta().Labels.Add("b", "2")
	r2.GetMeta().Labels.Add("c", "3")

	resMap := zebra.NewResourceMap(f)
	assert.Nil(resMap.Add(r1))
	assert.Nil(resMap.Add(r2))

	ls := store.NewLabelStore(resMap)
	assert.NotNil(ls)
	assert.Nil(ls.Initialize())

	query1 := zebra.Query{Op: zebra.MatchIn, Key: "a", Values: []string{"i", "1"}}
	query2 := zebra.Query{Op: zebra.MatchNotIn, Key: "a", Values: []string{"i"}}

	resources := ls.Query(query1)
	assert.Equal(1, len(resources.Resources))
	assert.Equal(2, len(resources.Resources["dummy-1"].Resources))

	query1.Key = "b"
	resources = ls.Query(query1)
	assert.Equal(0, len(resources.Resources))

	resources = ls.Query(query2)
	assert.Equal(1, len(resources.Resources))
}
