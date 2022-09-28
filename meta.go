package zebra

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

const (
	ShortIDSize        = 7
	DefaultMaxDuration = 4
)

var (
	ErrNameEmpty = errors.New("name is empty")
	ErrIDEmpty   = errors.New("id is empty")
	ErrIDShort   = errors.New("id must be at least 7 characters long")
)

// Meta defines the resource data that is required to be maintained by all the
// resource types.
type Meta struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	Type             Type      `json:"type"`
	Owner            string    `json:"owner"`
	CreationTime     time.Time `json:"creationTime"`
	ModificationTime time.Time `json:"modificationTime"`
	Labels           Labels    `json:"labels"`
}

func NewMeta(resType Type, name string, group string, owner string) Meta {
	t := time.Now()
	labels := Labels{}
	id := uuid.New().String()

	labels.Add("system.group", group)

	if name == "" {
		// Make the short ID as the defaul name
		name = id[:7]
	}

	return Meta{
		ID:               id,
		Name:             name,
		Type:             resType,
		Owner:            owner,
		CreationTime:     t,
		ModificationTime: t,
		Labels:           labels,
	}
}

// Validate returns an error if the given BaseResource object has incorrect values.
// Else, it returns nil.
func (r Meta) Validate() error {
	switch {
	case r.ID == "":
		return ErrIDEmpty
	case len(r.ID) < ShortIDSize:
		return ErrIDShort
	case r.Name == "":
		return ErrNameEmpty
	}

	if err := r.Type.Validate(); err != nil {
		return err
	}

	if err := r.Labels.Validate(); err != nil {
		return err
	}

	return nil
}
