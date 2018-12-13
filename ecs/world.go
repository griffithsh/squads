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

// Component of type for Entity.
func (mgr *World) Component(e Entity, ty string) Component {
	if v, ok := mgr.components[ty]; ok {
		if v, ok := v[e]; ok {
			return v
		}
	}
	return nil
}

// NewEntity creates an Entity
func (mgr *World) NewEntity() Entity {
	try := Entity(rand.Int63())
	// is the try already in use?
	if _, ok := mgr.entities[try]; ok {
		// recurse to try again.
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

// RemoveComponent from Entity.
func (mgr *World) RemoveComponent(e Entity, c Component) {
	delete(mgr.components[c.Type()], e)
	if len(mgr.components[c.Type()]) == 0 {
		delete(mgr.components, c.Type())
	}
}

// Clear all entities in the world.
func (mgr *World) Clear() {
	mgr.entities = map[Entity]struct{}{}
	mgr.components = map[string]map[Entity]Component{}
}
