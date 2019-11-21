package overworld

import (
	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/geom"
)

type TokenMoved struct {
	E    ecs.Entity
	From geom.Key
	To   geom.Key
}

// Type of the Event.
func (TokenMoved) Type() event.Type {
	return "overworld.TokenMoved"
}

type TokensCollided struct {
	E1, E2 ecs.Entity
	At     geom.Key
}

// Type of the Event.
func (TokensCollided) Type() event.Type {
	return "overworld.TokensCollided"
}

// Complete happens when the player escapes the current overworld through the portal.
type Complete event.Type

// CombatInitiated occurs when the player has met another squad for combat.
type CombatInitiated struct {
	Squads []ecs.Entity
	// info about the terrain?
}

// Type of the Component.
func (CombatInitiated) Type() event.Type {
	return "overworld.CombatInitiated"
}

// TODO GoingShopping event
