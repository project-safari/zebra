package lease_test

import (
	"testing"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/lease"
	"github.com/stretchr/testify/assert"
)

func TestNewQueue(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	g := zebra.NewGroup("dev")
	q := lease.NewQueue(g)
	assert.NotNil(q)
}

func TestQueueAdd(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	g := zebra.NewGroup("dev")
	q := lease.NewQueue(g)
	assert.NotNil(q)

	q.Enqueue(new(lease.Lease))
	assert.Equal(1, q.Size())
}

func TestQueueRemove(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	g := zebra.NewGroup("dev")
	q := lease.NewQueue(g)
	assert.NotNil(q)

	q.Enqueue(new(lease.Lease))
	assert.Equal(1, q.Size())

	val, err := q.Dequeue()
	assert.Nil(err)
	assert.NotNil(val)
	assert.Equal(0, q.Size())

	_, err = q.Dequeue()
	assert.NotNil(err)
}

func TestQueueEmpty(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	g := zebra.NewGroup("dev")
	q := lease.NewQueue(g)
	assert.NotNil(q)
	assert.True(q.IsEmpty())
}
