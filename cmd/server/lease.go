package main

import (
	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/lease"
)

type LeaseQueue map[*zebra.Group]*lease.Queue

// Initialize lease queue with resources from store.
func (api *ResourceAPI) InitializeLeaseQueue() {
	groups := make(map[string]*zebra.Group)

	api.LeaseQueue = make(LeaseQueue)

	// Go through all resources from store, add to proper group
	for _, l := range api.Store.Query().Resources {
		for _, r := range l.Resources {
			gName := r.GetLabels()["group"]

			// If group hasn't been created, create and add lease slice value
			g, ok := groups[gName]
			if !ok {
				g = zebra.NewGroup(gName)
				api.LeaseQueue[g] = lease.NewQueue(g)
				groups[gName] = g
			}

			g.Add(r)
		}
	}
}

// Add lease to all necessary group queues.
func (api *ResourceAPI) AddLease(l *lease.Lease) {
	for g := range l.Request() {
		api.LeaseQueue[g].Enqueue(l)
	}
}
