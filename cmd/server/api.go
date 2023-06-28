package main

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/go-logr/logr"
	"github.com/julienschmidt/httprouter"
	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/model"
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

		resources := api.Store.Query()
		// Get resources based on primary key (ID, Type, or Label)
		// Label Query currently doesn't work using url parameters
		if parames := req.URL.Query(); len(parames) != 0 {
			switch {
			case len(parames["ids"]) != 0:
				resources = api.Store.QueryUUID(strings.Split(parames["ids"][0], ","))
			case len(parames["types"]) != 0:
				resources = api.Store.QueryType(strings.Split(parames["types"][0], ","))
			}
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

		resMap := zebra.NewResourceMap(model.Factory())

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

		id := params.ByName("id")
		if id == "" {
			res.WriteHeader(http.StatusBadRequest)
			log.Info("resources could not be deleted, found invalid resource(s)")

			return
		}

		ids := []string{id}
		resMap := api.Store.QueryUUID(ids)

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
