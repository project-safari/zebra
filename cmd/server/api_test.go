package main //nolint:testpackage

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/dc"
	"github.com/project-safari/zebra/network"
	"github.com/project-safari/zebra/store"
	"github.com/stretchr/testify/assert"
)

func makeQueryRequest(assert *assert.Assertions, q *QueryRequest) *http.Request {
	req, err := http.NewRequest("GET", "/api/v1/resources", nil)
	assert.Nil(err)
	assert.NotNil(req)

	b, e := json.Marshal(q)
	assert.Nil(e)

	req.Body = ioutil.NopCloser(bytes.NewBuffer(b))

	return req
}

func TestQuery(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "testquery"

	t.Cleanup(func() { os.RemoveAll(root) })

	api := NewResourceAPI(store.DefaultFactory())
	assert.Nil(api.Initialize(root))

	h := handleQuery(context.Background(), api)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r, nil)
	})

	qr := new(QueryRequest)

	req := makeQueryRequest(assert, qr)
	handler.ServeHTTP(rr, req)
	assert.Equal(rr.Code, http.StatusOK)

	qr.IDs = []string{"0100000001"}
	req = makeQueryRequest(assert, qr)
	handler.ServeHTTP(rr, req)
	assert.Equal(rr.Code, http.StatusOK)

	qr.IDs = []string{}
	qr.Types = []string{"VLANPool"}
	req = makeQueryRequest(assert, qr)
	handler.ServeHTTP(rr, req)
	assert.Equal(rr.Code, http.StatusOK)

	qr.Types = []string{}
	qr.Labels = []zebra.Query{
		{Op: zebra.MatchEqual, Key: "test", Values: []string{"test"}},
		{Op: zebra.MatchIn, Key: "test2", Values: []string{"test1", "test2"}},
	}
	req = makeQueryRequest(assert, qr)
	handler.ServeHTTP(rr, req)
	assert.Equal(rr.Code, http.StatusOK)
}

func TestBadQuery(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "testbadquery"

	t.Cleanup(func() { os.RemoveAll(root) })

	api := NewResourceAPI(store.DefaultFactory())
	assert.Nil(api.Initialize(root))

	h := handleQuery(context.Background(), api)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r, nil)
	})

	// Cannot have both IDs and Types
	qr := new(QueryRequest)
	qr.IDs = []string{"0100000001"}
	qr.Types = []string{"VLANPool"}
	req := makeQueryRequest(assert, qr)
	handler.ServeHTTP(rr, req)
	assert.Equal(rr.Code, http.StatusBadRequest)

	// Cannot have both Types and Properties
	qr.IDs = []string{}
	qr.Properties = []zebra.Query{{Op: zebra.MatchEqual, Key: "test", Values: []string{"test"}}}
	req = makeQueryRequest(assert, qr)
	handler.ServeHTTP(rr, req)
	assert.Equal(rr.Code, http.StatusBadRequest)

	// Cannot have Labels with anything else
	qr.Properties = []zebra.Query{}
	qr.Labels = []zebra.Query{
		{Op: zebra.MatchEqual, Key: "test", Values: []string{"test"}},
		{Op: zebra.MatchEqual, Key: "blah", Values: []string{"blah", "blah2"}},
	}
	req = makeQueryRequest(assert, qr)
	handler.ServeHTTP(rr, req)
	assert.Equal(rr.Code, http.StatusBadRequest)

	// Must have valid label queries
	qr.Types = []string{}
	req = makeQueryRequest(assert, qr)
	handler.ServeHTTP(rr, req)
	assert.Equal(rr.Code, http.StatusBadRequest)

	// Must have valid property queries
	qr.Properties = qr.Labels
	qr.Labels = nil
	req = makeQueryRequest(assert, qr)
	handler.ServeHTTP(rr, req)
	assert.Equal(rr.Code, http.StatusBadRequest)
}

