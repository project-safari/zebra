package main //nolint:testpackage

import (
	"testing"

	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/assert"
)

// tests for routeHandler().
func TestRoutes(t *testing.T) {
	t.Parallel()

	assert := assert.New(t)

	router := routeHandler()
	assert.NotNil(router)

	r, ok := routeHandler().(*httprouter.Router)
	assert.NotNil(r)
	assert.True(ok)
}
