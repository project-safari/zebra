package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"
	"github.com/julienschmidt/httprouter"
	"github.com/rchamarthy/zebra"
	"github.com/rchamarthy/zebra/api"
	"github.com/rchamarthy/zebra/network"
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
		}(), "zebra-server.json"),
		"config file (default: $PWD/zebra-server.json",
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

	return logr.NewContext(ctx, logger.WithName("helper"))
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

	factory := initTypes()

	resAPI := api.NewResourceAPI(factory)
	if e := resAPI.Initialize(storeCfg.Root); e != nil {
		log.Error(e, "api initialization failed")
		panic(e)
	}

	router := httprouter.New()
	router.GET("/api/v1/resources", handle(resAPI))

	return router
}

func initTypes() zebra.ResourceFactory {
	factory := zebra.Factory()
	factory.Add("VLANPool", func() zebra.Resource {
		return new(network.VLANPool)
	})

	// Need to add all the known types here
	return factory
}

func handle(resAPI *api.ResourceAPI) httprouter.Handle {
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
