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
	"github.com/project-safari/zebra/dc"
	"github.com/project-safari/zebra/store"
	"github.com/stretchr/testify/assert"
)

func makeLabelRequest(assert *assert.Assertions, labels ...string) *http.Request {
	req, err := http.NewRequest("GET", "/api/v1/labels", nil)
	assert.Nil(err)
	assert.NotNil(req)

	v := map[string][]string{"labels": labels}
	b, e := json.Marshal(v)
	assert.Nil(e)

	req.Body = ioutil.NopCloser(bytes.NewBuffer(b))

	return req
}

func TestBadLabelReq(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	h := handleLabels(context.Background(), nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r, nil)
	})

	req, err := http.NewRequest("GET", "/api/v1/labels", nil)
	assert.Nil(err)
	assert.NotNil(req)

	v := "{....}" // bad json
	req.Body = ioutil.NopCloser(bytes.NewBufferString(v))

	handler.ServeHTTP(rr, req)
	assert.Equal(rr.Code, http.StatusBadRequest)
}

func TestAllLabels(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	store := makeStore(assert, "test_all_labels")
	assert.NotNil(store)

	defer func() {
		assert.Nil(os.RemoveAll("./test_all_labels"))
	}()

	h := handleLabels(context.Background(), store)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r, nil)
	})

	req := makeLabelRequest(assert)
	handler.ServeHTTP(rr, req)

	assert.Equal(rr.Code, http.StatusOK)
}

func TestLabels(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	store := makeStore(assert, "test_labels")
	assert.NotNil(store)

	defer func() {
		assert.Nil(os.RemoveAll("./test_labels"))
	}()

	h := handleLabels(context.Background(), store)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h(w, r, nil)
	})

	req := makeLabelRequest(assert, "label5", "label7")
	handler.ServeHTTP(rr, req)

	assert.Equal(rr.Code, http.StatusOK)
}

func makeStore(assert *assert.Assertions, root string) zebra.Store {
	s := store.NewResourceStore(root, store.DefaultFactory())
	assert.Nil(s.Initialize())

	labels := zebra.Labels{}

	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("label%d", i)
		value := fmt.Sprintf("value%d", i)
		labels[key] = value
	}

	// create 100 labs
	for i := 0; i < 100; i++ {
		l := new(dc.Lab)
		l.Name = fmt.Sprintf("lab-%d", i+1)
		l.Type = "Lab"
		l.BaseResource = *zebra.NewBaseResource("Lab", labels)

		assert.Nil(s.Create(l))
	}

	return s
}
