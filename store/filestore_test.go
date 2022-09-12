package store_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/store"
	"github.com/stretchr/testify/assert"
)

func factory() zebra.ResourceFactory {
	f := zebra.Factory()

	for i := 0; i < 10; i++ {
		t := zebra.Type{
			Name:        fmt.Sprintf("dummy-%d", i),
			Description: fmt.Sprintf("dummy type %d", i),
		}

		f.Add(t, func() zebra.Resource {
			return zebra.NewBaseResource(t, t.Name, t.Name, t.Name)
		})
	}

	return f
}

func TestFileStoreInitialize(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	fs := store.NewFileStore("test_file_store_initialize", factory())
	assert.Nil(fs.Initialize())

	defer func() { os.RemoveAll("test_file_store_initialize") }()
}

func TestFileStoreCreate(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore_fs_create"

	defer func() { os.RemoveAll(root) }()

	f := factory()
	fs := store.NewFileStore(root, f)

	// Initialize store
	assert.Nil(fs.Initialize())

	resource := f.New("dummy-1")
	// Store object
	assert.Nil(fs.Create(resource))

	// Check that object is indeed stored
	_, err := os.Stat(getPath(root, resource))
	assert.Nil(err)

	// Store object again, should update
	r, ok := resource.(*zebra.BaseResource)
	assert.True(ok)

	r.Meta.Name = "d1"
	assert.Nil(fs.Create(r))
}

func TestFSLoad(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore_fs_load"

	defer func() { os.RemoveAll(root) }()

	f := factory()
	resource := f.New("dummy-1")

	fs := store.NewFileStore(root, f)

	// Initialize store
	assert.Nil(fs.Initialize())

	// Store object
	assert.Nil(fs.Create(resource))

	// Check that object is indeed stored
	_, err := os.Stat(getPath(root, resource))
	assert.Nil(err)

	resources, err := fs.Load()
	assert.Nil(err)
	assert.NotNil(resources)
}

func TestFSDelete(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "test_fs_delete"

	defer func() { os.RemoveAll(root) }()

	f := factory()
	resource := f.New("dummy-1")

	fs := store.NewFileStore(root, f)

	// Initialize store
	assert.Nil(fs.Initialize())

	// Store object
	assert.Nil(fs.Create(resource))

	// Check that object is indeed stored
	_, err := os.Stat(getPath(root, resource))
	assert.Nil(err)

	// Delete object and check it is deleted
	assert.Nil(fs.Delete(resource))

	_, err = os.Stat(getPath(root, resource))
	assert.True(os.IsNotExist(err))

	// Try to delete res that doesn't exist, should succeed anyways
	assert.NotNil(fs.Delete(resource))
}

func TestFSClearStore(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "test_fs_clear_store"

	defer func() { os.RemoveAll(root) }()

	f := factory()

	resource1 := f.New("dummy-1")
	resource2 := f.New("dummy-1")

	fs := store.NewFileStore(root, f)

	// Initialize store
	assert.Nil(fs.Initialize())

	// Store object
	assert.Nil(fs.Create(resource1))
	assert.Nil(fs.Create(resource2))

	// Check that object is indeed stored
	_, err := os.Stat(getPath(root, resource1))
	assert.Nil(err)

	_, err = os.Stat(getPath(root, resource2))
	assert.Nil(err)

	// Delete object and check it is deleted
	assert.Nil(fs.Clear())

	_, err = os.Stat(getPath(root, resource1))
	assert.True(os.IsNotExist(err))

	_, err = os.Stat(getPath(root, resource2))
	assert.True(os.IsNotExist(err))
}

func TestFSWipeStore(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore_fs_wipe_store"

	defer func() { os.RemoveAll(root) }()

	fs := store.NewFileStore(root, nil)

	assert.Nil(fs.Initialize())
	assert.Nil(fs.Wipe())

	_, err := os.Stat(root + "/resources")
	assert.True(os.IsNotExist(err))
}

func TestFSBadLoad1(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore_fs_badtest1"

	defer func() { os.RemoveAll(root) }()

	fs := store.NewFileStore(root, nil)

	vals, err := fs.Load()
	assert.Nil(vals)
	assert.NotNil(err)

	assert.Nil(os.MkdirAll(root+"/resources/01", os.ModePerm))

	fd, err := os.OpenFile(root+"/resources/01/00", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
	assert.Nil(err)

	_, err = fmt.Fprintf(fd, `{"id":"0100000001","type":"VLANPool","rangeStart":0, "rangeEnd":10}`)
	assert.Nil(err)
	fd.Close()

	_, err = fs.Load()
	assert.NotNil(err)
}

func TestFSBadLoad2(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore_fs_badtest1a"

	defer func() { os.RemoveAll(root) }()

	f := factory()
	fs := store.NewFileStore(root, f)

	assert.Nil(os.MkdirAll(root+"/resources/01", os.ModePerm))

	fileDes, err := os.OpenFile(root+"/resources/01/01", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
	assert.Nil(err)

	_, err = fmt.Fprintf(fileDes, `{"id":"00", "type":123}`)
	assert.Nil(err)

	_, err = fs.Load()
	assert.NotNil(err)

	fileDes, err = os.OpenFile(root+"/resources/01/02", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
	assert.Nil(err)

	_, err = fmt.Fprintf(fileDes, `{"id":"0100000001","type":"VLANPool","rangeStart":0}`)
	assert.Nil(err)

	_, err = fs.Load()
	assert.NotNil(err)

	fileDes, err = os.OpenFile(root+"/resources/01/04", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
	assert.Nil(err)

	_, err = fmt.Fprintf(fileDes, `{"id":"0100000002","type":"invalid","rangeStart":0, "rangeEnd":10}`)
	assert.Nil(err)

	_, err = fs.Load()
	assert.NotNil(err)

	fileDes, err = os.OpenFile(root+"/resources/01/05", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
	assert.Nil(err)

	_, err = fmt.Fprintf(fileDes, `{"id":"0100000002","type":"Switch","rangeStart":0, "rangeEnd":10}`)
	assert.Nil(err)

	_, err = fs.Load()
	assert.NotNil(err)

	fileDes, err = os.OpenFile(root+"/resources/01/03", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
	assert.Nil(err)

	_, err = fmt.Fprintf(fileDes, "blahblah")
	assert.Nil(err)
	fileDes.Close()

	_, err = fs.Load()
	assert.NotNil(err)
}

func getPath(root string, res zebra.Resource) string {
	resID := res.GetMeta().ID

	return root + "/resources/" + resID[:2] + "/" + resID[2:]
}
