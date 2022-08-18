package store_test

import (
	"testing"

	"github.com/project-safari/zebra/store"
	"github.com/stretchr/testify/assert"
)

// Test for generating a resource factory with all the known types.
func TestDefaultFactory(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	types := store.DefaultFactory()
	assert.NotNil(types)

	typeList := types.Types()
	assert.NotEmpty(typeList)

	for _, t := range typeList {
		assert.NotNil(t.New())
	}
}
