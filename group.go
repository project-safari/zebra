package zebra

import "errors"

type Group struct {
	Name      string
	Resources *ResourceMap
	FreePool  *ResourceMap
	Count     uint
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

	if res.GetStatus().Lease == Free {
		g.FreePool.Add(res, res.GetType())
	}

	g.Count++
}

// Delete a resource from group.
func (g *Group) Delete(res Resource) {
	count := len(g.Resources.Resources[res.GetType()].Resources)

	g.Resources.Delete(res, res.GetType())
	g.FreePool.Add(res, res.GetType())

	newCount := len(g.Resources.Resources[res.GetType()].Resources)

	// Only update count if resource was really deleted
	if count != newCount {
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

	if count != int(g.Count) {
		return ErrWrongCount
	}

	return nil
}
