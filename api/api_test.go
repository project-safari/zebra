package api_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/rchamarthy/zebra/api"
	"github.com/stretchr/testify/assert"
	"gojini.dev/web"
)

//nolint:gochecknoglobals
var (
	handlerAll      = http.HandlerFunc(api.GetResources)
	handlerID       = http.HandlerFunc(api.GetResourcesByID)
	handlerType     = http.HandlerFunc(api.GetResourcesByType)
	handlerProperty = http.HandlerFunc(api.GetResourcesByProperty)
	handlerLabel    = http.HandlerFunc(api.GetResourcesByLabel)
)

//nolint:gochecknoglobals
var (
	resource1 = `{"VLANPool":[{"id":"0100000001","type":"VLANPool","labels":{"owner":"shravya"},` +
		`"rangeStart":0,"rangeEnd":10}]}` + "\n"
	resource2 = `{"VLANPool":[{"id":"0100000002","type":"VLANPool","labels":{"owner":"nandyala"},` +
		`"rangeStart":1,"rangeEnd":5}]}` + "\n"
	resources = `{"VLANPool":[{"id":"0100000001","type":"VLANPool","labels":{"owner":"shravya"},` +
		`"rangeStart":0,"rangeEnd":10},{"id":"0100000002","type":"VLANPool","labels":{"owner":"nandyala"},` +
		`"rangeStart":1,"rangeEnd":5}]}` + "\n"
	otherResources = `{"VLANPool":[{"id":"0100000002","type":"VLANPool","labels":{"owner":"nandyala"},` +
		`"rangeStart":1,"rangeEnd":5},{"id":"0100000001","type":"VLANPool","labels":{"owner":"shravya"},` +
		`"rangeStart":0,"rangeEnd":10}]}` + "\n"
	noResources = "{}\n"
)

func TestSetUp(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	assert.Nil(api.SetUp("teststore"))
}

func TestGetResources(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	assert.Nil(api.SetUp("teststore"))

	cfg := &web.Config{
		Address: web.NewAddress("127.0.0.1:9999"),
		TLS:     nil,
	}

	server := web.NewServer(cfg, handlerAll)
	assert.NotNil(server)

	ctx := context.Background()

	go func() {
		assert.NotNil(server.Start(ctx))
	}()
	time.Sleep(time.Second)

	resp, err := http.Get(fmt.Sprintf("http://%s/", cfg.Address))
	assert.True(err == nil && resp != nil)

	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(err)

	assert.True(string(body) == resources || string(body) == otherResources)
	assert.Nil(resp.Body.Close())

	assert.Nil(server.Stop(ctx, nil))
}

func TestGetResourcesByID(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	assert.Nil(api.SetUp("teststore"))

	cfg := &web.Config{
		Address: web.NewAddress("127.0.0.1:9998"),
		TLS:     nil,
	}

	server := web.NewServer(cfg, handlerID)
	assert.NotNil(server)

	ctx := context.Background()

	go func() {
		assert.NotNil(server.Start(ctx))
	}()
	time.Sleep(time.Second)

	// GET resource1
	resp, err := http.Get(fmt.Sprintf("http://%s?id=0100000001", cfg.Address))
	assert.True(err == nil && resp != nil)

	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(err)

	assert.Equal(string(body), resource1)
	assert.Nil(resp.Body.Close())

	// GET resource2
	resp, err = http.Get(fmt.Sprintf("http://%s?id=0100000002", cfg.Address))
	assert.True(err == nil && resp != nil)

	body, err = ioutil.ReadAll(resp.Body)
	assert.Nil(err)

	assert.Equal(string(body), resource2)
	assert.Nil(resp.Body.Close())

	// GET resources
	resp, err = http.Get(fmt.Sprintf("http://%s?id=0100000001,0100000002", cfg.Address))
	assert.True(err == nil && resp != nil)

	body, err = ioutil.ReadAll(resp.Body)
	assert.Nil(err)

	assert.Equal(string(body), resources)
	assert.Nil(resp.Body.Close())

	assert.Nil(server.Stop(ctx, nil))
}

func TestGetResourcesByType(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	assert.Nil(api.SetUp("teststore"))

	cfg := &web.Config{
		Address: web.NewAddress("127.0.0.1:9997"),
		TLS:     nil,
	}

	server := web.NewServer(cfg, handlerType)
	assert.NotNil(server)

	ctx := context.Background()

	go func() {
		assert.NotNil(server.Start(ctx))
	}()
	time.Sleep(time.Second)

	// GET resources
	resp, err := http.Get(fmt.Sprintf("http://%s?type=VLANPool", cfg.Address))
	assert.True(err == nil && resp != nil)

	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(err)

	assert.True(string(body) == resources || string(body) == otherResources)
	assert.Nil(resp.Body.Close())

	// GET no resources
	resp, err = http.Get(fmt.Sprintf("http://%s?type=IPAddressPool", cfg.Address))
	assert.True(err == nil && resp != nil)

	body, err = ioutil.ReadAll(resp.Body)
	assert.Nil(err)

	assert.Equal(string(body), noResources)
	assert.Nil(resp.Body.Close())

	assert.Nil(server.Stop(ctx, nil))
}

