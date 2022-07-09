package zebra

import "errors"

var (
	ErrNotFound        = errors.New("resource not found in store")
	ErrInvalidResource = errors.New("create/delete on invalid resource")
)

// Store interface requires basic store functionalities.
type Store interface {
	Initialize() error
	Wipe() error
	Clear() error
	Load() (*ResourceMap, error)
	Create(res Resource) error
	Delete(res Resource) error
}
