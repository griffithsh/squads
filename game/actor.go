package game

import (
	"time"
)

// Actor is a component that can be commanded to do things. Or maybe it's just an animator?
type Actor struct {
	M, N  int
	Moves []*Hex
	Busy  bool
	Free  time.Duration // Free is how long until this Actor is no longer Busy
}

// Type of this Component.
func (a *Actor) Type() string {
	return "Actor"
}

// Move adds movement Acts to this Actor.
func (a *Actor) Move(steps []*Hex) {
	if len(a.Moves) > 0 || a.Busy {
		return
	}

	a.Moves = steps
}
