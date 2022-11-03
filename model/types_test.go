package model_test

import (
	"testing"

	"github.com/project-safari/zebra/model"
	"github.com/stretchr/testify/assert"
)

// Test function that tests known types in a given resource factory.
func TestDefaultFactory(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)
	types := model.Factory()
	assert.NotNil(types)

	typeList := types.Types()
	assert.NotEmpty(typeList)

	for _, t := range typeList {
		assert.NotNil(t.Name)
		assert.NotNil(types.New(t.Name))
	}
}
