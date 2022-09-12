package zebra

import (
	"encoding/json"
)

type ResourceList struct {
	ctr       TypeConstructor
	Resources []Resource
}

func NewResourceList(ctr TypeConstructor) *ResourceList {
	return &ResourceList{
		ctr:       ctr,
		Resources: []Resource{},
	}
}

func (r *ResourceList) Add(res Resource) error {
	r.Resources = append(r.Resources, res)

	return nil
}

func (r *ResourceList) Delete(res Resource) error {
	listLen := len(r.Resources)

	for i, val := range r.Resources {
		if val.GetMeta().ID == res.GetMeta().ID {
			r.Resources[i] = r.Resources[listLen-1]
			r.Resources = r.Resources[:listLen-1]

			return nil
		}
	}

	return ErrNotFound
}

func CopyResourceList(dest *ResourceList, src *ResourceList) {
	if dest == nil || src == nil {
		return
	}

	dest.ctr = src.ctr
	dest.Resources = make([]Resource, len(src.Resources))
	copy(dest.Resources, src.Resources)
}

func (r *ResourceList) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Resources)
}

func (r *ResourceList) UnmarshalJSON(data []byte) error {
	// unmarshal the data as a list of maps with string as key so that
	// we can then create a resource object and parse each resource.
	values := []json.RawMessage{}
	if e := json.Unmarshal(data, &values); e != nil {
		return e
	}

	// Convert each value back into a byte array. Use the type to create the
	// actual resource object and then unmarshal the byte array into that
	// resource object.
	for _, value := range values {
		resource := r.ctr()

		// We have the []byte representation and we have the target object
		// so now can do the final unmarshal and add this object into the
		// resource list
		if e := json.Unmarshal(value, resource); e != nil {
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
		ctr, _ := src.factory.Constructor(key)
		dest.Resources[key] = NewResourceList(ctr)
		CopyResourceList(dest.Resources[key], val)
	}
}

func (r *ResourceMap) Factory() ResourceFactory {
	return r.factory
}

func (r *ResourceMap) Add(res Resource) error {
	key := res.GetMeta().Type.Name
	rl := r.Resources[key]

	if rl == nil {
		if ctr, ok := r.factory.Constructor(key); ok {
			rl = NewResourceList(ctr)
			r.Resources[key] = rl
		} else {
			return ErrTypeEmpty
		}
	}

	return rl.Add(res)
}

func (r *ResourceMap) Delete(res Resource) error {
	key := res.GetMeta().Type.Name
	if r.Resources[key] == nil {
		return ErrNotFound
	}

	if err := r.Resources[key].Delete(res); err != nil {
		return err
	}

	if len(r.Resources[key].Resources) == 0 {
		delete(r.Resources, key)
	}

	return nil
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

	// for each type create the resource list and parse the resource list
	for vType, rData := range values {
		ctr, ok := r.factory.Constructor(vType)
		if !ok {
			return ErrTypeEmpty
		}

		rList := NewResourceList(ctr)
		if e := json.Unmarshal(rData, rList); e != nil {
			return e
		}

		// Add this list to the Resource map
		r.Resources[vType] = rList
	}

	return nil
}
