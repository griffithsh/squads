package overworld

import (
	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/geom"
)

//go:generate stringer -type=TokenType

// TokenType represents a category of owner for the Token. The values are
// unusual increments so that their sum can show which types collided.
type TokenType int

const (
	// SquadToken represents the Token for a Squad.
	SquadToken TokenType = 1
	// GateToken represents the Token for an exit gate.
	GateToken = 2
	// MerchantToken represents a stop where the player can buy stuff.
	MerchantToken = 4
)

// Token represents a presence on the overworld map. This might be an enemy
// Squad, the exit gate, a merchant, or rest stop.
type Token struct {
	Key      geom.Key
	Presence ecs.Entity
	// type? Squad, Gate, Merchant, etc?
	Category TokenType
}

// Type of this Component.
func (Token) Type() string {
	return "Token"
}
