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
	"store": {"rootDir": "test"},
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
	e := cfgStore.LoadFromStr(ctx, `{"store": {"rootDir": "test"}}`)
	assert.Nil(e)

	// No authKey
	assert.Panics(func() {
		setupAdapter(ctx, cfgStore)
	})

	cfgStore = config.New()
	e = cfgStore.LoadFromStr(ctx, storeCfg)
	assert.Nil(e)

	defer func() {
		assert.Nil(os.RemoveAll("test"))
	}()

	assert.NotNil(setupAdapter(ctx, cfgStore))
}

func TestSetupAdapter(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	cfgStore := config.New()
	assert.Nil(cfgStore.LoadFromStr(context.Background(), storeCfg))

	defer func() {
		assert.Nil(os.RemoveAll("test"))
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
