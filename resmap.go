package zebra

import (
	"encoding/json"
)

type ResourceFactory interface {
	New(resourceType string) Resource
	Add(resourceType string, factory func() Resource) ResourceFactory
}

type typeMap map[string]func() Resource

func (t typeMap) New(resourceType string) Resource {
	factory, ok := t[resourceType]
	if !ok {
		return nil
	}

	return factory()
}

// Add adds a type and its factory method to the resource factory and returns the resource factory.
// The returned resource factory object can be used to add more types in a chained fashion.
func (t typeMap) Add(resourceType string, factory func() Resource) ResourceFactory {
	t[resourceType] = factory

	return t
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

func CopyResourceList(dest *ResourceList, src *ResourceList) {
	dest.factory = src.factory
	dest.Resources = make([]Resource, len(src.Resources))
	copy(dest.Resources, src.Resources)
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
			return ErrTypeEmpty
		}

		// Make a new resource for this value based on the embedded type field
		vType, ok := vAny.(string)
		if !ok {
			return ErrTypeEmpty
		}

		resource := r.factory.New(vType)

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
	dest.Resources = make(map[string]*ResourceList, len(src.Resources))

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
	listLen := len(r.Resources[key].Resources)

	for i, val := range r.Resources[key].Resources {
		if val.GetID() == res.GetID() {
			r.Resources[key].Resources[i] = r.Resources[key].Resources[listLen-1]
			r.Resources[key].Resources = r.Resources[key].Resources[:listLen-1]

			// If all values from key have been deleted, delete key entry
			if len(r.Resources[key].Resources) == 0 {
				delete(r.Resources, key)
			}

			return
		}
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
