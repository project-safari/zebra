package main //nolint:testpackage

import (
	"testing"

	"github.com/project-safari/zebra/auth"
	"github.com/project-safari/zebra/cmd/herd/pkg"
	"github.com/project-safari/zebra/compute"
	"github.com/project-safari/zebra/dc"
	"github.com/project-safari/zebra/network"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func dcData() (map[string]*dc.Datacenter, *dc.Datacenter) {
	this := make(map[string]*dc.Datacenter)
	valDC := &dc.Datacenter{} //nolint:exhaustruct,exhaustivestruct

	return this, valDC
}

func test() *cobra.Command {
	testCmd := NewShow()

	return testCmd
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

func TestNewNetCmd(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	netC := NewNetCmd(test())
	assert.NotNil(netC)
}

func TestNewSrvCmd(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	srvC := NewSrvCmd(test())
	assert.NotNil(srvC)
}

func TestNewDcCmd(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	dcC := NewDCCmd(test())
	assert.NotNil(dcC)
}

func TestShowServer(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"servers", "test-case"}

	srvCmd := NewSrvCmd(test())
	rootCmd := New()
	rootCmd.AddCommand(srvCmd)
	serv := ShowServ(rootCmd, args)

	assert.NotNil(serv)

	toPrint := make(map[string]*compute.Server)

	name := args[0]
	val := &compute.Server{} //nolint:exhaustruct,exhaustivestruct

	toPrint[name] = val

	printed := printServer(toPrint)

	assert.NotNil(printed)
}

func TestShowVC(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"vcenters", "test-case"}

	vc := NewSrvCmd(test())
	rootCmd := New()

	rootCmd.AddCommand(vc)

	vcShow := ShowVC(rootCmd, args)
	assert.NotNil(vcShow)

	toPrint := make(map[string]*compute.VCenter)

	name := args[0]
	val := &compute.VCenter{} //nolint:exhaustruct,exhaustivestruct

	toPrint[name] = val

	printed := printVC(toPrint)

	assert.NotNil(printed)
}

func TestShowVlan(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"vlans", "test-case"}

	v := NewNetCmd(test())
	rootCmd := New()

	rootCmd.AddCommand(v)

	vlan := ShowVlan(rootCmd, args)

	assert.NotNil(vlan)

	toPrint := make(map[string]*network.VLANPool)

	name := args[0]
	val := &network.VLANPool{} //nolint:exhaustruct,exhaustivestruct

	toPrint[name] = val

	printed := printNets(toPrint)

	assert.NotNil(printed)
}

//

func TestShowSw(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"sws", "test-case"}

	netCmd := NewNetCmd(test())
	rootCmd := New()

	rootCmd.AddCommand(netCmd)
	sw := ShowSw(rootCmd, args)

	assert.NotNil(sw)

	toPrint := make(map[string]*network.Switch)

	name := args[0]
	val := &network.Switch{} //nolint:exhaustruct,exhaustivestruct

	toPrint[name] = val

	printed := printSwitch(toPrint)

	assert.NotNil(printed)
}

func TestShowRack(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"racks", "test-case"}

	rackCmd := NewDCCmd(test())
	rootCmd := New()

	rootCmd.AddCommand(rackCmd)
	rack := ShowRack(rootCmd, args)

	assert.NotNil(rack)

	toPrint := make(map[string]*dc.Rack)

	name := args[0]
	val := &dc.Rack{} //nolint:exhaustruct,exhaustivestruct

	toPrint[name] = val

	printed := printRack(toPrint)

	assert.NotNil(printed)
}

func TestShowLab(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"labs", "test-case"}

	labCmd := NewDCCmd(test())
	rootCmd := New()

	rootCmd.AddCommand(labCmd)

	lab := ShowLab(rootCmd, args)

	assert.NotNil(lab)

	toPrint := make(map[string]*dc.Lab)

	name := args[0]
	val := &dc.Lab{} //nolint:exhaustruct,exhaustivestruct

	toPrint[name] = val

	printed := printLab(toPrint)

	assert.NotNil(printed)
}

//
//
func TestShowESX(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"esx", "test-case"}

	esxCmd := NewSrvCmd(test())
	rootCmd := New()

	rootCmd.AddCommand(esxCmd)

	esx := ShowESX(rootCmd, args)

	assert.NotNil(esx)

	toPrint := make(map[string]*compute.ESX)

	name := args[0]
	val := &compute.ESX{} //nolint:exhaustruct,exhaustivestruct

	toPrint[name] = val

	printed := printESX(toPrint)

	assert.NotNil(printed)
}

func TestShowDC(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"dc", "test-case"}

	dcCmd := NewDCCmd(test())
	rootCmd := New()

	rootCmd.AddCommand(dcCmd)

	dc := ShowDC(rootCmd, args)

	assert.NotNil(dc)

	name := args[0]

	this, valDC := dcData()

	this[name] = valDC

	printed := printDC(this)

	assert.NotNil(printed)
}

func TestShowUser(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"users", "test-case"}

	rootCmd := New()
	rootCmd.AddCommand(test())

	user := ShowUsr(rootCmd, args)

	assert.NotNil(user)

	toPrint := make(map[string]*auth.User)

	name := args[0]
	val := &auth.User{} //nolint:exhaustruct,exhaustivestruct

	val.Key = &auth.RsaIdentity{}
	val.PasswordHash = pkg.Password("user1")

	val.Role = new(auth.Role)

	val.Labels = pkg.CreateLabels()
	val.Labels = pkg.GroupLabels(val.Labels, "group1")

	val.Email = "sample@yahoo.com"

	toPrint[name] = val

	printed := printUser(toPrint)

	assert.NotNil(printed)
}

func TestShowReg(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"registrations", "test-case"}

	rootCmd := New()
	rootCmd.AddCommand(test())

	reg := ShowReg(rootCmd, args)

	assert.NotNil(reg)

	toPrint := make(map[string]*auth.User)

	name := args[0]
	val := &auth.User{} //nolint:exhaustruct,exhaustivestruct

	val.Key = &auth.RsaIdentity{}
	val.PasswordHash = pkg.Password("user2")

	val.Role = new(auth.Role)

	val.Labels = pkg.CreateLabels()
	val.Labels = pkg.GroupLabels(val.Labels, "group1")

	val.Email = "sample@yahoo.com"

	toPrint[name] = val

	toPrint[name] = val

	printed := printUser(toPrint)

	assert.NotNil(printed)
}

func TestShowIP(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"ip", "test-case"}

	rootCmd := New()
	rootCmd.AddCommand(test())

	addr := ShowIP(rootCmd, args)

	assert.NotNil(addr)

	toPrint := make(map[string]*network.IPAddressPool)

	name := args[0]
	val := &network.IPAddressPool{} //nolint:exhaustruct,exhaustivestruct

	toPrint[name] = val

	printed := printIP(toPrint)

	assert.NotNil(printed)
}

//

func TestShowVM(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	args := []string{"vms", "test-case"}

	rootCmd := New()
	rootCmd.AddCommand(test())

	machine := ShowVM(rootCmd, args)

	assert.NotNil(machine)

	toPrint := make(map[string]*compute.VM)

	name := args[0]
	val := &compute.VM{} //nolint:exhaustruct,exhaustivestruct

	toPrint[name] = val

	printed := printVM(toPrint)

	assert.NotNil(printed)
}
