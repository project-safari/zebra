package filestore_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/filestore"
	"github.com/project-safari/zebra/network"
	"github.com/stretchr/testify/assert"
)

const vlan string = "VLANPool"

func TestInitialize(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	t.Cleanup(func() { os.RemoveAll("teststore") })

	fs := filestore.NewFileStore("teststore", nil)

	assert.Nil(fs.Initialize())
	assert.Nil(fs.Initialize())
}

func TestCreate(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore1"

	t.Cleanup(func() { os.RemoveAll(root) })

	// Create VLANPool resource
	resource := getVLAN()
	resource.Labels = zebra.Labels{"key": "value"}

	types := zebra.Factory()
	types.Add(vlan, func() zebra.Resource { return new(network.VLANPool) })

	fs := filestore.NewFileStore(root, types)

	// Initialize store
	assert.Nil(fs.Initialize())

	// Store object
	assert.Nil(fs.Create(resource))

	// Check that object is indeed stored
	_, err := os.Stat(getPath(root, resource))
	assert.Nil(err)

	// Store object again, should update
	resource.RangeStart = 1
	assert.Nil(fs.Create(resource))
}

func TestLoad(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore2"

	t.Cleanup(func() { os.RemoveAll(root) })

	// Create VLANPool resource
	resource := getVLAN()

	types := zebra.Factory()
	types.Add(vlan, func() zebra.Resource { return new(network.VLANPool) })

	fs := filestore.NewFileStore(root, types)

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

	list := resources.Resources[vlan].Resources
	assert.Equal(1, len(list))
	assert.Equal(vlan, list[0].GetType())
}

func TestDelete(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore3"

	t.Cleanup(func() { os.RemoveAll(root) })

	// Create VLANPool resource
	resource := getVLAN()

	types := zebra.Factory()
	types.Add(vlan, func() zebra.Resource { return new(network.VLANPool) })

	fs := filestore.NewFileStore(root, types)

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
}

func TestClearStore(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore4"

	t.Cleanup(func() { os.RemoveAll(root) })

	// Create first VLANPool resource
	resource1 := getVLAN()

	// Create second VLANPool resource
	resource2 := getVLAN()

	types := zebra.Factory()
	types.Add(vlan, func() zebra.Resource { return new(network.VLANPool) })

	fs := filestore.NewFileStore(root, types)

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

func TestWipeStore(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "teststore5"

	t.Cleanup(func() { os.RemoveAll(root) })

	fs := filestore.NewFileStore(root, nil)

	assert.Nil(fs.Initialize())
	assert.Nil(fs.Wipe())

	_, err := os.Stat(root + "/resources")
	assert.True(os.IsNotExist(err))
}

func TestBadLoad1(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "badtest1"

	t.Cleanup(func() { os.RemoveAll(root) })

	fs := filestore.NewFileStore(root, nil)

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

func TestBadLoad2(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "badtest1a"

	t.Cleanup(func() { os.RemoveAll(root) })

	types := zebra.Factory()
	types.Add(vlan, func() zebra.Resource { return new(network.VLANPool) })
	types.Add("Switch", func() zebra.Resource { return new(network.Switch) })

	fs := filestore.NewFileStore(root, types)

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

func TestBadDelete(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "badtest2"

	t.Cleanup(func() { os.RemoveAll(root) })

	fs := filestore.NewFileStore(root, nil)

	assert.Nil(fs.Initialize())

	resource := &network.VLANPool{
		BaseResource: zebra.BaseResource{
			ID:     "010",
			Type:   "VLANPool",
			Labels: nil,
		},
		RangeStart: 0,
		RangeEnd:   10,
	}
	assert.NotNil(fs.Delete(resource))
}

func getVLAN() *network.VLANPool {
	return &network.VLANPool{
		BaseResource: *zebra.NewBaseResource("VLANPool", nil),
		RangeStart:   0,
		RangeEnd:     1,
	}
}

func getPath(root string, res zebra.Resource) string {
	resID := res.GetID()

	return root + "/resources/" + resID[:2] + "/" + resID[2:]
}
