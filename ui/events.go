package ui

import (
	"github.com/griffithsh/squads/event"
)

// UIInteract happens when the player interacts with the game by clicking with
// the mouse or tapping on the screen.
type UIInteract struct {
	X, Y                 float64
	AbsoluteX, AbsoluteY float64
}

// Type of the Event.
func (UIInteract) Type() event.Type {
	return "ui.UIInteract"
}

// Interact happens when a UIInteract event is unhandled by any UI.
type Interact struct {
	X, Y                 float64
	AbsoluteX, AbsoluteY float64
}

// Type of the Event.
func (Interact) Type() event.Type {
	return "ui.Interact"
}
