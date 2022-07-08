package api_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/api"
	"github.com/project-safari/zebra/network"
	"github.com/stretchr/testify/assert"
	"gojini.dev/web"
)

//nolint:gochecknoglobals
var (
	resource1 = `{"VLANPool":[{"id":"0100000001","type":"VLANPool","labels":{"owner":"shravya"},` +
		`"rangeStart":0,"rangeEnd":10}]}`
	resource2 = `{"VLANPool":[{"id":"0100000002","type":"VLANPool","labels":{"owner":"nandyala"},` +
		`"rangeStart":1,"rangeEnd":5}]}`
	resources = `{"VLANPool":[{"id":"0100000001","type":"VLANPool","labels":{"owner":"shravya"},` +
		`"rangeStart":0,"rangeEnd":10},{"id":"0100000002","type":"VLANPool","labels":{"owner":"nandyala"},` +
		`"rangeStart":1,"rangeEnd":5}]}`
	otherResources = `{"VLANPool":[{"id":"0100000002","type":"VLANPool","labels":{"owner":"nandyala"},` +
		`"rangeStart":1,"rangeEnd":5},{"id":"0100000001","type":"VLANPool","labels":{"owner":"shravya"},` +
		`"rangeStart":0,"rangeEnd":10}]}`
	noResources = "{}"
)

func TestNew(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	assert.NotNil(api.NewResourceAPI(nil))
}

func TestInitialize(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	f := zebra.Factory().Add("VLANPool", func() zebra.Resource { return new(network.VLANPool) })

	api := api.NewResourceAPI(f)
	assert.Nil(api.Initialize("teststore"))

	m := zebra.Factory().Add(" ", func() zebra.Resource { return new(network.VLANPool) })
	assert.NotEqual(m, t)

}

func TestGetResources(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	f := zebra.Factory().Add("VLANPool", func() zebra.Resource { return new(network.VLANPool) })

	myAPI := api.NewResourceAPI(f)
	assert.Nil(myAPI.Initialize("teststore"))

	cfg := &web.Config{
		Address: web.NewAddress("127.0.0.1:9999"),
		TLS:     nil,
	}

	server := web.NewServer(cfg, http.HandlerFunc(myAPI.GetResources))
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

	f := zebra.Factory().Add("VLANPool", func() zebra.Resource { return new(network.VLANPool) })
	myAPI := api.NewResourceAPI(f)
	assert.Nil(myAPI.Initialize("teststore"))

	cfg := &web.Config{
		Address: web.NewAddress("127.0.0.1:9998"),
		TLS:     nil,
	}

	server := web.NewServer(cfg, http.HandlerFunc(myAPI.GetResourcesByID))
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

	assert.Equal(resource1, string(body))
	assert.Nil(resp.Body.Close())

	// GET resource2
	resp, err = http.Get(fmt.Sprintf("http://%s?id=0100000002", cfg.Address))
	assert.True(err == nil && resp != nil)

	body, err = ioutil.ReadAll(resp.Body)
	assert.Nil(err)

	assert.Equal(resource2, string(body))
	assert.Nil(resp.Body.Close())

	// GET resources
	resp, err = http.Get(fmt.Sprintf("http://%s?id=0100000001,0100000002", cfg.Address))
	assert.True(err == nil && resp != nil)

	body, err = ioutil.ReadAll(resp.Body)
	assert.Nil(err)

	assert.Equal(resources, string(body))
	assert.Nil(resp.Body.Close())

	//GET noResources
	resp, err = http.Get(fmt.Sprintf("http://%s?id=0100000002", cfg.Address))
	assert.True(err == nil && resp != nil)

	body, err = ioutil.ReadAll(resp.Body)
	assert.Nil(err)

	assert.Equal(noResources, "{}")

	assert.NotEqual(noResources, string(body))
	assert.Nil(resp.Body.Close())

	//done
	assert.Nil(server.Stop(ctx, nil))
}

