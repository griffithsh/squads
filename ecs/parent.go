package ecs

// Parent allows hierarchical relationships between Entites. If an Entity has a
// Parent Component, then it is destroyed when its Parent is.
type Parent struct {
	Value Entity
}

// Type of this Component.
func (*Parent) Type() string {
	return "Parent"
}

// ParentSystem manages Parent Components.
type ParentSystem struct{}

// Update the ParentSystem.
func (s *ParentSystem) Update(mgr *World) {
	// Remove orphans.
	for _, e := range mgr.Get([]string{"Parent"}) {
		p := mgr.Component(e, "Parent").(*Parent)

		if mgr.Exists(p.Value) {
			continue
		}

		mgr.DestroyEntity(e)
	}

	// Remove references to Children that have been destroyed from the Parents
	// that claim them.
	for _, e := range mgr.Get([]string{"Children"}) {
		c := mgr.Component(e, "Children").(*Children)
		val := []Entity{}
		for _, child := range c.Value {
			if !mgr.Exists(child) {
				continue
			}
			val = append(val, child)
		}
		c.Value = val
	}
}

// Children allows hierarchical relationships between entities. Some Entities have Children.
type Children struct {
	Value []Entity
}

// Type of this Component.
func (*Children) Type() string {
	return "Children"
}
