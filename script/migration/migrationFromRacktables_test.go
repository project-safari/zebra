package main //nolint:testpackage

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/script"
	"github.com/stretchr/testify/assert"
)

func TestDetermineType(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// compute category
	means := "Compute"
	resName := "esxServer"

	result := script.DetermineType(means, resName)
	assert.Equal(result, "compute.esx")

	resName = "JENKINS"
	result = script.DetermineType(means, resName)
	assert.Equal(result, "compute.server")

	resName = "BLD123"
	result = script.DetermineType(means, resName)
	assert.Equal(result, "dc.datacenter")

	resName = "VLAN"
	result = script.DetermineType(means, resName)
	assert.Equal(result, "network.vlanPool")

	resName = "switchA"
	result = script.DetermineType(means, resName)
	assert.Equal(result, "network.switch")

	resName = "capic-1"
	result = script.DetermineType(means, resName)
	assert.Equal(result, "compute.vm")

	resName = "xYvapic/122"
	result = script.DetermineType(means, resName)
	assert.Equal(result, "compute.vcenter")

	resName = "Ipc"
	result = script.DetermineType(means, resName)
	assert.Equal(result, "network.ipAddressPool")

	// larger other category
	means = "Other"
	resName = "ixia"

	result = script.DetermineType(means, resName)
	assert.Equal(result, "dc.rack")

	resName = "nexus"

	result = script.DetermineType(means, resName)
	assert.Equal(result, "network.switch")

	// no category
	means = ""
	resName = ""

	result = script.DetermineType(means, resName)
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
	assert.Equal(result, "")

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

// Test for reading json.

func TestRead(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	req := makeLabelRequest(assert, nil, "a", "b", "c")

	labelReq := &struct {
		Labels []string `json:"labels"`
	}{Labels: []string{}}

	assert.Nil(script.ReadJSON(context.Background(), req, labelReq))

	// Bad IO reader
	req.Body = ioutil.NopCloser(mockReader{err: true})
	assert.NotNil(script.ReadJSON(context.Background(), req, nil))

	// Empty Body
	req.Body = ioutil.NopCloser(mockReader{err: false})
	assert.NotNil(script.ReadJSON(context.Background(), req, nil))
}

type mockReader struct {
	err bool
}

func (f mockReader) Read(b []byte) (int, error) {
	if f.err {
		return 0, errors.New("mock error") //nolint:goerr113
	}

	return 0, io.EOF
}

func makeLabelRequest(assert *assert.Assertions, resources *script.ResourceAPI, labels ...string) *http.Request {
	ctx := context.WithValue(context.Background(), script.ResourcesCtxKey, resources)
	ctx = context.WithValue(ctx, script.AuthCtxKey, "test-label")

	req, err := http.NewRequestWithContext(ctx, "GET", "/api/v1/labels", nil)
	assert.Nil(err)
	assert.NotNil(req)

	v := map[string][]string{"labels": labels}
	b, e := json.Marshal(v)
	assert.Nil(e)

	req.Body = ioutil.NopCloser(bytes.NewBuffer(b))

	return req
}

func TestAllData(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	var rack Racktables

	rackArr := []Racktables{}

	rack.ID = "123"
	rack.Name = "test-rack"

	rackArr = append(rackArr, rack)

	assert.NotNil((rackArr))
}

//nolint:funlen
func TestCreateRes(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	var rt Racktables

	// test for creating an empty resource
	testEmpty1, testEmpty2 := createResFromData(rt)

	assert.Nil(testEmpty1)

	assert.Equal(testEmpty2, "")

	rt.AssetNo = "1"
	rt.ID = "123"
	rt.IP = "1.1.1.1"
	rt.Name = "test-switch"
	rt.ObjtypeID = "8"

	// test for creating a switch
	rt.Type = "network.switch"
	testCreateSwitch, _ := createResFromData(rt)
	assert.NotNil(testCreateSwitch)

	// test for creating a dc
	rt.Type = "dc.datacenetr"
	testCreateDC, _ := createResFromData(rt)
	assert.NotNil(testCreateDC)

	// test for creating a lab
	rt.Type = "dc.lab"
	testCreateLab, _ := createResFromData(rt)
	assert.NotNil(testCreateLab)

	// test for creating a rack with shelf type
	rt.Type = "dc.shelf"
	testCreateShelf, _ := createResFromData(rt)
	assert.NotNil(testCreateShelf)

	// test for creating a vm
	rt.Type = "compute.vm"
	testCreateVM, _ := createResFromData(rt)
	assert.NotNil(testCreateVM)

	// test for creating a vc
	rt.Type = "compute.vcenetr"
	testCreateVC, _ := createResFromData(rt)
	assert.NotNil(testCreateVC)

	// test for creating a server
	rt.Type = "compute.server"
	testCreateSrv, _ := createResFromData(rt)
	assert.NotNil(testCreateSrv)

	// test for creating an esx server
	rt.Type = "compute.esx"
	testCreateESX, _ := createResFromData(rt)
	assert.NotNil(testCreateESX)

	// test for creating a IPAddressPool
	rt.Type = "network.ipaddresspool"
	testCreateIP, _ := createResFromData(rt)
	assert.NotNil(testCreateIP)

	// test for creating a vlanPool
	rt.Type = "network.vlanpool"
	testCreateVP, _ := createResFromData(rt)
	assert.NotNil(testCreateVP)
}

var errFake = errors.New("fake error")

type fakeReader struct {
	err bool
}

func (f fakeReader) Read(b []byte) (int, error) {
	if f.err {
		return 0, errFake
	}

	return 0, io.EOF
}

func TestReadJSON(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	req := makeLabelRequest(assert, nil, "a", "b", "c")

	labelReq := &struct {
		Labels []string `json:"labels"`
	}{Labels: []string{}}

	assert.Nil(script.ReadJSON(context.Background(), req, labelReq))

	// Bad IO reader
	req.Body = ioutil.NopCloser(fakeReader{err: true})
	assert.NotNil(script.ReadJSON(context.Background(), req, nil))

	// Empty Body
	req.Body = ioutil.NopCloser(fakeReader{err: false})
	assert.NotNil(script.ReadJSON(context.Background(), req, nil))
}

func TestDetermineParentType(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	parID := "4"
	childID := "1504"

	testParent := script.DetermineParentType(parID, childID, "ipc")

	assert.Equal("network.ipAddressPool", testParent)

	childID = "1507"
	testParent = script.DetermineParentType(parID, childID, "vshield")

	assert.Equal("infrastructure", testParent)

	parID = "8"
	testParent = script.DetermineParentType(parID, childID, "nexus")

	assert.Equal("network.switch", testParent)
}

func TestGetParent(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	testParent := script.GetParent("compute.esx")
	assert.NotNil(testParent)
	assert.Equal(testParent, "compute.server")

	testParent = script.GetParent("compute.vm")
	assert.NotNil(testParent)
	assert.Equal(testParent, "compute.esx")

	testParent = script.GetParent("dc.rack")
	assert.NotNil(testParent)
	assert.Equal(testParent, "dc.lab")

	testParent = script.GetParent("dc.lab")
	assert.NotNil(testParent)
	assert.Equal(testParent, "dc.datacenter")

	testParent = script.GetParent("compute.server")
	assert.NotNil(testParent)
	assert.Equal(testParent, "dc.rack")
}

//nolint:lll
func TestCreateRequests(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	factory := zebra.Factory()

	myAPI := script.NewResourceAPI(factory)

	body := `{"lab":[{"id":` + "123" + `,"type":` + "test-type" + `,"name":` + "test-name" + `,"owner":` + "test-owner" + "}]}"

	assert.NotNil(createRequests("POST", "/resources", body, myAPI))
}
