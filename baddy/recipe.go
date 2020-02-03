package baddy

import (
	"github.com/griffithsh/squads/game"
)

// RecipeID identifies a baddy recipe.
type RecipeID int

const (
	Skellington RecipeID = iota
	Wolf
	Giant
)

// Recipe describes a way to construct a baddy.
type Recipe struct {
	ActionPoints       int
	Preparation        int
	Name               string
	Sex                game.CharacterSex
	Profession         game.CharacterProfession
	SmallIcon, BigIcon game.Sprite
}

// Construct a baddy from a Recipe.
func (recipe Recipe) Construct() *game.Character {
	return &game.Character{
		Name:                 recipe.Name,
		Sex:                  recipe.Sex,
		Profession:           recipe.Profession,
		PreparationThreshold: recipe.Preparation,
		ActionPoints:         recipe.ActionPoints,
		SmallIcon:            recipe.SmallIcon,
		BigIcon:              recipe.BigIcon,
	}
}

// Recipes is the reference of baddy recipes in the game.
var Recipes map[RecipeID]Recipe = map[RecipeID]Recipe{
	Skellington: Recipe{
		ActionPoints: 60,
		Preparation:  1650,
		Name:         "Dumble",
		Sex:          game.Male,
		Profession:   game.Skeleton,
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
		Preparation:  1500,
		Name:         "Hustle",
		Sex:          game.Male,
		Profession:   game.Wolf,
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
	Giant: Recipe{
		ActionPoints: 60,
		Preparation:  1050,
		Name:         "Icarion",
		Sex:          game.Male,
		Profession:   game.Giant,
		SmallIcon: game.Sprite{
			Texture: "hud.png",
			X:       104,
			Y:       76,
			W:       26,
			H:       26,
		},
		BigIcon: game.Sprite{
			Texture: "hud.png",
			X:       104,
			Y:       24,
			W:       52,
			H:       52,
		},
	},
}
