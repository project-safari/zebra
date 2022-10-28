package lease_test

import (
	"net"
	"testing"
	"time"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/model"
	"github.com/project-safari/zebra/model/compute"
	"github.com/project-safari/zebra/model/lease"
	"github.com/project-safari/zebra/store"
	"github.com/stretchr/testify/assert"
)

func createLease(user, dur string) *lease.Lease {
	d, err := time.ParseDuration(dur)
	if err != nil {
		return nil
	}

	rr := []*lease.ResourceReq{
		{
			Type:  "compute.server",
			Group: "san-jose-building-14",
			Name:  "linux blah blah",
			Count: 1,
		},
		{
			Type:  "compute.vm",
			Group: "san-jose-building-18",
			Name:  "virtual",
			Count: 1,
		},
	}

	lease := lease.NewLease(user, d, rr)

	return lease
}

func makeStore(root string, factory zebra.ResourceFactory, assert *assert.Assertions) *store.ResourceStore {
	store := store.NewResourceStore(root, factory)
	err := store.Initialize()
	assert.Nil(err)

	server1 := compute.NewServer("1234", "new", "testserver1", "newuser@ciso.com", "san-jose-building-14")
	server1.BoardIP = net.IP{1, 1, 1, 1}
	server1.Credentials = zebra.NewCredentials("tester")
	assert.Nil(server1.Credentials.Add("password", "thisIsAGoodPassword!123"))

	server2 := compute.NewServer("1235", "new", "testserver2", "newuser@ciso.com", "san-jose-building-14")
	server2.BoardIP = net.IP{1, 1, 1, 1}
	server2.Credentials = zebra.NewCredentials("tester")
	assert.Nil(server2.Credentials.Add("password", "thisIsAGoodPassword!123"))

	vm1 := compute.NewVM("newesx", "testvm1", "newuser@cisco.com", "san-jose-building-18")
	vm1.ManagementIP = net.IP{1, 1, 1, 1}
	vm1.Credentials = zebra.NewCredentials("tester")
	assert.Nil(vm1.Credentials.Add("password", "thisIsAGoodPassword!123"))
	assert.Nil(store.Create(server1))
	assert.Nil(store.Create(server2))
	assert.Nil(store.Create(vm1))

	return store
}

func TestNewQueue(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "test_new_queue"

	factory := model.Factory()
	assert.NotNil(factory)

	store := store.NewResourceStore(root, factory)
	err := store.Initialize()
	assert.Nil(err)

	q := lease.NewQueue(store)
	assert.NotNil(q)
}

func TestProcess(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	root := "test_process"

	factory := model.Factory()
	assert.NotNil(factory)

	store := makeStore(root, factory, assert)
	q := lease.NewQueue(store)
	assert.NotNil(q)

	key1 := lease.Key{
		Type:  "compute.server",
		Group: "san-jose-building-14",
	}

	key2 := lease.Key{
		Type:  "compute.vm",
		Group: "san-jose-building-18",
	}

	q.Enqueue(createLease("testuser@cisco.com", "3h"))
	assert.NotEmpty(q.Queues)

	q.Enqueue(createLease("testuser1@cisco.com", "4h"))
	q.Process()

	assert.Empty(q.Queues[key1].Len())
	assert.NotEmpty(q.Queues[key2].Len())
}
