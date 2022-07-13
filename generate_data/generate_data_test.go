package generate_data_test

import (
	"net"
	"testing"

	"github.com/project-safari/zebra/generate_data"
	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	user := generate_data.User()
	assert.NotNil(user)
}

func TestPass(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	pwd := generate_data.Password()
	assert.NotNil(pwd)
}

func TestRange(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	num := generate_data.Range()

	assert.NotNil(num)
}

func TestCreateLabels(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	labels := generate_data.CreateLabels()

	assert.NotNil(labels)
}

func TestPorts(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	prt := generate_data.Ports()

	assert.NotNil(prt)
}

func TestModel(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	model := generate_data.Models()

	assert.NotNil(model)

	assert.NotEqual(model, " ")
}

func TestSerials(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ser := generate_data.Serials()

	assert.NotNil(ser)
	assert.NotEqual(ser, " ")

}

func TestCreateIPArr(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	IParr := generate_data.CreateIPArr(2)
	assert.NotEmpty(IParr)
}

func TestCreateVlanPool(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	VlanPool := generate_data.CreateVlanPool("VlanPool")

	assert.NotNil(VlanPool)

	if VlanPool.RangeStart > VlanPool.RangeEnd {
		assert.True(VlanPool.RangeStart > VlanPool.RangeEnd)
	} else if VlanPool.RangeStart < VlanPool.RangeEnd {
		assert.True(VlanPool.RangeStart < VlanPool.RangeEnd)
	} else if VlanPool.RangeStart == VlanPool.RangeEnd {
		assert.Equal(VlanPool.RangeStart, VlanPool.RangeEnd)
	}

}

func TestCreateSwitch(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	Switch := generate_data.CreateSwitch("Switch", net.IP("192.222.004"))

	assert.NotNil(Switch)

}

func TestCreateIPAddressPool(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	IPs := generate_data.CreateIPAddressPool("IPAddressPool", generate_data.CreateIPArr(2))

	assert.NotNil(IPs)

	assert.NotEmpty(IPs)
}

func TestCreateDatacenter(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	DataCenter := generate_data.CreateDatacenter()

	assert.NotNil(DataCenter)
}

func TestCreateLab(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	Lab := generate_data.CreateLab()

	assert.NotNil(Lab)
}

func TestCreateRack(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	Rack := generate_data.CreateRack()

	assert.NotNil(Rack)

	assert.NotNil(Rack.BaseResource)
}
