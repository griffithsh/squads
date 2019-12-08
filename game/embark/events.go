package embark

import "github.com/griffithsh/squads/event"

// SquadSelected occurs when the player has finished configuring their squad and
// is ready to embark on a new run.
type SquadSelected struct {
	// TODO: slice of characters?
}

// Type of the Event.
func (SquadSelected) Type() event.Type {
	return "embark.SquadSelected"
}

// Embarked is the event that occurs when the embark Manager has finished fading out.
type Embarked struct {
}

// Type of the Event.
func (Embarked) Type() event.Type {
	return "embark.Embarked"
}
