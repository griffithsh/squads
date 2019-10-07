package game

import (
	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/geom"
)

type StatType int

//go:generate stringer -type=StatType
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

// EndTurnRequested occurs when the player indicates that they are finished
// commanding the current Actor.
type EndTurnRequested struct{}

// Type of the Event.
func (EndTurnRequested) Type() event.Type {
	return "game.EndTurnRequested"
}

// MoveModeRequested occurs when the player indicates that they wish to move the
// actor awaiting input.
type MoveModeRequested struct{}

// Type of the Event.
func (MoveModeRequested) Type() event.Type {
	return "game.MoveModeRequested"
}

// CancelSkillRequested occurs when the player indicates they want to cancel
// targeting of the skill they selected.
type CancelSkillRequested struct{}

// Type of the Event.
func (CancelSkillRequested) Type() event.Type {
	return "game.CancelSkillRequested"
}

// CombatActorMoving occurs when an actor has begun their movement.
type CombatActorMoving struct {
	Entity               ecs.Entity
	NewSpeed, OldSpeed   float64
	OldFacing, NewFacing geom.DirectionType
}

// Type of the Event.
func (CombatActorMoving) Type() event.Type {
	return "game.CombatActorMoving"
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

// WindowSizeChanged occurs when the size of the window the game is running in
// changes.
type WindowSizeChanged struct {
	OldW, OldH int
	NewW, NewH int
}

// Type of the Event.
func (WindowSizeChanged) Type() event.Type {
	return "game.WindowSizeChanged"
}
