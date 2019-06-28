package game

import (
	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
)

type StatType int

const (
	HPStat StatType = iota
	EnergyStat
	ActionStat
	PrepStat
)

// CombatBegan occurs when a Combat has been initialised.
type CombatBegan struct{}

// Type of the Event.
func (CombatBegan) Type() event.Type {
	return "game.CombatBegan"
}

// CombatStatModified occurs when an Actor's current stats changed.
type CombatStatModified struct {
	Entity ecs.Entity
	Stat   StatType
	Amount int
}

// Type of the Event.
func (CombatStatModified) Type() event.Type {
	return "game.CombatStatModified"
}

// EndTurnRequested occurs when the player indicates that they are finished
// commanding the current Actor.
type EndTurnRequested struct{}

// Type of the Event.
func (EndTurnRequested) Type() event.Type {
	return "game.EndTurnRequested"
}

// CombatActorMovementConcluded occurs when an actor has finished their movement.
type CombatActorMovementConcluded struct {
	Entity ecs.Entity
}

// Type of the Event.
func (CombatActorMovementConcluded) Type() event.Type {
	return "game.CombatActorMovementConcluded"
}

// CombatAwaitingPlayerInput occurs when an actor needs a command.
type CombatAwaitingPlayerInput struct{}

// Type of the Event.
func (CombatAwaitingPlayerInput) Type() event.Type {
	return "game.CombatAwaitingPlayerInput"
}
