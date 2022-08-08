package zebra

import (
	"encoding/json"
	"log"
)

type Type struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Constructor func() Resource `json:"-"`
}

func (t *Type) New() Resource {
	return t.Constructor()
}

type ResourceFactory interface {
	New(string) Resource
	Add(Type) ResourceFactory
	Types() []Type
	Type(string) (Type, bool)
}

type typeMap map[string]Type

func (t typeMap) New(resourceType string) Resource {
	aType, ok := t[resourceType]
	if !ok {
		return nil
	}

	return aType.New()
}

// Add adds a type and its factory method to the resource factory and returns the resource factory.
// The returned resource factory object can be used to add more types in a chained fashion.
func (t typeMap) Add(aType Type) ResourceFactory {
	t[aType.Name] = aType

	return t
}

func (t typeMap) Types() []Type {
	types := make([]Type, 0, len(t))
	for _, aType := range t {
		types = append(types, aType)
	}

	return types
}

func (t typeMap) Type(name string) (Type, bool) {
	aType, ok := t[name]

	return aType, ok
}

func Factory() ResourceFactory {
	return typeMap{}
}

type ResourceList struct {
	factory   ResourceFactory
	Resources []Resource
}

func NewResourceList(f ResourceFactory) *ResourceList {
	return &ResourceList{
		factory:   f,
		Resources: []Resource{},
	}
}

func (r *ResourceList) Delete(res Resource) {
	listLen := len(r.Resources)

	for i, val := range r.Resources {
		if val.GetID() == res.GetID() {
			r.Resources[i] = r.Resources[listLen-1]
			r.Resources = r.Resources[:listLen-1]
		}
	}
}

func CopyResourceList(dest *ResourceList, src *ResourceList) {
	if dest == nil || src == nil {
		return
	}

	dest.factory = src.factory
	dest.Resources = make([]Resource, len(src.Resources))
	copy(dest.Resources, src.Resources)
}

// Checks if a given interface is a `Type`` struct and returns that `Type` if true and an error if not.
func TypeChecker(data interface{}) (Type, error) {
	empty := Type{
		Name:        "",
		Description: "",
		Constructor: func() Resource { return nil },
	}

	vType, err := json.Marshal(data)

	if err != nil {
		return empty, err
	}

	if err = json.Unmarshal(vType, &empty); err != nil {
		return empty, err
	}

	return empty, nil
}

func (r *ResourceList) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Resources)
}

func (r *ResourceList) UnmarshalJSON(data []byte) error {
	// unmarshal the data as a list of maps with string as key so that
	// we can look up the type value of the resource.
	values := []map[string]interface{}{}
	if e := json.Unmarshal(data, &values); e != nil {
		return e
	}

	// For each value find out the type, the convert the value back into a
	// byte array. Use the type to create the actual resource object and the
	// unmarshal the byte array into that resource object.
	for _, value := range values {
		// all resources must have type
		vAny, ok := value["type"]
		if !ok {
			log.Default().Print("Could not find type")

			return ErrTypeEmpty
		}

		vType, err := TypeChecker(vAny)
		if err != nil {
			return ErrTypeEmpty
		}

		resource := r.factory.New(vType.Name)

		if resource == nil {
			// Type factory doesnt know this type return error
			return ErrTypeEmpty
		}

		// Capture the []byte of just this value
		resData, err := json.Marshal(value)
		if err != nil {
			return err
		}

		// We have the []byte representation and we have the target object
		// so now can do the final unmarshal and add this object into the
		// resource list
		if e := json.Unmarshal(resData, resource); e != nil {
			return e
		}

		r.Resources = append(r.Resources, resource)
	}

	return nil
}

type ResourceMap struct {
	factory   ResourceFactory
	Resources map[string]*ResourceList
}

func NewResourceMap(f ResourceFactory) *ResourceMap {
	return &ResourceMap{
		factory:   f,
		Resources: map[string]*ResourceList{},
	}
}

func CopyResourceMap(dest *ResourceMap, src *ResourceMap) {
	if dest == nil || src == nil {
		return
	}

	dest.factory = src.factory
	dest.Resources = make(map[string]*ResourceList)

	for key, val := range src.Resources {
		dest.Resources[key] = NewResourceList(dest.factory)
		CopyResourceList(dest.Resources[key], val)
	}
}

func (r *ResourceMap) GetFactory() ResourceFactory {
	return r.factory
}

func (r *ResourceMap) Add(res Resource, key string) {
	if r.Resources[key] == nil {
		r.Resources[key] = NewResourceList(r.factory)
	}

	r.Resources[key].Resources = append(r.Resources[key].Resources, res)
}

func (r *ResourceMap) Delete(res Resource, key string) {
	if r.Resources[key] == nil {
		return
	}

	r.Resources[key].Delete(res)

	if len(r.Resources[key].Resources) == 0 {
		delete(r.Resources, key)
	}
}

func (r *ResourceMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Resources)
}

func (r *ResourceMap) UnmarshalJSON(data []byte) error {
	// unmarshal the data as a map[string][]byte to extract each resourcelist
	// against a type.
	values := map[string]json.RawMessage{}
	if e := json.Unmarshal(data, &values); e != nil {
		return e
	}

	//
	for vType, rData := range values {
		rList := NewResourceList(r.factory)
		if e := json.Unmarshal(rData, rList); e != nil {
			return e
		}

		// Add this list to the Resource map
		r.Resources[vType] = rList
	}

	return nil
}
