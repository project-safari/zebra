package pkg_test

import (
	"context"
	"net"
	"reflect"
	"testing"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/auth"
	"github.com/project-safari/zebra/cmd/herd/pkg"
	"github.com/project-safari/zebra/compute"
	"github.com/project-safari/zebra/dc"
	"github.com/project-safari/zebra/network"
	"github.com/stretchr/testify/assert"
)

func verifyType(assert *assert.Assertions, t string, resources []zebra.Resource) {
	for _, r := range resources {
		assert.NotNil(r)
		assert.Equal(t, r.GetType())
		assert.Nil(r.Validate(context.Background()))
	}
}

// Tests for the address file - IP generation.
func TestAddresses(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	addr := pkg.Addresses()

	assert.NotNil(addr)
	assert.NotEmpty(addr)
}

// Tests for setting the IP address.
func TestSetIP(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	samples := pkg.IPsamples()

	assert.NotNil(samples)
	assert.NotEmpty(samples)
}

// Tests for creating an array of IP addresses.
func TestCreateIPArr(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	num := 5
	arr := pkg.CreateIPArr(num)

	assert.NotNil(arr)
	assert.NotEmpty(arr)

	assert.Equal(num, len(arr))
}

// Tests for the device_info file.
func TestPorts(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	prt := pkg.Ports()

	assert.NotNil(prt)
}

// Tests for generating a model name.
func TestModel(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	model := pkg.Models()

	assert.NotNil(model)

	assert.NotEqual(model, " ")
}

// Tests for generating a serial code.
func TestSerials(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ser := pkg.Serials()

	assert.NotNil(ser)
	assert.NotEqual(ser, " ")
}

// Tests for generation of row.
func TestRows(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	row := pkg.Rows()

	assert.NotNil(row)
	assert.NotEqual(row, " ")
}

// Tests for generation of detailed information.
func TestSelectServer(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	srv := pkg.SelectServer()
	assert.NotNil(srv)
}

// Tests for generating / selecting an ID for esx servers.
func TestSelectESX(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	esx := pkg.SelectESX()
	assert.NotNil(esx)
}

// Tests for generating / selecting an ID for vcenters.
func TestSelectVCenter(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	center := pkg.SelectVcenter()
	assert.NotNil(center)
}

// Tests for the person file.
func TestUser(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	user := pkg.User()
	assert.NotNil(user)
}

// Tests for gereration of default passwords.
func TestPass(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	pwd := pkg.Password("person")
	assert.NotNil(pwd)
	assert.Equal("person123", pwd)
}

// Test for users' names.
func TestName(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	name := pkg.Name()

	assert.NotNil(name)
}

// Tests for the labels file.
func TestCreateLabels(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	labels := pkg.CreateLabels()

	assert.NotNil(labels)
}

// Tests for the rand file.
func TestRandData(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	testArr := []string{"this", "is", "a", "test", "array"}

	selected := pkg.RandData(testArr)

	assert.NotNil(selected)
	assert.NotEmpty(selected)

	test2Arr := []string{}

	assert.Empty(test2Arr)
}

// Tests for generation of random data for sample executions.
func TestRandNum(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	assert.NotNil(reflect.TypeOf(pkg.RandNum(5)))
}

// Tests for Vlan generation.
//
// #1 - creating the vlan.
func TestCreateVlanPool(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	var start uint16 = 111

	var end uint16 = 333

	VlanPool := network.NewVlanPool(start, end, pkg.CreateLabels())

	assert.NotNil(VlanPool)
	assert.NotEmpty(VlanPool)
}

// Tests for Vlan generation.
//
// #2 - generating the vlan.
func TestGenerateVlan(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	testVlans := pkg.GenerateVlanPool(3)

	assert.NotNil(testVlans)
	assert.Equal(len(testVlans), 3)
	verifyType(assert, "VLANPool", testVlans)
}

// Tests for switch generation.
//
// #1 - creating the switch.
func TestCreateSwitch(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	arr := []string{pkg.Serials(), pkg.Models(), pkg.Name()}
	sw := network.NewSwitch(arr, pkg.Ports(), net.IP("192.222.004"), pkg.CreateLabels())

	assert.NotNil(sw)
	assert.NotEmpty(sw)
}

// Tests for switch generation.
//
// #2 - generating the switch.
func TestGenerateSwitch(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	testSwitch := pkg.GenerateSwitch(2)

	assert.NotNil(testSwitch)
	assert.Equal(len(testSwitch), 2)
	verifyType(assert, "Switch", testSwitch)
}

// Tests for IPAddress generation.
//
// #1 - creating the IPAddressPool.
func TestCreateIPPool(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	pool := network.NewIPAddressPool(pkg.CreateIPArr(2), pkg.CreateLabels())

	assert.NotNil(pool)
	assert.NotEmpty(pool)
}

// Tests for IPAddress generation.
//
// #2 - generating the IPAddressPool.
func TestGenerateIPAddressPool(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	testPool := pkg.GenerateIPPool(2)

	assert.NotNil(testPool)
	verifyType(assert, "IPAddressPool", testPool)
}

