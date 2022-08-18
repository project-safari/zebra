package main //nolint:testpackage

import (
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

var argLock sync.Mutex //nolint:gochecknoglobals

// Testing execution of main.
func TestMain(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	argLock.Lock()
	defer argLock.Unlock()

	os.Args = append([]string{"zebra"}, "bad args")

	e := execRootCmd()
	assert.NotNil(e)

	os.Args = append([]string{"zebra"}, "--help")

	assert.Nil(execRootCmd())
}
