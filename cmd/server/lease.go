package main

import (
	"net/http"
	"time"

	"github.com/go-logr/logr"
	"github.com/julienschmidt/httprouter"
	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/model/lease"
	"github.com/project-safari/zebra/store"
)

type ResourceReq struct {
	Type         string        `json:"type"`
	Group        string        `json:"group"`
	Name         string        `json:"name"`
	Count        int           `json:"count"`
	Filters      []zebra.Query `json:"filters,omitempty"`
	ResourcesIds []string      `json:"resources,omitempty"`
}

func handleLease() httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		ctx := req.Context()
		log := logr.FromContextOrDiscard(ctx)
		api, ok := ctx.Value(ResourcesCtxKey).(*ResourceAPI)

		if !ok {
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		store := api.Store

		leaseReq := &struct {
			Email    string         `json:"email"`
			Duration time.Duration  `json:"duration"`
			Request  []*ResourceReq `json:"request"`
		}{}

		if err := readJSON(ctx, req, leaseReq); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			log.Info("lease could not be create, request was not valid")

			return
		}

		leasereqs := resourceList(api, leaseReq.Request)

		newlease := lease.NewLease(leaseReq.Email, leaseReq.Duration, leasereqs)

		if err := newlease.Validate(ctx); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			log.Info("lease could not be create, request was invalid")

			return
		}

		if err := store.Create(newlease); err != nil {
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		log.Info("successfully created lease")

		res.WriteHeader(http.StatusCreated)
	}
}

func resourceList(api *ResourceAPI, requests []*ResourceReq) []*lease.ResourceReq {
	leaseReqs := make([]*lease.ResourceReq, 0)

	resMap := api.Store.Query()
	idStore := store.NewIDStore(resMap)

	for _, rr := range requests {
		reqMap := idStore.Query(rr.ResourcesIds)
		resources := make([]zebra.Resource, 0)

		for _, resList := range reqMap.Resources {
			resources = append(resources, resList.Resources...)
		}

		leaseReq := &lease.ResourceReq{
			Type:      rr.Type,
			Group:     rr.Group,
			Name:      rr.Name,
			Count:     rr.Count,
			Filters:   rr.Filters,
			Resources: resources,
		}
		leaseReqs = append(leaseReqs, leaseReq)
	}

	return leaseReqs
}
