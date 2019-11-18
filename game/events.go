package game

import (
	"github.com/griffithsh/squads/ecs"
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

type CombatResult int

//go:generate stringer -type=CombatResult
const (
	Victorious CombatResult = iota
	Defeated
	Escaped
)

// CombatBegan occurs when a Combat has been initialised.
type CombatBegan struct{}

// Type of the Event.
func (CombatBegan) Type() event.Type {
	return "game.CombatBegan"
}

// CombatConcluded occurs when a Combat is over.
type CombatConcluded struct {
	Results map[ecs.Entity]CombatResult
}

// Type of the Event.
func (CombatConcluded) Type() event.Type {
	return "game.CombatConcluded"
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
