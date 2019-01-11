package game

import (
	"time"
)

// Mover is a component that can move.
type Mover struct {
	Moves   []*Hex
	Free    time.Duration // Free is how long until this Mover is done with its movement.
	Elapsed time.Duration // Elapsed since last step Move.
	Speed   float64       // Speed is how fast we're moving
}

// Type of this Component.
func (a *Mover) Type() string {
	return "Mover"
}
