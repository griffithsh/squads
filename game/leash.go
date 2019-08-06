package game

import (
	"time"

	"github.com/griffithsh/squads/ecs"
)

// Leash is a Component that makes the Entity Position follow another Entity's
// Position.
type Leash struct {
	Owner ecs.Entity

	LayerOffset int
}

// Type of this Component.
func (f *Leash) Type() string {
	return "Leash"
}

// LeashSystem synchronises the Position of Entities with a Leash Component to
// their Owner's Positions.
type LeashSystem struct{}

// Update the system.
func (sys *LeashSystem) Update(mgr *ecs.World, elapsed time.Duration) {
	for _, e := range mgr.Get([]string{"Leash"}) {
		leash := mgr.Component(e, "Leash").(*Leash)
		pos, ok := mgr.Component(leash.Owner, "Position").(*Position)
		if !ok {
			mgr.RemoveComponent(e, &Position{})
			continue
		}

		newPos := *pos
		newPos.Layer += leash.LayerOffset
		mgr.AddComponent(e, &newPos)
	}
}
