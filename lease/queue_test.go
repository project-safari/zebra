package lease //nolint:testpackage

import (
	"testing"

	"github.com/project-safari/zebra"
	"github.com/stretchr/testify/assert"
)

func TestQueue(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	q := getQueue()
	assert.NotNil(q)

	req := &ResourceReq{
		Type:      "NamedResource",
		Group:     "test-group",
		Name:      "name",
		Count:     1,
		Filters:   nil,
		Resources: nil,
	}

	q.Enqueue(req)
	assert.Equal(req, q.Dequeue())

	assert.True(q.Empty())
}

func TestQueueProcess(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	q := getQueue()
	assert.NotNil(q)

	q.Enqueue(getReq())
	q.Enqueue(getReq())

	stalled := make(chan bool)

	go q.Process(stalled)

	ret := <-stalled
	assert.Equal(true, ret)

	res := zebra.NewNamedResource("name", "NamedResource",
		map[string]string{"system.group": "test-group"})
	q.Group.Add(res)

	stalled <- false
}

func getQueue() *Queue {
	res := zebra.NewNamedResource("name", "NamedResource",
		map[string]string{"system.group": "test-group"})
	g := zebra.NewGroup("test-group")
	g.Add(res)

	resType := zebra.Type{
		Name:        "NamedResource",
		Description: "",
		Constructor: func() zebra.Resource {
			return zebra.NewNamedResource("name", "NamedResource", nil)
		},
	}

	return NewQueue(g, resType)
}

func getReq() *ResourceReq {
	return &ResourceReq{
		Type:      "NamedResource",
		Group:     "test-group",
		Name:      "name",
		Count:     1,
		Filters:   []zebra.Query{},
		Resources: nil,
	}
}
