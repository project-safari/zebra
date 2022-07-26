package zebra

import (
	"errors"

	"github.com/project-safari/zebra/status"
)

type Group struct {
	Name      string
	Resources *ResourceMap
	FreePool  *ResourceMap
	Count     int
}

var (
	ErrResourceMap = errors.New("resource map is nil")
	ErrWrongCount  = errors.New("group count differs from associated resource map count")
)

func NewGroup(name string) *Group {
	return &Group{
		Name:      name,
		Resources: NewResourceMap(nil), // have to do something about this factory
		FreePool:  NewResourceMap(nil), // considering we don't turn these into json, does this matter?
		Count:     0,
	}
}

// Add a resource to group.
func (g *Group) Add(res Resource) {
	g.Resources.Add(res, res.GetType())

	if res.GetStatus().Lease() == status.Free {
		g.FreePool.Add(res, res.GetType())
	}

	g.Count++
}

// Delete a resource from group.
func (g *Group) Delete(res Resource) {
	l, ok := g.Resources.Resources[res.GetType()]

	// Nothing to delete
	if !ok {
		return
	}

	// Keep track of how many resources there are before delete
	count := len(l.Resources)

	// Delete
	g.Resources.Delete(res, res.GetType())
	g.FreePool.Delete(res, res.GetType())

	countAfter := 0

	// Check if list still exists before updating countAfter
	if l, ok = g.Resources.Resources[res.GetType()]; ok {
		countAfter = len(l.Resources)
	}

	// Only update count if resource was really deleted
	if count != countAfter {
		g.Count--
	}
}

func (g *Group) Validate() error {
	// Check that name is valid
	if g.Name == "" {
		return ErrNameEmpty
	}

	if g.Resources == nil {
		return ErrResourceMap
	}

	if g.FreePool == nil {
		return ErrResourceMap
	}

	count := 0
	for _, l := range g.Resources.Resources {
		count += len(l.Resources)
	}

	if count != g.Count {
		return ErrWrongCount
	}

	return nil
}