func TestGetResourcesByProperty(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	assert.Nil(api.SetUp("teststore"))

	cfg := &web.Config{
		Address: web.NewAddress("127.0.0.1:9996"),
		TLS:     nil,
	}
	server := web.NewServer(cfg, handlerProperty)

	assert.NotNil(server)

	ctx := context.Background()

	go func() {
		assert.NotNil(server.Start(ctx))
	}()
	time.Sleep(time.Second)

	resp, err := http.Get(fmt.Sprintf("http://%s?property=Type-in-VLANPool,IPAddressPool", cfg.Address))
	assert.True(err == nil && resp != nil)

	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(err)

	assert.True(string(body) == resources || string(body) == otherResources)
	assert.Nil(resp.Body.Close())

	resp, err = http.Get(fmt.Sprintf("http://%s?property=Type-notin-VLANPool,IPAddressPool", cfg.Address))
	assert.True(err == nil && resp != nil)

	body, err = ioutil.ReadAll(resp.Body)
	assert.Nil(err)

	assert.True(string(body) == noResources)
	assert.Nil(resp.Body.Close())

	resp, err = http.Get(fmt.Sprintf("http://%s?property=Type-equal-VLANPool", cfg.Address))
	assert.True(err == nil && resp != nil)

	body, err = ioutil.ReadAll(resp.Body)
	assert.Nil(err)

	assert.True(string(body) == resources || string(body) == otherResources)
	assert.Nil(resp.Body.Close())

	resp, err = http.Get(fmt.Sprintf("http://%s?property=Type-notequal-VLANPool", cfg.Address))
	assert.True(err == nil && resp != nil)

	body, err = ioutil.ReadAll(resp.Body)
	assert.Nil(err)

	assert.True(string(body) == noResources)
	assert.Nil(resp.Body.Close())

	resp, err = http.Get(fmt.Sprintf("http://%s?property=Type-notequal", cfg.Address))
	assert.True(err == nil && resp != nil)
	assert.True(resp.StatusCode == 400)
	assert.Nil(server.Stop(ctx, nil))
}

func TestGetResourcesByLabel(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	assert.Nil(api.SetUp("teststore"))

	cfg := &web.Config{
		Address: web.NewAddress("127.0.0.1:9995"),
		TLS:     nil,
	}
	server := web.NewServer(cfg, handlerLabel)
	assert.NotNil(server)

	ctx := context.Background()

	go func() {
		assert.NotNil(server.Start(ctx))
	}()
	time.Sleep(time.Second)

	resp, err := http.Get(fmt.Sprintf("http://%s?label=owner-in-shravya,nandyala", cfg.Address))
	assert.Nil(err)
	assert.NotNil(resp)

	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(err)

	assert.True(string(body) == resources || string(body) == otherResources)
	assert.Nil(resp.Body.Close())

	resp, err = http.Get(fmt.Sprintf("http://%s?label=owner-notin-shravya,nandyala", cfg.Address))
	assert.True(err == nil && resp != nil)

	body, err = ioutil.ReadAll(resp.Body)
	assert.Nil(err)

	assert.True(string(body) == noResources)
	assert.Nil(resp.Body.Close())

	resp, err = http.Get(fmt.Sprintf("http://%s?label=owner-equal-shravya", cfg.Address))
	assert.True(err == nil && resp != nil)

	body, err = ioutil.ReadAll(resp.Body)
	assert.Nil(err)

	assert.True(string(body) == resource1)
	assert.Nil(resp.Body.Close())

	resp, err = http.Get(fmt.Sprintf("http://%s?label=owner-notequal-shravya", cfg.Address))
	assert.True(err == nil && resp != nil)

	body, err = ioutil.ReadAll(resp.Body)
	assert.Nil(err)

	assert.True(string(body) == resource2)
	assert.Nil(resp.Body.Close())

	resp, err = http.Get(fmt.Sprintf("http://%s?label=owner-notequal", cfg.Address))
	assert.True(err == nil && resp != nil)
	assert.True(resp.StatusCode == 400)
	assert.Nil(server.Stop(ctx, nil))
}
