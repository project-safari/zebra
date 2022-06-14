package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"
	"syscall"

	"github.com/rchamarthy/zebra"
)

// Store interface requires basic store functionalities.
type Store interface {
	Initialize() error
	Wipe() error
	Clear() error
	Load() (map[string]zebra.Resource, error)
	Create(res zebra.Resource) error
	Delete(res zebra.Resource) error
}

// FileStore implements Store.
type FileStore struct {
	storageRoot string
	types       map[string]zebra.Resource
}

var ErrTypeInvalid = errors.New("resource type invalid")

var ErrFolderInvalid = errors.New("folder invalid")

var ErrFileInvalid = errors.New("file invalid")

var ErrTypeUnpack = errors.New("unpack failed, resource type error")

// StoredResource allows information to be loaded from a store and thus queried
// by type and/or other attributes.
type storedResource struct {
	Type     string `json:"type"`
	Resource []byte `json:"resource"`
}

// Initialize store given path. Path is relative to current file location.
// If folders already exist, do nothing (existing store is unchanged).
func (f *FileStore) Initialize() error {
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
	err := os.RemoveAll(f.filestoreResourcesPath())
	if err != nil {
		return err
	}

	return nil
}

// Clear store given path (i.e. delete all resource objects). Path is relative
// to current file location. If store does not exist, create store.
func (f *FileStore) Clear() error {
	if err := f.Wipe(); err != nil {
		return err
	}

	if err := f.Initialize(); err != nil {
		return err
	}

	return nil
}

// Load object given storage root path and UUID.
// If object does not exist, return empty resource.
func (f *FileStore) Load() (map[string]zebra.Resource, error) {
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

// Store object given storage root path and resource pointer.
// If object already exists, overwrite.
func (f *FileStore) Create(res zebra.Resource) error {
	dir := f.resourcesFolderPath(res)

	object, err := json.Marshal(res)
	if err != nil {
		return err
	}

	fullName := strings.Split(reflect.TypeOf(res).String(), ".")
	resType := fullName[len(fullName)-1]

	storedRes, err := json.Marshal(storedResource{Type: resType, Resource: object})
	if err != nil {
		return err
	}

	file, err := ioutil.TempFile(dir, "temp_")
	if err != nil {
		return err
	}

	defer func() {
		if _, err := os.Stat(file.Name()); err == nil {
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

// Delete object given storage root path and UUID.
// If object does not exist, do nothing.
func (f *FileStore) Delete(res zebra.Resource) error {
	path := f.resourcesFilePath(res)

	if err := syscall.Unlink(path); err != nil {
		return err
	}

	return nil
}

// Unpack storedRes.Resource into correct type of resource and return zebra.Resource
// along with error if occurred.
func (f *FileStore) unpackResource(storedRes *storedResource) (zebra.Resource, error) {
	resType := reflect.TypeOf(f.types[storedRes.Type]).Elem()

	res, typeOK := reflect.New(resType).Interface().(zebra.Resource)
	if !typeOK {
		return nil, ErrTypeUnpack
	}

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

// Set storage root.
func (f *FileStore) SetStorageRoot(root string) {
	f.storageRoot = root
}

// Set types.
func (f *FileStore) SetTypes(types map[string]zebra.Resource) {
	f.types = types
}
