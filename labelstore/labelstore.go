package labelstore

import (
	"github.com/project-safari/zebra"
)

type LabelStore struct {
	factory   zebra.ResourceFactory
	uuids     map[string]zebra.Resource
	resources map[string]*zebra.ResourceMap
}

// Return new label store pointer given resource map.
func NewLabelStore(resources *zebra.ResourceMap) *LabelStore {
	labelstore := &LabelStore{
		factory: resources.GetFactory(),
		uuids: func() map[string]zebra.Resource {
			ret := make(map[string]zebra.Resource)

			for _, l := range resources.Resources {
				for _, res := range l.Resources {
					ret[res.GetID()] = res
				}
			}

			return ret
		}(),
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

	for _, res := range ls.uuids {
		for label, val := range res.GetLabels() {
			key := label + " = " + val
			retMap.Add(res, key)
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

	delete(ls.uuids, res.GetID())

	return nil
}

// Return all resources of given label - label value pairs in a ResourceMap.
func (ls *LabelStore) Query(query zebra.Query) *zebra.ResourceMap {
	if query.Op == zebra.MatchEqual || query.Op == zebra.MatchIn {
		return ls.labelMatch(query, true)
	}

	return ls.labelMatch(query, false)
}

func (ls *LabelStore) labelMatch(query zebra.Query, inVals bool) *zebra.ResourceMap {
	results := zebra.NewResourceMap(ls.factory)

	if inVals {
		for _, val := range query.Values {
			if ls.resources[query.Key] == nil || ls.resources[query.Key].Resources[val] == nil {
				continue
			}

			for _, res := range ls.resources[query.Key].Resources[val].Resources {
				results.Add(res, res.GetType().Name)
			}
		}

		return results
	}

	for val, valMap := range ls.resources[query.Key].Resources {
		if !zebra.IsIn(val, query.Values) {
			for _, res := range valMap.Resources {
				results.Add(res, res.GetType().Name)
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
