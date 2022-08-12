package main //nolint:testpackage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go tests for lease commands.
func TestLease(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	argLock.Lock()
	defer argLock.Unlock()

	os.Args = append([]string{"zebra"}, "lease")

	assert.NotNil(execRootCmd())

	os.Args = append([]string{"zebra"}, "-c", "../../simulator/admin.yaml",
		"lease", "blah")

	assert.NotNil(execRootCmd())

	os.Args = append([]string{"zebra"}, "-c", "junk.yaml",
		"lease", "Server")

	assert.NotNil(execRootCmd())
}
