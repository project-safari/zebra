package store

import (
	"context"
	"errors"
	"reflect"
	"sync"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/filestore"
	"github.com/project-safari/zebra/idstore"
	"github.com/project-safari/zebra/labelstore"
	"github.com/project-safari/zebra/typestore"
)

var ErrNilResource = errors.New("nil resource not allowed")

type ResourceStore struct {
	lock        sync.RWMutex
	StorageRoot string
	Factory     zebra.ResourceFactory
	fs          *filestore.FileStore
	ids         *idstore.IDStore
	ls          *labelstore.LabelStore
	ts          *typestore.TypeStore
}

func NewResourceStore(root string, factory zebra.ResourceFactory) *ResourceStore {
	return &ResourceStore{
		lock:        sync.RWMutex{},
		StorageRoot: root,
		Factory:     factory,
		fs:          nil,
		ids:         nil,
		ls:          nil,
		ts:          nil,
	}
}

func (rs *ResourceStore) Initialize() error {
	rs.lock.Lock()
	defer rs.lock.Unlock()

	rs.fs = filestore.NewFileStore(rs.StorageRoot, rs.Factory)
	if err := rs.fs.Initialize(); err != nil {
		return err
	}

	resources, err := rs.fs.Load()
	if err != nil {
		return err
	}

	rs.ids = idstore.NewIDStore(resources)
	rs.ls = labelstore.NewLabelStore(resources)
	rs.ts = typestore.NewTypeStore(resources)

	return nil
}

func (rs *ResourceStore) Wipe() error {
	rs.lock.Lock()
	defer rs.lock.Unlock()

	rs.fs = nil
	rs.ids = nil
	rs.ls = nil
	rs.ts = nil

	return nil
}

func (rs *ResourceStore) Clear() error {
	rs.lock.Lock()
	defer rs.lock.Unlock()

	if err := rs.fs.Clear(); err != nil {
		return err
	}

	if err := rs.ids.Clear(); err != nil {
		return err
	}

	if err := rs.ls.Clear(); err != nil {
		return err
	}

	if err := rs.ts.Clear(); err != nil {
		return err
	}

	return nil
}

// Return ResourceMap with resource type as key and list of resources as val.
func (rs *ResourceStore) Load() (*zebra.ResourceMap, error) {
	rs.lock.RLock()
	defer rs.lock.RUnlock()

	return rs.ts.Load()
}

func (rs *ResourceStore) Create(res zebra.Resource) error {
	if res == nil {
		return ErrNilResource
	}

	if err := res.Validate(context.Background()); err != nil {
		return err
	}

	rs.lock.Lock()
	defer rs.lock.Unlock()

	err := rs.fs.Create(res)
	if err != nil {
		return err
	}

	err = rs.ids.Create(res)
	if err != nil {
		return err
	}

	err = rs.ls.Create(res)
	if err != nil {
		return err
	}

	err = rs.ts.Create(res)
	if err != nil {
		return err
	}

	return nil
}

func (rs *ResourceStore) Delete(resource zebra.Resource) error {
	if resource == nil || resource.Validate(context.Background()) != nil {
		return zebra.ErrInvalidResource
	}

	rs.lock.Lock()
	defer rs.lock.Unlock()

	err := rs.fs.Delete(resource)
	if err != nil {
		return err
	}

	err = rs.ids.Delete(resource)
	if err != nil {
		return err
	}

	err = rs.ls.Delete(resource)
	if err != nil {
		return err
	}

	err = rs.ts.Delete(resource)
	if err != nil {
		return err
	}

	return nil
}

// Return all resources in a ResourceMap.
func (rs *ResourceStore) Query() *zebra.ResourceMap {
	rs.lock.RLock()
	defer rs.lock.RUnlock()

	resMap, err := rs.ts.Load()
	if err != nil {
		return nil
	}

	retMap := zebra.NewResourceMap(resMap.GetFactory())

	zebra.CopyResourceMap(retMap, resMap)

	return retMap
}

// Return resources with matching UUIDs.
func (rs *ResourceStore) QueryUUID(uuids []string) *zebra.ResourceMap {
	rs.lock.RLock()
	defer rs.lock.RUnlock()

	resMap := rs.ids.Query(uuids)
	retMap := zebra.NewResourceMap(resMap.GetFactory())

	zebra.CopyResourceMap(retMap, resMap)

	return retMap
}

// Return resources with matching types.
func (rs *ResourceStore) QueryType(types []string) *zebra.ResourceMap {
	rs.lock.RLock()
	defer rs.lock.RUnlock()

	resMap := rs.ts.Query(types)
	retMap := zebra.NewResourceMap(resMap.GetFactory())

	zebra.CopyResourceMap(retMap, resMap)

	return retMap
}

// Return resources with matching label.
func (rs *ResourceStore) QueryLabel(query zebra.Query) (*zebra.ResourceMap, error) {
	if err := query.Validate(); err != nil {
		return nil, err
	}

	rs.lock.RLock()
	defer rs.lock.RUnlock()

	resMap := rs.ls.Query(query)
	retMap := zebra.NewResourceMap(resMap.GetFactory())

	zebra.CopyResourceMap(retMap, resMap)

	return retMap, nil
}

// Return resources which match given property/value(s).
// Naive search implementation, >= O(n) for n resources.
func (rs *ResourceStore) QueryProperty(query zebra.Query) (*zebra.ResourceMap, error) {
	if err := query.Validate(); err != nil {
		return nil, err
	}

	rs.lock.RLock()
	defer rs.lock.RUnlock()

	if query.Op == zebra.MatchEqual || query.Op == zebra.MatchIn {
		return rs.propertyMatch(query, true)
	}

	return rs.propertyMatch(query, false)
}

func (rs *ResourceStore) propertyMatch(query zebra.Query, inVals bool) (*zebra.ResourceMap, error) {
	resMap, err := rs.ts.Load()
	if err != nil {
		return nil, err
	}

	retMap := zebra.NewResourceMap(rs.Factory)

	for t, l := range resMap.Resources {
		for _, res := range l.Resources {
			val := zebra.FieldByName(reflect.ValueOf(res).Elem(), query.Key).String()
			inList := zebra.IsIn(val, query.Values)

			if inVals && inList {
				retMap.Add(res, t)
			} else if !inVals && !inList {
				retMap.Add(res, t)
			}
		}
	}

	return retMap, nil
}
