package zebra_test

import (
	"fmt"
	"testing"

	"github.com/project-safari/zebra"
	"github.com/stretchr/testify/assert"
)

// Mock function that returns a new zebra.Resource to be used in tests.
// It crates a new zebra.BaseResource of type dummy.
func dummyCtr() zebra.Resource {
	r := new(zebra.BaseResource)
	r.Meta.Type.Name = "dummy"

	return r
}

// Test function for creating resource lists.
func TestNewResourceList(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	_, ctr := dummyType()
	assert.NotNil(zebra.NewResourceList(ctr))
}

// Test function for copying resource lists.
func TestCopyResourceList(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	_, ctr := dummyType()
	resA := zebra.NewResourceList(ctr)
	assert.NotNil(resA)

	resA.Resources = append(resA.Resources, ctr())

	resB := zebra.NewResourceList(ctr)
	assert.NotNil(resB)
	assert.Empty(len(resB.Resources))

	zebra.CopyResourceList(resB, resA)
	assert.Equal(1, len(resB.Resources))

	zebra.CopyResourceList(nil, nil)
}

func TestListMarshalUnmarshal(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	dummy, ctr := dummyType()
	funMap := zebra.Factory()
	funMap.Add(dummy, ctr)

	resA := zebra.NewResourceList(ctr)
	assert.NotNil(resA)

	d := ctr()
	assert.Nil(resA.Add(d))

	bytes, err := resA.MarshalJSON()
	assert.Nil(err)
	assert.NotNil(bytes)

	resB := zebra.NewResourceList(ctr)
	assert.NotNil(resB)

	err = resB.UnmarshalJSON(bytes)
	assert.Nil(err)

	d1 := ctr()
	resA.Resources = []zebra.Resource{d1}

	bytes, err = resA.MarshalJSON()
	assert.Nil(err)
	assert.NotNil(bytes)

	resB = zebra.NewResourceList(ctr)
	assert.NotNil(resB)

	fmt.Println(string(bytes))
	err = resB.UnmarshalJSON(bytes)
	assert.Nil(err)
	assert.Equal(1, len(resB.Resources))
}

func TestErrorMarshalUnmarshal(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	funMap := zebra.Factory()
	funMap.Add(zebra.Type{"dummy", "dummy type"}, dummyCtr)

	resList := zebra.NewResourceList(dummyCtr)
	assert.NotNil(resList.UnmarshalJSON(nil))
	assert.NotNil(resList.UnmarshalJSON([]byte(`[{"id":"0100000001", "meta":123}]`)))

	resMap := zebra.NewResourceMap(funMap)
	assert.NotNil(resMap.UnmarshalJSON(nil))
	assert.NotNil(resMap.UnmarshalJSON([]byte(`{"VLANPool":[{"id":"0100000001", "type":123}]}`)))
}

// Test function for creating resource maps.
func TestNewResourceMap(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	assert.NotNil(zebra.NewResourceMap(nil))
}

// Test function for copying resource maps.
func TestCopyResourceMap(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	funMap := zebra.Factory()
	dummy, ctr := dummyType()
	funMap.Add(dummy, ctr)

	resA := zebra.NewResourceMap(funMap)
	assert.NotNil(resA)
	assert.NotNil(resA.Factory())

	assert.Nil(resA.Add(ctr()))

	resB := zebra.NewResourceMap(funMap)
	assert.NotNil(resB)

	zebra.CopyResourceMap(resB, nil)

	zebra.CopyResourceMap(resB, resA)
	assert.Equal(1, len(resB.Resources))
	assert.Equal(1, len(resB.Resources["dummy"].Resources))
}

// Test function for deliting resource lists and resource maps.
func TestDelete(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	funMap := zebra.Factory()
	funMap.Add(dummyType())

	resA := zebra.NewResourceMap(funMap)
	assert.NotNil(resA)

	d := funMap.New("dummy")

	assert.Nil(resA.Add(d))
	assert.NotNil(len(resA.Resources["dummy"].Resources) == 1)

	assert.Nil(resA.Delete(d))

	assert.NotNil(resA.Delete(d))

	_, ok := resA.Resources["dummy"]
	assert.NotNil(ok)
}

func TestMapMarshalUnMarshal(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	funMap := zebra.Factory()
	funMap.Add(zebra.Type{"dummy", "dummy type"}, dummyCtr)

	resA := zebra.NewResourceMap(funMap)
	assert.NotNil(resA)

	d := funMap.New("dummy")
	assert.Nil(resA.Add(d))

	bytes, err := resA.MarshalJSON()
	assert.Nil(err)
	assert.NotNil(bytes)

	resB := zebra.NewResourceMap(funMap)
	assert.NotNil(resB)

	err = resB.UnmarshalJSON(bytes)
	assert.Nil(err)
}
