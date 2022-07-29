package main //nolint:testpackage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	os.Args = append([]string{"zebra"}, "bad args")

	e := execRootCmd()
	assert.NotNil(e)

	os.Args = append([]string{"zebra"}, "--help")

	assert.Nil(execRootCmd())
}
