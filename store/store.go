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

// StoredResource allows information to be loaded from a store and thus queried
// by type and/or other attributes.
type storedResource struct {
	Type     string `json:"type"`
	Resource []byte `json:"resource"`
}

// FileStore implements Store.
type FileStore struct {
	lock        sync.Mutex
	storageRoot string
	types       map[string]func() zebra.Resource
}

var ErrTypeInvalid = errors.New("resource type invalid")

var ErrFolderInvalid = errors.New("folder invalid")

var ErrFileInvalid = errors.New("file invalid")

var ErrTypeUnpack = errors.New("unpack failed, resource type error")

var ErrFileExists = errors.New("called create on resource that already exists")

var ErrFileDoesNotExist = errors.New("called update on resource that does not exist")

// Return new FileStore pointer set with storageRoot root, lock, and map of type
// name keys with corresponding constructor function values.
func NewFileStore(root string, types map[string]func() zebra.Resource) *FileStore {
	return &FileStore{
		lock:        sync.Mutex{},
		storageRoot: root,
		types: func() map[string]func() zebra.Resource {
			typeMap := make(map[string]func() zebra.Resource, len(types))

			// Make a copy of types so that they are not mutated after the store has
			// been created and initialized.
			for t, r := range types {
				typeMap[t] = r
			}

			return typeMap
		}(),
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
// Return list of resources.
func (f *FileStore) Load() (map[string]zebra.Resource, error) {
	f.lock.Lock()
	defer f.lock.Unlock()

	rootDir := f.filestoreResourcesPath()

	resources := make(map[string]zebra.Resource)

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

			storedRes := new(storedResource)

			err = json.Unmarshal(contents, storedRes)
			if err != nil {
				return nil, err
			}

			resID := subdir.Name() + file.Name()

			resources[resID], err = f.unpackResource(storedRes)
			if err != nil {
				return nil, err
			}
		}
	}

	return resources, nil
}

// Store new object given storage root path and resource pointer.
// If object already exists, return error.
func (f *FileStore) Create(res zebra.Resource) error {
	f.lock.Lock()
	defer f.lock.Unlock()

	if _, err := os.Stat(f.resourcesFilePath(res)); err == nil {
		return ErrFileExists
	}

	dir := f.resourcesFolderPath(res)

	object, err := json.Marshal(res)
	if err != nil {
		return err
	}

	resType := res.GetType()

	storedRes, err := json.Marshal(storedResource{Type: resType, Resource: object})
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

	if _, err := file.Write(storedRes); err != nil {
		return err
	}

	if err := file.Chmod(0o644); err != nil { //nolint:gomnd
		return err
	}

	return os.Rename(file.Name(), f.resourcesFilePath(res))
}

// Update existing object. If object does not exist, return error.
func (f *FileStore) Update(res zebra.Resource) error {
	filepath := f.resourcesFilePath(res)
	if _, err := os.Stat(filepath); err == nil {
		os.Remove(filepath)

		return f.Create(res)
	}

	return ErrFileDoesNotExist
}

// Delete object given storage root path and UUID.
// If object does not exist, do nothing.
func (f *FileStore) Delete(res zebra.Resource) error {
	f.lock.Lock()
	defer f.lock.Unlock()

	path := f.resourcesFilePath(res)

	if err := syscall.Unlink(path); err != nil {
		return err
	}

	return nil
}

// Unpack storedRes.Resource into correct type of resource and return zebra.Resource
// along with error if occurred.
func (f *FileStore) unpackResource(storedRes *storedResource) (zebra.Resource, error) {
	creator, ok := f.types[storedRes.Type]
	if !ok {
		return nil, ErrTypeUnpack
	}

	res := creator()
	if err := json.Unmarshal(storedRes.Resource, res); err != nil {
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
