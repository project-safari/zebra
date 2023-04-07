package store

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"sync"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/model/lease"
)

var ErrNilResource = errors.New("nil resource not allowed")

type ResourceStore struct {
	lock        sync.RWMutex
	StorageRoot string
	Factory     zebra.ResourceFactory
	fs          *FileStore
	ids         *IDStore
	ls          *LabelStore
	ts          *TypeStore
	qu          *Queue
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
		qu:          nil,
	}
}

func (rs *ResourceStore) Initialize() error {
	rs.lock.Lock()
	defer rs.lock.Unlock()

	rs.fs = NewFileStore(rs.StorageRoot, rs.Factory)
	if err := rs.fs.Initialize(); err != nil {
		return err
	}

	resources, err := rs.fs.Load()
	if err != nil {
		return err
	}

	rs.ids = NewIDStore(resources)
	rs.ls = NewLabelStore(resources)
	rs.ts = NewTypeStore(resources)
	rs.qu = NewQueue(rs)

	return nil
}

func (rs *ResourceStore) Wipe() error {
	rs.lock.Lock()
	defer rs.lock.Unlock()

	rs.fs = nil
	rs.ids = nil
	rs.ls = nil
	rs.ts = nil
	rs.qu = nil

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

	defer rs.qu.LeaseSatisfied()
	defer rs.qu.Process()
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

	if l, ok := res.(*lease.Lease); ok {
		rs.qu.Enqueue(l)
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

	retMap := zebra.NewResourceMap(resMap.Factory())

	zebra.CopyResourceMap(retMap, resMap)

	return retMap
}

// Return resources with matching UUIDs.
func (rs *ResourceStore) QueryUUID(uuids []string) *zebra.ResourceMap {
	rs.lock.RLock()
	defer rs.lock.RUnlock()

	resMap := rs.ids.Query(uuids)
	retMap := zebra.NewResourceMap(resMap.Factory())

	zebra.CopyResourceMap(retMap, resMap)

	return retMap
}

// Return resources with matching types.
func (rs *ResourceStore) QueryType(types []string) *zebra.ResourceMap {
	rs.lock.RLock()
	defer rs.lock.RUnlock()

	resMap := rs.ts.Query(types)
	retMap := zebra.NewResourceMap(resMap.Factory())

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
	retMap := zebra.NewResourceMap(resMap.Factory())

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

	for _, l := range resMap.Resources {
		for _, res := range l.Resources {
			val := FieldByName(reflect.ValueOf(res).Elem(), query.Key).String()
			inList := zebra.IsIn(val, query.Values)

			if (inVals && inList) || (!inVals && !inList) {
				if e := retMap.Add(res); e != nil {
					return nil, e
				}
			}
		}
	}

	return retMap, nil
}

// Filter given map by uuids.
func FilterUUID(uuids []string, resMap *zebra.ResourceMap) (*zebra.ResourceMap, error) {
	retMap := zebra.NewResourceMap(resMap.Factory())

	for _, l := range resMap.Resources {
		for _, res := range l.Resources {
			if zebra.IsIn(res.GetMeta().ID, uuids) {
				if e := retMap.Add(res); e != nil {
					return nil, e
				}
			}
		}
	}

	return retMap, nil
}

// Filter given map by types.
func FilterType(types []string, resMap *zebra.ResourceMap) (*zebra.ResourceMap, error) {
	f := resMap.Factory()
	retMap := zebra.NewResourceMap(f)

	for _, t := range types {
		l, ok := resMap.Resources[t]
		if !ok {
			continue
		}

		c, ok := f.Constructor(t)
		if !ok {
			return nil, zebra.ErrNotFound
		}

		copyL := zebra.NewResourceList(c)

		zebra.CopyResourceList(copyL, l)
		retMap.Resources[t] = copyL
	}

	return retMap, nil
}

// Filter given map by label name and val.
func FilterLabel(query zebra.Query, resMap *zebra.ResourceMap) (*zebra.ResourceMap, error) { //nolint:cyclop
	if err := query.Validate(); err != nil {
		return resMap, err
	}

	retMap := zebra.NewResourceMap(resMap.Factory())

	inVals := false

	if query.Op == zebra.MatchEqual || query.Op == zebra.MatchIn {
		inVals = true
	}

	for _, l := range resMap.Resources {
		for _, res := range l.Resources {
			labels := res.GetMeta().Labels
			matchIn := labels.MatchIn(query.Key, query.Values...)

			if (inVals && matchIn) || (!inVals && !matchIn) {
				if e := retMap.Add(res); e != nil {
					return nil, e
				}
			}
		}
	}

	return retMap, nil
}

// Filter given map by property name (case insensitive) and val.
func FilterProperty(query zebra.Query, resMap *zebra.ResourceMap) (*zebra.ResourceMap, error) { //nolint:cyclop
	if err := query.Validate(); err != nil {
		return resMap, err
	}

	retMap := zebra.NewResourceMap(resMap.Factory())

	inVals := false

	if query.Op == zebra.MatchEqual || query.Op == zebra.MatchIn {
		inVals = true
	}

	for _, l := range resMap.Resources {
		for _, res := range l.Resources {
			val := FieldByName(reflect.ValueOf(res).Elem(), query.Key).String()
			matchIn := zebra.IsIn(val, query.Values)

			if (inVals && matchIn) || (!inVals && !matchIn) {
				if e := retMap.Add(res); e != nil {
					return nil, e
				}
			}
		}
	}

	return retMap, nil
}

// Set Lease state of all resources in Resource Request.
func (rs *ResourceStore) FreeResources(reslist []string) {
	for _, id := range reslist {
		res, ok := rs.ids.resources[id]
		if !ok {
			continue
		}

		res.UpdateStatus().UpdateLeaseState(0)

		if err := rs.Create(res); err != nil {
			continue
		}
	}
}

// Ignore case in returning value of given field.
func FieldByName(v reflect.Value, field string) reflect.Value {
	field = strings.ToLower(field)

	return v.FieldByNameFunc(
		func(found string) bool {
			return strings.ToLower(found) == field
		})
}
