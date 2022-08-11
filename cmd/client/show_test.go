package main //nolint:testpackage

import (
	"testing"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/auth"
	"github.com/project-safari/zebra/cmd/herd/pkg"
	"github.com/project-safari/zebra/compute"
	"github.com/project-safari/zebra/dc"
	"github.com/project-safari/zebra/network"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

// to make it easier.
func test() *cobra.Command {
	showCmd := NewShow()

	return showCmd
}

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

// tests for server resource types.
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

// tests for dc resource types.
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

// tests for network resources.
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

	fact := new(zebra.ResourceFactory)

	resMap := zebra.NewResourceMap(*fact)

	assert.NotNil(resMap)

	rack := new(dc.Rack)
	rack.Status = new(zebra.Status)

	rack.Status.UsedBy = pkg.Name()

	resMap.Add(rack, "Rack")

	printResources(resMap)

	// test with many resources.

	bigMap := zebra.NewResourceMap(*fact)

	assert.NotNil(bigMap)

	addr := new(network.IPAddressPool)
	addr.Status = new(zebra.Status)

	addr.Status.UsedBy = pkg.Name()

	bigMap.Add(addr, "IPAddressPool")

	vlan := new(network.VLANPool)

	vlan.Status = new(zebra.Status)

	vlan.Status.UsedBy = pkg.Name()

	bigMap.Add(vlan, "VLANPool")

	sw := new(network.Switch)

	sw.Status = new(zebra.Status)

	sw.Status.UsedBy = pkg.Name()

	bigMap.Add(sw, "Switch")

	printResources(bigMap)

	// test with all resources.

	allMap := zebra.NewResourceMap(*fact)

	assert.NotNil(allMap)

	addr2 := new(network.IPAddressPool)
	addr2.Status = new(zebra.Status)

	addr2.Status.UsedBy = pkg.Name()

	allMap.Add(addr2, "IPAddressPool")

	vlan2 := new(network.VLANPool)
	vlan2.Status = new(zebra.Status)

	vlan2.Status.UsedBy = pkg.Name()

	allMap.Add(vlan2, "VLANPool")

	sw2 := new(network.Switch)
	sw2.Status = new(zebra.Status)

	sw2.Status.UsedBy = pkg.Name()

	allMap.Add(sw2, "Switch")

	center := new(dc.Datacenter)
	center.Status = new(zebra.Status)

	center.Status.UsedBy = pkg.Name()

	allMap.Add(center, "Datacenter")

	rack2 := new(dc.Rack)
	rack2.Status = new(zebra.Status)

	rack2.Status.UsedBy = pkg.Name()

	allMap.Add(rack, "Rack")

	lab := new(dc.Lab)
	lab.Status = new(zebra.Status)

	lab.Status.UsedBy = pkg.Name()

	allMap.Add(lab, "Lab")

	vc := new(compute.VCenter)
	vc.Status = new(zebra.Status)

	vc.Status.UsedBy = pkg.Name()

	allMap.Add(vc, "VCenter")

	vm := new(compute.VM)
	vm.Status = new(zebra.Status)

	vm.Status.UsedBy = pkg.Name()

	allMap.Add(vm, "VM")

	srv := new(compute.Server)
	srv.Status = new(zebra.Status)

	srv.Status.UsedBy = pkg.Name()

	allMap.Add(srv, "S")

	eserver := new(compute.ESX)
	eserver.Status = new(zebra.Status)

	eserver.Status.UsedBy = pkg.Name()

	allMap.Add(eserver, "esx")

	usr := new(auth.User)
	usr.Role = new(auth.Role)

	usr.Status = new(zebra.Status)

	usr.Status.UsedBy = pkg.Name()

	allMap.Add(usr, "person")

	printResources(allMap)
}

func TestPrintServers(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rootCmd := New()
	assert.NotNil(rootCmd)

	rootCmd.AddCommand(test())

	fact := new(zebra.ResourceFactory)

	resMap := zebra.NewResourceMap(*fact)

	server := new(compute.Server)

	server.Status = new(zebra.Status)

	server.Status.UsedBy = pkg.Name()

	resMap.Add(server, "Server")

	listed := resMap.Resources["Server"].Resources

	assert.NotNil(listed)

	printServers(listed)
}

func TestPrintESX(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rootCmd := New()
	assert.NotNil(rootCmd)

	rootCmd.AddCommand(test())

	fact := new(zebra.ResourceFactory)

	resMap := zebra.NewResourceMap(*fact)

	eserver := new(compute.ESX)
	eserver.Status = new(zebra.Status)

	eserver.Status.UsedBy = pkg.Name()

	resMap.Add(eserver, "ESX")

	listed := resMap.Resources["ESX"].Resources

	assert.NotNil(listed)

	printESX(listed)
}

