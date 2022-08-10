package main //nolint:testpackage

import (
	"os"
	"testing"

	"github.com/project-safari/zebra/auth"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	argLock.Lock()
	defer argLock.Unlock()

	cfg := createConfig(assert)

	os.Args = append([]string{"zebra"}, "registration")

	assert.NotNil(execRootCmd())

	os.Args = append([]string{"zebra"}, "-c", "../../simulator/admin.yaml",
		"registration", "blah")

	assert.NotNil(execRootCmd())

	os.Args = append([]string{"zebra"}, "-c", "junk.yaml",
		"registration", "Server")

	assert.NotNil(execRootCmd())

	assert.Nil(cfg.Save(testCfgFile))

	defer func() { assert.Nil(os.Remove(testCfgFile)) }()

	os.Args = append([]string{"zebra"}, "-c", testCfgFile, "registration", "blah")

	assert.NotNil(execRootCmd())
}

func createConfig(assert *assert.Assertions) *Config {
	cfg := NewConfig()
	assert.NotNil(cfg)

	cfg.Email = "test@zebra.project-safafi.io"
	key, err := auth.Load(testUserKeyFile)
	assert.Nil(err)

	cfg.Key = key
	cfg.CACert = testCACertFile
	cfg.Key = key
	c, err := NewClient(cfg)
	assert.Nil(err)
	assert.NotNil(c)

	return cfg
}
