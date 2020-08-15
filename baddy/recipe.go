package baddy

import (
	"math/rand"

	"github.com/griffithsh/squads/game"
)

// RecipeID identifies a baddy recipe.
type RecipeID int

const (
	Skellington RecipeID = iota
	Wolf
	Necro
)

// Recipe describes a way to construct a baddy.
type Recipe struct {
	ActionPoints int
	Preparation  int
	Name         string
	Sex          game.CharacterSex
	Profession   string
	Hair, Skin   string
}

// Construct a baddy from a Recipe.
func (recipe Recipe) Construct(rng *rand.Rand) *game.Character {
	return &game.Character{
		Name:                 recipe.Name,
		Hair:                 recipe.Hair,
		Skin:                 recipe.Skin,
		Sex:                  recipe.Sex,
		Profession:           recipe.Profession,
		InherantPreparation:  recipe.Preparation,
		InherantActionPoints: recipe.ActionPoints,
	}
}

// Recipes is the reference of baddy recipes in the game.
var Recipes map[RecipeID]Recipe = map[RecipeID]Recipe{
	Skellington: Recipe{
		ActionPoints: 60,
		Preparation:  50,
		Name:         "Dumble",
		Sex:          game.Male,
		Profession:   "Skeleton",
		Hair:         "black",
		Skin:         "pale",
	},
	Wolf: Recipe{
		ActionPoints: 60,
		Preparation:  50,
		Name:         "Hustle",
		Sex:          game.Male,
		Profession:   "Wolf",
		Hair:         "black",
		Skin:         "pale",
	},
	Necro: Recipe{
		ActionPoints: 60,
		Preparation:  50,
		Name:         "Pabst",
		Sex:          game.Male,
		Profession:   "Necromancer",
		Hair:         "black",
		Skin:         "pale",
	},
}
