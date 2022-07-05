package zebra

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
