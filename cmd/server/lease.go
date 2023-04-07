package main

import (
	"net/http"
	"time"

	"github.com/go-logr/logr"
	"github.com/julienschmidt/httprouter"
	"github.com/project-safari/zebra/model/lease"
)

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
			Email    string               `json:"email"`
			Duration time.Duration        `json:"duration"`
			Request  []*lease.ResourceReq `json:"request"`
		}{}

		if err := readJSON(ctx, req, leaseReq); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			log.Info("lease could not be create, request was not valid")

			return
		}

		newlease := lease.NewLease(leaseReq.Email, leaseReq.Duration, leaseReq.Request)

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

		res.WriteHeader(http.StatusOK)
	}
}
