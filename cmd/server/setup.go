package main

import (
	"context"
	"net/http"
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"
	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/model"
	"github.com/project-safari/zebra/model/user"
	"github.com/rs/zerolog"
	"gojini.dev/config"
	"gojini.dev/web"
)

// Function to set up the logger.
//
// It takes in a pointer to config.Store and returns a context.Context.
func setupLogger(cfgStore *config.Store) context.Context {
	ctx := context.Background()
	zl := zerolog.New(os.Stderr).Level(zerolog.DebugLevel)
	logger := zerologr.New(&zl)

	return logr.NewContext(ctx, logger.WithName("zebra"))
}

// Function that sets up the adapter.
//
// It takes in a context.Context and a pointer to config.Store and returns a web.Adapter.
func setupAdapter(ctx context.Context, cfgStore *config.Store) web.Adapter {
	log := logr.FromContextOrDiscard(ctx)

	storeCfg := struct {
		Root string `json:"rootDir"`
	}{Root: ""}

	if e := cfgStore.Get("store", &storeCfg); e != nil {
		panic(e)
	}

	authKey := "key"

	if e := cfgStore.Get("authKey", &authKey); e != nil {
		panic(e)
	}

	factory := model.Factory()

	resAPI := NewResourceAPI(factory)
	if e := resAPI.Initialize(storeCfg.Root); e != nil {
		panic(e)
	}

	log.Info("zebra store initialized")

	if e := initAdminUser(log, resAPI.Store, cfgStore); e != nil {
		panic(e)
	}

	return func(nextHandler http.Handler) http.Handler {
		return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			if nextHandler == nil {
				panic("setup MUST have a next handler")
			}

			// Create a new request with logger in its context.
			log := logr.FromContextOrDiscard(ctx)
			ctx = req.Context()
			ctx = logr.NewContext(ctx, log)
			ctx = context.WithValue(ctx, AuthCtxKey, authKey)
			ctx = context.WithValue(ctx, ResourcesCtxKey, resAPI)

			newReq := req.Clone(ctx)

			// Call the next handler in the chain with the request with logger
			nextHandler.ServeHTTP(res, newReq)
		})
	}
}

// Function to initialize an admin user.
//
// It takes in a logr.Logger, zebra.Store, and a pointer to config.Store.
//
// It returns an error or nil in the absence thereof.
func initAdminUser(log logr.Logger, store zebra.Store, cfgStore *config.Store) error {
	admin := new(user.User)

	if err := cfgStore.Get("admin", admin); err != nil {
		return err
	}

	if findUser(store, admin.Email) == nil {
		log.Info("creating admin user")

		return store.Create(admin)
	}

	return nil
}
