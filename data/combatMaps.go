package data

import (
	"math/rand"

	"github.com/griffithsh/squads/game"
)

var internalCombatMaps = []game.CombatMapRecipe{}

// GetCombatMap for use in a combat.
func (a *Archive) GetCombatMap() *game.CombatMapRecipe {

	switch len(a.combatMaps) {
	case 0:
		panic("no combat maps available")
	case 1:
		return &a.combatMaps[0]
	default:
		return &a.combatMaps[rand.Intn(len(a.combatMaps))]
	}
}
