package lease

import (
	"github.com/edwingeng/deque"
	"github.com/project-safari/zebra"
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

		q, ok := q.Queues[key]
		if !ok {
			q = deque.NewDeque()
		}

		q.Enqueue(req)
	}
}

func (q *Queue) Process() {
	for _, sq := range q.Queues {
		go process(q, sq)
	}
}

func process(lq *Queue, sq deque.Deque) {
	for {
		for !sq.Empty() {
			req, _ := sq.Peek(0).(*ResourceReq) // we know type for sure

			query := zebra.Query{Key: "system.group", Op: zebra.MatchEqual, Values: []string{req.Group}}

			resources, err := zebra.FilterLabel(query, lq.Store.QueryType([]string{req.Type}))
			if err != nil {
				panic(err)
			}

			freePool, ok := resources.Resources[req.Type]
			// No resources of this (group, type) right now, try again.
			if !ok {
				continue
			}

			// Try to assign as many resources as possible to reach count.
			count := 0
			for i := 0; i < len(freePool.Resources) && count < req.Count; i++ {
				res := freePool.Resources[i]

				// Try to lease resource
				if req.Assign(res) == nil {
					count++
				}
			}

			// If count was met, remove from queue. Else, try again.
			if count == req.Count {
				_ = sq.Dequeue() // throw away result, we don't care about it
			}
		}
	}
}
