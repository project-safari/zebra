package lease

import (
	"github.com/edwingeng/deque"
	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/store"
)

type Queue struct {
	Store  zebra.Store         `json:"store"`
	Queues map[Key]deque.Deque `json:"queues"`
}

type Key struct {
	Group string `json:"group"`
	Type  string `json:"type"`
}

func NewQueue(store zebra.Store) *Queue {
	return &Queue{
		Store:  store,
		Queues: make(map[Key]deque.Deque),
	}
}

func (q *Queue) Enqueue(l *Lease) {
	for _, req := range l.RequestList() {
		key := Key{
			Group: req.Group,
			Type:  req.Type,
		}

		queue, ok := q.Queues[key]
		if !ok {
			queue = deque.NewDeque()
		}

		queue.Enqueue(req)
		q.Queues[key] = queue
	}
}

func (q *Queue) Process() {
	for _, sq := range q.Queues {
		req, _ := sq.Peek(0).(*ResourceReq)
		query := zebra.Query{Key: "system.group", Op: zebra.MatchEqual, Values: []string{req.Group}}
		resmap := q.Store.QueryType([]string{req.Type})

		resources, err := store.FilterLabel(query, resmap)
		if err != nil {
			panic(err)
		}

		freePool, ok := resources.Resources[req.Type]
		// No resources of this (group, type) right now, try again.
		if !ok {
			continue
		}

		qch := make(chan deque.Deque)

		go process(freePool.Resources, sq, qch, req)
		sq = <-qch //nolint:staticcheck
	}
}

func process(freePool []zebra.Resource, sq deque.Deque, qch chan deque.Deque, req *ResourceReq) {
	for {
		for !sq.Empty() && len(freePool) != 0 {
			// Try to assign as many resources as possible to reach count.
			count := 0
			for i := 0; i < len(freePool) && count < req.Count; i++ {
				res := freePool[0]
				// Try to lease resource
				if req.Assign(res) == nil && res.GetStatus().LeaseStatus == 0 {
					count++
				}

				freePool = freePool[1:]
			}
			// If count was met, remove from queue. Else, try again.
			if count == req.Count {
				_ = sq.Dequeue() // throw away result, we don't care about it
			}
		}
		qch <- sq
	}
}
