package store_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/rchamarthy/zebra"
	"github.com/rchamarthy/zebra/network"
	"github.com/rchamarthy/zebra/store"
	"github.com/stretchr/testify/assert"
)

const vlan string = "VLANPool"

//nolint:gochecknoglobals
var (
	resource1 = &network.VLANPool{
		BaseResource: zebra.BaseResource{
			ID:     "0100000001",
			Type:   vlan,
			Labels: nil,
		},
		RangeStart: 0,
		RangeEnd:   10,
	}
	resource2 = &network.VLANPool{
		BaseResource: zebra.BaseResource{
			ID:     "0200000001",
			Type:   vlan,
			Labels: nil,
		},
		RangeStart: 1,
		RangeEnd:   5,
	}
)

func TestInitialize(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	assert.Nil(os.RemoveAll("teststore"))
	t.Cleanup(func() { os.RemoveAll("teststore") })

	filestore := store.NewFileStore("teststore", nil)

	assert.Nil(filestore.Initialize())
	assert.Nil(filestore.Initialize())
}

func TestCreate(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	assert.Nil(os.RemoveAll("teststore1"))
	t.Cleanup(func() { os.RemoveAll("teststore1") })

	// Create VLANPool resource
	resource := new(network.VLANPool)
	resource.ID = "0100000001"
	resource.Type = vlan
	resource.Labels = make(map[string]string)
	resource.Labels["key"] = "value"
	resource.RangeStart = 0
	resource.RangeEnd = 10

	types := make(map[string]func() zebra.Resource)

	types[vlan] = func() zebra.Resource { return new(network.VLANPool) }

	filestore := store.NewFileStore("teststore1", types)

	// Initialize store
	assert.Nil(filestore.Initialize())

	// Store object
	assert.Nil(filestore.Create(resource))
}

func TestUpdate(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	assert.Nil(os.RemoveAll("teststore1"))
	t.Cleanup(func() { os.RemoveAll("teststore1") })

	// Create VLANPool resource
	resource := new(network.VLANPool)
	resource.ID = "0100000001"
	resource.Type = vlan
	resource.Labels = make(map[string]string)
	resource.Labels["key"] = "value"
	resource.RangeStart = 0
	resource.RangeEnd = 10

	types := make(map[string]func() zebra.Resource)

	types[vlan] = func() zebra.Resource { return new(network.VLANPool) }

	filestore := store.NewFileStore("teststore1", types)

	// Initialize store
	assert.Nil(filestore.Initialize())

	assert.NotNil(filestore.Update(resource))

	// Store object
	assert.Nil(filestore.Create(resource))

	// Create VLANPool resource
	resource = new(network.VLANPool)
	resource.ID = "0100000001"
	resource.Type = vlan
	resource.Labels = make(map[string]string)
	resource.Labels["key1"] = "value1"
	resource.RangeStart = 0
	resource.RangeEnd = 10

	assert.NotNil(filestore.Create(resource))
	assert.Nil(filestore.Update(resource))
}

func TestLoad(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	assert.Nil(os.RemoveAll("teststore2"))
	t.Cleanup(func() { os.RemoveAll("teststore2") })

	// Create VLANPool resource
	resource := new(network.VLANPool)
	types := make(map[string]func() zebra.Resource)
	types[vlan] = func() zebra.Resource { return new(network.VLANPool) }
	types["IPAddressPool"] = func() zebra.Resource { return new(network.IPAddressPool) }
	filestore := store.NewFileStore("teststore2", types)

	// Initialize store
	assert.Nil(filestore.Initialize())

	// Store objects
	assert.Nil(filestore.Create(resource1))
	assert.NotNil(filestore.Create(resource))

	// Check that object is indeed stored
	_, err := os.Stat("teststore2/resources/01/00000001")
	assert.Nil(err)

	resources, err := filestore.Load()
	assert.Nil(err)

	assert.True(resources != nil && resources["0100000001"] != nil)
	assert.True(resources["0100000001"].GetType() == vlan)

	// Create invalid file and try loading
	fileName := "teststore2/resources/02/00000001"
	data := []byte("{\"test\":\"val\"}")

	assert.Nil(ioutil.WriteFile(fileName, data, 0o600))

	resources1, err := filestore.Load()
	assert.NotNil(err)
	assert.True(resources1["0100000001"] != nil)

	os.Remove("teststore2/resources/02/00000001")

	fileName = "teststore2/resources/03/00000001"
	data = []byte("{\"type\":\"IPAddressPool\"}")

	assert.Nil(ioutil.WriteFile(fileName, data, 0o600))

	_, err = filestore.Load()
	assert.NotNil(err)

	os.Remove("teststore2/resources/03/00000001")

	fileName = "teststore2/resources/04/00000001"

	data = []byte("{\"type\":\"invalid type\"}")

	assert.Nil(ioutil.WriteFile(fileName, data, 0o600))

	_, err = filestore.Load()
	assert.NotNil(err)
}

func TestDelete(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	assert.Nil(os.RemoveAll("teststore3"))
	t.Cleanup(func() { os.RemoveAll("teststore3") })

	// Create VLANPool resource
	resource := new(network.VLANPool)
	resource.ID = "0100000001"
	resource.Type = vlan
	resource.Labels = make(map[string]string)
	resource.RangeStart = 0
	resource.RangeEnd = 10

	types := make(map[string]func() zebra.Resource)

	types[vlan] = func() zebra.Resource { return new(network.VLANPool) }

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

	assert.Nil(os.RemoveAll("teststore4"))
	t.Cleanup(func() { os.RemoveAll("teststore4") })

	types := make(map[string]func() zebra.Resource)

	types[vlan] = func() zebra.Resource { return new(network.VLANPool) }

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

	assert.Nil(os.RemoveAll("teststore5"))
	t.Cleanup(func() { os.RemoveAll("teststore5") })

	filestore := store.NewFileStore("teststore5", nil)

	assert.Nil(filestore.Initialize())
	assert.Nil(filestore.Wipe())

	_, err := os.Stat("teststore5/resources")
	assert.True(os.IsNotExist(err))
}
