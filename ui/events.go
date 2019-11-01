package ui

import (
	"github.com/griffithsh/squads/event"
)

// Interact happens when the player interacts with the game by clicking with the
// mouse or tapping on the screen.
type Interact struct {
	X, Y     float64
	Absolute bool
}

// Type of the Event.
func (Interact) Type() event.Type {
	return "ui.Interact"
}
