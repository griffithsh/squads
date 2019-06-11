package ecs

import (
	"math/rand"
)

// NewWorld creates an Entity Component System World.
func NewWorld() *World {
	return &World{
		entities:   map[Entity]struct{}{},
		components: map[string]map[Entity]Component{},
	}
}

// World is an instance of an Entity Component System.
type World struct {
	entities   map[Entity]struct{}
	components map[string]map[Entity]Component
}

// Get returns the list of entities that have all of the provided types.
func (mgr *World) Get(types []string) []Entity {
	result := []Entity{}
	// Check every Entity ...
	for e := range mgr.entities {

		// For every type to check ...
		count := 0
		for _, ty := range types {

			// If there are no components for that type, break from checking this entity further.
			if _, ok := mgr.components[ty]; !ok {
				break
			}

			// If this entity is present in the components then add it to the result.
			if _, ok := mgr.components[ty][e]; ok {
				count++
			}
		}
		if len(types) == count {
			result = append(result, e)
		}
	}
	return result
}

// Exists returns whether an entity exists in this World.
func (mgr *World) Exists(e Entity) bool {
	_, ok := mgr.entities[e]
	return ok
}

// Component retrieves the Component of Type t for Entity.
func (mgr *World) Component(e Entity, t string) Component {
	if v, ok := mgr.components[t]; ok {
		if v, ok := v[e]; ok {
			return v
		}
	}
	return nil
}

// NewEntity creates an Entity
func (mgr *World) NewEntity() Entity {
	try := Entity(rand.Int63())
	// Is the try already in use? This is hugely unlikely given there 2^63 possibilities.
	if _, ok := mgr.entities[try]; ok {
		// Recurse to try again.
		return mgr.NewEntity()
	}
	mgr.entities[try] = struct{}{}
	return try
}

// DestroyEntity removes an Entity and all its Components.
func (mgr *World) DestroyEntity(e Entity) {
	for _, entities := range mgr.components {
		delete(entities, e)
	}
	delete(mgr.entities, e)
}

// AddComponent to Entity.
func (mgr *World) AddComponent(e Entity, c Component) {
	if _, ok := mgr.components[c.Type()]; !ok {
		mgr.components[c.Type()] = map[Entity]Component{}
	}
	mgr.components[c.Type()][e] = c
}

// RemoveType removes the Component of Type t from Entity e.
func (mgr *World) RemoveType(e Entity, t string) {
	delete(mgr.components[t], e)
	if len(mgr.components[t]) == 0 {
		delete(mgr.components, t)
	}
}

// RemoveComponent from an Entity.
func (mgr *World) RemoveComponent(e Entity, c Component) {
	delete(mgr.components[c.Type()], e)
	if len(mgr.components[c.Type()]) == 0 {
		delete(mgr.components, c.Type())
	}
}

// Clear all Entities and their Components from the World, resetting it to an empty state.
func (mgr *World) Clear() {
	mgr.entities = map[Entity]struct{}{}
	mgr.components = map[string]map[Entity]Component{}
}

// Tag an Entity with an arbitrary string.
func (mgr *World) Tag(e Entity, tag string) {
	c := mgr.Component(e, "Tags")
	if c == nil {
		mgr.AddComponent(e, &Tags{tag})
		return
	}
	t := c.(*Tags)

	t2 := append(*t, tag)
	mgr.AddComponent(e, &t2)
}

// AnyTagged returns any Entity tagged with tag. It returns 0 when there are no
// Entities tagged with tag.
func (mgr *World) AnyTagged(tag string) Entity {
	for _, e := range mgr.Get([]string{"Tags"}) {
		t := mgr.Component(e, "Tags").(*Tags)

		for _, v := range *t {
			if tag == v {
				return e
			}
		}
	}
	return 0
}

// Tagged returns all Entities tagged with tag.
func (mgr *World) Tagged(tag string) []Entity {
	var result []Entity
	for _, e := range mgr.Get([]string{"Tags"}) {
		t := mgr.Component(e, "Tags").(*Tags)

		for _, v := range *t {
			if tag == v {
				result = append(result, e)
			}
		}
	}
	return result
}
