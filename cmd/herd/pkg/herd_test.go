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

// tests for the address file - IP generation.
func TestAddresses(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	addr := pkg.Addresses()

	assert.NotNil(addr)
	assert.NotEmpty(addr)
}

func TestSetIP(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	samples := pkg.IPsamples()

	assert.NotNil(samples)
	assert.NotEmpty(samples)
}

func TestCreateIPArr(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	num := 5
	arr := pkg.CreateIPArr(num)

	assert.NotNil(arr)
	assert.NotEmpty(arr)

	assert.Equal(num, len(arr))
}

// tests for the device_info file.
func TestPorts(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	prt := pkg.Ports()

	assert.NotNil(prt)
}

func TestModel(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	model := pkg.Models()

	assert.NotNil(model)

	assert.NotEqual(model, " ")
}

func TestSerials(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	ser := pkg.Serials()

	assert.NotNil(ser)
	assert.NotEqual(ser, " ")
}

func TestRows(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	row := pkg.Rows()

	assert.NotNil(row)
	assert.NotEqual(row, " ")
}

// tests for generation of detailed information.
func TestSelectServer(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	srv := pkg.SelectServer()
	assert.NotNil(srv)
}

func TestSelectESX(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	esx := pkg.SelectESX()
	assert.NotNil(esx)
}

func TestSelectVCenter(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	center := pkg.SelectVcenter()
	assert.NotNil(center)
}

// tests for the person file.
func TestUser(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	user := pkg.User()
	assert.NotNil(user)
}

func TestPass(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	pwd := pkg.Password("person")
	assert.NotNil(pwd)
	assert.Equal("person123", pwd)
}

func TestName(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	name := pkg.Name()

	assert.NotNil(name)
}

// tests for the labels file.
func TestCreateLabels(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	labels := pkg.CreateLabels()

	assert.NotNil(labels)
}

// tests for the rand file.
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

func TestRandNum(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	assert.NotNil(reflect.TypeOf(pkg.RandNum(5)))
}

// tests for Vlan generation.
func TestCreateVlanPool(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	var start uint16 = 111

	var end uint16 = 333

	VlanPool := network.NewVlanPool(start, end, pkg.CreateLabels())

	assert.NotNil(VlanPool)
	assert.NotEmpty(VlanPool)
}

func TestGenerateVlan(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	testVlans := pkg.GenerateVlanPool(3)

	assert.NotNil(testVlans)
	assert.Equal(len(testVlans), 3)
	verifyType(assert, "VLANPool", testVlans)
}

// tests for switch generation.
func TestCreateSwitch(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	arr := []string{pkg.Serials(), pkg.Models(), pkg.Name()}
	sw := network.NewSwitch(arr, pkg.Ports(), net.IP("192.222.004"), pkg.CreateLabels())

	assert.NotNil(sw)
	assert.NotEmpty(sw)
}

func TestGenerateSwitch(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	testSwitch := pkg.GenerateSwitch(2)

	assert.NotNil(testSwitch)
	assert.Equal(len(testSwitch), 2)
	verifyType(assert, "Switch", testSwitch)
}

// tests for IPAddress generation.
func TestCreateIPPool(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	pool := network.NewIPAddressPool(pkg.CreateIPArr(2), pkg.CreateLabels())

	assert.NotNil(pool)
	assert.NotEmpty(pool)
}

func TestGenerateIPAddressPool(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	testPool := pkg.GenerateIPPool(2)

	assert.NotNil(testPool)
	verifyType(assert, "IPAddressPool", testPool)
}

// tests for datacenter generation.
func TestCreateDatacenter(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	datacent := dc.NewDatacenter(pkg.Addresses(), pkg.Name(), pkg.CreateLabels())

	assert.NotNil(datacent)
	assert.NotEmpty(datacent)
}

func TestGenerateDatacenterl(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	testDC := pkg.GenerateDatacenter(4)

	assert.NotNil(testDC)
	assert.Equal(len(testDC), 4)
	verifyType(assert, "Datacenter", testDC)
}

// tests for lab generation.
func TestCreateLab(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	lab := dc.NewLab(pkg.Name(), pkg.CreateLabels())

	assert.NotNil(lab)
	assert.NotEmpty(lab)
}

func TestGenerateLab(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	testLab := pkg.GenerateLab(2)

	assert.Equal(len(testLab), 2)
	verifyType(assert, "Lab", testLab)
}

// tests for rack generation.
func TestCreateRack(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rack := dc.NewRack(pkg.Name(), pkg.Rows(), pkg.CreateLabels())

	assert.NotEmpty(rack)
}

func TestGenerateRack(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	testRack := pkg.GenerateRack(1)

	assert.Equal(len(testRack), 1)
	verifyType(assert, "Rack", testRack)
}

// tests for vcenter generation.
func TestCreateVCenter(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	vc := compute.NewVCenter(pkg.Name(), net.IP("192.222.004"), pkg.CreateLabels())

	assert.NotEmpty(vc)
	assert.Equal("VCenter", vc.GetType())
}

func TestGenerateVC(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	testVC := pkg.GenerateVCenter(3)

	assert.Equal(len(testVC), 3)
	verifyType(assert, "VCenter", testVC)
	assert.NotNil(testVC)
}

// tests for esx generation.
func TestCreateESX(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	vc := compute.NewESX(pkg.Name(), pkg.SelectServer(), net.IP("192.222.004"), pkg.CreateLabels())

	assert.NotEmpty(vc)
}

func TestGenerateESX(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	testESX := pkg.GenerateESX(2)

	assert.Equal(len(testESX), 2)
	assert.NotNil((testESX))
	verifyType(assert, "ESX", testESX)
}

// tests for server generation.
func TestCreateServer(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	arr := []string{pkg.Serials(), pkg.Models(), pkg.Name()}
	server := compute.NewServer(arr, net.IP("192.222.004"), pkg.CreateLabels())

	assert.NotEmpty(server)
}

func TestGenerateServer(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	testServer := pkg.GenerateServer(10)

	assert.Equal(len(testServer), 10)
	assert.NotNil((testServer))
	verifyType(assert, "Server", testServer)
}

// tests for vm generation.
func TestCreateVM(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	arr := []string{pkg.Name(), pkg.SelectESX(), pkg.SelectVcenter()}
	vm := compute.NewVM(arr, net.IP("192.222.004"), pkg.CreateLabels())

	assert.NotEmpty(vm)
}

func TestGenerateVM(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	testVM := pkg.GenerateVM(5)

	assert.NotNil(testVM)
	verifyType(assert, "VM", testVM)
}

// tests for user info generation.
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

func TestGenerateUser(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	usr := pkg.GenerateUser(3)

	assert.NotEmpty(usr)
	verifyType(assert, "User", usr)
}

// tests for credential info generation.
func TestCreateCredential(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	credentials := zebra.NewCredential(pkg.Name(), pkg.CreateLabels())

	assert.NotEmpty(credentials)
}

func TestGenerateCredential(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	crd := pkg.GenerateCredential(2)

	assert.NotEmpty(crd)
}
