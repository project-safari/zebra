package typestore

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
			dest := zebra.NewResourceMap(resources.GetFactory())
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
	ts.resources = zebra.NewResourceMap(ts.resources.GetFactory())

	return nil
}

// Return all resources in a ResourceMap.
func (ts *TypeStore) Load() (*zebra.ResourceMap, error) {
	resources := zebra.NewResourceMap(ts.resources.GetFactory())

	zebra.CopyResourceMap(resources, ts.resources)

	return resources, nil
}

// Create a resource. If a resource with this ID already exists, update.
func (ts *TypeStore) Create(res zebra.Resource) error {
	// Check if resource already exists
	if oldRes, err := ts.find(res.GetID(), res.GetType()); err == nil {
		return ts.update(oldRes, res)
	}

	// If it doesn't, add to resource map
	ts.resources.Add(res, res.GetType())

	return nil
}

// Update a resource.
func (ts *TypeStore) update(oldRes zebra.Resource, res zebra.Resource) error {
	if err := ts.Delete(oldRes); err != nil {
		return err
	}

	return ts.Create(res)
}

// Delete a resource.
func (ts *TypeStore) Delete(res zebra.Resource) error {
	ts.resources.Delete(res, res.GetType())

	return nil
}

// Return all resources of given types in a ResourceMap.
func (ts *TypeStore) Query(types []string) *zebra.ResourceMap {
	factory := ts.resources.GetFactory()
	retMap := zebra.NewResourceMap(factory)

	for _, t := range types {
		retMap.Resources[t] = ts.resources.Resources[t]
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
		if val.GetID() == resID {
			return val, nil
		}
	}

	return nil, zebra.ErrNotFound
}
