package api_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/api"
	"github.com/project-safari/zebra/dc"
	"github.com/project-safari/zebra/network"
	"github.com/stretchr/testify/assert"
)

//nolint:gochecknoglobals
var (
	resource1 = `{"VLANPool":[{"id":"0100000001","type":"VLANPool","labels":{"owner":"shravya"},` +
		`"rangeStart":0,"rangeEnd":10}]}`
	resource2 = `{"VLANPool":[{"id":"0100000002","type":"VLANPool","labels":{"owner":"nandyala"},` +
		`"rangeStart":1,"rangeEnd":5}]}`
	resources = `{"VLANPool":[{"id":"0100000001","type":"VLANPool","labels":{"owner":"shravya"},` +
		`"rangeStart":0,"rangeEnd":10},{"id":"0100000002","type":"VLANPool","labels":{"owner":"nandyala"},` +
		`"rangeStart":1,"rangeEnd":5}]}`
	otherResources = `{"VLANPool":[{"id":"0100000002","type":"VLANPool","labels":{"owner":"nandyala"},` +
		`"rangeStart":1,"rangeEnd":5},{"id":"0100000001","type":"VLANPool","labels":{"owner":"shravya"},` +
		`"rangeStart":0,"rangeEnd":10}]}`
	noResources = "{}"
)

func TestNew(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	assert.NotNil(api.NewResourceAPI(nil))
}

func TestInitialize(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	f := zebra.Factory().Add("VLANPool", func() zebra.Resource { return new(network.VLANPool) })

	api := api.NewResourceAPI(f)
	assert.Nil(api.Initialize("teststore"))
}

func TestGetResources(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	f := zebra.Factory().Add("VLANPool", func() zebra.Resource { return new(network.VLANPool) })
	myAPI := api.NewResourceAPI(f)
	assert.Nil(myAPI.Initialize("teststore"))

	handler := http.HandlerFunc(myAPI.GetResources)

	req := makeRequest(assert, "GET", "/resources", "")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusOK, rr.Code)

	resBody := rr.Body.String()
	assert.True(resBody == resources || resBody == otherResources)
}

func TestGetResourcesByID(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	f := zebra.Factory().Add("VLANPool", func() zebra.Resource { return new(network.VLANPool) })
	myAPI := api.NewResourceAPI(f)
	assert.Nil(myAPI.Initialize("teststore"))

	handler := http.HandlerFunc(myAPI.GetResourcesByID)

	// GET resource1
	req := makeRequest(assert, "GET", "/resources?id=0100000001", "")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusOK, rr.Code)

	resBody := rr.Body.String()
	assert.Equal(resource1, resBody)

	// GET resource2
	req = makeRequest(assert, "GET", "/resources?id=0100000002", "")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusOK, rr.Code)

	resBody = rr.Body.String()
	assert.Equal(resource2, resBody)

	// GET resources
	req = makeRequest(assert, "GET", "/resources?id=0100000001,0100000002", "")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusOK, rr.Code)

	resBody = rr.Body.String()
	assert.True(resBody == resources || resBody == otherResources)
}

func TestGetResourcesByType(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	f := zebra.Factory().Add("VLANPool", func() zebra.Resource { return new(network.VLANPool) })
	myAPI := api.NewResourceAPI(f)
	assert.Nil(myAPI.Initialize("teststore"))

	handler := http.HandlerFunc(myAPI.GetResourcesByType)

	// GET resources
	req := makeRequest(assert, "GET", "/resources?type=VLANPool", "")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusOK, rr.Code)

	resBody := rr.Body.String()
	assert.True(resBody == resources || resBody == otherResources)

	// GET no resources
	req = makeRequest(assert, "GET", "/resources?type=IPAddressPool", "")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusOK, rr.Code)

	resBody = rr.Body.String()
	assert.Equal(noResources, resBody)
}

func TestGetResourcesByProperty(t *testing.T) { //nolint:funlen
	t.Parallel()
	assert := assert.New(t)

	f := zebra.Factory().Add("VLANPool", func() zebra.Resource { return new(network.VLANPool) })
	myAPI := api.NewResourceAPI(f)
	assert.Nil(myAPI.Initialize("teststore"))

	handler := http.HandlerFunc(myAPI.GetResourcesByProperty)

	// GET resources
	req := makeRequest(assert, "GET", "/resources?property=Type-in-VLANPool,IPAddressPool", "")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusOK, rr.Code)

	resBody := rr.Body.String()
	assert.True(resBody == resources || resBody == otherResources)

	// GET no resources
	req = makeRequest(assert, "GET", "/resources?property=Type-notin-VLANPool,IPAddressPool", "")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusOK, rr.Code)

	resBody = rr.Body.String()
	assert.Equal(noResources, resBody)

	// GET resources
	req = makeRequest(assert, "GET", "/resources?property=Type-equal-VLANPool", "")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusOK, rr.Code)

	resBody = rr.Body.String()
	assert.True(resBody == resources || resBody == otherResources)

	// GET no resources
	req = makeRequest(assert, "GET", "/resources?property=Type-notequal-VLANPool", "")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusOK, rr.Code)

	resBody = rr.Body.String()
	assert.Equal(noResources, resBody)

	// Invalid request
	req = makeRequest(assert, "GET", "/resources?property=Type-notequal", "")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusBadRequest, rr.Code)

	req = makeRequest(assert, "GET", "/resources?property=Type-notequal-VLANPool,Lab", "")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusBadRequest, rr.Code)

	req = makeRequest(assert, "GET", "/resources?property=Type-blahblah-VLANPool,Lab", "")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusBadRequest, rr.Code)
}

