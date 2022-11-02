package main //nolint:testpackage

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/auth"
	"github.com/stretchr/testify/assert"
)

// Constant values for paths to files that contain user credential info.
const (
	testCACertFile  = "../../simulator/zebra-ca.crt"
	testUserKeyFile = "../../simulator/user.key"
)

// Test function for a new client.
func TestNewClient(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	c, err := NewClient(nil)
	assert.Nil(c)
	assert.Equal(ErrNoConfig, err)

	cfg := new(Config)
	c, err = NewClient(cfg)
	assert.Nil(c)
	assert.Equal(ErrNoEmail, err)

	cfg.Email = "test@zebra.project-safafi.io"
	c, err = NewClient(cfg)
	assert.Nil(c)
	assert.Equal(ErrNoPrivateKey, err)

	key, err := auth.Load(testUserKeyFile)
	assert.Nil(err)
	assert.NotNil(key)

	cfg.Key = key
	c, err = NewClient(cfg)
	assert.Equal(ErrNoCACert, err)
	assert.Nil(c)

	cfg.Key = key.Public()
	c, err = NewClient(cfg)
	assert.Nil(c)
	assert.Equal(auth.ErrNoPrivateKey, err)

	cfg.CACert = testCACertFile
	cfg.Key = key
	c, err = NewClient(cfg)
	assert.Nil(err)
	assert.NotNil(c)

	_, e := c.Get("/blah", cfg, cfg)
	assert.NotNil(e)

	_, e = c.Delete("/blah", cfg, cfg)
	assert.NotNil(e)

	_, e = c.Post("/blah", cfg, cfg)
	assert.NotNil(e)
}

// Test function fot the TSL client.
func TestTLSClient(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	c, e := tlsClient(nil)
	assert.Nil(c)
	assert.Equal(ErrNoConfig, e)

	cfg := new(Config)

	c, e = tlsClient(cfg)
	assert.Nil(c)
	assert.Equal(ErrNoCACert, e)

	cfg.CACert = "random_file_doesnt_exist"
	c, e = tlsClient(cfg)
	assert.Nil(c)
	assert.NotNil(e)

	cfg.CACert = testCACertFile
	c, e = tlsClient(cfg)
	assert.NotNil(c)
	assert.Nil(e)
}

// Test function for client.
func TestClientDo(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	server := makeServer(assert)
	assert.NotNil(server)

	defer server.Close()

	key, err := auth.Load(testUserKeyFile)
	assert.Nil(err)
	assert.NotNil(key)

	cfg := &Config{
		ServerAddress: server.URL,
		Key:           key,
		User:          "loki",
		Email:         "loki@asgard.io",
		CACert:        testCACertFile,
		Defaults:      ConfigDefaults{Duration: zebra.DefaultMaxDuration},
	}

	client, err := NewClient(cfg)
	assert.Nil(err)
	assert.NotNil(client)

	code, err := client.do(context.Background(), "GET", "/test", make(chan int), nil)
	assert.Equal(0, code)
	assert.NotNil(err)

	code, err = client.do(nil, "GET", "/test", nil, nil) //nolint:staticcheck
	assert.Equal(0, code)
	assert.NotNil(err)

	code, err = client.do(context.Background(), "GET", "test", nil, nil)
	assert.Nil(err)
	assert.Equal(http.StatusOK, code)

	code, err = client.do(context.Background(), "GET", "bad_test", nil, nil)
	assert.NotNil(err)
	assert.Equal(http.StatusNotFound, code)

	code, err = client.do(context.Background(), "GET", "test", nil, &struct {
		A int `json:"a"`
	}{})
	assert.NotNil(err)
	assert.Equal(http.StatusOK, code)
}

// Mock function to create/start a server to be used in tests.
// Returns a pointer to a httptest.
func makeServer(assert *assert.Assertions) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		fmt.Println(req.URL.Path)
		if req.URL.Path == "/test" {
			rw.WriteHeader(http.StatusOK)
		} else {
			rw.WriteHeader(http.StatusNotFound)

			return
		}

		_, e := rw.Write([]byte(`OK`))
		assert.Nil(e)
	}))
}
