package game

import (
	"time"

	"github.com/griffithsh/squads/ecs"
)

// Choreographer is the System for Actors.
type Choreographer struct {
}

// Update Actors.
func (c *Choreographer) Update(mgr *ecs.World, elapsed time.Duration) {
	entities := mgr.Get([]string{"Actor"})

	for _, e := range entities {
		actor := mgr.Component(e, "Actor").(*Actor)
		// If actor has nothing to do ...
		if len(actor.Moves) == 0 {
			actor.Busy = false
			actor.Free = 0
			continue
		}

		// ... or is currently performing an action ...
		if actor.Busy {
			actor.Free -= elapsed
			if actor.Free < elapsed {
				actor.Busy = false
				actor.Free = 0
			}
			continue
		}

		// ... otherwise deal with the front of the command queue.
		step := actor.Moves[0]
		actor.M = step.M
		actor.N = step.N
		actor.Busy = true
		actor.Free = time.Millisecond * 100

		pos := mgr.Component(e, "Position").(*Position)

		pos.Center.X = actor.Moves[0].X()
		pos.Center.Y = actor.Moves[0].Y()

		// Take the Actor's Obstacle with them as they move.
		obstacle := mgr.Component(e, "Obstacle").(*Obstacle)
		obstacle.M = actor.Moves[0].M
		obstacle.N = actor.Moves[0].N

		actor.Moves = actor.Moves[1:]
	}
}
