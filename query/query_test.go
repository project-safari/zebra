package query_test

import (
	"net"
	"testing"

	"github.com/rchamarthy/zebra"
	"github.com/rchamarthy/zebra/network"
	"github.com/rchamarthy/zebra/query"
	"github.com/stretchr/testify/assert"
)

func TestNewQueryStore(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// Create first VLANPool resource
	resource1 := new(network.VLANPool)
	resource1.ID = "0100000001"
	resource1.RangeStart = 0
	resource1.RangeEnd = 10

	// Create second VLANPool resource
	resource2 := new(network.VLANPool)
	resource2.ID = "0200000001"
	resource2.RangeStart = 0
	resource2.RangeEnd = 10

	// Add resources to map
	resources := make(map[string]zebra.Resource)
	resources["0100000001"] = resource1
	resources["0200000001"] = resource2

	qs := query.NewQueryStore(resources)
	assert.NotNil(qs)
}

func TestQuery(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// Create VLANPool resource
	resource1 := new(network.VLANPool)
	resource1.ID = "0100000001"
	resource1.RangeStart = 0
	resource1.RangeEnd = 10

	// Create IPAddressPool resource
	resource2 := new(network.IPAddressPool)
	resource2.ID = "0200000001"
	ip := net.ParseIP("10.0.0.1")
	mask := ip.DefaultMask()
	resource2.Subnets = []net.IPNet{{IP: ip, Mask: mask}}

	// Add resources to map
	resources := make(map[string]zebra.Resource)
	resources["0100000001"] = resource1
	resources["0200000001"] = resource2

	qs := query.NewQueryStore(resources)
	assert.NotNil(qs)

	all := qs.Query()
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
	resource1.RangeStart = 0
	resource1.RangeEnd = 10

	// Create IPAddressPool resource
	resource2 := new(network.IPAddressPool)
	resource2.ID = "0200000001"
	ip := net.ParseIP("10.0.0.1")
	mask := ip.DefaultMask()
	resource2.Subnets = []net.IPNet{{IP: ip, Mask: mask}}

	// Add resources to map
	resources := make(map[string]zebra.Resource)
	resources["0100000001"] = resource1
	resources["0200000001"] = resource2

	querystore := query.NewQueryStore(resources)
	assert.NotNil(querystore)

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
	resource1.RangeStart = 0
	resource1.RangeEnd = 10

	// Create IPAddressPool resource
	resource2 := new(network.IPAddressPool)
	resource2.ID = "0200000001"
	ip := net.ParseIP("10.0.0.1")
	mask := ip.DefaultMask()
	resource2.Subnets = []net.IPNet{{IP: ip, Mask: mask}}

	// Add resources to map
	resources := make(map[string]zebra.Resource)
	resources["0100000001"] = resource1
	resources["0200000001"] = resource2

	querystore := query.NewQueryStore(resources)
	assert.NotNil(querystore)

	vlanpools := querystore.QueryType("VLANPool")
	assert.True(len(vlanpools) == 1)
	assert.True(vlanpools[0].GetID() == "0100000001")

	ippools := querystore.QueryType("IPAddressPool")
	assert.True(len(ippools) == 1)
	assert.True(ippools[0].GetID() == "0200000001")
}

func TestInvalidLabelQuery(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// Create VLANPool resource
	resource1 := new(network.VLANPool)
	resource1.ID = "0100000001"
	resource1.Labels = make(map[string]string)
	resource1.Labels["product-owner"] = "owner"
	resource1.RangeStart = 0
	resource1.RangeEnd = 10

	// Add resources to map
	resources := make(map[string]zebra.Resource)
	resources["0100000001"] = resource1

	querystore := query.NewQueryStore(resources)
	assert.NotNil(querystore)

	query1 := query.LabelQuery{
		Op:     7,
		Key:    "",
		Values: nil,
	}

	// Should fail on invalid query.
	_, err := querystore.QueryLabelsMatchAll([]query.LabelQuery{query1})
	assert.NotNil(err)

	// Should fail on invalid query.
	_, err = querystore.QueryLabelsMatchOne([]query.LabelQuery{query1})
	assert.NotNil(err)
}

