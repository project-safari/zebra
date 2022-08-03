package zebra

import (
	"errors"

	"github.com/project-safari/zebra/status"
)

type Group struct {
	Name      string
	Resources *ResourceMap
}

var (
	ErrResourceMap = errors.New("resource map is nil")
	ErrWrongGroup  = errors.New("resource does not have group name matching this group")
)

func NewGroup(name string) *Group {
	return &Group{
		Name:      name,
		Resources: NewResourceMap(nil),
	}
}

// Add a resource to group.
func (g *Group) Add(res Resource) {
	g.Resources.Add(res, res.GetType())
}

// Delete a resource from group.
func (g *Group) Delete(res Resource) {
	// Nothing to delete
	if _, ok := g.Resources.Resources[res.GetType()]; !ok {
		return
	}

	// Delete
	g.Resources.Delete(res, res.GetType())
}

// Given a resource, update lease status.
// If resource is not part of this group or is already free, throw an error.
func (g *Group) Free(res Resource) error {
	// If res is not in this group, return error
	if name, ok := res.GetLabels()["system.group"]; !ok || name != g.Name {
		return ErrWrongGroup
	}

	if err := res.GetStatus().SetFree(); err != nil {
		return err
	}

	return nil
}

// Given a resource, remove from free pool and update lease status.
// If resource is not part of this group or is already free, throw an error.
func (g *Group) Lease(res Resource) error {
	// If res is not in this group, return error
	if name, ok := res.GetLabels()["system.group"]; !ok || name != g.Name {
		return ErrWrongGroup
	}

	if err := res.GetStatus().SetLeased(); err != nil {
		return err
	}

	return nil
}

// Construct resource map with up-to-date free pool information.
func (g *Group) FreePool() *ResourceMap {
	freePool := NewResourceMap(nil)

	for t, l := range g.Resources.Resources {
		for _, r := range l.Resources {
			if r.GetStatus().Lease() == status.Free {
				freePool.Add(r, t)
			}
		}
	}

	return freePool
}

func (g *Group) Validate() error {
	// Check that name is valid
	if g.Name == "" {
		return ErrNameEmpty
	}

	if g.Resources == nil {
		return ErrResourceMap
	}

	return nil
}
