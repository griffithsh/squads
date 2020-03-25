package squad

import (
	"math/rand"

	"github.com/griffithsh/squads/baddy"
)

// RecipeID identifies squad recipes.
type RecipeID int

const (
	SoloSkellington RecipeID = iota
	WolfPack1
)

type Candidate struct {
	Chance float64
	ID     baddy.RecipeID
}

// Recipe describes how to construct an enemy squad.
type Recipe []Candidate

// Construct a squad of baddies from a Recipe.
func (recipe Recipe) Construct(rng *rand.Rand) []baddy.RecipeID {
	result := make([]baddy.RecipeID, 0, len(recipe))
	for _, candidate := range recipe {
		roll := rng.Float64()

		if roll < candidate.Chance {
			result = append(result, candidate.ID)
		}
	}
	return result
}

var Recipes = map[RecipeID]Recipe{
	SoloSkellington: Recipe{
		{1.0, baddy.Skellington},
	},
	WolfPack1: Recipe{
		{1.0, baddy.Wolf},
		{0.5, baddy.Wolf},
		{0.5, baddy.Wolf},
	},
}