func TestQueryLabelMatchAll(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// Create VLANPool resource
	resource1 := new(network.VLANPool)
	resource1.ID = "0100000001"
	resource1.Labels = make(map[string]string)
	resource1.Labels["product-owner"] = "shravya"
	resource1.RangeStart = 0
	resource1.RangeEnd = 10

	// Create IPAddressPool resource
	resource2 := new(network.IPAddressPool)
	resource2.ID = "0200000001"
	resource2.Labels = make(map[string]string)
	resource2.Labels["product-owner"] = "nandyala"
	resource2.Labels["team"] = "cloud networking"
	ip := net.ParseIP("10.0.0.1")
	mask := ip.DefaultMask()
	resource2.Subnets = []net.IPNet{{IP: ip, Mask: mask}}

	// Add resources to map
	resources := make(map[string]zebra.Resource)
	resources["0100000001"] = resource1
	resources["0200000001"] = resource2

	querystore := query.NewQueryStore(resources)
	assert.NotNil(querystore)

	query1 := query.LabelQuery{Op: query.MatchEqual, Key: "product-owner", Values: []string{"shravya", "nandyala"}}
	query2 := query.LabelQuery{Op: query.MatchIn, Key: "product-owner", Values: []string{"shravya", "nandyala"}}
	query3 := query.LabelQuery{Op: query.MatchNotEqual, Key: "product-owner", Values: []string{"shravya", "nandyala"}}
	query4 := query.LabelQuery{Op: query.MatchNotIn, Key: "product-owner", Values: []string{"shravya", "nandyala"}}
	query5 := query.LabelQuery{Op: query.MatchEqual, Key: "product-owner", Values: []string{"shravya"}}

	// Should fail on query 1 and query 3.
	_, err := querystore.QueryLabelsMatchAll([]query.LabelQuery{query1})
	assert.NotNil(err)

	_, err = querystore.QueryLabelsMatchAll([]query.LabelQuery{query3})
	assert.NotNil(err)

	// Should succeed on query 2, return both resources.
	pos, err := querystore.QueryLabelsMatchAll([]query.LabelQuery{query2})
	assert.True(err == nil && len(pos) == 2 && pos[0] != pos[1])

	// Should succeed on query 4, return no resources.
	pos, err = querystore.QueryLabelsMatchAll([]query.LabelQuery{query4})
	assert.True(err == nil && len(pos) == 0)

	// Should succeed on queries 2 and 5, return one resource.
	pos, err = querystore.QueryLabelsMatchAll([]query.LabelQuery{query2, query5})
	assert.True(err == nil && len(pos) == 1 && pos[0].GetID() == "0100000001")

	// Update query 3 to be valid, return 1 resource.
	query3.Values = []string{"shravya"}
	pos, err = querystore.QueryLabelsMatchAll([]query.LabelQuery{query3})
	assert.True(err == nil && len(pos) == 1 && pos[0].GetID() == "0200000001")
}

func TestQueryLabelMatchOne(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// Create VLANPool resource
	resource1 := new(network.VLANPool)
	resource1.ID = "0100000001"
	resource1.Labels = make(map[string]string)
	resource1.Labels["product-owner"] = "shravya"
	resource1.RangeStart = 0
	resource1.RangeEnd = 10

	// Create IPAddressPool resource
	resource2 := new(network.IPAddressPool)
	resource2.ID = "0200000001"
	resource2.Labels = make(map[string]string)
	resource2.Labels["product-owner"] = "nandyala"
	resource2.Labels["team"] = "cloud networking"
	ip := net.ParseIP("10.0.0.1")
	mask := ip.DefaultMask()
	resource2.Subnets = []net.IPNet{{IP: ip, Mask: mask}}

	// Add resources to map
	resources := make(map[string]zebra.Resource)
	resources["0100000001"] = resource1
	resources["0200000001"] = resource2

	querystore := query.NewQueryStore(resources)
	assert.NotNil(querystore)

	query1 := query.LabelQuery{Op: query.MatchIn, Key: "product-owner", Values: []string{"shravya", "nandyala"}}
	query2 := query.LabelQuery{Op: query.MatchNotIn, Key: "product-owner", Values: []string{"shravya", "nandyala"}}
	query3 := query.LabelQuery{Op: query.MatchEqual, Key: "product-owner", Values: []string{"shravya"}}
	query4 := query.LabelQuery{Op: query.MatchEqual, Key: "team", Values: []string{"cloud networking"}}

	// Should succeed on query 1, return both resources.
	pos, err := querystore.QueryLabelsMatchOne([]query.LabelQuery{query1})
	assert.Nil(err)
	assert.True(len(pos) == 2 && pos[0] != pos[1])

	// Should succeed on query 2, return no resources.
	pos, err = querystore.QueryLabelsMatchOne([]query.LabelQuery{query2})
	assert.Nil(err)
	assert.True(len(pos) == 0)

	// Should succeed on queries 1 and 3, return both resources.
	pos, err = querystore.QueryLabelsMatchOne([]query.LabelQuery{query1, query3})
	assert.Nil(err)
	assert.True(len(pos) == 2 && pos[0] != pos[1])

	// Should succeed on queries 3 and 4, return both resources.
	pos, err = querystore.QueryLabelsMatchOne([]query.LabelQuery{query3, query4})
	assert.Nil(err)
	assert.True(len(pos) == 2 && pos[0] != pos[1])
}
