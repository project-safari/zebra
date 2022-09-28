package zebra

import (
	"errors"
)

var (
	ErrTypeDescriptionEmpty = errors.New("type description is empty")
	ErrTypeNameEmpty        = errors.New("type name is empty")
	ErrTypeEmpty            = errors.New("type is empty")
	ErrWrongType            = errors.New("type mismatch")
)

// A Type struct represents the type for the resources, with name and descriotion.
type Type struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type TypeConstructor func() Resource

func (t Type) Validate() error {
	if t.Name == "" {
		return ErrTypeNameEmpty
	}

	if t.Description == "" {
		return ErrTypeDescriptionEmpty
	}

	return nil
}

// A ResourceFactory interface represents the series of operations that can be performed on resources.
type ResourceFactory interface {
	New(string) Resource
	Add(Type, TypeConstructor) ResourceFactory
	Types() []Type
	Type(string) (Type, bool)
	Constructor(string) (TypeConstructor, bool)
}

type typeData struct {
	Type        Type
	Constructor TypeConstructor
}

type typeMap map[string]typeData

func (t typeMap) New(resourceType string) Resource {
	aType, ok := t[resourceType]
	if !ok {
		return nil
	}

	return aType.Constructor()
}

// Add adds a type and its factory method to the resource factory and returns the resource factory.
// The returned resource factory object can be used to add more types in a chained fashion.
func (t typeMap) Add(aType Type, constructor TypeConstructor) ResourceFactory {
	t[aType.Name] = typeData{aType, constructor}

	return t
}

func (t typeMap) Types() []Type {
	types := make([]Type, 0, len(t))
	for _, aType := range t {
		types = append(types, aType.Type)
	}

	return types
}

func (t typeMap) Type(name string) (Type, bool) {
	if aType, ok := t[name]; ok {
		return aType.Type, true
	}

	return Type{}, false
}

func (t typeMap) Constructor(name string) (TypeConstructor, bool) {
	if aType, ok := t[name]; ok {
		return aType.Constructor, true
	}

	return nil, false
}

func Factory() ResourceFactory {
	return typeMap{}
}