func TestPrintVCenter(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rootCmd := New()
	assert.NotNil(rootCmd)

	rootCmd.AddCommand(test())

	fact := new(zebra.ResourceFactory)

	resMap := zebra.NewResourceMap(*fact)

	v := new(compute.VCenter)
	v.Status = new(zebra.Status)

	v.Status.UsedBy = pkg.Name()

	resMap.Add(v, "VCenter")

	listed := resMap.Resources["VCenter"].Resources

	assert.NotNil(listed)

	printVCenters(listed)
}

func TestPrintVM(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rootCmd := New()
	assert.NotNil(rootCmd)

	rootCmd.AddCommand(test())

	fact := new(zebra.ResourceFactory)

	resMap := zebra.NewResourceMap(*fact)

	machine := new(compute.VM)
	machine.Status = new(zebra.Status)

	machine.Status.UsedBy = pkg.Name()

	resMap.Add(machine, "VM")

	listed := resMap.Resources["VM"].Resources

	assert.NotNil(listed)

	printVM(listed)
}

func TestPrintVlan(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rootCmd := New()
	assert.NotNil(rootCmd)

	rootCmd.AddCommand(test())

	fact := new(zebra.ResourceFactory)

	resMap := zebra.NewResourceMap(*fact)

	vlan := new(network.VLANPool)

	vlan.Status = new(zebra.Status)

	vlan.Status.UsedBy = pkg.Name()

	resMap.Add(vlan, "VLANPool")

	listed := resMap.Resources["VLANPool"].Resources

	assert.NotNil(listed)

	printVlans(listed)
}

func TestPrintSwitches(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rootCmd := New()
	assert.NotNil(rootCmd)

	rootCmd.AddCommand(test())

	fact := new(zebra.ResourceFactory)

	resMap := zebra.NewResourceMap(*fact)

	sw := new(network.Switch)
	sw.Status = new(zebra.Status)

	sw.Status.UsedBy = pkg.Name()

	resMap.Add(sw, "Switch")

	listed := resMap.Resources["Switch"].Resources

	assert.NotNil(listed)

	printSwitches(listed)
}

func TestPrintIPPools(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rootCmd := New()
	assert.NotNil(rootCmd)

	rootCmd.AddCommand(test())

	fact := new(zebra.ResourceFactory)

	resMap := zebra.NewResourceMap(*fact)

	pool := new(network.IPAddressPool)

	pool.Status = new(zebra.Status)

	pool.Status.UsedBy = pkg.Name()

	resMap.Add(pool, "ips")

	listed := resMap.Resources["ips"].Resources

	assert.NotNil(listed)

	printIPs(listed)
}

func TestPrintDC(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rootCmd := New()
	assert.NotNil(rootCmd)

	rootCmd.AddCommand(test())

	fact := new(zebra.ResourceFactory)

	resMap := zebra.NewResourceMap(*fact)

	dc := new(dc.Datacenter)
	dc.Status = new(zebra.Status)

	dc.Status.UsedBy = pkg.Name()

	resMap.Add(dc, "Datacenter")

	listed := resMap.Resources["Datacenter"].Resources

	assert.NotNil(listed)

	printDatacenters(listed)
}

func TestPrintlabs(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rootCmd := New()
	assert.NotNil(rootCmd)

	rootCmd.AddCommand(test())

	fact := new(zebra.ResourceFactory)

	resMap := zebra.NewResourceMap(*fact)

	lb := new(dc.Lab)
	lb.Status = new(zebra.Status)

	lb.Status.UsedBy = pkg.Name()

	resMap.Add(lb, "Lab")

	listed := resMap.Resources["Lab"].Resources

	assert.NotNil(listed)

	printLabs(listed)
}

func TestPrintRacks(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rootCmd := New()
	assert.NotNil(rootCmd)

	rootCmd.AddCommand(test())

	fact := new(zebra.ResourceFactory)

	resMap := zebra.NewResourceMap(*fact)

	r := new(dc.Rack)
	r.Status = new(zebra.Status)

	r.Status.UsedBy = pkg.Name()

	resMap.Add(r, "Rack")

	listed := resMap.Resources["Rack"].Resources

	assert.NotNil(listed)

	printRacks(listed)
}

func TestPrintLeases(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rootCmd := New()
	assert.NotNil(rootCmd)

	rootCmd.AddCommand(test())

	fact := new(zebra.ResourceFactory)

	resMap := zebra.NewResourceMap(*fact)

	assert.NotNil(resMap)

	l := new([]zebra.Resource)

	assert.NotNil(l)

	printLeases(*l)
}

func TestPrintUsers(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rootCmd := New()
	assert.NotNil(rootCmd)

	rootCmd.AddCommand(test())

	fact := new(zebra.ResourceFactory)

	resMap := zebra.NewResourceMap(*fact)

	usr := new(auth.User)

	usr.Role = new(auth.Role)

	usr.Status = new(zebra.Status)

	resMap.Add(usr, "User")

	listed := resMap.Resources["User"].Resources

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
