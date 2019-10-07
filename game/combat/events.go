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

// EndTurnRequested occurs when the player indicates that they are finished
// commanding the current Actor.
type EndTurnRequested struct{}

// Type of the Event.
func (EndTurnRequested) Type() event.Type {
	return "combat.EndTurnRequested"
}

// MoveModeRequested occurs when the player indicates that they wish to move the
// actor awaiting input.
type MoveModeRequested struct{}

// Type of the Event.
func (MoveModeRequested) Type() event.Type {
	return "combat.MoveModeRequested"
}

// CancelSkillRequested occurs when the player indicates they want to cancel
// targeting of the skill they selected.
type CancelSkillRequested struct{}

// Type of the Event.
func (CancelSkillRequested) Type() event.Type {
	return "combat.CancelSkillRequested"
}

// ActorMoving occurs when an actor has begun their movement.
type ActorMoving struct {
	Entity               ecs.Entity
	NewSpeed, OldSpeed   float64
	OldFacing, NewFacing geom.DirectionType
}

// Type of the Event.
func (ActorMoving) Type() event.Type {
	return "combat.ActorMoving"
}
