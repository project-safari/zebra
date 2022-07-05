package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sync"
	"syscall"

	"github.com/rchamarthy/zebra"
)

// FileStore implements Store.
type FileStore struct {
	lock        sync.Mutex
	storageRoot string
	factory     zebra.ResourceFactory
}

var ErrTypeInvalid = errors.New("resource type invalid")

var ErrFolderInvalid = errors.New("folder invalid")

var ErrFileInvalid = errors.New("file invalid")

var ErrTypeUnpack = errors.New("unpack failed, resource type error")

var ErrFileExists = errors.New("called create on resource that already exists")

var ErrFileDoesNotExist = errors.New("called update on resource that does not exist")

var ErrNoType = errors.New("resource has no type field")

// Return new FileStore pointer set with storageRoot root, lock, and map of type
// name keys with corresponding constructor function values.
func NewFileStore(root string, resourceFactory zebra.ResourceFactory) *FileStore {
	return &FileStore{
		lock:        sync.Mutex{},
		storageRoot: root,
		factory:     resourceFactory,
	}
}

// Initialize store given path. Path is relative to current file location.
// If folders already exist, do nothing (existing store is unchanged).
func (f *FileStore) Initialize() error {
	f.lock.Lock()
	defer f.lock.Unlock()

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
	f.lock.Lock()
	defer f.lock.Unlock()

	return os.RemoveAll(f.filestoreResourcesPath())
}

// Clear store given path (i.e. delete all resource objects). Path is relative
// to current file location. If store does not exist, create store.
func (f *FileStore) Clear() error {
	f.lock.Lock()
	defer f.lock.Unlock()

	if err := os.RemoveAll(f.filestoreResourcesPath()); err != nil {
		return err
	}

	return f.init() // lock is held
}

// Load objects from filestore storageRoot.
// Return resources as ResourceMap where keys are types.
func (f *FileStore) Load() (*zebra.ResourceMap, error) {
	f.lock.Lock()
	defer f.lock.Unlock()

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

			resType, ok := object["type"].(string)
			if !ok {
				retErr = ErrNoType

				continue
			}

			newRes, err := f.unpackResource(contents, resType)
			if err != nil {
				retErr = err

				continue
			}

			resources.Add(newRes, resType)
		}
	}

	return resources, retErr
}

// Store new object given storage root path and resource pointer.
// If object already exists, return error.
func (f *FileStore) Create(res zebra.Resource) error {
	f.lock.Lock()
	defer f.lock.Unlock()

	return f.create(res)
}

// Should not be called without holding the write lock.
func (f *FileStore) create(res zebra.Resource) error {
	if _, err := os.Stat(f.resourcesFilePath(res)); err == nil {
		return ErrFileExists
	}

	dir := f.resourcesFolderPath(res)

	object, err := json.Marshal(res)
	if err != nil {
		return err
	}

	file, err := ioutil.TempFile(dir, "temp_")
	if err != nil {
		return err
	}

	defer func() {
		if _, err := os.Stat(file.Name()); err != nil {
			os.Remove(file.Name())
		}
	}()

	defer file.Close()

	if _, err := file.Write(object); err != nil {
		return err
	}

	if err := file.Chmod(0o644); err != nil { //nolint:gomnd
		return err
	}

	return os.Rename(file.Name(), f.resourcesFilePath(res))
}

// Update existing object. If object does not exist, return error.
func (f *FileStore) Update(res zebra.Resource) error {
	f.lock.Lock()
	defer f.lock.Unlock()

	filepath := f.resourcesFilePath(res)
	if _, err := os.Stat(filepath); err == nil {
		_ = f.delete(res)

		return f.create(res)
	}

	return ErrFileDoesNotExist
}

// Delete object given storage root path and UUID.
// If object does not exist, do nothing.
func (f *FileStore) Delete(res zebra.Resource) error {
	f.lock.Lock()
	defer f.lock.Unlock()

	return f.delete(res)
}

// Should not be called without holding the write lock.
func (f *FileStore) delete(res zebra.Resource) error {
	path := f.resourcesFilePath(res)

	if err := syscall.Unlink(path); err != nil {
		return err
	}

	return nil
}

// Unpack storedRes.Resource into correct type of resource and return zebra.Resource
// along with error if occurred.
func (f *FileStore) unpackResource(contents []byte, resType string) (zebra.Resource, error) {
	res := f.factory.New(resType)
	if res == nil {
		return nil, ErrTypeUnpack
	}

	if err := json.Unmarshal(contents, res); err != nil {
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
