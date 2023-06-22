// nolint: funlen // Some functions need to be longer.
package main //nolint:testpackage

import (
	"testing"

	"github.com/project-safari/zebra"
	"github.com/stretchr/testify/assert"
)

func TestDetermineType(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// compute category
	means := "Compute"
	resName := "esxServer"

	result := determineType(means, resName)
	assert.Equal(result, "compute.esx")

	resName = "JENKINS"
	result = determineType(means, resName)
	assert.Equal(result, "compute.server")

	resName = "BLD123"
	result = determineType(means, resName)
	assert.Equal(result, "dc.datacenter")

	resName = "VLAN"
	result = determineType(means, resName)
	assert.Equal(result, "network.vlanPool")

	resName = "switchA"
	result = determineType(means, resName)
	assert.Equal(result, "network.switch")

	resName = "capic-1"
	result = determineType(means, resName)
	assert.Equal(result, "compute.vm")

	resName = "xYvapic/122"
	result = determineType(means, resName)
	assert.Equal(result, "compute.vcenter")

	resName = "Ipc"
	result = determineType(means, resName)
	assert.Equal(result, "network.ipAddressPool")

	// larger other category
	means = "Other"
	resName = "ixia"

	result = determineType(means, resName)
	assert.Equal(result, "dc.rack")

	resName = "nexus"

	result = determineType(means, resName)
	assert.Equal(result, "network.switch")

	// no category
	means = ""
	resName = ""

	result = determineType(means, resName)
	assert.Equal(result, "")
}

//nolint:funlen
func TestDetermineIDMeaning(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// test for vm.
	id := "2"
	name := "VM"

	result := determineIDMeaning(id, name)
	assert.Equal(result, "compute.vm")

	// test for rack with name shelf.
	id = "30"
	name = "Shelf"
	result = determineIDMeaning(id, name)
	assert.Equal(result, "dc.rack")

	// test for rack with name rack.
	name = "Rack"
	result = determineIDMeaning(id, name)
	assert.Equal(result, "dc.rack")

	// test for vc.
	id = "38"
	name = "VC"
	result = determineIDMeaning(id, name)
	assert.Equal(result, "compute.vcenter")

	// test for server.
	id = "4"
	name = "server"
	result = determineIDMeaning(id, name)
	assert.Equal(result, "compute.server")

	// test for sw.
	id = "8"
	name = "sw"
	result = determineIDMeaning(id, name)
	assert.Equal(result, "network.switch")

	// tests for compute's id.
	id = "1504"
	name = "sw"
	result = determineIDMeaning(id, name)
	assert.Equal(result, "network.switch")

	id = "1504"
	name = "/"
	result = determineIDMeaning(id, name)
	assert.Equal("N/A", result)

	// test for other's id.
	id = "1503"
	name = "chasis"
	result = determineIDMeaning(id, name)
	assert.Equal(result, "dc.rack")

	// test for wrong id.
	id = "0"
	name = "chasis"
	result = determineIDMeaning(id, name)
	assert.Equal(result, "unclassified")
}

//nolint:funlen
func TestCreateRes(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	var rt Racktables

	// test for creating an empty resource
	testEmpty1 := CreateResFromData(rt)

	assert.Nil(testEmpty1)

	rt.AssetNo = "1"
	rt.ID = "123"
	rt.IP = "1.1.1.1"
	rt.Name = "test-switch"
	rt.ObjtypeID = "8"

	// test for creating a switch
	rt.Type = zebra.Type{
		Name:        "network.switch",
		Description: "tester network.switch",
	}

	testCreateSwitch := CreateResFromData(rt)
	assert.NotNil(testCreateSwitch)

	// test for creating a dc
	rt.Type = zebra.Type{
		Name:        "dc.dataceneter",
		Description: "tester dc.dataceneter",
	}

	testCreateDC := CreateResFromData(rt)
	assert.NotNil(testCreateDC)

	// test for creating a lab
	rt.Type = zebra.Type{
		Name:        "dc.lab",
		Description: "tester lab",
	}

	testCreateLab := CreateResFromData(rt)
	assert.NotNil(testCreateLab)

	// test for creating a rack with shelf type
	rt.Type = zebra.Type{
		Name:        "dc.shelf",
		Description: "tester dc.dataceneter",
	}

	testCreateShelf := CreateResFromData(rt)
	assert.NotNil(testCreateShelf)

	// test for creating a vm
	rt.Type = zebra.Type{
		Name:        "compute.vm",
		Description: "tester compute.vm",
	}

	testCreateVM := CreateResFromData(rt)
	assert.NotNil(testCreateVM)

	// test for creating a vc
	rt.Type = zebra.Type{
		Name:        "compute.vceneter",
		Description: "tester compute.vceneter",
	}

	testCreateVC := CreateResFromData(rt)
	assert.NotNil(testCreateVC)

	// test for creating a server
	rt.Type = zebra.Type{
		Name:        "compute.server",
		Description: "tester ompute.server",
	}

	testCreateSrv := CreateResFromData(rt)
	assert.NotNil(testCreateSrv)

	// test for creating an esx server
	rt.Type = zebra.Type{
		Name:        "compute.esx",
		Description: "tester compute.esx",
	}

	testCreateESX := CreateResFromData(rt)
	assert.NotNil(testCreateESX)

	// test for creating a IPAddressPool
	rt.Type = zebra.Type{
		Name:        "network.ipaddresspool",
		Description: "tester network.ipaddresspool",
	}

	testCreateIP := CreateResFromData(rt)
	assert.NotNil(testCreateIP)

	// test for creating a vlanPool
	rt.Type = zebra.Type{
		Name:        "network.vlanpool",
		Description: "network.vlanpool",
	}

	testCreateVP := CreateResFromData(rt)
	assert.NotNil(testCreateVP)
}

func TestFiller(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	var rt Racktables

	rt.AssetNo = "1"
	rt.ID = "123"
	rt.IP = "1.1.1.1"
	rt.Name = "test-switch"
	rt.ObjtypeID = "8"

	rt.Type = zebra.Type{
		Name:        "network.switch",
		Description: "network server",
	}
	testSwitchFiller := switchFiller(rt)
	assert.NotNil(testSwitchFiller)

	rt.Name = "test-server"
	rt.ObjtypeID = "8"

	rt.Type = zebra.Type{
		Name:        "compute.server",
		Description: "data center server",
	}

	testServerFiller := serverFiller(rt)
	assert.NotNil(testServerFiller)

	rt.Name = "test-esx"
	rt.ObjtypeID = "9"

	rt.Type = zebra.Type{
		Name:        "compute.esx",
		Description: "data center esx",
	}

	testESXfiller := esxFiller(rt)
	assert.NotNil(testESXfiller)

	rt.Name = "test-vc"
	rt.ObjtypeID = "9"

	rt.Type = zebra.Type{
		Name:        "compute.vcenter",
		Description: "VMWare vcenter",
	}

	testVCfiller := vcenterFiller(rt)
	assert.NotNil(testVCfiller)

	rt.Name = "test-vm"
	rt.ObjtypeID = "9"

	rt.Type = zebra.Type{
		Name:        "compute.vm",
		Description: "virtual machine",
	}

	testVMfiller := vmFiller(rt)
	assert.NotNil(testVMfiller)

	rt.Name = "network.vlanPool"

	rt.ObjtypeID = "10"

	rt.Type = zebra.Type{
		Name:        "network.vlanPool",
		Description: "tester",
	}

	testVLANfiller := vlanFiller(rt)
	assert.NotNil(testVLANfiller)
}
