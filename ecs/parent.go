package ecs

// ParentSystem manages Parent Components.
type ParentSystem struct {
	mgr *World
}

// NewParentSystem constructs a new ParentSystem.
func NewParentSystem(mgr *World) *ParentSystem {
	return &ParentSystem{
		mgr: mgr,
	}
}

// Update the ParentSystem.
func (s *ParentSystem) Update() {
	// Remove references to Children that have been destroyed from the Parents
	// that claim them.
	for _, e := range s.mgr.Get([]string{"Children"}) {
		c := s.mgr.Component(e, "Children").(*Children)
		val := []Entity{}
		for _, child := range c.Value {
			if !s.mgr.Exists(child) {
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
