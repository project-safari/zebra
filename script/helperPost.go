package script

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/go-logr/logr"
	"github.com/julienschmidt/httprouter"
	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/model"
)

var ErrEmptyBody = errors.New("invalid GET query request body")

func ReadJSON(ctx context.Context, req *http.Request, data interface{}) error {
	log := logr.FromContextOrDiscard(ctx)

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	log.Info("request", "body", string(body))

	if len(body) > 0 {
		err = json.Unmarshal(body, data)
	} else {
		err = ErrEmptyBody
	}

	return err
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

func NewResourceAPI(factory zebra.ResourceFactory) *ResourceAPI {
	return &ResourceAPI{
		factory: factory,
		Store:   nil,
	}
}

type ResourceAPI struct {
	factory zebra.ResourceFactory
	Store   zebra.Store
}

type CtxKey string

const (
	ResourcesCtxKey = CtxKey("resources")
	AuthCtxKey      = CtxKey("authKey")
	ClaimsCtxKey    = CtxKey("claims")
)

func HandlePost() httprouter.Handle {
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
		if err := ReadJSON(ctx, req, resMap); err != nil {
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
