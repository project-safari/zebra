package zebra

import "errors"

var (
	ErrCreateExists    = errors.New("called create on a resource that already exists")
	ErrUpdateNoExist   = errors.New("called update on a resource that does not exist")
	ErrNotFound        = errors.New("resource not found in store")
	ErrInvalidResource = errors.New("create/update/delete on invalid resource")
)

// Store interface requires basic store functionalities.
type Store interface {
	Initialize() error
	Wipe() error
	Clear() error
	Load() (*ResourceMap, error)
	Create(res Resource) error
	Update(res Resource) error
	Delete(res Resource) error
}
