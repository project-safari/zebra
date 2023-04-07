package store

import (
	"time"

	"github.com/edwingeng/deque"
	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/model/lease"
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

func (q *Queue) Enqueue(l *lease.Lease) {
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
		if sq.Len() < 1 {
			continue
		}

		req, _ := sq.Peek(0).(*lease.ResourceReq)
		query := zebra.Query{Key: "system.group", Op: zebra.MatchEqual, Values: []string{req.Group}}
		resmap := q.Store.QueryType([]string{req.Type})

		resources, err := FilterLabel(query, resmap)
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

func process(freePool []zebra.Resource, sq deque.Deque, qch chan deque.Deque, req *lease.ResourceReq) {
	for {
		for !sq.Empty() && len(freePool) != 0 {
			// Try to assign as many resources as possible to reach count.
			for i := 0; i < len(freePool) && len(req.Resources) < req.Count; i++ {
				res := freePool[0]
				// Try to lease resource
				if res.GetStatus().LeaseStatus == 0 && req.Assign(res) == nil {
					res.UpdateStatus().UpdateLeaseState(1)
				}

				freePool = freePool[1:]
			}
			// If count was met, remove from queue. Else, try again.
			if len(req.Resources) == req.Count {
				_ = sq.Dequeue() // throw away result, we don't care about it

				if sq.Len() < 1 {
					continue
				}

				req, _ = sq.Peek(0).(*lease.ResourceReq)
			}
		}
		qch <- sq
	}
}

func (q *Queue) LeaseSatisfied() {
	resMap := q.Store.QueryType([]string{"system.lease"})
	for k, v := range resMap.Resources {
		if k == "system.lease" {
			for _, leslist := range v.Resources {
				if l, ok := leslist.(*lease.Lease); ok {
					if l.IsSatisfied() && l.Status.State == zebra.Inactive {
						if l.Activate() == nil {
							go q.isLeaseActive(l)

							continue
						}
					}
				}
			}
		}
	}
}

func (q *Queue) isLeaseActive(lease *lease.Lease) {
	for !lease.IsExpired() {
		time.Sleep(1 * time.Minute)
	}

	lease.Deactivate()

	for _, reslist := range lease.RequestList() {
		list := reslist.Resources
		reslist.Resources = nil

		q.Store.FreeResources(list)
	}
}
