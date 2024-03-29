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

	"github.com/julienschmidt/httprouter"
	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/model"
	"github.com/project-safari/zebra/model/dc"
	"github.com/stretchr/testify/assert"
)

func makeQueryRequest(assert *assert.Assertions, resources *ResourceAPI, q *QueryRequest) *http.Request {
	ctx := context.WithValue(context.Background(), ResourcesCtxKey, resources)
	req, err := http.NewRequestWithContext(ctx, "GET", "/api/v1/resources", nil)
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

	root := "test_query"

	defer func() { os.RemoveAll(root) }()

	api := NewResourceAPI(model.Factory())
	assert.Nil(api.Initialize(root))

	h := handleQuery()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r, nil)
	})

	qr := new(QueryRequest)

	req := makeQueryRequest(assert, api, qr)
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusOK, rr.Code)

	qr.IDs = []string{"0100000001"}
	req = makeQueryRequest(assert, api, qr)
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusOK, rr.Code)

	qr.IDs = []string{}
	qr.Types = []string{"VLANPool"}
	req = makeQueryRequest(assert, api, qr)
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusOK, rr.Code)

	qr.Types = []string{}
	qr.Labels = []zebra.Query{
		{Op: zebra.MatchEqual, Key: "test", Values: []string{"test"}},
		{Op: zebra.MatchIn, Key: "test2", Values: []string{"test1", "test2"}},
	}
	req = makeQueryRequest(assert, api, qr)
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusOK, rr.Code)
}

func TestEmptyQuery(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "test_empty_query"

	defer func() { os.RemoveAll(root) }()

	api := NewResourceAPI(model.Factory())
	assert.Nil(api.Initialize(root))

	h := handleQuery()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r, nil)
	})

	ctx := context.WithValue(context.Background(), ResourcesCtxKey, api)
	req, err := http.NewRequestWithContext(ctx, "GET", "/api/v1/resources", nil)
	req.Body = ioutil.NopCloser(bytes.NewBuffer([]byte("")))

	assert.Nil(err)
	assert.NotNil(req)
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusOK, rr.Code)
}

func TestBadQuery(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "test_bad_query"

	defer func() { os.RemoveAll(root) }()

	api := NewResourceAPI(model.Factory())
	assert.Nil(api.Initialize(root))

	h := handleQuery()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r, nil)
	})

	// Cannot have both IDs and Types
	qr := new(QueryRequest)
	qr.IDs = []string{"0100000001"}
	qr.Types = []string{"VLANPool"}
	req := makeQueryRequest(assert, api, qr)
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusBadRequest, rr.Code)

	// Cannot have both Types and Properties
	qr.IDs = []string{}
	qr.Properties = []zebra.Query{{Op: zebra.MatchEqual, Key: "test", Values: []string{"test"}}}
	req = makeQueryRequest(assert, api, qr)
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusBadRequest, rr.Code)

	// Cannot have Labels with anything else
	qr.Properties = []zebra.Query{}
	qr.Labels = []zebra.Query{
		{Op: zebra.MatchEqual, Key: "test", Values: []string{"test"}},
		{Op: zebra.MatchEqual, Key: "blah", Values: []string{"blah", "blah2"}},
	}
	req = makeQueryRequest(assert, api, qr)
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusBadRequest, rr.Code)

	// Must have valid label queries
	qr.Types = []string{}
	req = makeQueryRequest(assert, api, qr)
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusBadRequest, rr.Code)

	// Must have valid property queries
	qr.Properties = qr.Labels
	qr.Labels = nil
	req = makeQueryRequest(assert, api, qr)
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusBadRequest, rr.Code)
}

func TestInvalidQuery(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "test_invalid_query"

	defer func() { os.RemoveAll(root) }()

	api := NewResourceAPI(model.Factory())
	assert.Nil(api.Initialize(root))

	h := handleQuery()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r, nil)
	})

	// Invalid context
	req, err := http.NewRequest("GET", "/api/v1/resources", nil)
	assert.Nil(err)
	assert.NotNil(req)

	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusInternalServerError, rr.Code)

	// Invalid json request
	ctx := context.WithValue(context.Background(), ResourcesCtxKey, api)
	req, err = http.NewRequestWithContext(ctx, "GET", "/api/v1/resources", nil)
	assert.Nil(err)
	assert.NotNil(req)

	v := "{...}" // bad json
	req.Body = ioutil.NopCloser(bytes.NewBufferString(v))

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusBadRequest, rr.Code)
}

