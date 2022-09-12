package store_test

import (
	"testing"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/store"
	"github.com/stretchr/testify/assert"
)

func TestNewIDStore(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resMap := zebra.NewResourceMap(factory())
	assert.NotNil(store.NewIDStore(resMap))
}

func TestInitialize(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resMap := zebra.NewResourceMap(factory())

	rs := store.NewIDStore(resMap)
	assert.NotNil(rs)

	assert.Nil(rs.Initialize())
}

func TestWipe(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resMap := zebra.NewResourceMap(factory())

	rs := store.NewIDStore(resMap)
	assert.NotNil(rs)
	assert.Nil(rs.Initialize())

	assert.Nil(rs.Wipe())
}

func TestClear(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	f := factory()
	r1 := f.New("dummy-1")
	r2 := f.New("dummy-2")

	resMap := zebra.NewResourceMap(f)
	assert.Nil(resMap.Add(r1))
	assert.Nil(resMap.Add(r2))

	rs := store.NewIDStore(resMap)
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

	f := factory()
	r1 := f.New("dummy-1")
	r2 := f.New("dummy-2")

	resMap := zebra.NewResourceMap(f)
	assert.Nil(resMap.Add(r1))
	assert.Nil(resMap.Add(r2))

	rs := store.NewIDStore(resMap)
	assert.NotNil(rs)
	assert.Nil(rs.Initialize())

	resources, err := rs.Load()
	assert.Nil(err)
	assert.Len(resources.Resources, 2)
	assert.Len(resources.Resources["dummy-1"].Resources, 1)
	assert.Len(resources.Resources["dummy-2"].Resources, 1)
}

func TestCreate(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	f := factory()
	r1 := f.New("dummy-1")
	r2 := f.New("dummy-1")

	resMap := zebra.NewResourceMap(f)

	rs := store.NewIDStore(resMap)
	assert.NotNil(rs)
	assert.Nil(rs.Initialize())

	// Create new resource, should pass
	assert.Nil(rs.Create(r1))

	resources, err := rs.Load()
	assert.Nil(err)
	assert.Len(resources.Resources, 1)
	assert.Len(resources.Resources["dummy-1"].Resources, 1)

	// Create another new resource, should pass
	assert.Nil(rs.Create(r2))

	// Create duplicate resource, should update
	assert.Nil(rs.Create(r1))
}

func TestDelete(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	f := factory()
	r1 := f.New("dummy-1")
	r2 := f.New("dummy-1")

	resMap := zebra.NewResourceMap(f)

	rs := store.NewIDStore(resMap)
	assert.NotNil(rs)
	assert.Nil(rs.Initialize())

	// Create new resource, should pass
	assert.Nil(rs.Create(r1))

	// Delete resource, should pass
	assert.Nil(rs.Delete(r1))

	// Try to delete non-existent resource, should pass anyways
	assert.Nil(rs.Delete(r2))
}

func TestQuery(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	f := factory()
	r1 := f.New("dummy-1")
	r2 := f.New("dummy-1")

	resMap := zebra.NewResourceMap(f)
	assert.Nil(resMap.Add(r1))
	assert.Nil(resMap.Add(r2))

	rs := store.NewIDStore(resMap)
	assert.NotNil(rs)
	assert.Nil(rs.Initialize())

	resources := rs.Query([]string{})
	assert.Empty(len(resources.Resources))

	resources = rs.Query([]string{r1.GetMeta().ID})
	assert.Equal(1, len(resources.Resources))
	assert.Equal(1, len(resources.Resources["dummy-1"].Resources))

	resources = rs.Query([]string{r1.GetMeta().ID, r2.GetMeta().ID})
	assert.Equal(1, len(resources.Resources))
	assert.Equal(2, len(resources.Resources["dummy-1"].Resources))

	resources = rs.Query([]string{"random id"})
	assert.Empty(resources.Resources)
}
