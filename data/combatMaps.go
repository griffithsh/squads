package data

import (
	"math/rand"

	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/geom"
)

func (a *Archive) GetCombatMap() *game.CombatMapRecipe {
	// FIXME: query from Archive.combatMaps instead of generating a random one.
	recipe := game.CombatMapRecipe{
		Starts: []geom.Key{
			{M: 6, N: 18},
			{M: 2, N: 8},
		},
		TileW: 24,
		TileH: 16,
	}
	maxM := rand.Intn(3) + 6
	for _, key := range geom.MByN(maxM, rand.Intn(7)+20) {
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
