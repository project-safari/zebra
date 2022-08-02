package lease

import (
	"sync"

	"github.com/edwingeng/deque"
	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/store"
)

type Queue struct {
	lock  sync.RWMutex
	Group *zebra.Group `json:"group"`
	Type  zebra.Type   `json:"type"`
	Queue deque.Deque  `json:"queue"`
}

func NewQueue(g *zebra.Group, t zebra.Type) *Queue {
	return &Queue{
		lock:  sync.RWMutex{},
		Group: g,
		Type:  t,
		Queue: deque.NewDeque(),
	}
}

func (q *Queue) Enqueue(r *ResourceReq) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.Queue.Enqueue(r)
}

func (q *Queue) Dequeue() *ResourceReq {
	q.lock.Lock()
	defer q.lock.Unlock()

	val, ok := q.Queue.Dequeue().(*ResourceReq)
	if !ok {
		panic(ok)
	}

	return val
}

func (q *Queue) Empty() bool {
	q.lock.RLock()
	defer q.lock.RUnlock()

	return q.Queue.Empty()
}

// Process queue.
func (q *Queue) Process(stalled chan bool) {
	for {
		if !q.Empty() {
			req := q.Dequeue()

			q.lock.Lock()
			freePool := q.getFreePool(req)

			// Now we have a group of eligible resources to assign from.
			if freePool == nil || len(freePool.Resources) < req.Count {
				// Try to release inactive leases' resources.
				// For now, should just return stalled in the channel
				// Unlock queue so more stuff can be added in the meanwhile
				q.lock.Unlock()
				stalled <- true
				// Wait until resources have been released into the free pool
				<-stalled
				q.lock.Lock()
				// Refresh free pool to have latest resources
				freePool = q.getFreePool(req)
			}

			// Continue when there are free resources to use.
			for i := 0; i < req.Count; i++ {
				res := freePool.Resources[i]
				if err := q.Group.Lease(res); err != nil {
					panic(err)
				}
				// Assign resource that was just leased to current request
				_ = req.Assign(res)
				i++
			}

			q.lock.Unlock()
		}
	}
}

// Should be called with write lock.
func (q *Queue) getFreePool(req *ResourceReq) *zebra.ResourceList {
	freePool, err := store.FilterType([]string{req.Type}, q.Group.FreePool())
	if err != nil {
		panic(err)
	}

	query := zebra.Query{Op: zebra.MatchEqual, Key: "system.group", Values: []string{req.Group}}
	if freePool, err = store.FilterLabel(query, freePool); err != nil {
		panic(err)
	}

	query = zebra.Query{Op: zebra.MatchEqual, Key: "Name", Values: []string{req.Name}}

	// Add name query to list of filters
	filters := make([]zebra.Query, len(req.Filters)+1)
	copy(filters, req.Filters)
	filters[len(filters)-1] = query

	// Iterate through and apply all filters
	for _, f := range filters {
		if freePool, err = store.FilterProperty(f, freePool); err != nil {
			panic(err)
		}
	}

	return freePool.Resources[req.Type]
}
