package main //nolint:testpackage

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/model"
	"github.com/project-safari/zebra/model/dc"
	"github.com/project-safari/zebra/store"
	"github.com/stretchr/testify/assert"
)

// Mock function that creates a store with 100 labs to be used for tests.
func makeStore(assert *assert.Assertions, root string) zebra.Store {
	f := model.Factory()
	s := store.NewResourceStore(root, f)
	assert.Nil(s.Initialize())

	// create 100 labs
	for i := 0; i < 100; i++ {
		l, ok := f.New("dc.lab").(*dc.Lab)
		assert.True(ok)

		n := fmt.Sprintf("lab-%d", i+1)
		l.BaseResource = *zebra.NewBaseResource(l.Meta.Type, n, "test_owner", "test_group")

		for i := 0; i < 10; i++ {
			key := fmt.Sprintf("label%d", i)
			value := fmt.Sprintf("value%d", i)
			l.Meta.Labels.Add(key, value)
		}

		assert.Nil(s.Create(l))
	}

	return s
}

// Mock function that creates a request for labels, to be used for tests.
func makeLabelRequest(assert *assert.Assertions, resources *ResourceAPI, labels ...string) *http.Request {
	ctx := context.WithValue(context.Background(), ResourcesCtxKey, resources)
	ctx = context.WithValue(ctx, AuthCtxKey, authKey)

	req, err := http.NewRequestWithContext(ctx, "GET", "/api/v1/labels", nil)
	assert.Nil(err)
	assert.NotNil(req)

	v := map[string][]string{"labels": labels}
	b, e := json.Marshal(v)
	assert.Nil(e)

	req.Body = ioutil.NopCloser(bytes.NewBuffer(b))

	return req
}

// Test for a bad request for labels.
func TestBadLabelReq(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	h := handleLabels()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r, nil)
	})

	ctx := context.WithValue(context.Background(), ResourcesCtxKey, NewResourceAPI(model.Factory()))
	req, err := http.NewRequestWithContext(ctx, "GET", "/api/v1/labels", nil)
	assert.Nil(err)
	assert.NotNil(req)

	v := "{....}" // bad json
	req.Body = ioutil.NopCloser(bytes.NewBufferString(v))

	handler.ServeHTTP(rr, req)
	assert.Equal(rr.Code, http.StatusBadRequest)

	// Bad context
	req, err = http.NewRequest("GET", "/api/v1/labels", nil)
	assert.Nil(err)
	assert.NotNil(req)

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(http.StatusInternalServerError, rr.Code)
}

// Function to test a request for all labels.
func TestAllLabels(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	defer func() { os.RemoveAll("test_all_labels") }()

	resources := NewResourceAPI(model.Factory())
	resources.Store = makeStore(assert, "test_all_labels")
	assert.NotNil(resources.Store)

	h := handleLabels()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r, nil)
	})

	req := makeLabelRequest(assert, resources)
	handler.ServeHTTP(rr, req)

	assert.Equal(rr.Code, http.StatusOK)
}

// Function to test labels and label requests.
func TestLabels(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	defer func() { os.RemoveAll("test_labels") }()

	resources := NewResourceAPI(model.Factory())
	resources.Store = makeStore(assert, "test_labels")
	assert.NotNil(resources.Store)

	h := handleLabels()
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r, nil)
	})

	req := makeLabelRequest(assert, resources, "label5", "label7")
	handler.ServeHTTP(rr, req)

	assert.Equal(rr.Code, http.StatusOK)
}
