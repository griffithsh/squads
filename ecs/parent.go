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

func (s *ParentSystem) Update(mgr *World) {
	for _, e := range mgr.Get([]string{"Parent"}) {
		p := mgr.Component(e, "Parent").(*Parent)

		if mgr.Exists(p.Value) {
			continue
		}

		mgr.DestroyEntity(e)
	}
}
