package combat

import (
	"fmt"
	"reflect"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/geom"
)

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
	K *geom.Key
}

// Type of the Event.
func (DifferentHexSelected) Type() event.Type {
	return "combat.DifferentHexSelected"
}

// ActorTurnChanged occurs when the Actor whose turn it is changes.
type ActorTurnChanged struct {
	Entity ecs.Entity
}

// Type of the Event.
func (v ActorTurnChanged) Type() event.Type {
	t := reflect.TypeOf(v)
	return event.Type(fmt.Sprintf("%s.%s", t.PkgPath(), t.Name()))
}

// StatModified occurs when an Actor's current stats changed.
type StatModified struct {
	Entity ecs.Entity
	Stat   game.StatType
	Amount int
}

// Type of the Event.
func (StatModified) Type() event.Type {
	return "combat.StatModified"
}
