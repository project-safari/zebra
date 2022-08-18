package zebra_test

import (
	"testing"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/network"
	"github.com/project-safari/zebra/status"
	"github.com/stretchr/testify/assert"
)

func TestAddNew(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	f := zebra.Factory()
	assert.NotNil(f)

	f.Add(network.SwitchType())
	assert.NotNil(f.New("Switch"))
	assert.Nil(f.New("random"))
}

func TestNewResourceList(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	assert.NotNil(zebra.NewResourceList(nil))
}

func TestCopyResourceList(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resA := zebra.NewResourceList(nil)
	assert.NotNil(resA)

	resA.Resources = append(resA.Resources, new(network.IPAddressPool))

	resB := zebra.NewResourceList(nil)
	assert.NotNil(resB)
	assert.Empty(len(resB.Resources))

	zebra.CopyResourceList(resB, resA)
	assert.Equal(1, len(resB.Resources))

	zebra.CopyResourceList(nil, nil)
}

func TestListMarshalUnmarshal(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	funMap := zebra.Factory()
	funMap.Add(network.VLANPoolType())

	resA := zebra.NewResourceList(funMap)
	assert.NotNil(resA)

	vlan := &network.VLANPool{
		BaseResource: zebra.BaseResource{
			ID:     "0100001",
			Type:   "invalid",
			Labels: nil,
			Status: status.DefaultStatus(),
		},
		RangeStart: 0,
		RangeEnd:   10,
	}

	resA.Resources = append(resA.Resources, vlan)

	bytes, err := resA.MarshalJSON()
	assert.Nil(err)
	assert.NotNil(bytes)

	resB := zebra.NewResourceList(funMap)
	assert.NotNil(resB)

	err = resB.UnmarshalJSON(bytes)
	assert.NotNil(err)

	vlan.Type = "VLANPool"
	resA.Resources = []zebra.Resource{vlan}

	bytes, err = resA.MarshalJSON()
	assert.Nil(err)
	assert.NotNil(bytes)

	resB = zebra.NewResourceList(funMap)
	assert.NotNil(resB)

	err = resB.UnmarshalJSON(bytes)
	assert.Nil(err)
	assert.Equal(1, len(resB.Resources))
}

func TestErrorMarshalUnmarshal(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	funMap := zebra.Factory()
	funMap.Add(network.VLANPoolType())
	resList := zebra.NewResourceList(funMap)
	assert.NotNil(resList.UnmarshalJSON(nil))
	assert.NotNil(resList.UnmarshalJSON([]byte(`[{"id":"0100000001"}]`)))
	assert.NotNil(resList.UnmarshalJSON([]byte(`[{"id":"0100000001", "type":123}]`)))

	resMap := zebra.NewResourceMap(nil)
	assert.NotNil(resMap.UnmarshalJSON(nil))
	assert.NotNil(resMap.UnmarshalJSON([]byte(`{"VLANPool":[{"id":"0100000001", "type":123}]}`)))
}

func TestNewResourceMap(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	assert.NotNil(zebra.NewResourceMap(nil))
}

func TestCopyResourceMap(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resA := zebra.NewResourceMap(nil)
	assert.NotNil(resA)

	resA.Add(new(network.IPAddressPool), "IPAddressPool")

	resB := zebra.NewResourceMap(nil)
	assert.NotNil(resB)

	zebra.CopyResourceMap(resB, nil)

	zebra.CopyResourceMap(resB, resA)
	assert.Equal(1, len(resB.Resources))
	assert.Equal(1, len(resB.Resources["IPAddressPool"].Resources))
}

func TestGetFactory(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	f := zebra.Factory()
	f.Add(network.SwitchType())
	assert.NotEmpty(f.Types())
	aType, ok := f.Type("Switch")
	assert.True(ok)
	assert.NotNil(aType)

	resA := zebra.NewResourceMap(f)
	assert.NotNil(resA)

	assert.NotNil(resA.GetFactory())
}

func TestAdd(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	funMap := zebra.Factory()
	funMap.Add(network.SwitchType())

	resA := zebra.NewResourceMap(funMap)
	assert.NotNil(resA)

	switch1 := funMap.New("Switch")

	resA.Add(switch1, "Switch")
	assert.NotNil(len(resA.Resources["Switch"].Resources) == 1)
}

func TestDelete(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	funMap := zebra.Factory()
	funMap.Add(network.SwitchType())

	resA := zebra.NewResourceMap(funMap)
	assert.NotNil(resA)

	switch1 := funMap.New("Switch")

	resA.Add(switch1, "Switch")
	assert.NotNil(len(resA.Resources["Switch"].Resources) == 1)

	resA.Delete(switch1, "Switch")

	resA.Delete(switch1, "invalid_key")

	_, ok := resA.Resources["Switch"]
	assert.NotNil(ok)
}

func TestMapMarshalUnMarshal(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	funMap := zebra.Factory()
	funMap.Add(network.VLANPoolType())

	resA := zebra.NewResourceMap(funMap)
	assert.NotNil(resA)

	vlan := &network.VLANPool{
		BaseResource: zebra.BaseResource{
			ID:     "0100001",
			Type:   "VLANPool",
			Labels: nil,
			Status: status.DefaultStatus(),
		},
		RangeStart: 0,
		RangeEnd:   10,
	}

	resA.Add(vlan, "VLANPool")

	bytes, err := resA.MarshalJSON()
	assert.Nil(err)
	assert.NotNil(bytes)

	resB := zebra.NewResourceMap(funMap)
	assert.NotNil(resB)

	err = resB.UnmarshalJSON(bytes)
	assert.Nil(err)
}
