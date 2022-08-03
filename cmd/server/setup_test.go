package main //nolint:testpackage

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"gojini.dev/config"
)

//nolint:lll
const storeCfg = `
{
	"store": {"rootDir": "test_setup"},
	"authKey": "abracadabra",
	"admin": {
		"id": "d0ff79eb-e820-469d-b924-dbb78992727e",
		"type": "User",
		"labels": {
		  "system.group": "users"
		},
		"name": "ravi",
		"key": "-----BEGIN RSA PUBLIC KEY-----\nMIIBCgKCAQEA4u0E+fwbAzEnN5e2meHWaadtSdFT/lfcyxAcjIp0PZ8InE7HBYMI\nwA9lSlXpfz25gGIdV4Fae+24FpQtPw3Fo8S0hVa0a4Hlz3HO02IxjPRt2aMmyLKm\nZGisIRlD5R5KufQlHp7ZD3opUCJElcm50F6utjlBq0N6+t8YZ20HVGdjwf61EnV4\nti5Q6E33ENMjEm/hcg0JuIUY1sLYDLPVaxAG+ZIU5wYGlhdIUgzFRk88ijApepDq\nxyXGvD+lsyduwSGNB7+AAnmyyWEXy1PxdXqlaKJh3vj7o2OlV/SrkwIN1UlgMtvB\nEsYw62spavvXEvv5DB7LDu2ORAPqbJ4QGQIDAQAB\n-----END RSA PUBLIC KEY-----\n",
		"passwordHash": "$2a$10$EuRgZ5GrkZ0zhavaLsGKhuTeSClCEiAx80h3FqOGiBnV6Vvv8av4y",
		"role": {
		  "name": "admin",
		  "privileges": [
			":c,r,u,d"
		  ]
		},
		"email": "ravchama@cisco.com"
	}
}
`

//nolint:lll
const storeCfgAdapter = `
{
	"store": {"rootDir": "test_setup_adapter"},
	"authKey": "abracadabra",
	"admin": {
		"id": "d0ff79eb-e820-469d-b924-dbb78992727e",
		"type": "User",
		"labels": {
		  "system.group": "users"
		},
		"name": "ravi",
		"key": "-----BEGIN RSA PUBLIC KEY-----\nMIIBCgKCAQEA4u0E+fwbAzEnN5e2meHWaadtSdFT/lfcyxAcjIp0PZ8InE7HBYMI\nwA9lSlXpfz25gGIdV4Fae+24FpQtPw3Fo8S0hVa0a4Hlz3HO02IxjPRt2aMmyLKm\nZGisIRlD5R5KufQlHp7ZD3opUCJElcm50F6utjlBq0N6+t8YZ20HVGdjwf61EnV4\nti5Q6E33ENMjEm/hcg0JuIUY1sLYDLPVaxAG+ZIU5wYGlhdIUgzFRk88ijApepDq\nxyXGvD+lsyduwSGNB7+AAnmyyWEXy1PxdXqlaKJh3vj7o2OlV/SrkwIN1UlgMtvB\nEsYw62spavvXEvv5DB7LDu2ORAPqbJ4QGQIDAQAB\n-----END RSA PUBLIC KEY-----\n",
		"passwordHash": "$2a$10$EuRgZ5GrkZ0zhavaLsGKhuTeSClCEiAx80h3FqOGiBnV6Vvv8av4y",
		"role": {
		  "name": "admin",
		  "privileges": [
			":c,r,u,d"
		  ]
		},
		"email": "ravchama@cisco.com"
	}
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
