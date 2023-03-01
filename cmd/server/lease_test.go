package main //nolint:testpackage

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/model"
	"github.com/project-safari/zebra/model/dc"
	"github.com/project-safari/zebra/model/lease"
	"github.com/stretchr/testify/assert"
)

func makeBody(duration string, resource zebra.Resource) io.ReadCloser {
	d, err := time.ParseDuration(duration)
	if err != nil {
		return nil
	}

	ids := []string{resource.GetMeta().ID}
	resources := []*lease.ResourceReq{
		{
			Type:      "Server",
			Group:     "san-jose-building-14",
			Name:      "linux blah blah",
			Count:     2,
			Resources: ids,
		},
		{
			Type:  "VM",
			Group: "san-jose-building-18",
			Name:  "virtual",
			Count: 1,
		},
	}

	leaseReq := &struct {
		Email    string               `json:"email"`
		Duration time.Duration        `json:"duration"`
		Request  []*lease.ResourceReq `json:"request"`
	}{
		Email:    "testuser@cisco.com",
		Duration: d,
		Request:  resources,
	}

	v, err := json.Marshal(leaseReq)
	if err != nil {
		panic(err)
	}

	return ioutil.NopCloser(bytes.NewBuffer(v))
}

func makeRequest(assert *assert.Assertions, method string, url string,
	body string, api *ResourceAPI,
) *http.Request {
	ctx := context.WithValue(context.Background(), ResourcesCtxKey, api)
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	assert.Nil(err)
	assert.NotNil(req)

	lab := dc.NewLab("Dexter's Labotory", "dexlab@cisoc.com", "blue")
	req.Body = makeBody(body, lab)

	return req
}

func TestLease(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "test_lease"

	defer func() { os.RemoveAll(root) }()

	api := NewResourceAPI(model.Factory())
	assert.Nil(api.Initialize(root))

	h := handleLease()
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r, nil)
	})

	req := makeRequest(assert, "POST", "/lease", "4h", api)

	assert.NotNil(req)

	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusOK, rr.Code)
}

func TestBadLease(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "test_bad_lease"

	defer func() { os.RemoveAll(root) }()

	api := NewResourceAPI(model.Factory())
	assert.Nil(api.Initialize(root))

	h := handleLease()
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r, nil)
	})
	lab := dc.NewLab("Dexter's Labotory", "dexlab@cisoc.com", "blue")
	req, err := http.NewRequest("POST", "/lease", makeBody("4h", lab))
	assert.Nil(err)
	assert.NotNil(req)

	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusInternalServerError, rr.Code)

	req = makeRequest(assert, "POST", "/lease", "4h", api)
	assert.Nil(err)
	assert.NotNil(req)

	v := "{...}"
	req.Body = ioutil.NopCloser(bytes.NewBufferString(v))
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusBadRequest, rr.Code)

	req.Body = makeBody("7h", lab)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusBadRequest, rr.Code)
}