func TestInvalidQuery(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "testinvalidquery"

	t.Cleanup(func() { os.RemoveAll(root) })

	api := NewResourceAPI(store.DefaultFactory())
	assert.Nil(api.Initialize(root))

	h := handleQuery(context.Background(), api)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r, nil)
	})

	// Invalid json request
	req, err := http.NewRequest("GET", "/api/v1/resources", nil)
	assert.Nil(err)
	assert.NotNil(req)

	v := "{...}" // bad json
	req.Body = ioutil.NopCloser(bytes.NewBufferString(v))

	handler.ServeHTTP(rr, req)
	assert.Equal(rr.Code, http.StatusBadRequest)
}

func TestNew(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	assert.NotNil(NewResourceAPI(nil))
}

func TestInitialize(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "api_teststore"

	t.Cleanup(func() { os.RemoveAll(root) })

	f := zebra.Factory().Add(network.VLANPoolType())

	api := NewResourceAPI(f)
	assert.Nil(api.Initialize(root))
}

func TestPutResource(t *testing.T) { //nolint:funlen
	t.Parallel()
	assert := assert.New(t)

	root := "api_teststore1"

	t.Cleanup(func() { os.RemoveAll(root) })

	f := zebra.Factory().Add(network.VLANPoolType()).Add(dc.LabType())
	myAPI := NewResourceAPI(f)
	assert.Nil(myAPI.Initialize(root))

	handler := http.HandlerFunc(myAPI.PutResource)

	body := `{"id":"0100000003","type":"Lab","labels": {"owner": "shravya"},"name": "shravya's lab"}`

	// Test error handling
	req, err := http.NewRequest("PUT", "/resources", nil)
	assert.Nil(err)

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Create new resource
	req = createRequest(assert, "PUT", "/resources", body)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusCreated, rr.Code)

	// Update existing resource
	req = createRequest(assert, "PUT", "/resources", body)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusOK, rr.Code)

	// Create resource without a type
	body = `{"id":"0100000003","labels": {"owner": "shravya"},"name": "shravya's lab"}`
	req = createRequest(assert, "PUT", "/resources", body)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusBadRequest, rr.Code)

	// Create resource with an invalid type
	body = `{"id":"0100000003","type":"test","labels": {"owner": "shravya"},"name": "shravya's lab"}`
	req = createRequest(assert, "PUT", "/resources", body)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusBadRequest, rr.Code)

	// Create invalid resource
	body = `{"id":"0100000003","type":"Lab"}`
	req = createRequest(assert, "PUT", "/resources", body)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusBadRequest, rr.Code)

	// Trigger ioutil.ReadAll() panic
	body = ""
	req = createRequest(assert, "PUT", "/resources", body)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusBadRequest, rr.Code)

	// Create resource with no information
	body = " "
	req = createRequest(assert, "PUT", "/resources", body)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusBadRequest, rr.Code)
}

func TestDeleteResource(t *testing.T) { //nolint:funlen
	t.Parallel()
	assert := assert.New(t)

	root := "api_teststore2"

	t.Cleanup(func() { os.RemoveAll(root) })

	f := zebra.Factory().Add(dc.LabType())
	myAPI := NewResourceAPI(f)
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
	req := createRequest(assert, "DELETE", "/resources", body)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusOK, rr.Code)

	expected := "Invalid resource IDs: 10000003, 10000004\n"
	resBody := rr.Body.String()
	assert.Equal(expected, resBody)

	// DELETE resources
	body = `["10000001", "10000002"]`
	req = createRequest(assert, "DELETE", "/resources", body)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusOK, rr.Code)

	expected = "Deleted the following resources: 10000001, 10000002\n"
	resBody = rr.Body.String()
	assert.Equal(expected, resBody)

	// Trigger ioutil.ReadAll() panic
	body = ""
	req = createRequest(assert, "DELETE", "/resources", body)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusBadRequest, rr.Code)

	// Bad request, type of list
	body = `[1, 2]`
	req = createRequest(assert, "DELETE", "/resources", body)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusBadRequest, rr.Code)
}

func createRequest(assert *assert.Assertions, method string, url string, body string) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	assert.Nil(err)
	assert.NotNil(req)

	if body != "" {
		req.Body = ioutil.NopCloser(bytes.NewBufferString(body))
	}

	return req
}
