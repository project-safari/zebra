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

	t.Cleanup(func() { os.RemoveAll("teststore1") })

	// Create VLANPool resource
	resource := new(network.VLANPool)
	resource.ID = "0100000001"
	resource.Type = vlan
	resource.Labels = make(map[string]string)
	resource.Labels["key"] = "value"
	resource.RangeStart = 0
	resource.RangeEnd = 10

	types := zebra.Factory()
	types.Add(vlan, func() zebra.Resource { return new(network.VLANPool) })

	fs := filestore.NewFileStore("teststore1", types)

	// Initialize store
	assert.Nil(fs.Initialize())

	// Store object
	assert.Nil(fs.Create(resource))

	// Check that object is indeed stored
	_, err := os.Stat("teststore1/resources/01/00000001")
	assert.Nil(err)
}

func TestUpdate(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	t.Cleanup(func() { os.RemoveAll("teststore6") })

	// Create VLANPool resource
	resource := new(network.VLANPool)
	resource.ID = "0100000001"
	resource.Type = vlan
	resource.Labels = make(map[string]string)
	resource.Labels["key"] = "value"
	resource.RangeStart = 0
	resource.RangeEnd = 10

	types := zebra.Factory()
	types.Add(vlan, func() zebra.Resource { return new(network.VLANPool) })

	fs := filestore.NewFileStore("teststore6", types)

	// Initialize store
	assert.Nil(fs.Initialize())

	assert.NotNil(fs.Update(resource))

	// Store object
	assert.Nil(fs.Create(resource))

	// Check that object is indeed stored
	_, err := os.Stat("teststore6/resources/01/00000001")
	assert.Nil(err)

	// Create VLANPool resource
	resource = new(network.VLANPool)
	resource.ID = "0100000001"
	resource.Type = vlan
	resource.Labels = make(map[string]string)
	resource.Labels["key1"] = "value1"
	resource.RangeStart = 0
	resource.RangeEnd = 10

	err = fs.Create(resource)
	assert.NotNil(err)

	err = fs.Update(resource)
	assert.Nil(err)

	// Check that object is stored
	_, err = os.Stat("teststore6/resources/01/00000001")
	assert.Nil(err)
}

func TestLoad(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	t.Cleanup(func() { os.RemoveAll("teststore2") })

	// Create VLANPool resource
	resource := new(network.VLANPool)
	resource.ID = "0100000001"
	resource.Type = vlan
	resource.Labels = make(map[string]string)
	resource.RangeStart = 0
	resource.RangeEnd = 10

	types := zebra.Factory()
	types.Add(vlan, func() zebra.Resource { return new(network.VLANPool) })

	fs := filestore.NewFileStore("teststore2", types)

	// Initialize store
	assert.Nil(fs.Initialize())

	// Store object
	assert.Nil(fs.Create(resource))

	// Check that object is indeed stored
	_, err := os.Stat("teststore2/resources/01/00000001")
	assert.Nil(err)

	resources, err := fs.Load()
	assert.Nil(err)

	assert.True(resources != nil)

	list := resources.Resources["VLANPool"].Resources
	assert.True(len(list) == 1)
	assert.True(list[0].GetType() == vlan)
}

func TestDelete(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	t.Cleanup(func() { os.RemoveAll("teststore3") })

	// Create VLANPool resource
	resource := new(network.VLANPool)
	resource.ID = "0100000001"
	resource.Type = vlan
	resource.Labels = make(map[string]string)
	resource.RangeStart = 0
	resource.RangeEnd = 10

	types := zebra.Factory()
	types.Add(vlan, func() zebra.Resource { return new(network.VLANPool) })

	fs := filestore.NewFileStore("teststore3", types)

	// Initialize store
	assert.Nil(fs.Initialize())

	// Store object
	assert.Nil(fs.Create(resource))

	// Check that object is indeed stored
	_, err := os.Stat("teststore3/resources/01/00000001")
	assert.Nil(err)

	// Delete object and check it is deleted
	assert.Nil(fs.Delete(resource))

	_, err = os.Stat("teststore3/resources/01/00000001")
	assert.True(os.IsNotExist(err))
}

