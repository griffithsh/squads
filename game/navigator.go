package game

import (
	"time"

	"github.com/griffithsh/squads/ecs"
)

// Navigator is the System for Movers.
type Navigator struct {
}

// Update Actors.
func (nav *Navigator) Update(mgr *ecs.World, elapsed time.Duration) {
	entities := mgr.Get([]string{"Mover", "Actor", "Position", "Obstacle"})

	spd := float64(time.Millisecond * 150)
	for _, e := range entities {
		mover := mgr.Component(e, "Mover").(*Mover)
		actor := mgr.Component(e, "Actor").(*Actor)
		pos := mgr.Component(e, "Position").(*Position)
		obstacle := mgr.Component(e, "Obstacle").(*Obstacle)

		switch {
		// case "first move":
		case len(mover.Moves) > 0 && mover.Moves[0].M == actor.M && mover.Moves[0].N == actor.N:
			mover.Moves = mover.Moves[1:]
			mover.Speed = 0.5
			mover.Free = time.Duration(spd / mover.Speed)
			mover.Elapsed = 0

		// case "traversing":
		default:
			mover.Elapsed += elapsed

		// case "next move":
		case mover.Elapsed >= mover.Free && len(mover.Moves) > 0:
			mover.Elapsed -= mover.Free
			mover.Free = 0
			switch len(mover.Moves) {
			default:
				mover.Speed = 0.8
			case 2:
				mover.Speed = 0.6
			case 1:
				mover.Speed = 0.4
			}
			mover.Free = time.Duration(spd / mover.Speed)
			mover.Elapsed = 0

			step := mover.Moves[0]

			actor.M = step.M
			actor.N = step.N

			// Update the position of this Mover.
			pos.Center.X = mover.Moves[0].X()
			pos.Center.Y = mover.Moves[0].Y()

			// Take the Mover's Obstacle with them as they move.
			obstacle.M = step.M
			obstacle.N = step.N

			mover.Moves = mover.Moves[1:]

		// case "all done":
		case mover.Elapsed >= mover.Free && len(mover.Moves) == 0:
			mgr.RemoveComponent(e, mover)
		}
	}
}