// Tests for datacenter generation.
//
// #1 - creating the datacenter.
func TestCreateDatacenter(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	datacent := dc.NewDatacenter(pkg.Addresses(), pkg.Name(), pkg.CreateLabels())

	assert.NotNil(datacent)
	assert.NotEmpty(datacent)
}

// Tests for datacenter generation.
//
// #2 - generating the datacenter.
func TestGenerateDatacenterl(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	testDC := pkg.GenerateDatacenter(4)

	assert.NotNil(testDC)
	assert.Equal(len(testDC), 4)
	verifyType(assert, "Datacenter", testDC)
}

// Tests for lab generation.
//
// #1 - creating the lab.
func TestCreateLab(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	lab := dc.NewLab(pkg.Name(), pkg.CreateLabels())

	assert.NotNil(lab)
	assert.NotEmpty(lab)
}

// Tests for lab generation.
//
// #2 - generating the lab.
func TestGenerateLab(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	testLab := pkg.GenerateLab(2)

	assert.Equal(len(testLab), 2)
	verifyType(assert, "Lab", testLab)
}

// Tests for rack generation.
//
// #1 - creating the rack.
func TestCreateRack(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rack := dc.NewRack(pkg.Name(), pkg.Rows(), pkg.CreateLabels())

	assert.NotEmpty(rack)
}

// Tests for rack generation.
//
// #2 - generating the rack.
func TestGenerateRack(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	testRack := pkg.GenerateRack(1)

	assert.Equal(len(testRack), 1)
	verifyType(assert, "Rack", testRack)
}

// Tests for vcenter generation.
//
// #1 - creating the vcenter.
func TestCreateVCenter(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	vc := compute.NewVCenter(pkg.Name(), net.IP("192.222.004"), pkg.CreateLabels())

	assert.NotEmpty(vc)
	assert.Equal("VCenter", vc.GetType())
}

// Tests for vcenter generation.
//
// #2 - generating the vcenter.
func TestGenerateVC(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	testVC := pkg.GenerateVCenter(3)

	assert.Equal(len(testVC), 3)
	verifyType(assert, "VCenter", testVC)
	assert.NotNil(testVC)
}

// Tests for esx generation.
//
// #1 - creaing the esx.
func TestCreateESX(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	vc := compute.NewESX(pkg.Name(), pkg.SelectServer(), net.IP("192.222.004"), pkg.CreateLabels())

	assert.NotEmpty(vc)
}

// Tests for esx generation.
//
// #2 - generating the esx.
func TestGenerateESX(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	testESX := pkg.GenerateESX(2)

	assert.Equal(len(testESX), 2)
	assert.NotNil((testESX))
	verifyType(assert, "ESX", testESX)
}

// Tests for server generation.
//
// #1 - creating the server.
func TestCreateServer(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	arr := []string{pkg.Serials(), pkg.Models(), pkg.Name()}
	server := compute.NewServer(arr, net.IP("192.222.004"), pkg.CreateLabels())

	assert.NotEmpty(server)
}

// Tests for server generation.
//
// #2 - generating the server.
func TestGenerateServer(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	testServer := pkg.GenerateServer(10)

	assert.Equal(len(testServer), 10)
	assert.NotNil((testServer))
	verifyType(assert, "Server", testServer)
}

// Tests for vm generation.
//
// #1 - creating the vm.
func TestCreateVM(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	arr := []string{pkg.Name(), pkg.SelectESX(), pkg.SelectVcenter()}
	vm := compute.NewVM(arr, net.IP("192.222.004"), pkg.CreateLabels())

	assert.NotEmpty(vm)
}

// Tests for vm generation.
//
// #2 - generating the vm.
func TestGenerateVM(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	testVM := pkg.GenerateVM(5)

	assert.NotNil(testVM)
	verifyType(assert, "VM", testVM)
}

// Tests for user info generation.
//
// #1 - creating the user.
func TestCreateUser(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	key, err := auth.Generate()
	assert.Nil(err)

	name := pkg.Name()

	user := auth.NewUser(name, pkg.Email(name), pkg.Password(name), key, pkg.CreateLabels())
	assert.NotNil(user)
	assert.Nil(user.Validate(context.Background()))
}

// Tests for user info generation.
//
// #2 - generating the user.
func TestGenerateUser(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	usr := pkg.GenerateUser(3)

	assert.NotEmpty(usr)
	verifyType(assert, "User", usr)
}

// Tests for credential info generation.
//
// #1 - creating the credential info.
func TestCreateCredential(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	credentials := zebra.NewCredential(pkg.Name(), pkg.CreateLabels())

	assert.NotEmpty(credentials)
}

// Tests for credential info generation.
//
// #2 - generating the credential info.
func TestGenerateCredential(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	crd := pkg.GenerateCredential(2)

	assert.NotEmpty(crd)
}
