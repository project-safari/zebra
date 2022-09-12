package store

import (
	"github.com/project-safari/zebra"
)

type valMap map[string][]zebra.Resource

func (v valMap) Add(val string, r zebra.Resource) {
	l := v[val]
	l = append(l, r)
	v[val] = l
}

func (v valMap) Delete(val string, r zebra.Resource) error {
	l := v[val]
	if l == nil {
		return zebra.ErrNotFound
	}

	listLen := len(l)

	for i, res := range l {
		if r.GetMeta().ID == res.GetMeta().ID {
			l[i] = l[listLen-1]
			l = l[:listLen-1]

			if len(l) == 0 {
				delete(v, val) // no more resources for this value
			} else {
				v[val] = l
			}

			return nil
		}
	}

	return zebra.ErrNotFound
}

type LabelStore struct {
	factory   zebra.ResourceFactory
	uuids     map[string]zebra.Resource
	resources map[string]valMap
}

// Return new label store pointer given resource map.
func NewLabelStore(resources *zebra.ResourceMap) *LabelStore {
	labelstore := &LabelStore{
		factory: resources.Factory(),
		uuids: func() map[string]zebra.Resource {
			ret := make(map[string]zebra.Resource)

			for _, l := range resources.Resources {
				for _, res := range l.Resources {
					ret[res.GetMeta().ID] = res
				}
			}

			return ret
		}(),
		resources: makeLabelMap(resources),
	}

	return labelstore
}

func makeLabelMap(resources *zebra.ResourceMap) map[string]valMap {
	labelMap := make(map[string]valMap)

	for _, l := range resources.Resources {
		for _, res := range l.Resources {
			for key, val := range res.GetMeta().Labels {
				if labelMap[key] == nil {
					labelMap[key] = valMap{}
				}

				labelMap[key].Add(val, res)
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
	ls.resources = make(map[string]valMap)
	ls.uuids = make(map[string]zebra.Resource)

	return nil
}

// Return all resources in a ResourceMap where keys are labelName = labelVal.
func (ls *LabelStore) Load() (*zebra.ResourceMap, error) {
	retMap := zebra.NewResourceMap(ls.factory)

	for _, res := range ls.uuids {
		if e := retMap.Add(res); e != nil {
			return nil, e
		}
	}

	return retMap, nil
}

// Create a resource. If a resource with this ID already exists, update.
func (ls *LabelStore) Create(res zebra.Resource) error {
	// do a best effor delete so the latest res wins
	_ = ls.Delete(res)

	// Create a new resource
	ls.uuids[res.GetMeta().ID] = res

	for label, val := range res.GetMeta().Labels {
		if ls.resources[label] == nil {
			ls.resources[label] = valMap{}
		}

		ls.resources[label].Add(val, res)
	}

	return nil
}

// Delete a resource.
func (ls *LabelStore) Delete(res zebra.Resource) error {
	// If resource does not exist in store, just return without error
	if _, err := ls.find(res.GetMeta().ID); err != nil {
		return zebra.ErrNotFound
	}

	for label, val := range res.GetMeta().Labels {
		if ls.resources[label] != nil {
			_ = ls.resources[label].Delete(val, res)

			if len(ls.resources[label]) == 0 {
				delete(ls.resources, label)
			}
		}
	}

	delete(ls.uuids, res.GetMeta().ID)

	return nil
}

// Return all resources of given label - label value pairs in a ResourceMap.
func (ls *LabelStore) Query(query zebra.Query) *zebra.ResourceMap {
	if query.Op == zebra.MatchEqual || query.Op == zebra.MatchIn {
		return ls.labelMatch(query, true)
	}

	return ls.labelMatch(query, false)
}

func (ls *LabelStore) labelMatch(query zebra.Query, inVals bool) *zebra.ResourceMap { //nolint:cyclop
	results := zebra.NewResourceMap(ls.factory)

	if inVals {
		for _, val := range query.Values {
			if ls.resources[query.Key] == nil || ls.resources[query.Key][val] == nil {
				continue
			}

			for _, res := range ls.resources[query.Key][val] {
				if results.Add(res) != nil {
					return nil
				}
			}
		}

		return results
	}

	for val, valMap := range ls.resources[query.Key] {
		if !zebra.IsIn(val, query.Values) {
			for _, res := range valMap {
				if results.Add(res) != nil {
					return nil
				}
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
