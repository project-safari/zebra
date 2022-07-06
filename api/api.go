package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/rchamarthy/zebra"
	"github.com/rchamarthy/zebra/query"
	"github.com/rchamarthy/zebra/store"
)

type ResourceAPI struct {
	factory    zebra.ResourceFactory
	resStore   *store.FileStore
	queryStore *query.QueryStore
}

var ErrNumArgs = errors.New("wrong number of args")

func NewResourceAPI(factory zebra.ResourceFactory) *ResourceAPI {
	return &ResourceAPI{
		factory:    factory,
		resStore:   nil,
		queryStore: nil,
	}
}

// Set up store and query store given storage root.
func (api *ResourceAPI) Initialize(storageRoot string) error {
	api.resStore = store.NewFileStore(storageRoot, api.factory)

	if err := api.resStore.Initialize(); err != nil {
		return err
	}

	resMap, err := api.resStore.Load()
	if err != nil {
		return err
	}

	api.queryStore = query.NewQueryStore(resMap)

	if err = api.queryStore.Initialize(); err != nil {
		return err
	}

	return nil
}

func (api *ResourceAPI) GetResources(w http.ResponseWriter, req *http.Request) { //nolint:varnamelen
	results := api.queryStore.Query()

	bytes, err := json.Marshal(results)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes) // nolint:errcheck
}

func (api *ResourceAPI) GetResourcesByID(w http.ResponseWriter, req *http.Request) { //nolint:varnamelen
	uuids := strings.Split(req.URL.Query().Get("id"), ",")

	results := api.queryStore.QueryUUID(uuids)

	bytes, err := json.Marshal(results)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes) // nolint:errcheck
}

func (api *ResourceAPI) GetResourcesByType(w http.ResponseWriter, req *http.Request) { //nolint:varnamelen
	resTypes := strings.Split(req.URL.Query().Get("type"), ",")

	results := api.queryStore.QueryType(resTypes)

	bytes, err := json.Marshal(results)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes) // nolint:errcheck
}

func (api *ResourceAPI) GetResourcesByProperty(w http.ResponseWriter, req *http.Request) { //nolint:varnamelen
	query, err := buildQuery(req.URL.Query().Get("property"), true)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	results, err := api.queryStore.QueryProperty(query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	bytes, err := json.Marshal(results)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes) // nolint:errcheck
}

func (api *ResourceAPI) GetResourcesByLabel(w http.ResponseWriter, req *http.Request) { //nolint:varnamelen
	query, err := buildQuery(req.URL.Query().Get("label"), false)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	results, err := api.queryStore.QueryLabel(query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	bytes, err := json.Marshal(results)
	if err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(bytes) // nolint:errcheck
}

func buildQuery(url string, isProperty bool) (query.Query, error) {
	params := strings.Split(url, "-")

	if len(params) != 3 { //nolint:gomnd
		return query.Query{}, ErrNumArgs
	}

	key := params[0]

	var operator query.Operator

	switch params[1] {
	case "equal":
		operator = query.MatchEqual
	case "notequal":
		operator = query.MatchNotEqual
	case "in":
		operator = query.MatchIn
	case "notin":
		operator = query.MatchNotIn
	}

	values := strings.Split(params[2], ",")

	return query.Query{Op: operator, Key: key, Values: values}, nil
}
