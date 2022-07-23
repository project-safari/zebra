package lease

import (
	"errors"
	"sync"

	"github.com/edwingeng/deque"
	"github.com/project-safari/zebra"
)

type Queue struct {
	lock  sync.RWMutex
	group *zebra.Group
	queue deque.Deque
}

var ErrEmptyQueue = errors.New("tried to remove from an empty queue")

func NewQueue(g *zebra.Group) *Queue {
	return &Queue{
		lock:  sync.RWMutex{},
		group: g,
		queue: deque.NewDeque(),
	}
}

func (q *Queue) Enqueue(l *Lease) {
	q.lock.Lock()
	defer q.lock.Unlock()

	// Really, we should check to make sure lease has this group
	// But for now, move on

	q.queue.Enqueue(l)
}

func (q *Queue) Dequeue() (*Lease, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.IsEmpty() {
		return nil, ErrEmptyQueue
	}

	val, _ := q.queue.Dequeue().(*Lease)

	return val, nil
}

func (q *Queue) Size() int {
	return q.queue.Len()
}

func (q *Queue) IsEmpty() bool {
	return q.queue.Empty()
}

// Process lease at front
// Right now, this only assigns whatever is available and then releases from queue, need to fix this
// Idea: move lease to "waiting" area and then whenever we release resources from an expired lease,
// try to process waiting leases again.
func (q *Queue) Process() error {
	l, err := q.Dequeue()
	if err != nil {
		return err
	}

	// Try to process lease request for this queue's group
	request := l.Request()[q.group]

	// Does this group even have resources of this type available?
	resList, ok := q.group.FreePool.Resources[request.Type]
	if !ok {
		// Just exit without assigning any resources
		return nil
	}

	// Try to assign as many resources as possible
	end := int(request.Count)
	if end < len(resList.Resources) {
		end = len(resList.Resources)
	}

	for i := 0; i < end; i++ {
		l.Assign(resList.Resources[i])
	}

	// Update progress of lease
	l.Finish(q.group)

	return nil
}

// Naive implementation of wait time index
// We should store these somehow, maybe an overall wait time for the queue,
// that could be a property of the queue. But what if we want wait time for individual
// lease, then it's tricky.

/*
func (q *Queue) WaitTime() int {
	index := 0

	for i := 0; i < q.queue.Len(); i++ {
		l := q.queue.Peek(i).(*Lease)

		for _, req := range l.Request {
			index += int(req.Count)
		}
	}

	return index / int(q.group.Count)
}
*/
