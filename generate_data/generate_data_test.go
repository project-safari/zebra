/*
create 100 instances of each resource for some users
program to execute tests
*/

package generate_data_test

import (
	"net"
	"testing"

	"github.com/project-safari/zebra/generate_data"
	"github.com/stretchr/testify/assert"
)

func TestSetTypes(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	types := generate_data.AllResourceTypes()
	assert.NotNil(types)
}

func TestSetIP(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	samples := generate_data.IPsamples()
	assert.NotNil(samples)
}

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

func TestName(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	name := generate_data.Name()

	assert.NotNil(name)
}

func TestRange(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	num := generate_data.Range()

	assert.NotNil(num)
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

func TestRows(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	row := generate_data.Rows()

	assert.NotNil(row)
	assert.NotEqual(row, " ")
}

func TestAddresses(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	adr := generate_data.Addresses()

	assert.NotNil(adr)
	assert.NotEqual(adr, " ")
}

func TestOrder(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	var a uint16 = 20
	var b uint16 = 5

	one, two := generate_data.Order(a, b)
	assert.True(one < two)
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

	VlanPool := generate_data.NewVlanPool("VlanPool")

	assert.NotNil(VlanPool)
	assert.NotEmpty(VlanPool)
}

func TestCreateVcenter(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	Vcenter := generate_data.NewVCenter("VlanPool", net.IP("192.222.004"))

	assert.NotNil(Vcenter)
}

func TestCreateSwitch(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	Switch := generate_data.NewSwitch("Switch", net.IP("192.222.004"))

	assert.NotNil(Switch)
}

func TestCreateIPAddressPool(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	IPs := generate_data.NewIPAddressPool("IPAddressPool", generate_data.CreateIPArr(2))

	assert.NotNil(IPs)

	assert.NotEmpty(IPs)
}

func TestCreateDatacenter(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	DataCenter := generate_data.NewDatacenter("Datacenter")

	assert.NotNil(DataCenter)
}

func TestCreateLabels(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	labels := generate_data.CreateLabels()

	assert.NotNil(labels)
}

func TestCreateLab(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	Lab := generate_data.NewLab("Lab")

	assert.NotNil(Lab)
}

func TestCreateRack(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	Rack := generate_data.NewRack("Rack")

	assert.NotNil(Rack)

	assert.NotNil(Rack.BaseResource)
}

func TestIsGood(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	result := generate_data.IsGood(100)

	assert.NotNil(result)
	assert.False(result)

	errRes := generate_data.IsGood(0)
	assert.True(errRes)
}

func TestGeneration(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	num := 100
	result := generate_data.IsGood(num)

	assert.NotNil(generate_data.GenerateData(result, num))
	user, arr := generate_data.GenerateData(result, num)

	assert.NotNil(user)

	assert.NotEmpty(arr)
}