func TestGetResourcesByLabel(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	f := zebra.Factory().Add("VLANPool", func() zebra.Resource { return new(network.VLANPool) })
	myAPI := api.NewResourceAPI(f)
	assert.Nil(myAPI.Initialize("teststore"))

	handler := http.HandlerFunc(myAPI.GetResourcesByLabel)

	// GET resources
	req := makeRequest(assert, "GET", "/resources?label=owner-in-shravya,nandyala", "")
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusOK, rr.Code)

	resBody := rr.Body.String()
	assert.True(resBody == resources || resBody == otherResources)

	// GET no resources
	req = makeRequest(assert, "GET", "/resources?label=owner-notin-shravya,nandyala", "")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusOK, rr.Code)

	resBody = rr.Body.String()
	assert.Equal(noResources, resBody)

	// GET resource1
	req = makeRequest(assert, "GET", "/resources?label=owner-equal-shravya", "")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusOK, rr.Code)

	resBody = rr.Body.String()
	assert.Equal(resource1, resBody)

	// GET resource2
	req = makeRequest(assert, "GET", "/resources?label=owner-notequal-shravya", "")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusOK, rr.Code)

	resBody = rr.Body.String()
	assert.Equal(resource2, resBody)

	// Invalid requests
	req = makeRequest(assert, "GET", "/resources?label=owner-equal", "")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusBadRequest, rr.Code)

	req = makeRequest(assert, "GET", "/resources?label=owner-equal-shravya,nandyala", "")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusBadRequest, rr.Code)
}

func TestPutResource(t *testing.T) { //nolint:funlen
	t.Parallel()
	assert := assert.New(t)

	t.Cleanup(func() { os.Remove("teststore/resources/01/00000003") })

	f := zebra.Factory().Add("VLANPool", func() zebra.Resource { return new(network.VLANPool) })
	f.Add("Lab", func() zebra.Resource { return new(dc.Lab) })
	myAPI := api.NewResourceAPI(f)
	assert.Nil(myAPI.Initialize("teststore"))

	handler := http.HandlerFunc(myAPI.PutResource)

	body := `{"id":"0100000003","type":"Lab","labels": {"owner": "shravya"},"name": "shravya's lab"}`

	// Test error handling
	req, err := http.NewRequest("POST", "/resources", nil)
	assert.Nil(err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Create new resource
	req = makeRequest(assert, "POST", "/resources", body)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusCreated, rr.Code)

	// Update existing resource
	req = makeRequest(assert, "POST", "/resources", body)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusOK, rr.Code)

	// Create resource without a type
	body = `{"id":"0100000003","labels": {"owner": "shravya"},"name": "shravya's lab"}`
	req = makeRequest(assert, "POST", "/resources", body)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusBadRequest, rr.Code)

	// Create resource with an invalid type
	body = `{"id":"0100000003","type":"test","labels": {"owner": "shravya"},"name": "shravya's lab"}`
	req = makeRequest(assert, "POST", "/resources", body)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusBadRequest, rr.Code)

	// Create invalid resource
	body = `{"id":"0100000003","type":"Lab"}`
	req = makeRequest(assert, "POST", "/resources", body)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusBadRequest, rr.Code)

	// Trigger ioutil.ReadAll() panic
	body = ""
	req = makeRequest(assert, "POST", "/resources", body)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusBadRequest, rr.Code)

	// Create resource with no information
	body = " "
	req = makeRequest(assert, "POST", "/resources", body)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusBadRequest, rr.Code)
}

func TestDeleteResource(t *testing.T) { //nolint:funlen
	t.Parallel()
	assert := assert.New(t)

	root := "teststore1"

	t.Cleanup(func() { os.RemoveAll(root) })

	f := zebra.Factory().Add("Lab", func() zebra.Resource { return new(network.VLANPool) })
	myAPI := api.NewResourceAPI(f)
	assert.Nil(myAPI.Initialize(root))

	handler := http.HandlerFunc(myAPI.DeleteResource)

	lab1 := &dc.Lab{
		NamedResource: zebra.NamedResource{
			BaseResource: zebra.BaseResource{
				ID:     "10000001",
				Type:   "Lab",
				Labels: nil,
			},
			Name: "Lab1",
		},
	}

	lab2 := &dc.Lab{
		NamedResource: zebra.NamedResource{
			BaseResource: zebra.BaseResource{
				ID:     "10000002",
				Type:   "Lab",
				Labels: nil,
			},
			Name: "Lab2",
		},
	}

	assert.Nil(myAPI.Store.Create(lab1))
	assert.Nil(myAPI.Store.Create(lab2))

	// Invalid resources requested to be deleted
	body := `["10000003", "10000004"]`
	req := makeRequest(assert, "DELETE", "/resources", body)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusOK, rr.Code)

	expected := "Invalid resource IDs: 10000003, 10000004\n"
	resBody := rr.Body.String()
	assert.Equal(expected, resBody)

	// DELETE resources
	body = `["10000001", "10000002"]`
	req = makeRequest(assert, "DELETE", "/resources", body)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusOK, rr.Code)

	expected = "Deleted the following resources: 10000001, 10000002\n"
	resBody = rr.Body.String()
	assert.Equal(expected, resBody)

	// Trigger ioutil.ReadAll() panic
	body = ""
	req = makeRequest(assert, "DELETE", "/resources", body)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusBadRequest, rr.Code)

	// Bad request, type of list
	body = `[1, 2]`
	req = makeRequest(assert, "DELETE", "/resources", body)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusBadRequest, rr.Code)
}

func makeRequest(assert *assert.Assertions, method string, url string, body string) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	assert.Nil(err)
	assert.NotNil(req)

	if body != "" {
		req.Body = ioutil.NopCloser(bytes.NewBufferString(body))
	}

	return req
}
