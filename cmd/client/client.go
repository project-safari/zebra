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

var (
	ErrNoCACert     = errors.New("zebra CA certificate file is not conifugred")
	ErrNoConfig     = errors.New("zebra config file is not specified")
	ErrNoEmail      = errors.New("user email is not configured")
	ErrNoPrivateKey = errors.New("user private key is not configured")
)

// var timeLimitedClient = http.Client{Timeout: 18000 * time.Second}

type Client struct {
	cfg *Config
	c   *http.Client
	h   http.Header
}

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

func (c *Client) Get(path string, in, out interface{}) (int, error) {
	c.c.Timeout = time.Duration(1) * time.Minute
	// add 1 min. timeout for get request.
	return c.do(context.Background(), "GET", path, in, out)
}

func (c *Client) Delete(path string, in, out interface{}) (int, error) {
	c.c.Timeout = time.Duration(1) * time.Minute
	// add 1 min. timeout for delete request.
	return c.do(context.Background(), "DELETE", path, in, out)
}

func (c *Client) Post(path string, in, out interface{}) (int, error) {
	c.c.Timeout = time.Duration(1) * time.Minute
	// add 1 min. timeout for post request.
	return c.do(context.Background(), "POST", path, in, out)
}

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

	transport.ForceAttemptHTTP2 = true

	client := new(http.Client)
	client.Timeout = time.Duration(1) * time.Minute
	client.Transport = transport

	return client, nil
}
