package game

import (
	"time"

	"github.com/griffithsh/squads/geom"
)

type CombatMapRecipeHexFrame struct {
	Texture  string
	X, Y     int
	W, H     int
	Duration time.Duration
}

type CombatMapRecipeVisual struct {
	Frames           []CombatMapRecipeHexFrame
	XOffset, YOffset int
	Layer            int
	Obscures         []geom.Key
}

type CombatMapRecipeHex struct {
	Position geom.Key
	Obstacle ObstacleType
	Visuals  []CombatMapRecipeVisual
}
type CombatMapRecipe struct {
	Hexes        []CombatMapRecipeHex
	Starts       []geom.Key
	TileW, TileH int
}
