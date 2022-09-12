package store

import (
	"github.com/project-safari/zebra"
)

type TypeStore struct {
	resources *zebra.ResourceMap
}

// Return new type store pointer given resource map.
func NewTypeStore(resources *zebra.ResourceMap) *TypeStore {
	typestore := &TypeStore{
		resources: func() *zebra.ResourceMap {
			dest := zebra.NewResourceMap(resources.Factory())
			zebra.CopyResourceMap(dest, resources)

			return dest
		}(),
	}

	return typestore
}

func (ts *TypeStore) Initialize() error {
	return nil
}

func (ts *TypeStore) Wipe() error {
	ts.resources = nil

	return nil
}

func (ts *TypeStore) Clear() error {
	ts.resources = zebra.NewResourceMap(ts.resources.Factory())

	return nil
}

// Return all resources in a ResourceMap.
func (ts *TypeStore) Load() (*zebra.ResourceMap, error) {
	resources := zebra.NewResourceMap(ts.resources.Factory())

	zebra.CopyResourceMap(resources, ts.resources)

	return resources, nil
}

// Create a resource. If a resource with this ID already exists, update.
func (ts *TypeStore) Create(res zebra.Resource) error {
	// Check if resource already exists
	meta := res.GetMeta()
	if oldRes, err := ts.find(meta.ID, meta.Type.Name); err == nil {
		// Delete the old resource
		_ = ts.Delete(oldRes)
	}

	// If it doesn't, add to resource map
	return ts.resources.Add(res)
}

// Delete a resource.
func (ts *TypeStore) Delete(res zebra.Resource) error {
	return ts.resources.Delete(res)
}

// Return all resources of given types in a ResourceMap.
func (ts *TypeStore) Query(types []string) *zebra.ResourceMap {
	factory := ts.resources.Factory()
	retMap := zebra.NewResourceMap(factory)

	for _, t := range types {
		if ts.resources.Resources[t] != nil {
			retMap.Resources[t] = ts.resources.Resources[t]
		}
	}

	return retMap
}

// Find given resource in TypeStore. If not found, return nil and error.
// If found, return resource and nil.
func (ts *TypeStore) find(resID string, resType string) (zebra.Resource, error) {
	resMap := ts.resources.Resources[resType]
	if resMap == nil {
		return nil, zebra.ErrNotFound
	}

	for _, val := range resMap.Resources {
		if val.GetMeta().ID == resID {
			return val, nil
		}
	}

	return nil, zebra.ErrNotFound
}
