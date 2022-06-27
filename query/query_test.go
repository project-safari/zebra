package query_test

import (
	"net"
	"testing"

	"github.com/rchamarthy/zebra"
	"github.com/rchamarthy/zebra/network"
	"github.com/rchamarthy/zebra/query"
	"github.com/stretchr/testify/assert"
)

const (
	vlan  = "VLANPool"
	ipool = "IPAddressPool"
)

func TestNewQueryStore(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// Create first VLANPool resource
	resource1 := new(network.VLANPool)
	resource1.ID = "0100000001"
	resource1.Type = vlan
	resource1.RangeStart = 0
	resource1.RangeEnd = 10

	// Add resources to map
	resources := make(map[string]zebra.Resource)
	resources["0100000001"] = resource1

	qs := query.NewQueryStore(resources)
	assert.NotNil(qs)
}

func TestInitialize(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// Create first VLANPool resource
	resource1 := new(network.VLANPool)
	resource1.ID = "0100000001"
	resource1.Type = vlan
	resource1.RangeStart = 0
	resource1.RangeEnd = 10

	// Add resources to map
	resources := make(map[string]zebra.Resource, 1)
	resources["0100000001"] = resource1

	querystore := query.NewQueryStore(resources)
	assert.NotNil(querystore)

	err := querystore.Initialize()
	assert.Nil(err)

	ret, err := querystore.Load()
	assert.True(err == nil && len(ret) == 1)
	assert.True(ret["0100000001"] == resource1)
}

func TestWipe(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// Create first VLANPool resource
	resource1 := new(network.VLANPool)
	resource1.ID = "0100000001"
	resource1.Type = vlan
	resource1.RangeStart = 0
	resource1.RangeEnd = 10

	// Add resources to map
	resources := make(map[string]zebra.Resource)
	resources["0100000001"] = resource1

	querystore := query.NewQueryStore(resources)
	assert.NotNil(querystore)

	err := querystore.Initialize()
	assert.Nil(err)

	err = querystore.Wipe()
	assert.Nil(err)
	assert.NotNil(querystore)
}

func TestClear(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// Create first VLANPool resource
	resource1 := new(network.VLANPool)
	resource1.ID = "0100000001"
	resource1.Type = vlan
	resource1.RangeStart = 0
	resource1.RangeEnd = 10

	// Add resources to map
	resources := make(map[string]zebra.Resource)
	resources["0100000001"] = resource1

	querystore := query.NewQueryStore(resources)
	assert.NotNil(querystore)

	err := querystore.Initialize()
	assert.Nil(err)

	err = querystore.Clear()
	assert.Nil(err)

	ret, err := querystore.Load()
	assert.Nil(err)
	assert.True(len(ret) == 0)
}

func TestCreate(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// Create first VLANPool resource
	resource1 := new(network.VLANPool)
	resource1.ID = "0100000001"
	resource1.Type = vlan
	resource1.RangeStart = 0
	resource1.RangeEnd = 10

	// Create second VLANPool resource
	resource2 := new(network.VLANPool)
	resource2.ID = "0200000001"
	resource2.Type = vlan
	resource2.Labels = map[string]string{"stage": "prod"}
	resource2.RangeStart = 1
	resource2.RangeEnd = 5

	// Add resources to map
	resources := make(map[string]zebra.Resource)
	resources["0100000001"] = resource1

	querystore := query.NewQueryStore(resources)
	assert.NotNil(querystore)

	err := querystore.Initialize()
	assert.Nil(err)

	ret, err := querystore.Load()
	assert.True(err == nil && len(ret) == 1)
	assert.True(ret["0100000001"] == resource1)

	err = querystore.Create(resource2)
	assert.Nil(err)

	ret, err = querystore.Load()
	assert.True(err == nil && len(ret) == 2)
	assert.True(ret["0100000001"] == resource1)
	assert.True(ret["0200000001"] == resource2)

	// Create a third VLANPool resource with same ID as resource2
	resource3 := new(network.VLANPool)
	resource3.ID = "0200000001"
	resource3.Type = vlan
	resource3.Labels = map[string]string{"stagetest": "dev"}
	resource3.RangeStart = 1
	resource3.RangeEnd = 5

	err = querystore.Create(resource3)
	assert.Nil(err)

	res, err := querystore.QueryLabel(query.Query{Op: query.MatchEqual, Key: "stage", Values: []string{"prod"}})
	assert.True(err == nil && len(res) == 0)

	res, err = querystore.QueryLabel(query.Query{Op: query.MatchEqual, Key: "stagetest", Values: []string{"dev"}})
	assert.True(err == nil && len(res) == 1 && res[0].GetID() == "0200000001")
}

