package pkg_test

import (
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

	pwd := pkg.Password()
	assert.NotNil(pwd)
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
}

// tests for vcenter generation.
func TestCreateVCenter(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	vc := compute.NewVCenter(pkg.Name(), net.IP("192.222.004"), pkg.CreateLabels())

	assert.NotEmpty(vc)
}

func TestGenerateVC(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	testVC := pkg.GenerateVcenter(3)

	assert.Equal(len(testVC), 3)
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
}

// tests for user info generation.
func TestCreateUser(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	user := auth.NewUser(pkg.User(), pkg.Name(), pkg.Password(), pkg.CreateLabels())

	assert.NotEmpty(user)
}

func TestGenerateUser(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	usr := pkg.GenerateUser(3)

	assert.NotEmpty(usr)
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
