package ecs

import (
	"math/rand"
)

// NewWorld creates an Entity Component System World.
func NewWorld() *World {
	return &World{
		entities:     map[Entity]struct{}{},
		components:   map[string]map[Entity]Component{},
		dependencies: map[Entity][]Entity{},
	}
}

// World is an instance of an Entity Component System.
type World struct {
	entities   map[Entity]struct{}
	components map[string]map[Entity]Component

	dependencies map[Entity][]Entity
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

// Must returns the Entity when ok is true, otherwise it will panic.
func Must(e Entity, ok bool) Entity {
	if !ok {
		// FIXME: Better message.
		panic("ecs.Must not ok")
	}
	return e
}

// Single Entity that has all Components specified by types. Returns the Entity
// and a boolean indicating whether there was exactly one Entity that satisfies
// all types. When the second return value is false, the Entity returned is not
// valid.
func (mgr *World) Single(types []string) (Entity, bool) {
	es := mgr.Get(types)
	if len(es) != 1 {
		return 0, false
	}
	return es[0], true
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
	if dependants, ok := mgr.dependencies[e]; ok {
		for _, e := range dependants {
			mgr.DestroyEntity(e)
		}
		delete(mgr.dependencies, e)
	}

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

// ListComponents returns the Component types that are present on the Entity.
func (mgr *World) ListComponents(e Entity) []string {
	result := []string{}
	for ty, v := range mgr.components {
		if _, ok := v[e]; ok {
			result = append(result, ty)
		}
	}
	return result
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

// HasTag returns whether an Entity has the passed tag.
func (mgr *World) HasTag(e Entity, tag string) bool {
	comp := mgr.Component(e, "Tags")
	if comp == nil {
		return false
	}

	tags := comp.(*Tags)

	for _, t := range *tags {
		if tag == t {
			return true
		}
	}
	return false
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

// RemoveTag from the Entity.
func (mgr *World) RemoveTag(e Entity, tag string) {
	comp := mgr.Component(e, "Tags")
	if comp == nil {
		return
	}
	tags := comp.(*Tags)
	newTags := Tags{}
	for _, t := range *tags {
		if tag != t {
			newTags = append(newTags, t)
		}
	}
	if len(newTags) > 0 {
		mgr.AddComponent(e, &newTags)
	} else {
		mgr.RemoveComponent(e, tags)
	}
}

// Dependency adds a cascading destroy rule for a pair of Entities. When parent
// is destroyed, then child is also destroyed.
func (mgr *World) Dependency(parent, child Entity) {
	mgr.dependencies[parent] = append(mgr.dependencies[parent], child)
}

// Len returns the number of Entites currently in the world.
func (mgr *World) Len() int {
	return len(mgr.entities)
}