func TestGetResourcesByType(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	f := zebra.Factory().Add("VLANPool", func() zebra.Resource { return new(network.VLANPool) })
	myAPI := api.NewResourceAPI(f)
	assert.Nil(myAPI.Initialize("teststore"))

	cfg := &web.Config{
		Address: web.NewAddress("127.0.0.1:9997"),
		TLS:     nil,
	}

	server := web.NewServer(cfg, http.HandlerFunc(myAPI.GetResourcesByType))
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

	assert.Equal(noResources, string(body))
	assert.Nil(resp.Body.Close())

	assert.Nil(server.Stop(ctx, nil))
}


func TestGetResourcesByProperty(t *testing.T) { // nolint:funlen
	t.Parallel()
	assert := assert.New(t)

	f := zebra.Factory().Add("VLANPool", func() zebra.Resource { return new(network.VLANPool) })
	myAPI := api.NewResourceAPI(f)
	assert.Nil(myAPI.Initialize("teststore"))

	cfg := &web.Config{
		Address: web.NewAddress("127.0.0.1:9996"),
		TLS:     nil,
	}
	server := web.NewServer(cfg, http.HandlerFunc(myAPI.GetResourcesByProperty))

	assert.NotNil(server)

	ctx := context.Background()

	go func() {
		assert.NotNil(server.Start(ctx))
	}()
	time.Sleep(time.Second)

	//otherResources

	resp, err := http.Get(fmt.Sprintf("http://%s?property=Type-in-VLANPool,IPAddressPool", cfg.Address))
	assert.True(err == nil && resp != nil)

	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(err)

	assert.True(string(body) == resources || string(body) == otherResources)
	assert.Nil(resp.Body.Close())

	//noResources

	resp, err = http.Get(fmt.Sprintf("http://%s?property=Type-notin-VLANPool,IPAddressPool", cfg.Address))
	assert.True(err == nil && resp != nil)

	body, err = ioutil.ReadAll(resp.Body)
	assert.Nil(err)

	assert.Equal(noResources, string(body))
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

	assert.Equal(noResources, string(body))
	assert.Nil(resp.Body.Close())

	resp, err = http.Get(fmt.Sprintf("http://%s?property=Type-blahblah-test", cfg.Address))
	assert.True(err == nil && resp != nil)
	assert.True(resp.StatusCode == 400)
	assert.Nil(server.Stop(ctx, nil))
}

func TestGetResourcesByLabel(t *testing.T) { // nolint:funlen
	t.Parallel()
	assert := assert.New(t)

	f := zebra.Factory().Add("VLANPool", func() zebra.Resource { return new(network.VLANPool) })
	myAPI := api.NewResourceAPI(f)
	assert.Nil(myAPI.Initialize("teststore"))

	cfg := &web.Config{
		Address: web.NewAddress("127.0.0.1:9995"),
		TLS:     nil,
	}
	server := web.NewServer(cfg, http.HandlerFunc(myAPI.GetResourcesByLabel))
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

	assert.Equal(noResources, string(body))
	assert.Nil(resp.Body.Close())

	resp, err = http.Get(fmt.Sprintf("http://%s?label=owner-equal-shravya", cfg.Address))
	assert.True(err == nil && resp != nil)

	body, err = ioutil.ReadAll(resp.Body)
	assert.Nil(err)

	assert.Equal(resource1, string(body))
	assert.Nil(resp.Body.Close())

	resp, err = http.Get(fmt.Sprintf("http://%s?label=owner-notequal-shravya", cfg.Address))
	assert.True(err == nil && resp != nil)

	body, err = ioutil.ReadAll(resp.Body)
	assert.Nil(err)

	assert.Equal(resource2, string(body))
	assert.Nil(resp.Body.Close())

	resp, err = http.Get(fmt.Sprintf("http://%s?label=owner-notequal", cfg.Address))
	assert.True(err == nil && resp != nil)
	assert.True(resp.StatusCode == 400)
	assert.Nil(server.Stop(ctx, nil))
}
