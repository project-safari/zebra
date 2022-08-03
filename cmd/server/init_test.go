package main //nolint:testpackage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitCmd(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	os.Args = append([]string{"zebra-server"}, "-c",
		"../../simulator/zebra-simulator.json",
		"init")

	assert.NotNil(execRootCmd())

	os.Args = append([]string{"zebra-server"}, "-c",
		"../../simulator/zebra-simulator.json",
		"init", "--user", "../../simulator/admin.yaml", "--password",
		"blah", "--auth-key", "blee")

	assert.Nil(execRootCmd())
}
