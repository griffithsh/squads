package game

import (
	"time"

	"github.com/griffithsh/squads/geom"
)

type CombatMapRecipeHexFrame struct {
	Texture  string
	X, Y     int
	Duration time.Duration
}

type CombatMapRecipeVisual struct {
	Frames []CombatMapRecipeHexFrame
	Layer  int
}

type CombatMapRecipeHex struct {
	Position geom.Key
	Obstacle bool
	Visuals  []CombatMapRecipeVisual
}
type CombatMapRecipe struct {
	Hexes        []CombatMapRecipeHex
	Starts       []geom.Key
	TileW, TileH int
}
