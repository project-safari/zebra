package main //nolint:testpackage

import (
	"os"
	"testing"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/cmd/herd/pkg"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	stored := initStore("test_store")
	assert.NotNil(stored)
	assert.Nil(os.RemoveAll("test_store"))
}

func TestFileStore(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	fs := initStore("herd_file_store")

	defer os.RemoveAll("herd_file_store")

	resources := make([]zebra.Resource, 0, 2)
	res := storeResources(resources, fs)

	assert.Nil(res)

	labels := pkg.CreateLabels()
	resources = append(resources, zebra.NewBaseResource("", labels))
	ret := storeResources(resources, fs)

	assert.Nil(ret)
}

func TestInitStore(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	initial := initStore("./users")

	assert.NotNil(initial)
}
