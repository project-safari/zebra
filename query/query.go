package query

import (
	"errors"
	"reflect"
	"strings"

	"github.com/rchamarthy/zebra"
)

type Operator uint8

// Constants defined for QueryOperator type.
const (
	MatchEqual Operator = iota
	MatchNotEqual
	MatchIn
	MatchNotIn
)

// Command struct for label queries.
type LabelQuery struct {
	Op     Operator
	Key    string
	Values []string
}

// QueryStore keeps track of different maps for fast querying.
type QueryStore struct { //nolint:revive
	rAll   []zebra.Resource
	rUUID  map[string]zebra.Resource
	rType  map[string][]zebra.Resource
	rLabel map[(string)][]pair
}

// pair acts as a tuple for tracking label value and associated resource.
type pair struct {
	value    string
	resource zebra.Resource
}

var ErrQuery = errors.New("query not valid")

// Initialize new query store and set up indexes.
func NewQueryStore(resourcesUUID map[string]zebra.Resource) *QueryStore {
	qs := new(QueryStore)
	qs.rUUID = resourcesUUID
	qs.setUp()

	return qs
}

// Set up QueryStore object maps.
func (qs *QueryStore) setUp() {
	qs.rAll = make([]zebra.Resource, 0, len(qs.rUUID))
	qs.rType = make(map[string][]zebra.Resource)
	qs.rLabel = make(map[string][]pair)

	for _, val := range qs.rUUID {
		qs.rAll = append(qs.rAll, val)
		fullName := strings.Split(reflect.TypeOf(val).String(), ".")
		resType := fullName[len(fullName)-1]
		qs.rType[resType] = append(qs.rType[resType], val)

		for name, label := range val.GetLabels() {
			qs.rLabel[name] = append(qs.rLabel[name], pair{label, val})
		}
	}
}

// Return all resources.
func (qs *QueryStore) Query() []zebra.Resource {
	dest := make([]zebra.Resource, len(qs.rAll))
	copy(dest, qs.rAll)

	return dest
}

// Return resources with matching UUIDs.
func (qs *QueryStore) QueryUUID(uuids []string) []zebra.Resource {
	resources := make([]zebra.Resource, 0, len(uuids))

	for _, id := range uuids {
		res, ok := qs.rUUID[id]
		if ok {
			resources = append(resources, res)
		}
	}

	return resources
}

// Return resources with matching types.
func (qs *QueryStore) QueryType(types ...string) []zebra.Resource {
	resources := []zebra.Resource{}
	for _, t := range types {
		resources = append(resources, qs.rType[t]...)
	}

	return resources
}

// Return resources with matching type and names.
func (qs *QueryStore) QueryTypeName(resType string, names []string) []zebra.Resource {
	resources := qs.QueryType(resType)
	for i, res := range resources {
		if !isIn(res.GetName(), names) {
			resources[i] = resources[len(resources)-1]
			resources = resources[:len(resources)-1]
		}
	}

	return resources
}

// Return resources with all matching labels.
func (qs *QueryStore) QueryLabelsMatchAll(queries []LabelQuery) ([]zebra.Resource, error) {
	resources := make(map[string]zebra.Resource)

	for queryNum, query := range queries {
		results, err := qs.matchLabel(query)
		if err != nil {
			return nil, err
		}

		// If it's the first query, only add to resources.
		// Else, update resources.
		if queryNum == 0 {
			for _, res := range results {
				resources[res.GetID()] = res
			}
		} else {
			for _, res := range resources {
				id := res.GetID()
				if results[id] == nil {
					delete(resources, id)
				}
			}
		}
	}

	returnVal := make([]zebra.Resource, 0, len(resources))
	for _, val := range resources {
		returnVal = append(returnVal, val)
	}

	return returnVal, nil
}

// Return resources with at least one matching label.
func (qs *QueryStore) QueryLabelsMatchOne(queries []LabelQuery) ([]zebra.Resource, error) {
	resources := make(map[string]zebra.Resource)

	for _, query := range queries {
		results, err := qs.matchLabel(query)
		if err != nil {
			return nil, err
		}

		for _, res := range results {
			resources[res.GetID()] = res
		}
	}

	returnVal := make([]zebra.Resource, 0, len(resources))
	for _, val := range resources {
		returnVal = append(returnVal, val)
	}

	return returnVal, nil
}

// Return list of resources filtered by specific label's values.
func (qs *QueryStore) matchLabel(query LabelQuery) (map[string]zebra.Resource, error) {
	switch query.Op {
	case MatchEqual:
		if len(query.Values) != 1 {
			return nil, ErrQuery
		}

		fallthrough
	case MatchIn:
		return match(qs.rLabel[query.Key], query.Values), nil
	case MatchNotEqual:
		if len(query.Values) != 1 {
			return nil, ErrQuery
		}

		fallthrough
	case MatchNotIn:
		return notMatch(qs.rLabel[query.Key], query.Values), nil
	default:
		return nil, ErrQuery
	}
}

// Return resources that have label values matching something in vals.
func match(pairs []pair, vals []string) map[string]zebra.Resource {
	resources := make(map[string]zebra.Resource)

	for _, res := range pairs {
		if isIn(res.value, vals) {
			resources[res.resource.GetID()] = res.resource
		}
	}

	return resources
}

// Return resources that do not have label values matching something in vals.
func notMatch(pairs []pair, vals []string) map[string]zebra.Resource {
	resources := make(map[string]zebra.Resource)

	for _, res := range pairs {
		if !isIn(res.value, vals) {
			resources[res.resource.GetID()] = res.resource
		}
	}

	return resources
}

// Return if val is in string list.
func isIn(val string, list []string) bool {
	for _, v := range list {
		if val == v {
			return true
		}
	}

	return false
}
