package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/rchamarthy/zebra"
	"github.com/rchamarthy/zebra/network"
	"github.com/rchamarthy/zebra/query"
	"github.com/rchamarthy/zebra/store"
)

//nolint:gochecknoglobals
var (
	types      map[string]func() zebra.Resource
	resStore   *store.FileStore
	queryStore *query.QueryStore

	ErrNumArgs = errors.New("wrong number of args")
)

// Set up store and query store given storage root.
func SetUp(storageRoot string) error {
	types = map[string]func() zebra.Resource{
		"VLANPool": func() zebra.Resource { return new(network.VLANPool) },
	}
	resStore = store.NewFileStore(storageRoot, types)

	err := resStore.Initialize()
	if err != nil {
		return err
	}

	resMap, err := resStore.Load()
	if err != nil {
		return err
	}

	queryStore = query.NewQueryStore(resMap)

	err = queryStore.Initialize()
	if err != nil {
		return err
	}

	return nil
}

func GetResources(w http.ResponseWriter, req *http.Request) { //nolint:varnamelen
	results := queryStore.Query()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err := json.NewEncoder(w).Encode(results)
	if err != nil {
		panic(err)
	}
}

func GetResourcesByID(w http.ResponseWriter, req *http.Request) { //nolint:varnamelen
	uuids := strings.Split(req.URL.Query().Get("id"), ",")
	results := queryStore.QueryUUID(uuids)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err := json.NewEncoder(w).Encode(results)
	if err != nil {
		panic(err)
	}
}

func GetResourcesByType(w http.ResponseWriter, req *http.Request) { //nolint:varnamelen
	resTypes := strings.Split(req.URL.Query().Get("type"), ",")
	results := queryStore.QueryType(resTypes)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err := json.NewEncoder(w).Encode(results)
	if err != nil {
		panic(err)
	}
}

func GetResourcesByProperty(w http.ResponseWriter, req *http.Request) { //nolint:varnamelen
	query, err := buildQuery(req.URL.Query().Get("property"), true)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	results, err := queryStore.QueryProperty(query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		panic(err)
	}
}

func GetResourcesByLabel(w http.ResponseWriter, req *http.Request) { //nolint:varnamelen
	query, err := buildQuery(req.URL.Query().Get("label"), false)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	results, err := queryStore.QueryLabel(query)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(results)
	if err != nil {
		panic(err)
	}
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
