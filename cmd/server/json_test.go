package main //nolint:testpackage

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

// errFake is a mock error to be used in tests.
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

type fakeWriter struct {
	status int
	header http.Header
	err    bool
}

func (f *fakeWriter) WriteHeader(status int) {
	f.status = status
}

func (f *fakeWriter) Header() http.Header {
	return f.header
}

func (f *fakeWriter) Write(b []byte) (int, error) {
	if f.err {
		return 0, errFake
	}

	return len(b), nil
}

// Test function for reading json.
func TestReadJSON(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	req := makeLabelRequest(assert, nil, "a", "b", "c")

	labelReq := &struct {
		Labels []string `json:"labels"`
	}{Labels: []string{}}

	assert.Nil(readJSON(context.Background(), req, labelReq))

	// Bad IO reader
	req.Body = ioutil.NopCloser(fakeReader{err: true})
	assert.NotNil(readJSON(context.Background(), req, nil))

	// Empty Body
	req.Body = ioutil.NopCloser(fakeReader{err: false})
	assert.NotNil(readJSON(context.Background(), req, nil))
}

// Test function for writing json.
func TestWriteJSON(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	f := &fakeWriter{status: 0, header: http.Header{}, err: false}

	x := map[string]interface{}{
		"foo": make(chan int),
	}

	writeJSON(context.Background(), f, x)
	assert.Equal(http.StatusInternalServerError, f.status)

	f.err = true
	writeJSON(context.Background(), f, 10)
}
