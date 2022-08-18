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

// Resource api is a struct that contains a fzebra.ResourceFactory and a zebra.Store.
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

// ErrQueryRequest returns an error message if the GET request has an invalid body.
var ErrQueryRequest = errors.New("invalid GET query request body")

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

// Validate all resources in a resource map.
func validateResources(ctx context.Context, resMap *zebra.ResourceMap) error {
	// Check all resources to make sure they are valid
	for _, l := range resMap.Resources {
		for _, r := range l.Resources {
			if err := r.Validate(ctx); err != nil {
				return err
			}
		}
	}

	return nil
}

// Function to handle querries.
// This function helps log the info and write the response in a json file, or return if an error occurs.
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

func handlePost() httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		ctx := req.Context()
		log := logr.FromContextOrDiscard(ctx)
		api, ok := ctx.Value(ResourcesCtxKey).(*ResourceAPI)

		if !ok {
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		resMap := zebra.NewResourceMap(store.DefaultFactory())

		// Read request, return error if applicable
		if err := readJSON(ctx, req, resMap); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			log.Info("resources could not be created, could not read request")

			return
		}

		if validateResources(ctx, resMap) != nil {
			res.WriteHeader(http.StatusBadRequest)
			log.Info("resources could not be created, found invalid resource(s)")

			return
		}

		// Add all resources to store
		if applyFunc(resMap, api.Store.Create) != nil {
			res.WriteHeader(http.StatusInternalServerError)
			log.Info("internal server error while creating resources")

			return
		}

		log.Info("successfully created resources")

		res.WriteHeader(http.StatusOK)
	}
}

func handleDelete() httprouter.Handle {
	return func(res http.ResponseWriter, req *http.Request, params httprouter.Params) {
		ctx := req.Context()
		log := logr.FromContextOrDiscard(ctx)
		api, ok := ctx.Value(ResourcesCtxKey).(*ResourceAPI)

		if !ok {
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		resMap := zebra.NewResourceMap(store.DefaultFactory())

		// Read request, return error if applicable
		if err := readJSON(ctx, req, resMap); err != nil {
			res.WriteHeader(http.StatusBadRequest)
			log.Info("resources could not be deleted, could not read request")

			return
		}

		if validateResources(ctx, resMap) != nil {
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
