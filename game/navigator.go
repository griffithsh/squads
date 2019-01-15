package game

import (
	"time"

	"github.com/griffithsh/squads/ecs"
)

// Navigator is the System for Movers.
type Navigator struct {
}

// Update Movers.
func (nav *Navigator) Update(mgr *ecs.World, elapsed time.Duration) {

	entities := mgr.Get([]string{"Mover", "Actor", "Position"})

	speed := float64(time.Millisecond * 150)
	for _, e := range entities {
		mover := mgr.Component(e, "Mover").(*Mover)
		actor := mgr.Component(e, "Actor").(*Actor)
		pos := mgr.Component(e, "Position").(*Position)

		if len(mover.Moves) > 0 && mover.Moves[0].M == actor.M && mover.Moves[0].N == actor.N {
			// Pop the first move, because it's the current position.
			mover.Moves = mover.Moves[1:]

			// First move is a little slow.
			mover.Speed = 0.5

			// Start-of-move tasks...
			mover.Elapsed = 0
			mover.Duration = time.Duration(speed / mover.Speed)
			mover.dx = mover.Moves[0].X() - pos.Center.X
			mover.dy = mover.Moves[0].Y() - pos.Center.Y
			mover.x = pos.Center.X
			mover.y = pos.Center.Y

		} else if mover.Elapsed >= mover.Duration {
			// End-of-move tasks...
			dest := mover.Moves[0]
			pos.Center.X = dest.X()
			pos.Center.Y = dest.Y()
			actor.M = dest.M
			actor.N = dest.N
			if obstacle, ok := mgr.Component(e, "Obstacle").(*Obstacle); ok {
				obstacle.M = dest.M
				obstacle.N = dest.N
			}

			// Pop the move list to update the next destination.
			mover.Moves = mover.Moves[1:]

			// Are we done?
			if len(mover.Moves) == 0 {
				mgr.RemoveComponent(e, mover)
				continue
			}

			// The last few moves are slower than normal.
			switch len(mover.Moves) {
			default:
				mover.Speed = 0.75
			case 2:
				mover.Speed = 0.55
			case 1:
				mover.Speed = 0.30
			}

			// Start-of-move tasks...
			mover.Elapsed -= mover.Duration
			mover.Duration = time.Duration(speed / mover.Speed)
			mover.dx = mover.Moves[0].X() - pos.Center.X
			mover.dy = mover.Moves[0].Y() - pos.Center.Y
			mover.x = pos.Center.X
			mover.y = pos.Center.Y

		} else {
			// Traversing ....
			mover.Elapsed += elapsed

			pos.Center.X = mover.Elapsed.Seconds()/mover.Duration.Seconds()*mover.dx + mover.x
			pos.Center.Y = mover.Elapsed.Seconds()/mover.Duration.Seconds()*mover.dy + mover.y
		}
	}
}
