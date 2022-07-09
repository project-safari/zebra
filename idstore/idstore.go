package idstore

import (
	"sync"

	"github.com/project-safari/zebra"
)

type IDStore struct {
	lock      sync.RWMutex
	factory   zebra.ResourceFactory
	resources map[string]zebra.Resource
}

// Return new resource store pointer given resource map.
func NewIDStore(resources *zebra.ResourceMap) *IDStore {
	ids := &IDStore{
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

	return ids
}

func (ids *IDStore) Initialize() error {
	return nil
}

func (ids *IDStore) Wipe() error {
	ids.lock.Lock()
	defer ids.lock.Unlock()

	ids.resources = nil

	return nil
}

func (ids *IDStore) Clear() error {
	ids.lock.Lock()
	defer ids.lock.Unlock()

	ids.resources = make(map[string]zebra.Resource)

	return nil
}

// Return all resources in a ResourceMap with UUID key and res in list val.
func (ids *IDStore) Load() (*zebra.ResourceMap, error) {
	ids.lock.RLock()
	defer ids.lock.RUnlock()

	resMap := zebra.NewResourceMap(ids.factory)

	for key, val := range ids.resources {
		resMap.Add(val, key)
	}

	return resMap, nil
}

// Create a resource. If a resource with this ID already exists, return error.
func (ids *IDStore) Create(res zebra.Resource) error {
	ids.lock.Lock()
	defer ids.lock.Unlock()

	return ids.create(res)
}

// Should not be called without holding the write lock.
func (ids *IDStore) create(res zebra.Resource) error {
	// Check if resource already exists
	if _, err := ids.find(res.GetID()); err == nil {
		return zebra.ErrCreateExists
	}

	ids.resources[res.GetID()] = res

	return nil
}

// Update a resource. Return error if resource does not exist.
func (ids *IDStore) Update(res zebra.Resource) error {
	ids.lock.Lock()
	defer ids.lock.Unlock()

	oldRes, err := ids.find(res.GetID())
	// If resource does not exist, return error.
	if err != nil {
		return zebra.ErrUpdateNoExist
	}

	_ = ids.delete(oldRes)

	_ = ids.create(res)

	return nil
}

// Delete a resource.
func (ids *IDStore) Delete(res zebra.Resource) error {
	ids.lock.Lock()
	defer ids.lock.Unlock()

	return ids.delete(res)
}

// Should not be called without holding the write lock.
func (ids *IDStore) delete(res zebra.Resource) error {
	delete(ids.resources, res.GetID())

	return nil
}

// Return all resources of given UUIDs in a ResourceMap.
func (ids *IDStore) Query(uuids []string) *zebra.ResourceMap {
	retMap := zebra.NewResourceMap(ids.factory)

	for _, id := range uuids {
		res, ok := ids.resources[id]
		if !ok {
			return nil
		}

		retMap.Add(res, res.GetType())
	}

	return retMap
}

// Find given resource in IDStore. If not found, return nil and error.
// If found, return resource and nil.
func (ids *IDStore) find(resID string) (zebra.Resource, error) {
	res, ok := ids.resources[resID]
	if !ok {
		return nil, zebra.ErrNotFound
	}

	return res, nil
}
