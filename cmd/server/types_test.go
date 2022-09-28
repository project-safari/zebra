package main //nolint:testpackage

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Mock function that creates a type request to be used in tests.
func makeTypeRequest(assert *assert.Assertions, types ...string) *http.Request {
	req, err := http.NewRequest("GET", "/api/v1/types", nil)
	assert.Nil(err)
	assert.NotNil(req)

	v := map[string][]string{"types": types}
	b, e := json.Marshal(v)
	assert.Nil(e)

	req.Body = ioutil.NopCloser(bytes.NewBuffer(b))

	return req
}

// Test function for a bad request.
func TestBadReq(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	h := handleTypes()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r, nil)
	})

	req, err := http.NewRequest("GET", "/api/v1/types", nil)
	assert.Nil(err)
	assert.NotNil(req)

	v := "{....}" // bad json
	req.Body = ioutil.NopCloser(bytes.NewBufferString(v))

	handler.ServeHTTP(rr, req)
	assert.Equal(rr.Code, http.StatusBadRequest)
}

func TestDefaultFactory(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	h := handleTypes()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r, nil)
	})

	req := makeTypeRequest(assert)
	handler.ServeHTTP(rr, req)

	assert.Equal(rr.Code, http.StatusOK)
}

// Test function for type requests.
func TestTypes(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	h := handleTypes()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r, nil)
	})

	req := makeTypeRequest(assert, "compute.server", "network.switch")
	handler.ServeHTTP(rr, req)

	assert.Equal(rr.Code, http.StatusOK)
}
