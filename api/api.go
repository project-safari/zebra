package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io/ioutil"
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

func (api *ResourceAPI) PutResource(w http.ResponseWriter, req *http.Request) {
	if req.Body == nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	res := api.unpackResource(w, body)
	if res == nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	// Check if this is a create or an update.
	exists := len(api.Store.QueryUUID([]string{res.GetID()}).Resources) != 0

	if err := api.Store.Create(res); err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if exists {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

func (api *ResourceAPI) DeleteResource(w http.ResponseWriter, req *http.Request) {
	if req.Body == nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	var ids []string

	err = json.Unmarshal(body, &ids)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	resources := api.Store.QueryUUID(ids)
	status := make(map[string]int, len(ids))

	for _, l := range resources.Resources {
		for _, res := range l.Resources {
			if api.Store.Delete(res) != nil {
				status[res.GetID()] = -1
			} else {
				status[res.GetID()] = 1
			}
		}
	}

	httpStatus, response := createDeleteResponse(ids, status)

	w.WriteHeader(httpStatus)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(html.EscapeString(response))) //nolint:errcheck
}

func createDeleteResponse(ids []string, status map[string]int) (int, string) {
	httpStatus := http.StatusOK
	successID := make([]string, 0)
	failedID := make([]string, 0)
	invalidID := make([]string, 0)

	for _, id := range ids {
		switch status[id] {
		case -1:
			failedID = append(failedID, id)
		case 1:
			successID = append(successID, id)
		default:
			invalidID = append(invalidID, id)
		}
	}

	var response string

	if len(successID) > 0 {
		response = fmt.Sprintf("Deleted the following resources: %s\n", strings.Join(successID, ", "))
	}

	if len(failedID) > 0 {
		httpStatus = http.StatusMultiStatus

		response += fmt.Sprintf("Failed to delete the following resources: %s\n", strings.Join(failedID, ", "))
	}

	if len(invalidID) > 0 {
		response += fmt.Sprintf("Invalid resource IDs: %s\n", strings.Join(invalidID, ", "))
	}

	return httpStatus, response
}

func (api *ResourceAPI) unpackResource(w http.ResponseWriter, body []byte) zebra.Resource {
	object := make(map[string]interface{})
	if err := json.Unmarshal(body, &object); err != nil {
		return nil
	}

	resType, ok := object["type"].(string)
	if !ok {
		return nil
	}

	res := api.factory.New(resType)
	if err := json.Unmarshal(body, res); err != nil {
		return nil
	}

	return res
}
