package store_test

import (
	"os"
	"testing"

	"github.com/rchamarthy/zebra"
	"github.com/rchamarthy/zebra/network"
	"github.com/rchamarthy/zebra/store"
	"github.com/stretchr/testify/assert"
)

const vlan string = "VLANPool"

func TestInitialize(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	t.Cleanup(func() { os.RemoveAll("teststore") })

	filestore := store.NewFileStore("teststore", nil)

	assert.Nil(filestore.Initialize())
	assert.Nil(filestore.Initialize())
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

	filestore := store.NewFileStore("teststore1", types)

	// Initialize store
	assert.Nil(filestore.Initialize())

	// Store object
	assert.Nil(filestore.Create(resource))

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

	filestore := store.NewFileStore("teststore6", types)

	// Initialize store
	assert.Nil(filestore.Initialize())

	assert.NotNil(filestore.Update(resource))

	// Store object
	assert.Nil(filestore.Create(resource))

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

	err = filestore.Create(resource)
	assert.NotNil(err)

	err = filestore.Update(resource)
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

	filestore := store.NewFileStore("teststore2", types)

	// Initialize store
	assert.Nil(filestore.Initialize())

	// Store object
	assert.Nil(filestore.Create(resource))

	// Check that object is indeed stored
	_, err := os.Stat("teststore2/resources/01/00000001")
	assert.Nil(err)

	resources, err := filestore.Load()
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

	filestore := store.NewFileStore("teststore3", types)

	// Initialize store
	assert.Nil(filestore.Initialize())

	// Store object
	assert.Nil(filestore.Create(resource))

	// Check that object is indeed stored
	_, err := os.Stat("teststore3/resources/01/00000001")
	assert.Nil(err)

	// Delete object and check it is deleted
	assert.Nil(filestore.Delete(resource))

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

	filestore := store.NewFileStore("teststore4", types)

	// Initialize store
	assert.Nil(filestore.Initialize())

	// Store object
	assert.Nil(filestore.Create(resource1))
	assert.Nil(filestore.Create(resource2))

	// Check that object is indeed stored
	_, err := os.Stat("teststore4/resources/01/00000001")
	assert.Nil(err)

	_, err = os.Stat("teststore4/resources/02/00000001")
	assert.Nil(err)

	// Delete object and check it is deleted
	assert.Nil(filestore.Clear())

	_, err = os.Stat("teststore4/resources/01/00000001")
	assert.True(os.IsNotExist(err))

	_, err = os.Stat("teststore4/resources/02/00000001")
	assert.True(os.IsNotExist(err))
}

func TestWipeStore(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	t.Cleanup(func() { os.RemoveAll("teststore5") })

	filestore := store.NewFileStore("teststore5", nil)

	assert.Nil(filestore.Initialize())
	assert.Nil(filestore.Wipe())

	_, err := os.Stat("teststore5/resources")
	assert.True(os.IsNotExist(err))
}