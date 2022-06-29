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
	Types      map[string]func() zebra.Resource
	ResStore   *store.FileStore
	QueryStore *query.QueryStore

	ErrNumArgs = errors.New("wrong number of args")
)

// Set up store and query store given storage root.
func SetUp(storageRoot string) error {
	Types = map[string]func() zebra.Resource{
		"VLANPool": func() zebra.Resource { return new(network.VLANPool) },
	}
	ResStore = store.NewFileStore(storageRoot, Types)

	err := ResStore.Initialize()
	if err != nil {
		return err
	}

	resMap, err := ResStore.Load()
	if err != nil {
		return err
	}

	QueryStore = query.NewQueryStore(resMap)

	err = QueryStore.Initialize()
	if err != nil {
		return err
	}

	return nil
}

func GetResources(writer http.ResponseWriter, req *http.Request) {
	data := sortByType(QueryStore.Query())

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)

	err := json.NewEncoder(writer).Encode(data)
	if err != nil {
		panic(err)
	}
}

func GetResourcesByID(writer http.ResponseWriter, req *http.Request) {
	uuids := strings.Split(req.URL.Query().Get("id"), ",")
	data := sortByType(QueryStore.QueryUUID(uuids))

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)

	err := json.NewEncoder(writer).Encode(data)
	if err != nil {
		panic(err)
	}
}

func GetResourcesByType(writer http.ResponseWriter, req *http.Request) {
	types := strings.Split(req.URL.Query().Get("type"), ",")
	data := sortByType(QueryStore.QueryType(types))

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)

	err := json.NewEncoder(writer).Encode(data)
	if err != nil {
		panic(err)
	}
}

func GetResourcesByProperty(writer http.ResponseWriter, req *http.Request) {
	query, err := buildQuery(req.URL.Query().Get("property"), true)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	res, err := QueryStore.QueryProperty(query)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	data := sortByType(res)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(writer).Encode(data)
	if err != nil {
		panic(err)
	}
}

func GetResourcesByLabel(writer http.ResponseWriter, req *http.Request) {
	query, err := buildQuery(req.URL.Query().Get("label"), false)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	res, err := QueryStore.QueryLabel(query)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	data := sortByType(res)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(writer).Encode(data)
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

func sortByType(resources []zebra.Resource) map[string][]zebra.Resource {
	sortedRes := make(map[string][]zebra.Resource)

	for _, res := range resources {
		resType := res.GetType()
		sortedRes[resType] = append(sortedRes[resType], res)
	}

	return sortedRes
}
