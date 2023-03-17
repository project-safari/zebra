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
	"store": { "rootDir": "test_setup" },
	"authKey": "AvadaKedavra",
	"admin": {
	  "meta": {
		"id": "d770d80d-cd76-4b42-b617-e0f853d60c0a",
		"name": "admin",
		"type": {
		  "name": "system.user",
		  "description": "zebra system user"
		},
		"owner": "admin",
		"creationTime": "2022-09-10T15:52:17.245582145Z",
		"modificationTime": "2022-09-10T15:52:17.245582145Z",
		"labels": {
		  "system.group": "users"
		}
	  },
	  "status": {
		"state": "active"
	  },
	  "key": "-----BEGIN RSA PUBLIC KEY-----\nMIIBCgKCAQEAvv5CZovyjcIezMXcEFvsL5OyOHocMmzkac7Q+aBlEiGcEFEOlbgr\n8nR+8fuLvfjLHVwuYXRD+6D7OQaXRDTCCoPnYdEpwhc7fWUDjnLFfMb7xfn+yIjy\nzkarL2JO075CYMitT5o5Vqbad8EC05FVfFVSPKGaUg0eSpjuVnZyfKbW4Fhnf+LI\nIimd0dVSMHu5CFpLdTLQzXAHzUoSjyqBzhS9gHGiDHqfkVTqd326uJJLBHjNwTSj\nkf0rp8Ib43ldC+pB2BW6gbJbcg7mGjcFXmXY/rkdKA/qCSJkHVIzHMyoW02URz1O\nBDY5//Hr01R7+EsPNN/tpvQ4XBRJh7JlPwIDAQAB\n-----END RSA PUBLIC KEY-----\n",
	  "passwordHash": "$2a$10$BXs49DHWpXwKZlJb6dhriOyc2Ujl2vFj0Giaqt2itjb41qLR5/4mW",
	  "role": {
		"name": "admin",
		"privileges": [
		  ":c,r,u,d"
		]
	  },
	  "email": "admin@zebra.project-safari.io"
	}
  }
`

//nolint:lll
const storeCfgAdapter = `
{
	"store": { "rootDir": "test_setup_adapter" },
	"authKey": "AvadaKedavra",
	"admin": {
	  "meta": {
		"id": "d770d80d-cd76-4b42-b617-e0f853d60c0a",
		"name": "admin",
		"type": {
		  "name": "system.user",
		  "description": "zebra system user"
		},
		"owner": "admin",
		"creationTime": "2022-09-10T15:52:17.245582145Z",
		"modificationTime": "2022-09-10T15:52:17.245582145Z",
		"labels": {
		  "system.group": "users"
		}
	  },
	  "status": {
		"state": "active"
	  },
	  "key": "-----BEGIN RSA PUBLIC KEY-----\nMIIBCgKCAQEAvv5CZovyjcIezMXcEFvsL5OyOHocMmzkac7Q+aBlEiGcEFEOlbgr\n8nR+8fuLvfjLHVwuYXRD+6D7OQaXRDTCCoPnYdEpwhc7fWUDjnLFfMb7xfn+yIjy\nzkarL2JO075CYMitT5o5Vqbad8EC05FVfFVSPKGaUg0eSpjuVnZyfKbW4Fhnf+LI\nIimd0dVSMHu5CFpLdTLQzXAHzUoSjyqBzhS9gHGiDHqfkVTqd326uJJLBHjNwTSj\nkf0rp8Ib43ldC+pB2BW6gbJbcg7mGjcFXmXY/rkdKA/qCSJkHVIzHMyoW02URz1O\nBDY5//Hr01R7+EsPNN/tpvQ4XBRJh7JlPwIDAQAB\n-----END RSA PUBLIC KEY-----\n",
	  "passwordHash": "$2a$10$BXs49DHWpXwKZlJb6dhriOyc2Ujl2vFj0Giaqt2itjb41qLR5/4mW",
	  "role": {
		"name": "admin",
		"privileges": [
		  ":c,r,u,d"
		]
	  },
	  "email": "admin@zebra.project-safari.io"
	}
  }
`

// Test function for the setup.
func TestSetup(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	defer func() { os.RemoveAll("test_setup") }()

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

	assert.NotNil(setupAdapter(ctx, cfgStore))
}

// Test function for the setup adapter.
func TestSetupAdapter(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	defer func() { os.RemoveAll("test_setup_adapter") }()

	cfgStore := config.New()
	assert.Nil(cfgStore.LoadFromStr(context.Background(), storeCfgAdapter))

	ctx := setupLogger(cfgStore)
	a := setupAdapter(ctx, cfgStore)
	assert.NotNil(a)

	assert.Panics(func() {
		handler := a(nil)
		handler.ServeHTTP(nil, nil)
	})

	testForward(assert, a)
}
