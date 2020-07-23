package data

import (
	"math/rand"

	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/geom"
)

var internalCombatMaps = []game.CombatMapRecipe{}

// mbyn creates a slice of Keys that fill a rectangular field of m by n.
func mbyn(m, n int) []geom.Key {
	var result []geom.Key
	for ni := 0; ni < n; ni++ {
		for mi := 0; mi < m; mi++ {
			result = append(result, geom.Key{M: mi, N: ni})
		}
	}
	return result
}

func random() *game.CombatMapRecipe {
	recipe := game.CombatMapRecipe{
		Starts: []geom.Key{
			{M: 6, N: 18},
			{M: 2, N: 8},
		},
		TileW: 24,
		TileH: 16,
	}
	maxM := rand.Intn(6) + 12
	for _, key := range mbyn(maxM, rand.Intn(4)+10) {
		recipe.Hexes = append(recipe.Hexes, game.CombatMapRecipeHex{
			Position: key,
			Visuals: []game.CombatMapRecipeVisual{
				{
					Frames: []game.CombatMapRecipeHexFrame{
						{
							Texture: "terrain.png",
							X:       0,
							Y:       0,
						},
					},
					Layer: 1,
				},
			},
		})
		m, n := key.M, key.N
		i := m + n*maxM
		if (m != 4 || n != 14) && i%17 == 1 || i%23 == 1 {
			// add a tree!
			last := len(recipe.Hexes) - 1
			recipe.Hexes[last].Obstacle = true
			recipe.Hexes[last].Visuals = append(recipe.Hexes[last].Visuals, game.CombatMapRecipeVisual{
				Frames: []game.CombatMapRecipeHexFrame{
					{
						Texture: "trees.png",
						X:       0,
						Y:       0,
					},
				},
				Layer: 10,
			})
		}
	}
	return &recipe
}

// GetCombatMap for use in a combat.
func (a *Archive) GetCombatMap() *game.CombatMapRecipe {

	switch len(a.combatMaps) {
	case 0:
		return random()
	case 1:
		return &a.combatMaps[0]
	default:
		return &a.combatMaps[rand.Intn(len(a.combatMaps))]
	}
}
