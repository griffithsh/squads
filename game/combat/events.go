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
