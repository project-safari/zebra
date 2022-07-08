package resourcestore

import (
	"context"
	"sync"

	"github.com/project-safari/zebra"
)

type ResourceStore struct {
	lock      sync.RWMutex
	factory   zebra.ResourceFactory
	resources map[string]zebra.Resource
}

// Return new resource store pointer given resource map.
func NewResourceStore(resources *zebra.ResourceMap) *ResourceStore {
	rs := &ResourceStore{
		lock:    sync.RWMutex{},
		factory: resources.GetFactory(),
		resources: func() map[string]zebra.Resource {
			resMap := make(map[string]zebra.Resource)
			for _, l := range resources.Resources {
				for _, res := range l.Resources {
					resMap[res.GetID()] = res
				}
			}

			return resMap
		}(),
	}

	return rs
}

func (rs *ResourceStore) Initialize() error {
	return nil
}

func (rs *ResourceStore) Wipe() error {
	rs.lock.Lock()
	defer rs.lock.Unlock()

	rs.resources = nil

	return nil
}

func (rs *ResourceStore) Clear() error {
	rs.lock.Lock()
	defer rs.lock.Unlock()

	rs.resources = make(map[string]zebra.Resource)

	return nil
}

// Return all resources in a ResourceMap with UUID key and res in list val.
func (rs *ResourceStore) Load() (*zebra.ResourceMap, error) {
	rs.lock.RLock()
	defer rs.lock.RUnlock()

	resMap := zebra.NewResourceMap(rs.factory)

	for key, val := range rs.resources {
		resMap.Add(val, key)
	}

	return resMap, nil
}

// Create a resource. If a resource with this ID already exists, return error.
func (rs *ResourceStore) Create(res zebra.Resource) error {
	if err := res.Validate(context.Background()); err != nil {
		return err
	}

	rs.lock.Lock()
	defer rs.lock.Unlock()

	return rs.create(res)
}

// Should not be called without holding the write lock.
func (rs *ResourceStore) create(res zebra.Resource) error {
	// Check if resource already exists
	if _, err := rs.find(res.GetID()); err == nil {
		return zebra.ErrCreateExists
	}

	rs.resources[res.GetID()] = res

	return nil
}

// Update a resource. Return error if resource does not exist.
func (rs *ResourceStore) Update(res zebra.Resource) error {
	if err := res.Validate(context.Background()); err != nil {
		return err
	}

	rs.lock.Lock()
	defer rs.lock.Unlock()

	oldRes, err := rs.find(res.GetID())
	// If resource does not exist, return error.
	if err != nil {
		return zebra.ErrUpdateNoExist
	}

	_ = rs.delete(oldRes)

	_ = rs.create(res)

	return nil
}

// Delete a resource.
func (rs *ResourceStore) Delete(res zebra.Resource) error {
	if err := res.Validate(context.Background()); err != nil {
		return err
	}

	rs.lock.Lock()
	defer rs.lock.Unlock()

	return rs.delete(res)
}

// Should not be called without holding the write lock.
func (rs *ResourceStore) delete(res zebra.Resource) error {
	delete(rs.resources, res.GetID())

	return nil
}

// Return all resources of given UUIDs in a ResourceMap.
func (rs *ResourceStore) Query(ids []string) *zebra.ResourceMap {
	retMap := zebra.NewResourceMap(rs.factory)

	for _, id := range ids {
		res, ok := rs.resources[id]
		if !ok {
			return nil
		}

		retMap.Add(res, res.GetType())
	}

	return retMap
}

// Find given resource in ResourceStore. If not found, return nil and error.
// If found, return resource and nil.
func (rs *ResourceStore) find(resID string) (zebra.Resource, error) {
	res, ok := rs.resources[resID]
	if !ok {
		return nil, zebra.ErrNotFound
	}

	return res, nil
}
