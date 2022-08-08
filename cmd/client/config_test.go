package main //nolint:testpackage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testCfgFile = "./test_config.yaml"

func TestConfig(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	cfg := NewConfig()
	assert.NotNil(cfg)

	assert.Nil(cfg.Save(testCfgFile))

	defer func() { assert.Nil(os.Remove(testCfgFile)) }()

	argLock.Lock()
	defer argLock.Unlock()

	os.Args = append([]string{"zebra"}, "-c", testCfgFile, "config")

	assert.Nil(execRootCmd())

	os.Args = append([]string{"zebra"}, "-c", testCfgFile, "config", "init", "https://127.0.0.1:6666")

	assert.Nil(execRootCmd())

	os.Args = append([]string{"zebra"}, "-c", testCfgFile, "config", "user", "tester")

	assert.Nil(execRootCmd())

	os.Args = append([]string{"zebra"}, "-c", testCfgFile, "config", "email", "tester@zebra.safari.io")

	assert.Nil(execRootCmd())

	os.Args = append([]string{"zebra"}, "-c", testCfgFile, "config", "ca-cert", "cert_file")

	assert.Nil(execRootCmd())

	os.Args = append([]string{"zebra"}, "-c", testCfgFile, "config", "server", "https://zebra.safari.io")

	assert.Nil(execRootCmd())

	os.Args = append([]string{"zebra"}, "-c", testCfgFile, "config", "defaults", "--duration", "10")

	assert.NotNil(execRootCmd())

	os.Args = append([]string{"zebra"}, "-c", testCfgFile, "config", "defaults", "--duration", "2")

	assert.Nil(execRootCmd())

	os.Args = append([]string{"zebra"}, "-c", testCfgFile, "config", "public-key")

	assert.Nil(execRootCmd())

	cfg, err := Load("junk")
	assert.NotNil(err)
	assert.Nil(cfg)

	cfg, err = Load(testCfgFile)
	assert.Nil(err)
	assert.NotNil(cfg)
	assert.Equal("tester", cfg.User)
	assert.Equal("tester@zebra.safari.io", cfg.Email)
	assert.Equal(2, cfg.Defaults.Duration)
	assert.Equal("https://zebra.safari.io", cfg.ServerAddress)
}
