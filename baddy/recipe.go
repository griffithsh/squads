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
)

// Recipe describes a way to construct a baddy.
type Recipe struct {
	ActionPoints       int
	Preparation        int
	Name               string
	Sex                game.CharacterSex
	Profession         string
	SmallIcon, BigIcon game.Sprite
}

// Construct a baddy from a Recipe.
func (recipe Recipe) Construct(rng *rand.Rand) *game.Character {
	return &game.Character{
		Name:                 recipe.Name,
		Sex:                  recipe.Sex,
		Profession:           recipe.Profession,
		InherantPreparation:  recipe.Preparation,
		InherantActionPoints: recipe.ActionPoints,
		SmallIcon:            recipe.SmallIcon,
		BigIcon:              recipe.BigIcon,
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
		SmallIcon: game.Sprite{
			Texture: "hud.png",
			X:       0,
			Y:       154,
			W:       26,
			H:       26,
		},
		BigIcon: game.Sprite{
			Texture: "hud.png",
			X:       0,
			Y:       102,
			W:       52,
			H:       52,
		},
	},
	Wolf: Recipe{
		ActionPoints: 60,
		Preparation:  50,
		Name:         "Hustle",
		Sex:          game.Male,
		Profession:   "Wolf",
		SmallIcon: game.Sprite{
			Texture: "hud.png",
			X:       52,
			Y:       76,
			W:       26,
			H:       26,
		},
		BigIcon: game.Sprite{
			Texture: "hud.png",
			X:       52,
			Y:       24,
			W:       52,
			H:       52,
		},
	},
}
