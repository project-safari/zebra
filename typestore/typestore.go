package typestore

import (
	"context"
	"sync"

	"github.com/project-safari/zebra"
)

type TypeStore struct {
	lock      sync.RWMutex
	resources *zebra.ResourceMap
}

// Return new type store pointer given resource map.
func NewTypeStore(resources *zebra.ResourceMap) *TypeStore {
	typestore := &TypeStore{
		lock: sync.RWMutex{},
		resources: func() *zebra.ResourceMap {
			dest := zebra.NewResourceMap(nil)
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
	ts.lock.Lock()
	defer ts.lock.Unlock()

	ts.resources = nil

	return nil
}

func (ts *TypeStore) Clear() error {
	ts.lock.Lock()
	defer ts.lock.Unlock()

	ts.resources = zebra.NewResourceMap(ts.resources.GetFactory())

	return nil
}

// Return all resources in a ResourceMap.
func (ts *TypeStore) Load() (*zebra.ResourceMap, error) {
	ts.lock.RLock()
	defer ts.lock.RUnlock()

	resources := zebra.NewResourceMap(nil)

	zebra.CopyResourceMap(resources, ts.resources)

	return resources, nil
}

// Create a resource. If a resource with this ID already exists, return error.
func (ts *TypeStore) Create(res zebra.Resource) error {
	if err := res.Validate(context.Background()); err != nil {
		return err
	}

	ts.lock.Lock()
	defer ts.lock.Unlock()

	return ts.create(res)
}

// Should not be called without holding the write lock.
func (ts *TypeStore) create(res zebra.Resource) error {
	// Check if resource already exists
	if _, err := ts.find(res.GetID(), res.GetType()); err == nil {
		return zebra.ErrCreateExists
	}

	ts.resources.Add(res, res.GetType())

	return nil
}

// Update a resource. Return error if resource does not exist.
func (ts *TypeStore) Update(res zebra.Resource) error {
	if err := res.Validate(context.Background()); err != nil {
		return err
	}

	ts.lock.Lock()
	defer ts.lock.Unlock()

	oldRes, err := ts.find(res.GetID(), res.GetType())
	// If resource does not exist, return error.
	if err != nil {
		return zebra.ErrUpdateNoExist
	}

	_ = ts.delete(oldRes)

	_ = ts.create(res)

	return nil
}

// Delete a resource.
func (ts *TypeStore) Delete(res zebra.Resource) error {
	if err := res.Validate(context.Background()); err != nil {
		return err
	}

	ts.lock.Lock()
	defer ts.lock.Unlock()

	return ts.delete(res)
}

// Should not be called without holding the write lock.
func (ts *TypeStore) delete(res zebra.Resource) error {
	ts.resources.Delete(res, res.GetType())

	return nil
}

// Return all resources of given types in a ResourceMap.
func (ts *TypeStore) Query(types []string) *zebra.ResourceMap {
	factory := ts.resources.GetFactory()
	retMap := zebra.NewResourceMap(factory)

	for _, t := range types {
		list := zebra.NewResourceList(factory)
		zebra.CopyResourceList(list, ts.resources.Resources[t])
		retMap.Resources[t] = list
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
