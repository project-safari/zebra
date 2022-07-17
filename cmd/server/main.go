package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"
	"github.com/julienschmidt/httprouter"
	"github.com/project-safari/zebra/api"
	"github.com/project-safari/zebra/store"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"gojini.dev/config"
	"gojini.dev/web"
)

const version = "unknown"

func main() {
	name := filepath.Base(os.Args[0])
	rootCmd := &cobra.Command{ // nolint:exhaustruct,exhaustivestruct
		Use:          name,
		Short:        "zebra server",
		Version:      version + "\n",
		RunE:         run,
		SilenceUsage: true,
	}
	rootCmd.SetVersionTemplate(version + "\n")
	rootCmd.Flags().StringP("config", "c", path.Join(
		func() string {
			s, _ := os.Getwd()

			return s
		}(), "server.json"),
		"config file (default: $PWD/server.json",
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	// Load server configuration
	cfgFile := cmd.Flag("config").Value.String()

	cfgStore := config.New()
	if err := cfgStore.LoadFromFile(context.Background(), cfgFile); err != nil {
		return err
	}

	return startServer(cfgStore)
}

func setupLogger(cfgStore *config.Store) context.Context {
	ctx := context.Background()
	zl := zerolog.New(os.Stderr).Level(zerolog.DebugLevel)
	logger := zerologr.New(&zl)

	return logr.NewContext(ctx, logger.WithName("zebra"))
}

func startServer(cfgStore *config.Store) error {
	appCtx := setupLogger(cfgStore)

	log := logr.FromContextOrDiscard(appCtx)

	serverCfg := new(web.Config)
	if e := cfgStore.Get("server", serverCfg); e != nil {
		log.Error(e, "web server config missing")

		return e
	}

	handler := httpHandler(appCtx, cfgStore)
	webServer := web.NewServer(serverCfg, handler)

	return webServer.Start(appCtx)
}

func httpHandler(ctx context.Context, cfgStore *config.Store) http.Handler {
	log := logr.FromContextOrDiscard(ctx)
	storeCfg := struct {
		Root string `json:"rootDir"`
	}{Root: ""}

	if e := cfgStore.Get("store", &storeCfg); e != nil {
		log.Error(e, "store configuration missing")
		panic(e)
	}

	authKey := "key"

	if e := cfgStore.Get("authKey", &authKey); e != nil {
		log.Error(e, "auth key missing")
		panic(e)
	}

	factory := store.DefaultFactory()

	resAPI := api.NewResourceAPI(factory)
	if e := resAPI.Initialize(storeCfg.Root); e != nil {
		log.Error(e, "api initialization failed")
		panic(e)
	}

	router := httprouter.New()
	router.GET("/api/v1/resources", handleQuery(resAPI))
	router.GET("/api/v1/types", handleTypes(ctx))
	router.GET("/login", handleLogin(ctx, resAPI.Store, authKey))

	return router
}

func handleQuery(resAPI *api.ResourceAPI) httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		switch {
		case strings.HasPrefix(req.URL.RawQuery, "id"):
			resAPI.GetResourcesByID(res, req)

		case strings.HasPrefix(req.URL.RawQuery, "type"):
			resAPI.GetResourcesByType(res, req)

		case strings.HasPrefix(req.URL.RawQuery, "property"):
			resAPI.GetResourcesByProperty(res, req)

		case strings.HasPrefix(req.URL.RawQuery, "label"):
			resAPI.GetResourcesByLabel(res, req)

		default:
			resAPI.GetResources(res, req)
		}
	}
}

func writeJSON(ctx context.Context, res http.ResponseWriter, data interface{}) {
	log := logr.FromContextOrDiscard(ctx)

	bytes, err := json.Marshal(data)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusOK)
	if _, err := res.Write(bytes); err != nil {
		log.Error(err, "error writing response")
	}
}
