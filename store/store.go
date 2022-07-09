package store

import (
	"context"
	"sync"

	"github.com/project-safari/zebra"
	"github.com/project-safari/zebra/filestore"
	"github.com/project-safari/zebra/idstore"
	"github.com/project-safari/zebra/labelstore"
	"github.com/project-safari/zebra/typestore"
)

type ResourceStore struct {
	lock sync.RWMutex
	fs   *filestore.FileStore
	ids  *idstore.IDStore
	ls   *labelstore.LabelStore
	ts   *typestore.TypeStore
}

func NewResourceStore(root string, factory zebra.ResourceFactory) *ResourceStore {
	return &ResourceStore{
		lock: sync.RWMutex{},
		fs: func() *filestore.FileStore {
			fs := filestore.NewFileStore(root, factory)
			_ = fs.Initialize()

			return fs
		}(),
		ids: nil,
		ls:  nil,
		ts:  nil,
	}
}

func (rs *ResourceStore) Initialize() error {
	rs.lock.Lock()
	defer rs.lock.Unlock()

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
	if res == nil || res.Validate(context.Background()) != nil {
		return zebra.ErrInvalidResource
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

func (rs *ResourceStore) Delete(res zebra.Resource) error {
	if res == nil || res.Validate(context.Background()) != nil {
		return zebra.ErrInvalidResource
	}

	rs.lock.Lock()
	defer rs.lock.Unlock()

	err := rs.fs.Delete(res)
	if err != nil {
		return err
	}

	err = rs.ids.Delete(res)
	if err != nil {
		return err
	}

	err = rs.ls.Delete(res)
	if err != nil {
		return err
	}

	err = rs.ts.Delete(res)
	if err != nil {
		return err
	}

	return nil
}
