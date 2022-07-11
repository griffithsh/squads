package game

import (
	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
)

//go:generate go run github.com/dmarkham/enumer -output=./events_enumer.go -type=StatType,CombatResult

type StatType int

const (
	HPStat StatType = iota
	EnergyStat
	ActionStat
	PrepStat
)

type CombatResult int

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

// SomethingInteresting occurs when something interesting has occurred. The
// camera may want to focus on X,Y.
type SomethingInteresting struct {
	X, Y float64
}

// Type of the Event.
func (SomethingInteresting) Type() event.Type {
	return "game.SomethingInteresting"
}
