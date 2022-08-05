package zebra

import (
	"context"
	"errors"
	"log"
	"strings"
	"unicode"
)

// Resource interface is implemented by all resources and provides resource
// validation and label selection methods.
type Resource interface {
	Validate(ctx context.Context) error
	GetID() string
	GetType() Type
	GetLabels() Labels
}

var (
	ErrNameEmpty   = errors.New("name is empty")
	ErrIDEmpty     = errors.New("id is empty")
	ErrIDShort     = errors.New("id must be at least 3 characters long")
	ErrTypeEmpty   = errors.New("type is empty")
	ErrWrongType   = errors.New("type value is incorrect")
	ErrPassLen     = errors.New("password is less than 12 characters long")
	ErrPassCase    = errors.New("password does not contain both upper and lowercase")
	ErrPassNum     = errors.New("password does not contain a number")
	ErrPassSpecial = errors.New("password does not contain a special character")
	ErrNoKeys      = errors.New("keys is nil")
	ErrLabel       = errors.New("missing mandatory system label")
)

// BaseResource must be embedded in all resource structs, ensuring each resource is
// assigned an ID string.
type BaseResource struct {
	ID     string `json:"id"`
	Type   Type   `json:"type"`
	Labels Labels `json:"labels,omitempty"`
	Status Status `json:"status,omitempty"`
}

// Validate returns an error if the given BaseResource object has incorrect values.
// Else, it returns nil.
func (r *BaseResource) Validate(ctx context.Context) error {
	switch {
	case r.ID == "":
		return ErrIDEmpty
	case len(r.ID) < 3: // nolint:gomnd
		return ErrIDShort
	case r.Type.Name == "":
		return ErrTypeEmpty
	}

	if err := r.LabelsValidate(); err != nil {
		return err
	}

	return r.Status.Validate(ctx)
}

// Return ID of BaseResource r.
func (r *BaseResource) GetID() string {
	return r.ID
}

// BaseResource has no name. Return empty string.
func (r *BaseResource) GetType() Type {
	return r.Type
}

// Return labels of BaseResource r.
func (r *BaseResource) GetLabels() Labels {
	dest := make(Labels, len(r.Labels))
	for key, val := range r.Labels {
		dest[key] = val
	}

	return dest
}

// Special label validation to ensure all resources have group label.
func (r *BaseResource) LabelsValidate() error {
	log.Default().Printf("Checking Lables: %v", r.Labels)
	if _, ok := r.Labels["system.group"]; !ok {
		return ErrLabel
	}

	return nil
}

// NamedResource represents all resources assigned both a string ID and a name.
type NamedResource struct {
	BaseResource
	Name string `json:"name"`
}

// Validate returns an error if the given NamedResource object has incorrect values.
// Else, it returns nil.
func (r *NamedResource) Validate(ctx context.Context) error {
	if r.Name == "" {
		return ErrNameEmpty
	}

	return r.BaseResource.Validate(ctx)
}

// Credentials represents a named resource that has a set of keys (where each key is
// an authentication method) with corresponding values (where each value is the
// information to store for the authentication method).
type Credentials struct {
	NamedResource
	Keys map[string]string
}

func CredentialsType() Type {
	return Type{
		Name:        "Credentials",
		Description: "Credentials",
		Constructor: func() Resource { return new(Credentials) },
	}
}

// Validate returns an error if the given Credentials object has incorrect values.
// Else, it returns nil.
func (c *Credentials) Validate(ctx context.Context) error {
	keyValidators := map[string]func(string) error{"password": ValidatePassword, "ssh-key": ValidateSSHKey}

	for keyType, key := range c.Keys {
		v := keyValidators[keyType]
		if err := v(key); err != nil {
			return err
		}
	}

	if c.Keys == nil {
		return ErrNoKeys
	}

	return c.NamedResource.Validate(ctx)
}

// Check to make sure password follows rules.
// 1. At least 12 characters long.
// 2. Contains upper and lowercase letters.
// 3. Contains at least 1 number.
// 4. Contains at least 1 special character.
func ValidatePassword(password string) error { //nolint:cyclop
	minLength := 12

	var okLen, okLower, okUpper, okNum, okSpec bool

	for _, char := range password {
		switch {
		case unicode.IsLower(char) && unicode.IsLetter(char):
			okLower = true
		case unicode.IsUpper(char) && unicode.IsLetter(char):
			okUpper = true
		case unicode.IsDigit(char):
			okNum = true
		case strings.Contains("!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~", string(char)): //nolint:gocritic
			okSpec = true
		}
	}

	okLen = len(password) >= minLength

	switch {
	case !okLen:
		return ErrPassLen
	case (!okLower || !okUpper):
		return ErrPassCase
	case !okNum:
		return ErrPassNum
	case !okSpec:
		return ErrPassSpecial
	}

	return nil
}

// TO BE IMPLEMENTED AT A LATER TIME.
// Check to make sure SSH key follows specified rules.
func ValidateSSHKey(key string) error {
	return nil
}

func NewCredential(name string, labels Labels) *Credentials {
	namedRes := new(NamedResource)

	namedRes.BaseResource = *NewBaseResource(CredentialsType(), labels)

	namedRes.Name = name

	ret := &Credentials{
		NamedResource: *namedRes,
		// some labels.
		Keys: labels,
	}

	return ret
}
