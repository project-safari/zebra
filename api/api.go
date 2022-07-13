package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/store"
)

type ResourceAPI struct {
	factory zebra.ResourceFactory
	Store   zebra.Store
}

var ErrNumArgs = errors.New("wrong number of args")

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

func (api *ResourceAPI) GetResources(w http.ResponseWriter, req *http.Request) {
	results := api.Store.Query()

	bytes, err := json.Marshal(results)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes) // nolint:errcheck
}

func (api *ResourceAPI) GetResourcesByID(w http.ResponseWriter, req *http.Request) {
	uuids := strings.Split(req.URL.Query().Get("id"), ",")

	results := api.Store.QueryUUID(uuids)

	bytes, err := json.Marshal(results)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes) // nolint:errcheck
}

func (api *ResourceAPI) GetResourcesByType(w http.ResponseWriter, req *http.Request) {
	resTypes := strings.Split(req.URL.Query().Get("type"), ",")

	results := api.Store.QueryType(resTypes)

	bytes, err := json.Marshal(results)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes) // nolint:errcheck
}

func (api *ResourceAPI) GetResourcesByProperty(w http.ResponseWriter, req *http.Request) {
	query, err := buildQuery(req.URL.Query().Get("property"), true)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	results, err := api.Store.QueryProperty(query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	bytes, err := json.Marshal(results)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes) // nolint:errcheck
}

func (api *ResourceAPI) GetResourcesByLabel(w http.ResponseWriter, req *http.Request) {
	query, err := buildQuery(req.URL.Query().Get("label"), false)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	results, err := api.Store.QueryLabel(query)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	bytes, err := json.Marshal(results)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes) // nolint:errcheck
}

func buildQuery(url string, isProperty bool) (zebra.Query, error) {
	params := strings.Split(url, "-")

	if len(params) != 3 { //nolint:gomnd
		return zebra.Query{}, ErrNumArgs
	}

	key := params[0]

	var operator zebra.Operator

	switch params[1] {
	case "equal":
		operator = zebra.MatchEqual
	case "notequal":
		operator = zebra.MatchNotEqual
	case "in":
		operator = zebra.MatchIn
	case "notin":
		operator = zebra.MatchNotIn
	default:
		return zebra.Query{}, zebra.ErrInvalidQuery
	}

	values := strings.Split(params[2], ",")

	return zebra.Query{Op: operator, Key: key, Values: values}, nil
}
