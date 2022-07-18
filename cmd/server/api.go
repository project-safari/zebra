package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/store"
)

type ResourceAPI struct {
	factory zebra.ResourceFactory
	Store   zebra.Store
}

type QueryRequest struct {
	IDs        []string      `json:"ids,omitempty"`
	Types      []string      `json:"types,omitempty"`
	Labels     []zebra.Query `json:"labels,omitempty"`
	Properties []zebra.Query `json:"properties,omitempty"`
}

var ErrQueryRequest = errors.New("invalid GET query request body")

func handleQuery(ctx context.Context, api *ResourceAPI) httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		qr := new(QueryRequest)

		// Read request, return error if applicable
		if err := readReq(ctx, req, qr); err != nil {
			res.WriteHeader(http.StatusBadRequest)

			return
		}

		// Validate query request and label/property queries
		if err := qr.Validate(ctx); err != nil {
			res.WriteHeader(http.StatusBadRequest)

			return
		}

		var resources *zebra.ResourceMap

		// Get resources based on primary key (ID, Type, or Label)
		switch {
		case len(qr.IDs) != 0:
			resources = api.Store.QueryUUID(qr.IDs)
		case len(qr.Types) != 0:
			resources = api.Store.QueryType(qr.Types)
		case len(qr.Labels) != 0:
			q := qr.Labels[0]
			qr.Labels = qr.Labels[1:]
			// Can safely ignore error because we have already validated the query
			resources, _ = api.Store.QueryLabel(q)
		default:
			resources = api.Store.Query()
		}

		// Filter further based on label queries
		for _, q := range qr.Labels {
			// Can safely ignore error because we have already validated the query
			resources, _ = store.FilterLabel(q, resources)
		}

		// Write response body
		writeJSON(ctx, res, resources)
	}
}

func (qr *QueryRequest) Validate(ctx context.Context) error {
	id := len(qr.IDs) != 0
	t := len(qr.Types) != 0
	l := len(qr.Labels) != 0
	p := len(qr.Properties) != 0

	// Make sure only id (and labels), types (and labels), or labels are present
	if (id && (t || p)) || (p && (t || l)) {
		return ErrQueryRequest
	}

	// Check Labels queries are valid
	if err := checkQueries(qr.Labels); err != nil {
		return err
	}

	// Check Properties queries are valid
	return checkQueries(qr.Properties)
}

func checkQueries(queries []zebra.Query) error {
	for _, q := range queries {
		if err := q.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func handlePut(ctx context.Context, api *ResourceAPI) httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		resMap := zebra.NewResourceMap(store.DefaultFactory())

		// Read request, return error if applicable
		if err := readReq(ctx, req, resMap); err != nil {
			res.WriteHeader(http.StatusBadRequest)

			return
		}

		// Check all resources to make sure they are valid
		for _, l := range resMap.Resources {
			for _, r := range l.Resources {
				if r.Validate(ctx) != nil {
					res.WriteHeader(http.StatusBadRequest)

					return
				}
			}
		}

		// Add all resources to store
		for _, l := range resMap.Resources {
			for _, r := range l.Resources {
				if api.Store.Create(r) != nil {
					res.WriteHeader(http.StatusInternalServerError)

					return
				}
			}
		}

		// Return 200
		res.WriteHeader(http.StatusOK)
	}
}

func handleDelete(ctx context.Context, api *ResourceAPI) httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		resMap := zebra.NewResourceMap(store.DefaultFactory())

		// Read request, return error if applicable
		if err := readReq(ctx, req, resMap); err != nil {
			res.WriteHeader(http.StatusBadRequest)

			return
		}

		// Check all resources to make sure they are valid
		for _, l := range resMap.Resources {
			for _, r := range l.Resources {
				if r.Validate(ctx) != nil {
					res.WriteHeader(http.StatusBadRequest)

					return
				}
			}
		}

		// Delete all resources from store
		for _, l := range resMap.Resources {
			for _, r := range l.Resources {
				if api.Store.Delete(r) != nil {
					res.WriteHeader(http.StatusInternalServerError)

					return
				}
			}
		}

		// Return 200
		res.WriteHeader(http.StatusOK)
	}
}

func NewResourceAPI(factory zebra.ResourceFactory) *ResourceAPI {
	return &ResourceAPI{
		factory: factory,
		Store:   nil,
	}
}

// Set up store and query store given storage root.
func (api *ResourceAPI) Initialize(storageRoot string) error {
	api.Store = store.NewResourceStore(storageRoot, api.factory)

	return api.Store.Initialize()
}
