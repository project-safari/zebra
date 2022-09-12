package zebra_test

import (
	"testing"

	"github.com/project-safari/zebra"
	"github.com/stretchr/testify/assert"
)

func dummyType() (zebra.Type, zebra.TypeConstructor) {
	t := zebra.Type{
		Name:        "dummy",
		Description: "dummy type",
	}

	return t, func() zebra.Resource {
		return zebra.NewBaseResource(t, "dummy", "dummy", "dummy")
	}
}

func TestAddNew(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	d := zebra.Type{}
	assert.Equal(zebra.ErrTypeNameEmpty, d.Validate())
	d.Name = "dummy"
	assert.Equal(zebra.ErrTypeDescriptionEmpty, d.Validate())
	d.Description = "dummy test"
	assert.Nil(d.Validate())

	f := zebra.Factory()
	assert.NotNil(f)

	f.Add(dummyType())
	assert.NotNil(f.New("dummy"))
	assert.Nil(f.New("random"))

	_, ok := f.Type("random")
	assert.False(ok)
}

func TestFactory(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	f := zebra.Factory()
	f.Add(dummyType())

	assert.NotEmpty(f.Types())
	aType, ok := f.Type("dummy")
	assert.True(ok)
	assert.NotNil(aType)

	resA := zebra.NewResourceMap(f)
	assert.NotNil(resA)

	assert.NotNil(resA.Factory())
}

func TestAdd(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	funMap := zebra.Factory()
	funMap.Add(dummyType())

	resA := zebra.NewResourceMap(funMap)
	assert.NotNil(resA)

	d1 := funMap.New("dummy")

	assert.Nil(resA.Add(d1))
	assert.NotNil(len(resA.Resources["dummy"].Resources) == 1)
}
