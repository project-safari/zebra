package query

import (
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
	lock   sync.RWMutex
	rUUID  map[string]zebra.Resource
	rType  zebra.ResourceSet
	rLabel map[string]zebra.ResourceSet
}

var ErrOpVals = errors.New("number of values not valid for query operator")

var ErrOp = errors.New("query operator not valid")

var ErrResExists = errors.New("called create on resource that already exists")

var ErrResDoesNotExist = errors.New("called update on resource that does not exist")

// Return new query store pointer given resource map.
func NewQueryStore(resources zebra.ResourceSet) *QueryStore {
	querystore := &QueryStore{
		lock:  sync.RWMutex{},
		rUUID: nil,
		rType: func() zebra.ResourceSet {
			results := make(zebra.ResourceSet, len(resources))
			// Make a copy of resources so that they are not mutated after the store has
			// been created and initialized.
			for resType, resList := range resources {
				results[resType] = make(zebra.ResourceList, len(resList))
				copy(results[resType], resList)
			}

			return results
		}(),
		rLabel: nil,
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
	qs.rLabel = make(map[string]zebra.ResourceSet)

	for _, resList := range qs.rType {
		for _, res := range resList {
			qs.rUUID[res.GetID()] = res

			for labelName, labelVal := range res.GetLabels() {
				if qs.rLabel[labelName] == nil {
					qs.rLabel[labelName] = make(zebra.ResourceSet, 1)
				}

				qs.rLabel[labelName][labelVal] = append(qs.rLabel[labelName][labelVal], res)
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
	qs.rType = make(zebra.ResourceSet, 0)
	qs.rLabel = make(map[string]zebra.ResourceSet, 0)

	return nil
}

// Return all resources in a ResourceSet.
func (qs *QueryStore) Load() (zebra.ResourceSet, error) {
	qs.lock.RLock()
	defer qs.lock.RUnlock()

	resources := make(zebra.ResourceSet)

	for resType, resList := range qs.rType {
		resources[resType] = make(zebra.ResourceList, len(resList))
		copy(resources[resType], resList)
	}

	return resources, nil
}

// Create a resource. If a resource with this ID already exists, return error.
func (qs *QueryStore) Create(res zebra.Resource) error {
	qs.lock.RLock()
	resID := res.GetID()
	_, exists := qs.rUUID[resID] //nolint:ifshort
	qs.lock.RUnlock()

	// If resource already exists, return error.
	if exists {
		return ErrResExists
	}

	qs.lock.Lock()
	defer qs.lock.Unlock()

	resType := res.GetType()

	qs.rUUID[resID] = res
	qs.rType[resType] = append(qs.rType[resType], res)

	for labelName, labelVal := range res.GetLabels() {
		if qs.rLabel[labelName] == nil {
			qs.rLabel[labelName] = make(zebra.ResourceSet, 1)
		}

		qs.rLabel[labelName][labelVal] = append(qs.rLabel[labelName][labelVal], res)
	}

	return nil
}

// Update a resource. Return error if resource does not exist.
func (qs *QueryStore) Update(res zebra.Resource) error {
	qs.lock.RLock()
	resID := res.GetID()
	oldRes, exists := qs.rUUID[resID]
	qs.lock.RUnlock()

	// If resource does not exist, return error.
	if !exists {
		return ErrResDoesNotExist
	}

	_ = qs.Delete(oldRes)

	_ = qs.Create(res)

	return nil
}

// Delete a resource.
func (qs *QueryStore) Delete(res zebra.Resource) error {
	qs.lock.Lock()
	defer qs.lock.Unlock()

	resID := res.GetID()
	resType := res.GetType()

	delete(qs.rUUID, resID)

	for i, val := range qs.rType[resType] {
		listLen := len(qs.rType[resType])
		// To delete from types list, move last elem.
		if val.GetID() == resID {
			qs.rType[resType][i] = qs.rType[resType][listLen-1]
			qs.rType[resType] = qs.rType[resType][:listLen-1]
		}
	}

	for labelName, labelVal := range res.GetLabels() {
		length := len(qs.rLabel[labelName][labelVal])

		for i, val := range qs.rLabel[labelName][labelVal] {
			if val.GetID() == resID {
				qs.rLabel[labelName][labelVal][i] = qs.rLabel[labelName][labelVal][length-1]
				qs.rLabel[labelName][labelVal] = qs.rLabel[labelName][labelVal][:length-1]
			}
		}
	}

	return nil
}

// Return all resources in a slice.
func (qs *QueryStore) Query() zebra.ResourceSet {
	qs.lock.RLock()
	defer qs.lock.RUnlock()

	resources := make(zebra.ResourceSet, len(qs.rType))

	for resType, resList := range qs.rType {
		resources[resType] = make(zebra.ResourceList, len(resList))
		copy(resources[resType], resList)
	}

	return resources
}

// Return resources with matching UUIDs.
func (qs *QueryStore) QueryUUID(uuids []string) zebra.ResourceSet {
	qs.lock.RLock()
	defer qs.lock.RUnlock()

	resources := make(zebra.ResourceSet, 0)

	for _, id := range uuids {
		res, ok := qs.rUUID[id]
		if ok {
			resType := res.GetType()
			resources[resType] = append(resources[resType], res)
		}
	}

	return resources
}

// Return resources with matching types.
func (qs *QueryStore) QueryType(types []string) zebra.ResourceSet {
	qs.lock.RLock()
	defer qs.lock.RUnlock()

	resources := make(zebra.ResourceSet)

	for _, t := range types {
		resList := qs.rType[t]
		if len(resList) > 0 {
			resources[t] = append(resources[t], resList...)
		}
	}

	return resources
}

// Return resources with matching label.
func (qs *QueryStore) QueryLabel(query Query) (zebra.ResourceSet, error) {
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
func (qs *QueryStore) QueryProperty(query Query) (zebra.ResourceSet, error) {
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

func (qs *QueryStore) labelMatch(query Query, inVals bool) (zebra.ResourceSet, error) {
	if inVals {
		results := make(zebra.ResourceSet, 0)

		for _, val := range query.Values {
			for _, res := range qs.rLabel[query.Key][val] {
				resType := res.GetType()
				results[resType] = append(results[resType], res)
			}
		}

		return results, nil
	}

	results := make(zebra.ResourceSet, 0)

	for val, valMap := range qs.rLabel[query.Key] {
		if !isIn(val, query.Values) {
			for _, res := range valMap {
				resType := res.GetType()
				results[resType] = append(results[resType], res)
			}
		}
	}

	return results, nil
}

func (qs *QueryStore) propertyMatch(query Query, inVals bool) (zebra.ResourceSet, error) {
	results := make(zebra.ResourceSet, 0)

	for _, res := range qs.rUUID {
		val := reflect.ValueOf(res).Elem().FieldByName(query.Key).String()
		inList := isIn(val, query.Values)

		if inVals && inList {
			resType := res.GetType()
			results[resType] = append(results[resType], res)
		} else if !inVals && !inList {
			resType := res.GetType()
			results[resType] = append(results[resType], res)
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
