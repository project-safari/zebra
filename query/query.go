package query

import (
	"context"
	"errors"
	"reflect"
	"sync"

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
type Query struct {
	Op     Operator
	Key    string
	Values []string
}

// QueryStore keeps track of different maps for fast querying.
type QueryStore struct { //nolint:revive
	lock    sync.RWMutex
	rUUID   map[string]zebra.Resource
	rType   *zebra.ResourceMap
	rLabel  map[string]*zebra.ResourceMap
	factory zebra.ResourceFactory
}

var ErrOpVals = errors.New("number of values not valid for query operator")

var ErrOp = errors.New("query operator not valid")

var ErrResExists = errors.New("called create on resource that already exists")

var ErrResDoesNotExist = errors.New("called update on resource that does not exist")

// Return new query store pointer given resource map.
func NewQueryStore(resources *zebra.ResourceMap) *QueryStore {
	querystore := &QueryStore{
		lock:  sync.RWMutex{},
		rUUID: nil,
		rType: func() *zebra.ResourceMap {
			dest := zebra.NewResourceMap(nil)
			zebra.CopyResourceMap(dest, resources)

			return dest
		}(),
		rLabel:  nil,
		factory: resources.GetFactory(),
	}

	return querystore
}

// Initialize indexes for query store.
func (qs *QueryStore) Initialize() error {
	qs.lock.Lock()
	defer qs.lock.Unlock()

	return qs.init()
}

// init implements the store initialization. This function must never be called
// without holding the write lock.
func (qs *QueryStore) init() error {
	qs.rUUID = make(map[string]zebra.Resource)
	qs.rLabel = make(map[string]*zebra.ResourceMap)

	for _, resList := range qs.rType.Resources {
		for _, res := range resList.Resources {
			qs.rUUID[res.GetID()] = res

			for labelName, labelVal := range res.GetLabels() {
				if qs.rLabel[labelName] == nil {
					qs.rLabel[labelName] = zebra.NewResourceMap(qs.factory)
				}

				qs.rLabel[labelName].Add(res, labelVal)
			}
		}
	}

	return nil
}

// Delete all index maps for qs.
func (qs *QueryStore) Wipe() error {
	qs.lock.Lock()
	defer qs.lock.Unlock()

	qs.rUUID = nil
	qs.rType = nil
	qs.rLabel = nil

	return nil
}

// Clear all index maps for qs.
func (qs *QueryStore) Clear() error {
	qs.lock.Lock()
	defer qs.lock.Unlock()

	qs.rUUID = make(map[string]zebra.Resource, 0)
	qs.rType = zebra.NewResourceMap(nil)
	qs.rLabel = make(map[string]*zebra.ResourceMap, 0)

	return nil
}

// Return all resources in a ResourceSet.
func (qs *QueryStore) Load() (*zebra.ResourceMap, error) {
	qs.lock.RLock()
	defer qs.lock.RUnlock()

	resources := zebra.NewResourceMap(nil)

	zebra.CopyResourceMap(resources, qs.rType)

	return resources, nil
}

// Create a resource. If a resource with this ID already exists, return error.
func (qs *QueryStore) Create(res zebra.Resource) error {
	qs.lock.Lock()
	defer qs.lock.Unlock()

	return qs.create(res)
}

// Should not be called without holding the write lock.
func (qs *QueryStore) create(res zebra.Resource) error {
	resID := res.GetID()

	// If resource already exists, return error.
	if _, exists := qs.rUUID[resID]; exists {
		return ErrResExists
	}

	if err := res.Validate(context.Background()); err != nil {
		return err
	}

	resType := res.GetType()

	qs.rUUID[resID] = res
	qs.rType.Add(res, resType)

	for labelName, labelVal := range res.GetLabels() {
		if qs.rLabel[labelName] == nil {
			qs.rLabel[labelName] = zebra.NewResourceMap(qs.factory)
		}

		qs.rLabel[labelName].Add(res, labelVal)
	}

	return nil
}

// Update a resource. Return error if resource does not exist.
func (qs *QueryStore) Update(res zebra.Resource) error {
	qs.lock.Lock()
	defer qs.lock.Unlock()

	if err := res.Validate(context.Background()); err != nil {
		return err
	}

	resID := res.GetID()
	oldRes, exists := qs.rUUID[resID]

	// If resource does not exist, return error.
	if !exists {
		return ErrResDoesNotExist
	}

	_ = qs.delete(oldRes)

	_ = qs.create(res)

	return nil
}

// Delete a resource.
func (qs *QueryStore) Delete(res zebra.Resource) error {
	qs.lock.Lock()
	defer qs.lock.Unlock()

	if err := res.Validate(context.Background()); err != nil {
		return err
	}

	return qs.delete(res)
}

// Should not be called without holding the write lock.
func (qs *QueryStore) delete(res zebra.Resource) error {
	delete(qs.rUUID, res.GetID())
	qs.rType.Delete(res, res.GetType())

	for labelName, labelVal := range res.GetLabels() {
		qs.rLabel[labelName].Delete(res, labelVal)
	}

	return nil
}

// Return all resources in a ResourceMap.
func (qs *QueryStore) Query() *zebra.ResourceMap {
	qs.lock.RLock()
	defer qs.lock.RUnlock()

	resources := zebra.NewResourceMap(nil)

	zebra.CopyResourceMap(resources, qs.rType)

	return resources
}

// Return resources with matching UUIDs.
func (qs *QueryStore) QueryUUID(uuids []string) *zebra.ResourceMap {
	qs.lock.RLock()
	defer qs.lock.RUnlock()

	resources := zebra.NewResourceMap(qs.factory)

	for _, id := range uuids {
		res, ok := qs.rUUID[id]
		if ok {
			resources.Add(res, res.GetType())
		}
	}

	return resources
}

// Return resources with matching types.
func (qs *QueryStore) QueryType(types []string) *zebra.ResourceMap {
	qs.lock.RLock()
	defer qs.lock.RUnlock()

	resources := zebra.NewResourceMap(qs.factory)

	for _, t := range types {
		resList := qs.rType.Resources[t]
		if resList != nil {
			resources.Resources[t] = zebra.NewResourceList(qs.factory)
			zebra.CopyResourceList(resources.Resources[t], resList)
		}
	}

	return resources
}

// Return resources with matching label.
func (qs *QueryStore) QueryLabel(query Query) (*zebra.ResourceMap, error) {
	qs.lock.RLock()
	defer qs.lock.RUnlock()

	switch query.Op {
	case MatchEqual:
		if len(query.Values) != 1 {
			return nil, ErrOpVals
		}

		fallthrough
	case MatchIn:
		return qs.labelMatch(query, true)
	case MatchNotEqual:
		if len(query.Values) != 1 {
			return nil, ErrOpVals
		}

		fallthrough
	case MatchNotIn:
		return qs.labelMatch(query, false)
	default:
		return nil, ErrOp
	}
}

// Return resources which match given property/value(s).
// Naive search implementation, >= O(n) for n resources.
func (qs *QueryStore) QueryProperty(query Query) (*zebra.ResourceMap, error) {
	qs.lock.RLock()
	defer qs.lock.RUnlock()

	switch query.Op {
	case MatchEqual:
		if len(query.Values) != 1 {
			return nil, ErrOpVals
		}

		fallthrough
	case MatchIn:
		return qs.propertyMatch(query, true)
	case MatchNotEqual:
		if len(query.Values) != 1 {
			return nil, ErrOpVals
		}

		fallthrough
	case MatchNotIn:
		return qs.propertyMatch(query, false)
	default:
		return nil, ErrOp
	}
}

func (qs *QueryStore) labelMatch(query Query, inVals bool) (*zebra.ResourceMap, error) {
	results := zebra.NewResourceMap(qs.factory)

	if inVals {
		for _, val := range query.Values {
			for _, res := range qs.rLabel[query.Key].Resources[val].Resources {
				results.Add(res, res.GetType())
			}
		}

		return results, nil
	}

	for val, valMap := range qs.rLabel[query.Key].Resources {
		if !isIn(val, query.Values) {
			for _, res := range valMap.Resources {
				results.Add(res, res.GetType())
			}
		}
	}

	return results, nil
}

func (qs *QueryStore) propertyMatch(query Query, inVals bool) (*zebra.ResourceMap, error) {
	results := zebra.NewResourceMap(qs.factory)

	for _, res := range qs.rUUID {
		val := reflect.ValueOf(res).Elem().FieldByName(query.Key).String()
		inList := isIn(val, query.Values)

		if inVals && inList {
			results.Add(res, res.GetType())
		} else if !inVals && !inList {
			results.Add(res, res.GetType())
		}
	}

	return results, nil
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
