package main //nolint:testpackage

import (
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var argLock sync.Mutex //nolint:gochecknoglobals

// Function to test execution of main.
func TestMain(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	argLock.Lock()
	defer argLock.Unlock()

	defer func() { os.RemoveAll("herd_store") }()

	os.Args = []string{"herd"}

	e := execRootCmd()
	assert.Nil(e)

	os.Args = append([]string{"herd"}, "--help")

	assert.Nil(execRootCmd())
}
