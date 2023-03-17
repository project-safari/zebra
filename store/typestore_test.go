package store_test

import (
	"testing"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/store"
	"github.com/stretchr/testify/assert"
)

// Test function for creating a new type store.
func TestTSNewTestStore(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resMap := zebra.NewResourceMap(factory())
	assert.NotNil(store.NewTypeStore(resMap))
}

// Test function for initializing a type store.
func TestTSInitialize(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resMap := zebra.NewResourceMap(factory())

	ts := store.NewTypeStore(resMap)
	assert.NotNil(ts)
	assert.Nil(ts.Initialize())
}

// Test function for type store wipe.
func TestTSWipe(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resMap := zebra.NewResourceMap(factory())

	ts := store.NewTypeStore(resMap)
	assert.NotNil(ts)
	assert.Nil(ts.Initialize())

	assert.Nil(ts.Wipe())
}

// Test function for type store clear.
func TestTSClear(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	f := factory()

	resMap := zebra.NewResourceMap(f)
	assert.Nil(resMap.Add(f.New("dummy-1")))
	assert.Nil(resMap.Add(f.New("dummy-2")))

	ts := store.NewTypeStore(resMap)
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

// Test function for type store load.
func TestTSLoad(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	f := factory()

	resMap := zebra.NewResourceMap(f)
	assert.Nil(resMap.Add(f.New("dummy-1")))
	assert.Nil(resMap.Add(f.New("dummy-2")))

	ts := store.NewTypeStore(resMap)
	assert.NotNil(ts)
	assert.Nil(ts.Initialize())

	resources, err := ts.Load()
	assert.Nil(err)
	assert.Len(resources.Resources, 2)
	assert.Len(resources.Resources["dummy-1"].Resources, 1)
	assert.Len(resources.Resources["dummy-2"].Resources, 1)
}

// Test function for creating a resource in a type store.
func TestTSCreate(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	f := factory()

	t1 := f.New("dummy-1")
	t2 := f.New("dummy-1")

	resMap := zebra.NewResourceMap(f)

	ts := store.NewTypeStore(resMap)
	assert.NotNil(ts)
	assert.Nil(ts.Initialize())

	// Create new resource, should pass
	assert.Nil(ts.Create(t1))

	resources, err := ts.Load()
	assert.Nil(err)
	assert.Len(resources.Resources, 1)
	assert.Len(resources.Resources["dummy-1"].Resources, 1)

	// Create duplicate resource, should update
	t1.GetMeta().Labels.Add("system.group", "newT1")
	assert.Nil(ts.Create(t1))

	resources, err = ts.Load()
	assert.Nil(err)

	labels := resources.Resources["dummy-1"].Resources[0].GetMeta().Labels
	assert.False(labels.MatchEqual("test", "true"))
	assert.True(labels.MatchEqual("system.group", "newT1"))

	// Create another new resource, should pass
	assert.Nil(ts.Create(t2))

	resources, err = ts.Load()
	assert.Nil(err)
	assert.Len(resources.Resources["dummy-1"].Resources, 2)
}

// Test function to delete resource in a type store.
func TestTSDelete(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	f := factory()

	t1 := f.New("dummy-1")
	t2 := f.New("dummy-1")

	resMap := zebra.NewResourceMap(f)

	ts := store.NewTypeStore(resMap)
	assert.NotNil(ts)
	assert.Nil(ts.Initialize())

	// Create new resource, should pass
	assert.Nil(ts.Create(t1))

	// Delete resource, should pass
	assert.Nil(ts.Delete(t1))

	// Try to delete non-existent resource, should fail
	assert.NotNil(ts.Delete(t2))
}

// Test function for a type store query.
func TestTSQuery(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	f := factory()
	resMap := zebra.NewResourceMap(f)
	assert.Nil(resMap.Add(f.New("dummy-1")))
	assert.Nil(resMap.Add(f.New("dummy-2")))

	ts := store.NewTypeStore(resMap)
	assert.NotNil(ts)
	assert.Nil(ts.Initialize())

	resources := ts.Query([]string{})
	assert.Empty(len(resources.Resources))

	resources = ts.Query([]string{"dummy-1"})
	assert.Len(resources.Resources, 1)

	resources = ts.Query([]string{"dummy-1", "dummy-2"})
	assert.Len(resources.Resources, 2)
	assert.Len(resources.Resources["dummy-1"].Resources, 1)
	assert.Len(resources.Resources["dummy-2"].Resources, 1)
}