func TestDelete(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// Create first VLANPool resource
	resource1 := new(network.VLANPool)
	resource1.ID = "0100000001"
	resource1.Type = vlan
	resource1.RangeStart = 0
	resource1.RangeEnd = 10

	// Create second VLANPool resource
	resource2 := new(network.VLANPool)
	resource2.ID = "0200000001"
	resource2.Type = vlan
	resource2.Labels = map[string]string{"stage": "dev"}
	resource2.RangeStart = 1
	resource2.RangeEnd = 5

	// Add resources to map
	resources := make(map[string]zebra.Resource)
	resources["0100000001"] = resource1
	resources["0200000001"] = resource2

	querystore := query.NewQueryStore(resources)
	assert.NotNil(querystore)

	err := querystore.Initialize()
	assert.Nil(err)

	ret, err := querystore.Load()
	assert.True(err == nil && len(ret) == 2)

	err = querystore.Delete(resource2)
	assert.Nil(err)

	ret, err = querystore.Load()
	assert.True(err == nil && len(ret) == 1)
	assert.True(ret["0100000001"] == resource1)
}

func TestQuery(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// Create VLANPool resource
	resource1 := new(network.VLANPool)
	resource1.ID = "0100000001"
	resource1.Type = vlan
	resource1.RangeStart = 0
	resource1.RangeEnd = 10

	// Create IPAddressPool resource
	resource2 := new(network.IPAddressPool)
	resource2.ID = "0200000001"
	resource2.Type = "IPAddressPool"
	ip := net.ParseIP("10.0.0.1")
	mask := ip.DefaultMask()
	resource2.Subnets = []net.IPNet{{IP: ip, Mask: mask}}

	// Add resources to map
	resources := make(map[string]zebra.Resource)
	resources["0100000001"] = resource1
	resources["0200000001"] = resource2

	querystore := query.NewQueryStore(resources)
	assert.NotNil(querystore)

	err := querystore.Initialize()
	assert.Nil(err)

	all := querystore.Query()
	assert.True(len(all) == 2)
	assert.True(all[0].GetID() == "0100000001" || all[0].GetID() == "0200000001")
	assert.True(all[1].GetID() == "0100000001" || all[1].GetID() == "0200000001")
	assert.True(all[0].GetID() != all[1].GetID())
}

func TestQueryUUID(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// Create VLANPool resource
	resource1 := new(network.VLANPool)
	resource1.ID = "0100000001"
	resource1.Type = vlan
	resource1.RangeStart = 0
	resource1.RangeEnd = 10

	// Create IPAddressPool resource
	resource2 := new(network.IPAddressPool)
	resource2.ID = "0200000001"
	resource2.Type = ipool
	ip := net.ParseIP("10.0.0.1")
	mask := ip.DefaultMask()
	resource2.Subnets = []net.IPNet{{IP: ip, Mask: mask}}

	// Add resources to map
	resources := make(map[string]zebra.Resource)
	resources["0100000001"] = resource1
	resources["0200000001"] = resource2

	querystore := query.NewQueryStore(resources)
	assert.NotNil(querystore)

	err := querystore.Initialize()
	assert.Nil(err)

	results := querystore.QueryUUID([]string{"0100000001"})
	assert.True(len(results) == 1 && results[0].GetID() == "0100000001")

	results = querystore.QueryUUID([]string{"0100000001", "0200000001"})
	assert.True(len(results) == 2)

	results = querystore.QueryUUID([]string{"0100000001", "0300000001"})
	assert.True(len(results) == 1)
}

func TestQueryType(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// Create VLANPool resource
	resource1 := new(network.VLANPool)
	resource1.ID = "0100000001"
	resource1.Type = vlan
	resource1.RangeStart = 0
	resource1.RangeEnd = 10

	// Create IPAddressPool resource
	resource2 := new(network.IPAddressPool)
	resource2.ID = "0200000001"
	resource2.Type = ipool
	ip := net.ParseIP("10.0.0.1")
	mask := ip.DefaultMask()
	resource2.Subnets = []net.IPNet{{IP: ip, Mask: mask}}

	// Add resources to map
	resources := make(map[string]zebra.Resource)
	resources["0100000001"] = resource1
	resources["0200000001"] = resource2

	querystore := query.NewQueryStore(resources)

	err := querystore.Initialize()
	assert.Nil(err)

	vlanpools := querystore.QueryType("VLANPool")
	assert.True(len(vlanpools) == 1)
	assert.True(vlanpools[0].GetID() == "0100000001")

	ippools := querystore.QueryType(ipool)
	assert.True(len(ippools) == 1)
	assert.True(ippools[0].GetID() == "0200000001")
}

func TestInvalidLabelQuery(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// Create VLANPool resource
	resource1 := new(network.VLANPool)
	resource1.ID = "0100000001"
	resource1.Type = vlan
	resource1.Labels = make(map[string]string)
	resource1.Labels["product-owner"] = "owner"
	resource1.RangeStart = 0
	resource1.RangeEnd = 10

	// Add resources to map
	resources := make(map[string]zebra.Resource)
	resources["0100000001"] = resource1

	querystore := query.NewQueryStore(resources)

	err := querystore.Initialize()
	assert.Nil(err)

	query := query.Query{
		Op:     7,
		Key:    "",
		Values: nil,
	}

	// Should fail on invalid query.
	_, err = querystore.QueryLabel(query)
	assert.NotNil(err)
}

