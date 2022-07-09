package labelstore

import (
	"github.com/project-safari/zebra"
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

type LabelStore struct {
	factory   zebra.ResourceFactory
	uuids     map[string]zebra.Resource
	resources map[string]*zebra.ResourceMap
}

// Return new label store pointer given resource map.
func NewLabelStore(resources *zebra.ResourceMap) *LabelStore {
	labelstore := &LabelStore{
		factory:   resources.GetFactory(),
		uuids:     make(map[string]zebra.Resource),
		resources: makeLabelMap(resources),
	}

	return labelstore
}

func makeLabelMap(resources *zebra.ResourceMap) map[string]*zebra.ResourceMap {
	labelMap := make(map[string]*zebra.ResourceMap)

	for _, l := range resources.Resources {
		for _, res := range l.Resources {
			for label, val := range res.GetLabels() {
				if labelMap[label] == nil {
					labelMap[label] = zebra.NewResourceMap(resources.GetFactory())
				}

				labelMap[label].Add(res, val)
			}
		}
	}

	return labelMap
}

func (ls *LabelStore) Initialize() error {
	return nil
}

func (ls *LabelStore) Wipe() error {
	ls.resources = nil
	ls.uuids = nil

	return nil
}

func (ls *LabelStore) Clear() error {
	ls.resources = make(map[string]*zebra.ResourceMap)
	ls.uuids = make(map[string]zebra.Resource)

	return nil
}

// Return all resources in a ResourceMap where keys are labelName = labelVal.
func (ls *LabelStore) Load() (*zebra.ResourceMap, error) {
	retMap := zebra.NewResourceMap(ls.factory)

	for label, valMap := range ls.resources {
		for val, resList := range valMap.Resources {
			list := zebra.NewResourceList(nil)
			key := label + " = " + val

			zebra.CopyResourceList(list, resList)
			retMap.Resources[key] = list
		}
	}

	return retMap, nil
}

// Create a resource. If a resource with this ID already exists, update.
func (ls *LabelStore) Create(res zebra.Resource) error {
	// Check if resource already exists, update if so
	if oldRes, err := ls.find(res.GetID()); err == nil {
		return ls.update(oldRes, res)
	}

	// Create a new resource
	ls.uuids[res.GetID()] = res

	for label, val := range res.GetLabels() {
		if ls.resources[label] == nil {
			ls.resources[label] = zebra.NewResourceMap(ls.factory)
		}

		ls.resources[label].Add(res, val)
	}

	return nil
}

// Update a resource.
func (ls *LabelStore) update(oldRes zebra.Resource, res zebra.Resource) error {
	if err := ls.Delete(oldRes); err != nil {
		return err
	}

	return ls.Create(res)
}

// Delete a resource.
func (ls *LabelStore) Delete(res zebra.Resource) error {
	// If resource does not exist in store, just return without error
	if _, err := ls.find(res.GetID()); err != nil {
		return nil
	}

	for label, val := range res.GetLabels() {
		if ls.resources[label] != nil {
			ls.resources[label].Delete(res, val)

			if len(ls.resources[label].Resources) == 0 {
				delete(ls.resources, label)
			}
		}
	}

	ls.uuids[res.GetID()] = nil

	return nil
}

// Return all resources of given label - label value pairs in a ResourceMap.
func (ls *LabelStore) Query(query Query) *zebra.ResourceMap {
	switch query.Op {
	case MatchEqual:
		if len(query.Values) != 1 {
			return nil
		}

		fallthrough
	case MatchIn:
		return ls.labelMatch(query, true)
	case MatchNotEqual:
		if len(query.Values) != 1 {
			return nil
		}

		fallthrough
	case MatchNotIn:
		return ls.labelMatch(query, false)
	default:
		return nil
	}
}

func (ls *LabelStore) labelMatch(query Query, inVals bool) *zebra.ResourceMap {
	results := zebra.NewResourceMap(ls.factory)

	if inVals {
		for _, val := range query.Values {
			for _, res := range ls.resources[query.Key].Resources[val].Resources {
				results.Add(res, res.GetType())
			}
		}

		return results
	}

	for val, valMap := range ls.resources[query.Key].Resources {
		if !isIn(val, query.Values) {
			for _, res := range valMap.Resources {
				results.Add(res, res.GetType())
			}
		}
	}

	return results
}

// Find given resource in LabelStore. If not found, return nil and error.
// If found, return resource and nil.
func (ls *LabelStore) find(resID string) (zebra.Resource, error) {
	val, ok := ls.uuids[resID]
	if !ok {
		return nil, zebra.ErrNotFound
	}

	return val, nil
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
