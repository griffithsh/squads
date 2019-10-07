package game

import (
	"github.com/griffithsh/squads/event"
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
