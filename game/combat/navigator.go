package combat

import (
	"fmt"
	"time"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
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

var professionSpeeds = map[string]time.Duration{
	"Wolf":        150 * time.Millisecond,
	"Skeleton":    800 * time.Millisecond,
	"Necromancer": 500 * time.Millisecond,
}

// Update Movers.
func (nav *Navigator) Update(mgr *ecs.World, elapsed time.Duration) {
	entities := mgr.Get([]string{"Mover", "Participant", "Position"})

	for _, e := range entities {
		mover := mgr.Component(e, "Mover").(*Mover)
		pos := mgr.Component(e, "Position").(*game.Position)
		facer := mgr.Component(e, "Facer").(*game.Facer)
		participant := mgr.Component(e, "Participant").(*Participant)

		oldFace := facer.Face

		speed := float64(250 * time.Millisecond)
		if t, ok := professionSpeeds[participant.Profession]; ok {
			speed = float64(t)
		}

		if len(mover.Moves) > 0 && mover.Moves[0].X == pos.Center.X && mover.Moves[0].Y == pos.Center.Y {
			// Pop the first move, because it's the current position.
			mover.Moves = mover.Moves[1:]

			// First move is a little slow.
			mover.Speed = 0.75

			// Start-of-move tasks...
			mover.Elapsed = 0
			mover.Duration = time.Duration(speed / mover.Speed)
			mover.dx = mover.Moves[0].X - pos.Center.X
			mover.dy = mover.Moves[0].Y - pos.Center.Y
			mover.x = pos.Center.X
			mover.y = pos.Center.Y

			if dir, err := direction(mover.dx, mover.dy); err == nil {
				facer.Face = dir
			}
			if oldFace != facer.Face || 0 != mover.Speed {
				nav.Publish(&ParticipantMoving{
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
				nav.Publish(&ParticipantMoving{
					Entity:    e,
					NewSpeed:  0,
					OldSpeed:  mover.Speed,
					NewFacing: facer.Face,
					OldFacing: oldFace,
				})
				mgr.RemoveComponent(e, mover)
				nav.Publish(&ParticipantMovementConcluded{Entity: e})
				continue
			}

			oldSpeed := mover.Speed
			// The last few moves are slower than normal.
			switch len(mover.Moves) {
			default:
				// 1.00 = 100% of configured speed.
				mover.Speed = 1.00
			case 2:
				mover.Speed = 0.75
			case 1:
				mover.Speed = 0.50
			}

			// Start-of-move tasks...
			mover.Elapsed -= mover.Duration
			mover.Duration = time.Duration(speed / mover.Speed)
			mover.dx = mover.Moves[0].X - pos.Center.X
			mover.dy = mover.Moves[0].Y - pos.Center.Y
			mover.x = pos.Center.X
			mover.y = pos.Center.Y
			if dir, err := direction(mover.dx, mover.dy); err == nil {
				facer.Face = dir
			}
			if oldFace != facer.Face || oldSpeed != mover.Speed {
				nav.Publish(&ParticipantMoving{
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

// direction calculates which hexagonal direction the vector of x,y aligns with.
// FIXME: hard-coded assumption about hexes being 24 by 16
func direction(x, y float64) (geom.DirectionType, error) {
	switch {
	case x == 0 && y == 16:
		return geom.S, nil
	case x == 0 && y == -16:
		return geom.N, nil
	case x == 17 && y == 8:
		return geom.SE, nil
	case x == -17 && y == 8:
		return geom.SW, nil
	case x == 17 && y == -8:
		return geom.NE, nil
	case x == -17 && y == -8:
		return geom.NW, nil
	}
	return geom.DirectionType(0), fmt.Errorf("unhandled: %f,%f", x, y)
}
