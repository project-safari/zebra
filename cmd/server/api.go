package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-logr/logr"
	"github.com/julienschmidt/httprouter"
	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/store"
)

// ResourceAPI is a struct that contains resources for the apis.
type ResourceAPI struct {
	factory zebra.ResourceFactory
	Store   zebra.Store
}

// QueryRequest struct is a struct that contains data for possible requests.
type QueryRequest struct {
	IDs        []string      `json:"ids,omitempty"`
	Types      []string      `json:"types,omitempty"`
	Labels     []zebra.Query `json:"labels,omitempty"`
	Properties []zebra.Query `json:"properties,omitempty"`
}

// ErrQueryRequest is an error that occurs if the GET request is invalid.
var ErrQueryRequest = errors.New("invalid GET query request body")

// Validate function for the query request.
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
	if err := validateQueries(qr.Labels); err != nil {
		return err
	}

	// Check Properties queries are valid
	return validateQueries(qr.Properties)
}

// Function that returns a new struct with resources for the API.
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

// Apply given function f to each resource in resMap.
// Return error if it occurrs or nil if successful.
func applyFunc(resMap *zebra.ResourceMap, f func(zebra.Resource) error) error {
	for _, l := range resMap.Resources {
		for _, r := range l.Resources {
			if err := f(r); err != nil {
				return err
			}
		}
	}

	return nil
}

// Validate all queries in given slice.
func validateQueries(queries []zebra.Query) error {
	for _, q := range queries {
		if err := q.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// Function thathandles a query and retursn a httprouter.Handle.
func handleQuery() httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		ctx := req.Context()
		log := logr.FromContextOrDiscard(ctx)
		api, ok := ctx.Value(ResourcesCtxKey).(*ResourceAPI)

		if !ok {
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		qr := new(QueryRequest)

		// Read request, return error if applicable
		if err := readJSON(ctx, req, qr); err != nil && !errors.Is(err, ErrEmptyBody) {
			res.WriteHeader(http.StatusBadRequest)
			log.Info("resources could not be queried, could not read request")

			return
		}

		// Validate query request and label/property queries
		if err := qr.Validate(ctx); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			log.Info("resources could not be queried, found invalid quer(y/ies)")

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

		log.Info("successfully queried resources")

		// Write response body
		writeJSON(ctx, res, resources)
	}
}

// Function that handles a delete nad returns a httprouter.Handle.
func handleDelete() httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		ctx := req.Context()
		log := logr.FromContextOrDiscard(ctx)
		api, ok := ctx.Value(ResourcesCtxKey).(*ResourceAPI)

		if !ok {
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		idReq := &struct {
			IDs []string `json:"ids"`
		}{}
		resMap := api.Store.Query()

		// Read request, return error if applicable
		if err := readJSON(ctx, req, idReq); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			log.Info("resources could not be deleted, could not read request")

			return
		}

		idStore := store.NewIDStore(resMap)
		resMap = idStore.Query(idReq.IDs)

		if len(idReq.IDs) == 0 {
			res.WriteHeader(http.StatusBadRequest)
			log.Info("resources could not be deleted, found invalid resource(s)")

			return
		}

		// Delete all resources from store
		if applyFunc(resMap, api.Store.Delete) != nil {
			res.WriteHeader(http.StatusInternalServerError)
			log.Info("internal server error while deleting resources")

			return
		}

		log.Info("successfully deleted resources")

		res.WriteHeader(http.StatusOK)
	}
}