func TestNew(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	assert.NotNil(NewResourceAPI(nil))
}

func TestInitialize(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "test_initialize"

	defer func() { os.RemoveAll(root) }()

	f := model.Factory()

	api := NewResourceAPI(f)
	assert.Nil(api.Initialize(root))
}

func TestPostResource(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "test_post_resource"

	defer func() { os.RemoveAll(root) }()

	myAPI := NewResourceAPI(model.Factory())
	assert.Nil(myAPI.Initialize(root))

	h := handlePost()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r, nil)
	})

	body := `{"lab":[{"id":"0100000003","type":"Lab","labels": {"owner": "shravya"},"name": "shravya's lab"}]}`

	// Create new resource
	req := createRequest(assert, "POST", "/resources", body, myAPI)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.NotEqual(http.StatusOK, rr.Code)

	// Update existing resource
	req = createRequest(assert, "POST", "/resources", body, myAPI)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.NotEqual(http.StatusOK, rr.Code)

	// Create resource with an invalid type, won't read properly
	body = `{"lab":[{"id":"","type":"test","labels": {"owner": "shravya"},"name": "shravya's lab"}]}`
	req = createRequest(assert, "POST", "/resources", body, myAPI)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusBadRequest, rr.Code)

	// Create resource with an invalid ID
	body = `{"lab":[{"id":"","type":"Lab","labels": {"owner": "shravya"},"name": "shravya's lab"}]}`
	req = createRequest(assert, "POST", "/resources", body, myAPI)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusBadRequest, rr.Code)
}

func TestDeleteResource(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "test_delete_resource"

	defer func() { os.RemoveAll(root) }()

	myAPI := NewResourceAPI(model.Factory())
	assert.Nil(myAPI.Initialize(root))

	h := handleDelete()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r, nil)
	})

	lab1 := dc.NewLab("Lab1", "test_owner", "test_group")

	assert.Nil(myAPI.Store.Create(lab1))

	// Invalid resources requested to be deleted
	body := `{"type":[]}`
	req := createRequest(assert, "DELETE", "/resources", body, myAPI)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.NotEqual(http.StatusOK, rr.Code)

	body = `{"lab":[{"id":"","type":"","name": "shravya's lab"}]}`
	req = createRequest(assert, "DELETE", "/resources", body, myAPI)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusBadRequest, rr.Code)

	body = `{"lab":[{"id":"0","type":"Lab","name": "shravya's lab"}]}`
	req = createRequest(assert, "DELETE", "/resources", body, myAPI)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusBadRequest, rr.Code)

	// DELETE resources
	param := httprouter.Param{
		Key:   "id",
		Value: lab1.Meta.ID,
	}

	params := httprouter.Params{param}
	handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r, params)
	})

	req = createRequest(assert, "DELETE", "/resources", "", myAPI)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusOK, rr.Code)

	assert.Empty(myAPI.Store.Query().Resources)
}

func TestValidateQueries(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// Should fail
	qs := []zebra.Query{
		{Op: zebra.MatchIn, Key: "test", Values: []string{"blah", "blah2"}},
		{Op: zebra.MatchEqual, Key: "test", Values: []string{"blah", "blah2"}},
		{Op: 8, Key: "test", Values: []string{"blah", "blah2"}},
	}

	assert.Nil(qs[0].Validate())
	assert.NotNil(qs[1].Validate())
	assert.NotNil(qs[2].Validate())

	assert.NotNil(validateQueries(qs))
	assert.Nil(validateQueries(qs[:1]))
}

func TestApplyFunc(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resMap := zebra.NewResourceMap(model.Factory())

	f := func(r zebra.Resource) error {
		return r.Validate(context.Background())
	}

	invalidRes := dc.NewLab("", "", "")
	assert.Nil(resMap.Add(invalidRes))

	assert.NotNil(applyFunc(resMap, f))
}

func createRequest(assert *assert.Assertions, method string, url string,
	body string, api *ResourceAPI,
) *http.Request {
	ctx := context.WithValue(context.Background(), ResourcesCtxKey, api)
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	assert.Nil(err)
	assert.NotNil(req)

	if body != "" {
		req.Body = ioutil.NopCloser(bytes.NewBufferString(body))
	}

	return req
}
