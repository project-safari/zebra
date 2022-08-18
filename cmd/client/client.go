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

// Potential errors with the configuration.
var (
	ErrNoCACert     = errors.New("zebra CA certificate file is not conifugred")
	ErrNoConfig     = errors.New("zebra config file is not specified")
	ErrNoEmail      = errors.New("user email is not configured")
	ErrNoPrivateKey = errors.New("user private key is not configured")
)

// A struct for client data, it contains the configuration, the http client, and the http header.
type Client struct {
	cfg *Config
	c   *http.Client
	h   http.Header
}

// A function for a new client.
// Takes in the config.
// Returns a pointer to the client and (a) potential error(s).
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

// Operation function on the client: do.
// Uses the do function with "GET".
// Returns int and (a) potential error(s).
func (c *Client) Get(path string, in, out interface{}) (int, error) {
	return c.do(context.Background(), "GET", path, in, out)
}

// Operation function on the client: delete.
// Uses the do function with "DELETE".
// Returns int and (a) potential error(s).
func (c *Client) Delete(path string, in, out interface{}) (int, error) {
	return c.do(context.Background(), "DELETE", path, in, out)
}

// Operation function on the client: post.
// Uses the do function with "POST".
// Returns int and (a) potential error(s).
func (c *Client) Post(path string, in, out interface{}) (int, error) {
	return c.do(context.Background(), "POST", path, in, out)
}

// Function on the client: do.
// Implemented in the client's operation functions.
// Uses a server addr, a new buffer, a request with a given method, and a response body.
// Its purpose is to get the responce's status.
// Returns int and (a) potential error(s).
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

// Function for a tls client.
// Takes in a config.
// Returns a pointer to an http client and (a) potential error(s).
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
