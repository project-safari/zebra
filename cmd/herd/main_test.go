package main //nolint:testpackage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	stored := initStore("./")
	assert.NotNil(stored)
}
