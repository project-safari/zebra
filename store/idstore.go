package store

import (
	"github.com/project-safari/zebra"
)

type IDStore struct {
	factory   zebra.ResourceFactory
	resources map[string]zebra.Resource
}

// Return new resource store pointer given resource map.
func NewIDStore(resources *zebra.ResourceMap) *IDStore {
	ids := &IDStore{
		factory: resources.Factory(),
		resources: func() map[string]zebra.Resource {
			resMap := make(map[string]zebra.Resource)
			for _, l := range resources.Resources {
				for _, res := range l.Resources {
					resMap[res.GetMeta().ID] = res
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
	ids.resources = nil

	return nil
}

func (ids *IDStore) Clear() error {
	ids.resources = make(map[string]zebra.Resource)

	return nil
}

// Return all resources in a ResourceMap with UUID key and res in list val.
func (ids *IDStore) Load() (*zebra.ResourceMap, error) {
	resMap := zebra.NewResourceMap(ids.factory)

	for _, val := range ids.resources {
		if e := resMap.Add(val); e != nil {
			return nil, e
		}
	}

	return resMap, nil
}

// Create a resource. If a resource with this ID already exists, update.
func (ids *IDStore) Create(res zebra.Resource) error {
	// Check if resource already exists
	if oldRes, err := ids.find(res.GetMeta().ID); err == nil {
		return ids.update(oldRes, res)
	}

	ids.resources[res.GetMeta().ID] = res

	return nil
}

// Update a resource.
func (ids *IDStore) update(oldRes zebra.Resource, res zebra.Resource) error {
	if err := ids.Delete(oldRes); err != nil {
		return err
	}

	return ids.Create(res)
}

// to add any necessary funcs. for edit.

// Delete a resource.
func (ids *IDStore) Delete(res zebra.Resource) error {
	delete(ids.resources, res.GetMeta().ID)

	return nil
}

// Return all resources of given UUIDs in a ResourceMap.
// If resource with id cannot be found, skip.
func (ids *IDStore) Query(uuids []string) *zebra.ResourceMap {
	retMap := zebra.NewResourceMap(ids.factory)

	for _, id := range uuids {
		res, ok := ids.resources[id]
		if !ok {
			continue
		}

		if e := retMap.Add(res); e != nil {
			return nil
		}
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
