package main //nolint:testpackage

import (
	"os"
	"path"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"gojini.dev/config"
	"gojini.dev/web"
)

func TestRun(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	rootCmd := new(cobra.Command)

	rootCmd.Use = "named"
	rootCmd.Short = "zebra server"

	rootCmd.Version = version + "\n"
	rootCmd.RunE = run

	rootCmd.SilenceUsage = true
	rootCmd.SetVersionTemplate(version + "\n")

	rootCmd.Flags().StringP("config", "c", path.Join(
		func() string {
			s, _ := os.Getwd()

			return s
		}(), "server.json"),
		"config file (default: $PWD/server.json",
	)

	args := []string{"", "test"}
	running := run(rootCmd, args)

	assert.NotNil(running)
}

func TestStartSrv(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	serverCfg := new(web.Config)

	serverCfg.Address = new(web.Address)

	assert.NotNil(serverCfg)

	cfgStore := new(config.Store)
	assert.NotNil(cfgStore)

	appCtx := setupLogger(cfgStore)
	assert.NotNil(appCtx)

	stored := startServer(cfgStore)

	assert.NotNil(stored)

	log := loginAdapter()
	assert.NotNil(log)

	register := registerAdapter()
	assert.NotNil(register)

	auth := authAdapter()
	assert.NotNil(auth)

	refresh := refreshAdapter()
	assert.NotNil(refresh)

	routes := routeHandler()
	assert.NotNil(routes)

	s := makeStore(assert, "./")
	assert.NotNil(s)

	setup := new(web.Adapter)
	assert.NotNil(setup)
}
