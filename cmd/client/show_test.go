package main //nolint:testpackage

import (
	"testing"
	"time"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/auth"
	"github.com/project-safari/zebra/model"
	"github.com/project-safari/zebra/model/compute"
	"github.com/project-safari/zebra/model/dc"
	"github.com/project-safari/zebra/model/lease"
	"github.com/project-safari/zebra/model/network"
	"github.com/project-safari/zebra/model/user"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// to make it easier.
func test() *cobra.Command {
	showCmd := NewShow()

	return showCmd
}

// Function that tests the client.
func TestClient(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	c, err := NewClient(nil)
	assert.Nil(c)
	assert.Equal(ErrNoConfig, err)

	cfg := new(Config)
	c, err = NewClient(cfg)
	assert.Nil(c)
	assert.Equal(ErrNoEmail, err)

	cfg.Email = "test@zebra.project-safafi.io"
	c, err = NewClient(cfg)
	assert.Nil(c)
	assert.Equal(ErrNoPrivateKey, err)

	key, err := auth.Load(testUserKeyFile)
	assert.Nil(err)
	assert.NotNil(key)

	cfg.Key = key
	c, err = NewClient(cfg)
	assert.Equal(ErrNoCACert, err)
	assert.Nil(c)

	key.Public()
	cfg.Key = key.Public()
	c, err = NewClient(cfg)
	assert.Nil(c)
	assert.Equal(auth.ErrNoPrivateKey, err)

	cfg.CACert = testCACertFile
	cfg.Key = key
	c, err = NewClient(cfg)
	assert.Nil(err)
	assert.NotNil(c)

	cli, err := NewClient(cfg)

	assert.Nil(err)

	assert.NotNil(cli)
}

// tests for adding new show command(s).
func TestNewZebraCommand(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	cmd := NewShow()
	assert.NotNil(cmd)
}

// tests for resources.
func TestShowRes(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	argLock.Lock()
	defer argLock.Unlock()

	args := []string{"rack", "user", "test-case"}

	rootCmd := New()

	rootCmd.AddCommand(test())

	res := showResources(rootCmd, args)

	assert.NotNil(res)
}

// showing lease status and information.
func TestShowLease(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	argLock.Lock()
	defer argLock.Unlock()

	args := []string{"lease-status", "test-case"}

	rootCmd := New()

	rootCmd.AddCommand(test())

	res := showLeases(rootCmd, args)

	assert.NotNil(res)
}

// tests for server resource types (server, esx, vcenter, vm).
func TestShowServer(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	argLock.Lock()
	defer argLock.Unlock()

	args := []string{"server", "test-case"}

	rootCmd := New()

	rootCmd.AddCommand(test())

	res := showServers(rootCmd, args)

	assert.NotNil(res)
}

func TestShowESX(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	argLock.Lock()
	defer argLock.Unlock()

	args := []string{"esx", "test-case"}

	rootCmd := New()

	rootCmd.AddCommand(test())

	res := showESX(rootCmd, args)

	assert.NotNil(res)
}

func TestShowVC(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	argLock.Lock()
	defer argLock.Unlock()

	args := []string{"vcenter", "test-case"}

	rootCmd := New()

	rootCmd.AddCommand(test())

	res := showVCenters(rootCmd, args)

	assert.NotNil(res)
}

func TestShowVM(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	argLock.Lock()
	defer argLock.Unlock()

	args := []string{"vm", "test-case"}

	rootCmd := New()

	rootCmd.AddCommand(test())

	res := showVM(rootCmd, args)

	assert.NotNil(res)
}

// tests for dc resource types (datacenter, lab, rack).
func TestShowDatacenter(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	argLock.Lock()
	defer argLock.Unlock()

	args := []string{"datacenter", "test-case"}

	rootCmd := New()

	rootCmd.AddCommand(test())

	res := showDatacenters(rootCmd, args)

	assert.NotNil(res)
}

func TestShowLab(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	argLock.Lock()
	defer argLock.Unlock()

	args := []string{"lab", "test-case"}

	rootCmd := New()

	rootCmd.AddCommand(test())

	res := showLabs(rootCmd, args)

	assert.NotNil(res)
}

func TestShowRack(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	argLock.Lock()
	defer argLock.Unlock()

	args := []string{"rack", "test-case"}

	rootCmd := New()

	rootCmd.AddCommand(test())

	res := showRacks(rootCmd, args)

	assert.NotNil(res)
}

// tests for network resources (switch, vlan, ip-address).
func TestShowSwitches(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	argLock.Lock()
	defer argLock.Unlock()

	args := []string{"switch", "test-case"}

	rootCmd := New()

	rootCmd.AddCommand(test())

	res := showSwitches(rootCmd, args)

	assert.NotNil(res)
}

func TestShowVlans(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	argLock.Lock()
	defer argLock.Unlock()

	args := []string{"vlan", "test-case"}

	rootCmd := New()

	rootCmd.AddCommand(test())

	res := showVlans(rootCmd, args)

	assert.NotNil(res)
}

func TestShowPools(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	argLock.Lock()
	defer argLock.Unlock()

	args := []string{"ip", "test-case"}

	rootCmd := New()

	rootCmd.AddCommand(test())

	res := showIPs(rootCmd, args)

	assert.NotNil(res)
}

// tests for user resources (user data, registrations, key).
func TestShowUsers(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	argLock.Lock()
	defer argLock.Unlock()

	args := []string{"user", "test-case"}

	rootCmd := New()

	rootCmd.AddCommand(test())

	res := showUsers(rootCmd, args)

	assert.NotNil(res)
}

func TestShowRegistrations(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	argLock.Lock()
	defer argLock.Unlock()

	args := []string{"registration", "test-case"}

	rootCmd := New()

	rootCmd.AddCommand(test())

	res := showRegistrations(rootCmd, args)

	assert.NotNil(res)
}

func TestPrintResources(t *testing.T) { //nolint:funlen
	t.Parallel()
	assert := assert.New(t)

	rootCmd := New()
	assert.NotNil(rootCmd)

	rootCmd.AddCommand(test())

	// test with only one resource.

	fact := model.Factory()

	resMap := zebra.NewResourceMap(fact)

	assert.NotNil(resMap)

	rack := dc.NewRack("test_row", "test_rack", "test_owner", "test_group")
	assert.Nil(resMap.Add(rack))

	printResources(resMap)

	// test with many resources.

	bigMap := zebra.NewResourceMap(fact)

	assert.NotNil(bigMap)

	addr := network.NewIPAddressPool("test_ip_pool", "test_owner", "test_group")
	assert.Nil(bigMap.Add(addr))

	vlan := network.NewVLANPool("test_vlan_pool", "test_owner", "test_group")
	assert.Nil(bigMap.Add(vlan))

	sw := network.NewSwitch("test_switch", "test_owner", "test_group")
	assert.Nil(bigMap.Add(sw))

	printResources(bigMap)

	// test with all resources.

	allMap := zebra.NewResourceMap(fact)
	assert.NotNil(allMap)

	assert.Nil(allMap.Add(addr))
	assert.Nil(allMap.Add(vlan))
	assert.Nil(allMap.Add(sw))

	center := dc.NewDatacenter("test_dc_addr", "test_dc", "test_owner", "test_group")
	assert.Nil(allMap.Add(center))
	assert.Nil(allMap.Add(rack))

	lab := dc.NewLab("test_lab", "test_owner", "test_group")
	assert.Nil(allMap.Add(lab))

	vc := compute.NewVCenter("test_vcenter", "test_owner", "test_group")
	assert.Nil(allMap.Add(vc))

	vm := compute.NewVM("test_esx", "test_vm", "test_owner", "test_group")
	assert.Nil(allMap.Add(vm))

	srv := compute.NewServer("test_serial", "test_model", "test_server", "test_owner", "test_group")
	assert.Nil(allMap.Add(srv))

	eserver := compute.NewESX("test_server", "test_esx", "test_owner", "test_group")
	assert.Nil(allMap.Add(eserver))

	key, _ := auth.Generate()
	p, _ := auth.NewPriv("", false, true, false, false)
	role := &auth.Role{
		Name:       "user",
		Privileges: []*auth.Priv{p},
	}
	usr := user.NewUser("test_user", "test@zebra.io", "bigPassword1!!!", key.Public(), role)
	assert.Nil(allMap.Add(usr))

	printResources(allMap)
}

func TestPrintServers(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rootCmd := New()
	assert.NotNil(rootCmd)

	rootCmd.AddCommand(test())

	fact := model.Factory()
	resMap := zebra.NewResourceMap(fact)
	server := compute.NewServer("test_serial", "test_model", "test_server", "test_owner", "test_group")
	assert.Nil(resMap.Add(server))

	listed := resMap.Resources["compute.server"].Resources
	assert.NotNil(listed)

	printServers(listed)
}

func TestPrintESX(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rootCmd := New()
	assert.NotNil(rootCmd)

	rootCmd.AddCommand(test())

	fact := model.Factory()
	resMap := zebra.NewResourceMap(fact)
	esx := compute.NewESX("test_server", "test_esx", "test_owner", "test_group")
	assert.Nil(resMap.Add(esx))

	listed := resMap.Resources["compute.esx"].Resources
	assert.NotNil(listed)

	printESX(listed)
}

func TestPrintVCenter(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rootCmd := New()
	assert.NotNil(rootCmd)

	rootCmd.AddCommand(test())

	fact := model.Factory()
	resMap := zebra.NewResourceMap(fact)
	vc := compute.NewVCenter("test_vcenter", "test_owner", "test_group")
	assert.Nil(resMap.Add(vc))

	listed := resMap.Resources["compute.vcenter"].Resources

	assert.NotNil(listed)

	printVCenters(listed)
}

func TestPrintVM(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rootCmd := New()
	assert.NotNil(rootCmd)

	rootCmd.AddCommand(test())

	fact := model.Factory()
	resMap := zebra.NewResourceMap(fact)
	vm := compute.NewVM("test_esx", "test_vm", "test_owner", "test_group")
	assert.Nil(resMap.Add(vm))

	listed := resMap.Resources["compute.vm"].Resources

	assert.NotNil(listed)

	printVM(listed)
}

func TestPrintVlan(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rootCmd := New()
	assert.NotNil(rootCmd)

	rootCmd.AddCommand(test())

	fact := model.Factory()
	resMap := zebra.NewResourceMap(fact)
	vlan := network.NewVLANPool("test_vlan_pool", "test_owner", "test_group")
	assert.Nil(resMap.Add(vlan))

	listed := resMap.Resources["network.vlanPool"].Resources

	assert.NotNil(listed)

	printVlans(listed)
}

func TestPrintSwitches(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rootCmd := New()
	assert.NotNil(rootCmd)

	rootCmd.AddCommand(test())

	fact := model.Factory()
	resMap := zebra.NewResourceMap(fact)
	sw := network.NewSwitch("test_switch", "test_owner", "test_group")
	assert.Nil(resMap.Add(sw))

	listed := resMap.Resources["network.switch"].Resources

	assert.NotNil(listed)

	printSwitches(listed)
}

func TestPrintIPPools(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rootCmd := New()
	assert.NotNil(rootCmd)

	rootCmd.AddCommand(test())

	fact := model.Factory()
	resMap := zebra.NewResourceMap(fact)
	ip := network.NewIPAddressPool("test_ip_pool", "test_owner", "test_group")
	assert.Nil(resMap.Add(ip))

	listed := resMap.Resources["network.ipAddressPool"].Resources
	assert.NotNil(listed)

	printIPs(listed)
}

func TestPrintDC(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rootCmd := New()
	assert.NotNil(rootCmd)

	rootCmd.AddCommand(test())

	fact := model.Factory()
	resMap := zebra.NewResourceMap(fact)
	dc := dc.NewDatacenter("test_dc_addr", "test_dc", "test_owner", "test_group")
	assert.Nil(resMap.Add(dc))

	listed := resMap.Resources["dc.datacenter"].Resources
	assert.NotNil(listed)

	printDatacenters(listed)
}

func TestPrintlabs(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rootCmd := New()
	assert.NotNil(rootCmd)

	rootCmd.AddCommand(test())

	fact := model.Factory()
	resMap := zebra.NewResourceMap(fact)
	lab := dc.NewLab("test_lab", "test_owner", "test_group")
	assert.Nil(resMap.Add(lab))

	listed := resMap.Resources["dc.lab"].Resources
	assert.NotNil(listed)

	printLabs(listed)
}

func TestPrintRacks(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rootCmd := New()
	assert.NotNil(rootCmd)

	rootCmd.AddCommand(test())

	fact := model.Factory()
	resMap := zebra.NewResourceMap(fact)
	rack := dc.NewRack("test_row", "test_rack", "test_owner", "test_group")
	assert.Nil(resMap.Add(rack))

	listed := resMap.Resources["dc.rack"].Resources
	assert.NotNil(listed)

	printRacks(listed)
}

func TestPrintLeases(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rootCmd := New()
	assert.NotNil(rootCmd)

	rootCmd.AddCommand(test())

	fact := model.Factory()
	resMap := zebra.NewResourceMap(fact)
	lease := lease.NewLease("test@zebra.io", time.Hour, nil)
	assert.Nil(resMap.Add(lease))

	listed := resMap.Resources["system.lease"].Resources
	assert.NotNil(listed)

	printLeases(listed)
}

func TestPrintUsers(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rootCmd := New()
	assert.NotNil(rootCmd)

	rootCmd.AddCommand(test())

	fact := model.Factory()
	resMap := zebra.NewResourceMap(fact)
	key, _ := auth.Generate()
	p, _ := auth.NewPriv("", false, true, false, false)
	role := &auth.Role{
		Name:       "user",
		Privileges: []*auth.Priv{p},
	}
	usr := user.NewUser("test_user", "test@zebra.io", "bigPassword1!!!", key.Public(), role)
	assert.Nil(resMap.Add(usr))

	listed := resMap.Resources["system.user"].Resources
	assert.NotNil(listed)

	printUsers(listed)
}

func TestBehavior(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rootCmd := New()
	assert.NotNil(rootCmd)

	assert.NotNil(test())

	main()
}
