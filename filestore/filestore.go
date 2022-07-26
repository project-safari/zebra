package filestore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"syscall"

	"github.com/hashicorp/go-multierror"
	"github.com/project-safari/zebra"
)

const RWRR = os.FileMode(0o644)

// FileStore implements Store.
type FileStore struct {
	storageRoot string
	factory     zebra.ResourceFactory
}

var ErrTypeInvalid = errors.New("resource type invalid")

var ErrFolderInvalid = errors.New("folder invalid")

var ErrFileInvalid = errors.New("file invalid")

var ErrTypeUnpack = errors.New("unpack failed, resource type error")

var ErrNoType = errors.New("resource has no type field")

var ErrFactoryNil = errors.New("resource factory is nil for filestore")

// Return new FileStore pointer set with storageRoot root, lock, and map of type
// name keys with corresponding constructor function values.
func NewFileStore(root string, resourceFactory zebra.ResourceFactory) *FileStore {
	return &FileStore{
		storageRoot: root,
		factory:     resourceFactory,
	}
}

// Initialize store given path. Path is relative to current file location.
// If folders already exist, do nothing (existing store is unchanged).
func (f *FileStore) Initialize() error {
	return f.init()
}

// init implements the store initialization. This function must never be called
// without holding the write lock.
func (f *FileStore) init() error {
	location := f.filestoreResourcesPath()
	err := os.MkdirAll(location, os.ModePerm)

	if err != nil && !os.IsExist(err) {
		return err
	}

	for i := 0; i < 256; i++ {
		h := fmt.Sprintf("%02x", i)
		err := os.Mkdir(path.Join(location, h), os.ModePerm)

		if err != nil && !os.IsExist(err) {
			return err
		}
	}

	return nil
}

// Wipe store given path. Path is relative to current file location.
// If store does not exist, do nothing.
func (f *FileStore) Wipe() error {
	return os.RemoveAll(f.filestoreResourcesPath())
}

// Clear store given path (i.e. delete all resource objects). Path is relative
// to current file location. If store does not exist, create store.
func (f *FileStore) Clear() error {
	if err := os.RemoveAll(f.filestoreResourcesPath()); err != nil {
		return err
	}

	return f.init() // lock is held
}

// Load objects from filestore storageRoot.
// Return resources as ResourceMap where keys are types.
func (f *FileStore) Load() (*zebra.ResourceMap, error) {
	var retErr error

	rootDir := f.filestoreResourcesPath()

	resources := zebra.NewResourceMap(f.factory)

	dirs, err := os.ReadDir(rootDir)
	if err != nil {
		return nil, err
	}

	for _, subdir := range dirs {
		files, err := os.ReadDir(path.Join(rootDir, subdir.Name()))
		if err != nil {
			return nil, err
		}

		for _, file := range files {
			contents, err := os.ReadFile(path.Join(rootDir, subdir.Name(), file.Name()))
			if err != nil {
				return nil, err
			}

			object := make(map[string]interface{})

			err = json.Unmarshal(contents, &object)
			if err != nil {
				return nil, err
			}

			resType, ok := object["type"]
			if !ok {
				retErr = ErrNoType

				continue
			}

			vType, err := zebra.TypeChecker(resType)
			if err != nil {
				retErr = ErrNoType

				continue
			}

			newRes, err := f.unpackResource(contents, vType.Name)
			if err != nil {
				retErr = err

				continue
			}

			resources.Add(newRes, vType.Name)
		}
	}

	return resources, retErr
}

// Store new object given storage root path and resource pointer.
// If object already exists, update.
func (f *FileStore) Create(res zebra.Resource) error {
	if _, err := os.Stat(f.resourcesFilePath(res)); err == nil {
		return f.update(res)
	}

	dir := f.resourcesFolderPath(res)

	object, err := json.Marshal(res)
	if err != nil {
		return err
	}

	cleanup := func(f *os.File, err error) error {
		errs := multierror.Append(nil, err)

		if e := f.Close(); e != nil {
			errs = multierror.Append(errs, e)
		}

		if e := os.Remove(f.Name()); e != nil {
			errs = multierror.Append(errs, e)
		}

		return errs
	}

	file, err := ioutil.TempFile(dir, "temp_")
	if err != nil {
		return err
	}

	if _, err := file.Write(object); err != nil {
		return cleanup(file, err)
	}

	if err := file.Chmod(RWRR); err != nil {
		return cleanup(file, err)
	}

	if err := file.Close(); err != nil {
		return cleanup(file, err)
	}

	return os.Rename(file.Name(), f.resourcesFilePath(res))
}

// Update existing object.
func (f *FileStore) update(res zebra.Resource) error {
	if err := f.Delete(res); err != nil {
		return err
	}

	return f.Create(res)
}

// Delete object given storage root path and UUID.
// If object does not exist, return error.
func (f *FileStore) Delete(res zebra.Resource) error {
	path := f.resourcesFilePath(res)

	// attempt to delete resource that does not exist, just return nil
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return nil
	}

	if err := syscall.Unlink(path); err != nil {
		return err
	}

	return nil
}

// Unpack storedRes.Resource into correct type of resource and return zebra.Resource
// along with error if occurred.
func (f *FileStore) unpackResource(contents []byte, resType string) (zebra.Resource, error) {
	if f.factory == nil {
		return nil, ErrFactoryNil
	}

	res := f.factory.New(resType)
	if res == nil {
		return nil, ErrTypeUnpack
	}

	if err := json.Unmarshal(contents, res); err != nil {
		return nil, err
	}

	if err := res.Validate(context.Background()); err != nil {
		return nil, err
	}

	return res, nil
}

// Return file path given resource.
func (f *FileStore) resourcesFilePath(res zebra.Resource) string {
	resID := res.GetID()

	return path.Join(f.storageRoot, "resources", resID[:2], resID[2:])
}

// Return folder path given resource.
func (f *FileStore) resourcesFolderPath(res zebra.Resource) string {
	resID := res.GetID()

	return path.Join(f.storageRoot, "resources", resID[:2])
}

// Return path to filestore resources folder.
func (f *FileStore) filestoreResourcesPath() string {
	return path.Join(f.storageRoot, "resources")
}
