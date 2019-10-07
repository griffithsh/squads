package combat

import (
	"time"
)

type Waypoint struct {
	X, Y float64
}

// Mover is a component that can move.
type Mover struct {
	Moves    []Waypoint
	Duration time.Duration // Duration is how long this Mover will require to complete a move to the next Hex.
	Elapsed  time.Duration // Elapsed time since started the move to the next Hex.
	Speed    float64       // Speed is how fast we're moving

	// Delta to next Hex.
	dx, dy float64

	// Position of the Hex the mover has started the latest move from.
	x, y float64
}

// Type of this Component.
func (a *Mover) Type() string {
	return "Mover"
}
