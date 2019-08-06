package combat

import "github.com/griffithsh/squads/event"

// StateTransition occurs when the combat's state changes
type StateTransition struct {
	Old, New State
}

// Type of the Event.
func (StateTransition) Type() event.Type {
	return "combat.StateTransition"
}

// DifferentHexSelected occurs when the user has selected a different hex - i.e.
// via mousing over.
type DifferentHexSelected struct {
	M,N int
}

// Type of the Event.
func (DifferentHexSelected) Type() event.Type {
	return "combat.DifferentHexSelected"
}
