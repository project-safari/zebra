package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

// possible errors wit user credentials.
var (
	ErrNoCACert     = errors.New("zebra CA certificate file is not conifugred")
	ErrNoConfig     = errors.New("zebra config file is not specified")
	ErrNoEmail      = errors.New("user email is not configured")
	ErrNoPrivateKey = errors.New("user private key is not configured")
)

// Struct for client information,
// contans pointers to Config and http.Client and an http.Header.
type Client struct {
	cfg *Config
	c   *http.Client
	h   http.Header
}

// Function that initiates a new client.
//
// Takes in a pointer to the Config struct,
// returns a pointer ro Client and a potential error or nil if there is no error.
func NewClient(cfg *Config) (*Client, error) {
	if cfg == nil {
		return nil, ErrNoConfig
	}

	if cfg.Email == "" {
		return nil, ErrNoEmail
	}

	if cfg.Key == nil {
		return nil, ErrNoPrivateKey
	}

	t, err := cfg.Key.Sign([]byte(cfg.Email))
	token := base64.StdEncoding.EncodeToString(t)

	if err != nil {
		return nil, err
	}

	h := http.Header{}
	h.Add("Zebra-Auth-User", cfg.Email)
	h.Add("Zebra-Auth-Token", token)
	h.Add("User-Agent", "zebra-client")
	h.Add("Accept-Encoding", "application/json")
	h.Add("Content-Type", "application/json")

	c, err := tlsClient(cfg)
	if err != nil {
		return nil, err
	}

	return &Client{
		cfg: cfg,
		c:   c,
		h:   h,
	}, nil
}

// Operation function on Client, given a pointer to Client - get.
//
// Function takes in a string for path, and two interfaces: in and out.
// It executes the do operation function on the given pointer,
// and returns the resp.StatusCode as an int and an error, or nil if there is no error.
//
// Uses "GET".
func (c *Client) Get(path string, in, out interface{}) (int, error) {
	return c.do(context.Background(), "GET", path, in, out)
}

// Operation function on Client, given a pointer to Client - delete.
//
// Function takes in a string for path, and two interfaces: in and out.
// It executes the do operation function on the given pointer,
// and returns the resp.StatusCode as an int and an error, or nil if there is no error.
//
// Uses "DELETE".
func (c *Client) Delete(path string, in, out interface{}) (int, error) {
	return c.do(context.Background(), "DELETE", path, in, out)
}

// Operation function on Client, given a pointer to Client - post.
//
// Function takes in a string for path, and two interfaces: in and out.
// It executes the do operation function on the given pointer,
// and returns the resp.StatusCode as an int and an error, or nil if there is no error.
//
// Uses "POST".
func (c *Client) Post(path string, in, out interface{}) (int, error) {
	return c.do(context.Background(), "POST", path, in, out)
}

// Operation function on Client, given a pointer to Client - do.
//
// Function takes in a string for path, and two interfaces: in and out.
// It returns the resp.StatusCode as an int and an error, or nil if there is no error.
func (c *Client) do(ctx context.Context, method, path string, in, out interface{}) (int, error) {
	url := fmt.Sprintf("%s/%s", c.cfg.ServerAddress, path)
	buf := bytes.NewBuffer([]byte{})

	if in != nil {
		b, e := json.Marshal(in)
		if e != nil {
			return 0, e
		}

		buf = bytes.NewBuffer(b)
	}

	r, err := http.NewRequestWithContext(ctx, method, url, buf)
	if err != nil {
		return 0, err
	}

	r.Header = c.h

	resp, err := c.c.Do(r)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != http.StatusOK {
		return resp.StatusCode, fmt.Errorf("%s: %s", url, resp.Status) //nolint:goerr113
	}

	b, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return resp.StatusCode, e
	}

	if out != nil {
		if e := json.Unmarshal(b, out); e != nil {
			return resp.StatusCode, e
		}
	}

	return resp.StatusCode, nil
}

// Function to initiate a tsl client, given a pointer to the Config struct.
// It returns a pointer to the http.Client and an error, or nil if there is no error.
func tlsClient(cfg *Config) (*http.Client, error) {
	if cfg == nil {
		return nil, ErrNoConfig
	}

	if cfg.CACert == "" {
		return nil, ErrNoCACert
	}

	caCert, err := ioutil.ReadFile(cfg.CACert)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	transport := new(http.Transport)
	transport.TLSClientConfig = new(tls.Config)
	transport.TLSClientConfig.RootCAs = caCertPool
	transport.TLSClientConfig.MinVersion = tls.VersionTLS13

	client := new(http.Client)
	client.Timeout = time.Duration(1) * time.Minute
	client.Transport = transport

	return client, nil
}
