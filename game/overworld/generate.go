package overworld

import (
	"math/rand"

	"github.com/griffithsh/squads/baddy"
	"github.com/griffithsh/squads/game/overworld/procedural"
	"github.com/griffithsh/squads/squad"

	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/geom"
)

func generateProcedural(rng *rand.Rand, recipe *procedural.Generator, lvl int) Map {
	generated := recipe.Generate(rng.Int63(), lvl)

	nodes := map[geom.Key]*Node{}
	for key, placement := range generated.Paths.Nodes {
		connections := map[geom.DirectionType]geom.Key{}
		for dir := range placement.Connections {
			connections[dir] = key.ToDirection(dir)
		}
		nodes[key] = &Node{
			ID:        key,
			e:         0,
			Connected: connections,
		}
	}

	enemies := map[geom.Key][]*game.Character{}
	for key, recipeID := range generated.Opponents {
		characters := []*game.Character{}
		for _, recipeID := range squad.Recipes[recipeID].Construct(rng) {
			char := baddy.Recipes[recipeID].Construct(rng)
			char.Level = lvl
			vit := int(char.VitalityPerLevel * float64(char.Level))
			hp := char.BaseHealth
			char.CurrentHealth = game.MaxHealth(hp, vit)
			characters = append(characters, char)
		}
		enemies[key] = characters
	}

	d := Map{
		Terrain: generated.Terrain,
		Nodes:   nodes,
		Enemies: enemies,
		Start:   generated.Paths.Start,
		Gate:    generated.Paths.Goal,
	}

	return d
}
