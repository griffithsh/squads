package game

import (
	"time"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/geom"
)

// Navigator is the System for Movers.
type Navigator struct {
	*event.Bus
}

// NewNavigator constructs a Navigator.
func NewNavigator(bus *event.Bus) *Navigator {
	return &Navigator{
		Bus: bus,
	}
}

// Update Movers.
func (nav *Navigator) Update(mgr *ecs.World, elapsed time.Duration) {
	entities := mgr.Get([]string{"Mover", "Actor", "Position"})

	speed := float64(time.Millisecond * 150)
	for _, e := range entities {
		mover := mgr.Component(e, "Mover").(*Mover)
		pos := mgr.Component(e, "Position").(*Position)
		facer := mgr.Component(e, "Facer").(*Facer)
		oldFace := facer.Face

		if len(mover.Moves) > 0 && mover.Moves[0].X == pos.Center.X && mover.Moves[0].Y == pos.Center.X {
			// Pop the first move, because it's the current position.
			mover.Moves = mover.Moves[1:]

			// First move is a little slow.
			mover.Speed = 0.5

			// Start-of-move tasks...
			mover.Elapsed = 0
			mover.Duration = time.Duration(speed / mover.Speed)
			mover.dx = mover.Moves[0].X - pos.Center.X
			mover.dy = mover.Moves[0].Y - pos.Center.Y
			mover.x = pos.Center.X
			mover.y = pos.Center.Y
			if dir, err := geom.Direction(mover.dx, mover.dy); err == nil {
				facer.Face = dir
			}

			if oldFace == facer.Face || 0 != mover.Speed {
				nav.Publish(&CombatActorMoving{
					Entity:    e,
					NewSpeed:  mover.Speed,
					OldSpeed:  0,
					NewFacing: facer.Face,
					OldFacing: oldFace,
				})
			}

		} else if mover.Elapsed >= mover.Duration {
			// End-of-move tasks...
			dest := mover.Moves[0]
			pos.Center.X = float64(dest.X)
			pos.Center.Y = float64(dest.Y)

			// Pop the move list to update the next destination.
			mover.Moves = mover.Moves[1:]

			// Are we done?
			if len(mover.Moves) == 0 {
				nav.Publish(&CombatActorMoving{
					Entity:    e,
					NewSpeed:  0,
					OldSpeed:  mover.Speed,
					NewFacing: facer.Face,
					OldFacing: oldFace,
				})
				mgr.RemoveComponent(e, mover)
				nav.Publish(&CombatActorMovementConcluded{Entity: e})
				continue
			}

			oldSpeed := mover.Speed
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
			mover.dx = float64(mover.Moves[0].X) - pos.Center.X
			mover.dy = float64(mover.Moves[0].Y) - pos.Center.Y
			mover.x = pos.Center.X
			mover.y = pos.Center.Y
			if dir, err := geom.Direction(mover.dx, mover.dy); err == nil {
				facer.Face = dir
			}

			if oldFace == facer.Face || oldSpeed != mover.Speed {
				nav.Publish(&CombatActorMoving{
					Entity:    e,
					NewSpeed:  mover.Speed,
					OldSpeed:  oldSpeed,
					NewFacing: facer.Face,
					OldFacing: oldFace,
				})
			}

		} else {
			// Traversing ....
			mover.Elapsed += elapsed

			pos.Center.X = mover.Elapsed.Seconds()/mover.Duration.Seconds()*mover.dx + mover.x
			pos.Center.Y = mover.Elapsed.Seconds()/mover.Duration.Seconds()*mover.dy + mover.y
		}
	}
}