func TestClearStore(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	t.Cleanup(func() { os.RemoveAll("teststore4") })

	// Create first VLANPool resource
	resource1 := new(network.VLANPool)
	resource1.ID = "0100000001"
	resource1.Type = vlan
	resource1.Labels = make(map[string]string)
	resource1.RangeStart = 0
	resource1.RangeEnd = 10

	// Create second VLANPool resource
	resource2 := new(network.VLANPool)
	resource2.ID = "0200000001"
	resource2.Type = vlan
	resource2.Labels = make(map[string]string)
	resource2.RangeStart = 0
	resource2.RangeEnd = 10

	types := zebra.Factory()
	types.Add(vlan, func() zebra.Resource { return new(network.VLANPool) })

	fs := filestore.NewFileStore("teststore4", types)

	// Initialize store
	assert.Nil(fs.Initialize())

	// Store object
	assert.Nil(fs.Create(resource1))
	assert.Nil(fs.Create(resource2))

	// Check that object is indeed stored
	_, err := os.Stat("teststore4/resources/01/00000001")
	assert.Nil(err)

	_, err = os.Stat("teststore4/resources/02/00000001")
	assert.Nil(err)

	// Delete object and check it is deleted
	assert.Nil(fs.Clear())

	_, err = os.Stat("teststore4/resources/01/00000001")
	assert.True(os.IsNotExist(err))

	_, err = os.Stat("teststore4/resources/02/00000001")
	assert.True(os.IsNotExist(err))
}

func TestWipeStore(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	t.Cleanup(func() { os.RemoveAll("teststore5") })

	fs := filestore.NewFileStore("teststore5", nil)

	assert.Nil(fs.Initialize())
	assert.Nil(fs.Wipe())

	_, err := os.Stat("teststore5/resources")
	assert.True(os.IsNotExist(err))
}

func TestBadLoad1(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	t.Cleanup(func() { os.RemoveAll("badtest1") })

	fs := filestore.NewFileStore("badtest1", nil)

	vals, err := fs.Load()
	assert.Nil(vals)
	assert.NotNil(err)

	assert.Nil(os.MkdirAll("badtest1/resources/01", os.ModePerm))

	fd, err := os.OpenFile("badtest1/resources/01/00", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
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

	t.Cleanup(func() { os.RemoveAll("badtest1a") })

	types := zebra.Factory()
	types.Add(vlan, func() zebra.Resource { return new(network.VLANPool) })
	types.Add("Switch", func() zebra.Resource { return new(network.Switch) })

	fs := filestore.NewFileStore("badtest1a", types)

	assert.Nil(os.MkdirAll("badtest1a/resources/01", os.ModePerm))

	fileDes, err := os.OpenFile("badtest1a/resources/01/01", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
	assert.Nil(err)

	_, err = fmt.Fprintf(fileDes, `{"id":"00", "type":123}`)
	assert.Nil(err)

	_, err = fs.Load()
	assert.NotNil(err)

	fileDes, err = os.OpenFile("badtest1a/resources/01/02", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
	assert.Nil(err)

	_, err = fmt.Fprintf(fileDes, `{"id":"0100000001","type":"VLANPool","rangeStart":0}`)
	assert.Nil(err)

	_, err = fs.Load()
	assert.NotNil(err)

	fileDes, err = os.OpenFile("badtest1a/resources/01/04", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
	assert.Nil(err)

	_, err = fmt.Fprintf(fileDes, `{"id":"0100000002","type":"invalid","rangeStart":0, "rangeEnd":10}`)
	assert.Nil(err)

	_, err = fs.Load()
	assert.NotNil(err)

	fileDes, err = os.OpenFile("badtest1a/resources/01/05", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
	assert.Nil(err)

	_, err = fmt.Fprintf(fileDes, `{"id":"0100000002","type":"Switch","rangeStart":0, "rangeEnd":10}`)
	assert.Nil(err)

	_, err = fs.Load()
	assert.NotNil(err)

	fileDes, err = os.OpenFile("badtest1a/resources/01/03", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o666)
	assert.Nil(err)

	_, err = fmt.Fprintf(fileDes, "blahblah")
	assert.Nil(err)
	fileDes.Close()

	_, err = fs.Load()
	assert.NotNil(err)
}

func TestBadCreate(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	t.Cleanup(func() { os.RemoveAll("badtest2") })

	fs := filestore.NewFileStore("badtest2", nil)

	assert.Nil(fs.Initialize())
	assert.NotNil(fs.Create(new(network.VLANPool)))

	resource := &network.VLANPool{
		BaseResource: zebra.BaseResource{
			ID:     "01",
			Type:   "VLANPool",
			Labels: nil,
		},
		RangeStart: 0,
		RangeEnd:   10,
	}
	assert.NotNil(fs.Create(resource))
}

func TestBadUpdate(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	t.Cleanup(func() { os.RemoveAll("badtest3") })

	fs := filestore.NewFileStore("badtest3", nil)

	assert.Nil(fs.Initialize())
	assert.NotNil(fs.Update(new(network.VLANPool)))
}

func TestBadDelete(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	t.Cleanup(func() { os.RemoveAll("badtest4") })

	fs := filestore.NewFileStore("badtest4", nil)

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

	assert.NotNil(fs.Delete(new(network.VLANPool)))
	assert.NotNil(fs.Delete(resource))
}
