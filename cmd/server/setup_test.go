package main //nolint:testpackage

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gojini.dev/config"
)

const storeCfg = `
{
	"store": {"rootDir": "test_setup"},
	"authKey": "abracadabra"
}
`

const storeCfgAdapter = `
{
	"store": {"rootDir": "test_setup_adapter"},
	"authKey": "abracadabra"
}
`

func TestSetup(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	cfgStore := config.New()
	ctx := setupLogger(cfgStore)
	assert.NotNil(ctx)

	// No store
	assert.Panics(func() {
		setupAdapter(ctx, cfgStore)
	})

	cfgStore = config.New()
	e := cfgStore.LoadFromStr(ctx, `{"store": {"rootDir": "test_setup"}}`)
	assert.Nil(e)

	// No authKey
	assert.Panics(func() {
		setupAdapter(ctx, cfgStore)
	})

	cfgStore = config.New()
	e = cfgStore.LoadFromStr(ctx, storeCfg)
	assert.Nil(e)

	defer func() {
		assert.Nil(os.RemoveAll("test_setup"))
	}()

	assert.NotNil(setupAdapter(ctx, cfgStore))
}

func TestSetupAdapter(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	cfgStore := config.New()
	assert.Nil(cfgStore.LoadFromStr(context.Background(), storeCfgAdapter))

	defer func() {
		assert.Nil(os.RemoveAll("test_setup_adapter"))
	}()

	ctx := setupLogger(cfgStore)
	a := setupAdapter(ctx, cfgStore)
	assert.NotNil(a)

	assert.Panics(func() {
		handler := a(nil)
		handler.ServeHTTP(nil, nil)
	})

	testForward(assert, a)
}
