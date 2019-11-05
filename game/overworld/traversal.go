package overworld

import (
	"time"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/game"
)

// Traversal is a Component that changes the X and Y coordinates of an Entity over time.
type Traversal struct {
	age      time.Duration
	Duration time.Duration

	origin      *game.Center
	Destination game.Center
	diff        game.Center

	Complete func()
}

// Type of this Component.
func (Traversal) Type() string {
	return "Traversal"
}

// TraversalSystem manager Traversal Components.
type TraversalSystem struct{}

// Update all Traversal Components.
func (ts *TraversalSystem) Update(mgr *ecs.World, elapsed time.Duration) {
	for _, e := range mgr.Get([]string{"Traversal", "Position"}) {
		traversal := mgr.Component(e, "Traversal").(*Traversal)
		position := mgr.Component(e, "Position").(*game.Position)

		// First update of this Traversal.
		if traversal.origin == nil {
			traversal.origin = &game.Center{
				X: position.Center.X,
				Y: position.Center.Y,
			}
			traversal.diff.X = traversal.Destination.X - position.Center.X
			traversal.diff.Y = traversal.Destination.Y - position.Center.Y
		}

		// Traversal complete.
		if traversal.age+elapsed > traversal.Duration {
			position.Center = traversal.Destination
			if traversal.Complete != nil {
				traversal.Complete()
			}
			mgr.RemoveComponent(e, traversal)
			continue
		}

		// Otherwise ...
		traversal.age += elapsed

		// Update the position of the Entity.
		perc := float64(traversal.age) / float64(traversal.Duration)
		// TODO: This function generates linear velocity. I think I would prefer
		// a velocity that starts and ends slow, but speeds up in the middle.
		position.Center.X = traversal.origin.X + (traversal.diff.X * perc)
		position.Center.Y = traversal.origin.Y + (traversal.diff.Y * perc)
	}
}