func getResources() (*network.VLANPool, *network.IPAddressPool) {
	// Create VLANPool resource
	resource1 := &network.VLANPool{
		BaseResource: zebra.BaseResource{
			ID:     "0100000001",
			Type:   vlan,
			Labels: map[string]string{"product-owner": "shravya"},
		},
		RangeStart: 0,
		RangeEnd:   10,
	}

	// Create IPAddressPool resource
	ipaddress := net.ParseIP("10.0.0.1")
	mask := ipaddress.DefaultMask()
	resource2 := &network.IPAddressPool{
		BaseResource: zebra.BaseResource{
			ID:   "0200000001",
			Type: ipool,
			Labels: map[string]string{
				"product-owner": "nandyala",
				"team":          "cloud networking",
			},
		},
		Subnets: []net.IPNet{{IP: ipaddress, Mask: mask}},
	}

	return resource1, resource2
}

func TestQueryLabel(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resource1, resource2 := getResources()

	// Add resources to map
	resources := make(map[string]zebra.Resource)
	resources["0100000001"] = resource1
	resources["0200000001"] = resource2

	querystore := query.NewQueryStore(resources)

	err := querystore.Initialize()
	assert.Nil(err)

	query1 := query.Query{Op: query.MatchEqual, Key: "product-owner", Values: []string{"shravya", "nandyala"}}
	query2 := query.Query{Op: query.MatchIn, Key: "product-owner", Values: []string{"shravya", "nandyala"}}
	query3 := query.Query{Op: query.MatchNotEqual, Key: "product-owner", Values: []string{"shravya", "nandyala"}}
	query4 := query.Query{Op: query.MatchNotIn, Key: "product-owner", Values: []string{"shravya", "nandyala"}}

	// Should fail on query 1 and query 3.
	_, err = querystore.QueryLabel(query1)
	assert.NotNil(err)

	_, err = querystore.QueryLabel(query3)
	assert.NotNil(err)

	// Update query 1, should succeed.
	query1.Values = []string{"nandyala"}
	pos, err := querystore.QueryLabel(query1)
	assert.True(err == nil && len(pos) == 1 && pos[0].GetID() == resource2.ID)

	// Should succeed on query 2, return both resources.
	pos, err = querystore.QueryLabel(query2)
	assert.True(err == nil && len(pos) == 2 && pos[0] != pos[1])

	// Should succeed on query 4, return no resources.
	pos, err = querystore.QueryLabel(query4)
	assert.True(err == nil && len(pos) == 0)

	// Update query 3 to be valid, return 1 resource.
	query3.Values = []string{"shravya"}
	pos, err = querystore.QueryLabel(query3)
	assert.True(err == nil && len(pos) == 1 && pos[0].GetID() == "0200000001")
}

func TestQueryProperty(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	resource1, resource2 := getResources()

	// Add resources to map
	resources := make(map[string]zebra.Resource)
	resources["0100000001"] = resource1
	resources["0200000001"] = resource2

	querystore := query.NewQueryStore(resources)

	err := querystore.Initialize()
	assert.Nil(err)

	query1 := query.Query{Op: query.MatchEqual, Key: "Type", Values: []string{vlan, ipool}}
	query2 := query.Query{Op: query.MatchIn, Key: "Type", Values: []string{vlan}}
	query3 := query.Query{Op: query.MatchNotEqual, Key: "Type", Values: []string{vlan, ipool}}
	query4 := query.Query{Op: query.MatchNotIn, Key: "Type", Values: []string{vlan, ipool}}

	// Should fail on query 1 and query 3.
	_, err = querystore.QueryProperty(query1)
	assert.NotNil(err)

	_, err = querystore.QueryProperty(query3)
	assert.NotNil(err)

	// Update query 1, should succeed.
	query1.Values = []string{ipool}
	pos, err := querystore.QueryProperty(query1)
	assert.Nil(err)
	assert.True(len(pos) == 1)
	assert.True(pos[0].GetID() == resource2.ID)

	// Should succeed on query 2, return first resource.
	pos, err = querystore.QueryProperty(query2)
	assert.Nil(err)
	assert.True(len(pos) == 1)
	assert.True(pos[0].GetID() == resource1.ID)

	// Should succeed on query 4, return no resources.
	pos, err = querystore.QueryProperty(query4)
	assert.True(err == nil && len(pos) == 0)

	// Update query 3 to be valid, return 1 resource.
	query3.Values = []string{ipool}
	pos, err = querystore.QueryProperty(query3)
	assert.Nil(err)
	assert.True(len(pos) == 1)
	assert.True(pos[0].GetID() == resource1.ID)

	pos, err = querystore.QueryProperty(query.Query{Op: 0x7, Key: "", Values: []string{""}})
	assert.Nil(pos)
	assert.NotNil(err)
}
