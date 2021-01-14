package embark

import (
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/geom"
)

// Feature is something that appears on the embark screen, like a house or other
// doodad.
type Feature struct {
	Sprite game.Sprite

	// Coverage is a series of contiguous hexes that this feature also expands
	// over. One feature might not fit on a hex, so the way to describe which
	// other hexes it also requires to fit is by traversing the geometry
	// directions, and collecting the hexes traversed.
	Coverage []geom.DirectionType

	// PathConnect is the series of steps required to get from the anchor point
	// of this Feature (the lowest Y coordinate), to the hex that connects to a
	// roadway.
	PathConnect []geom.DirectionType
}

// Occupies calculates the geom.keys that this feature would take if it was
// placed at a given origin.
func (f *Feature) Occupies(origin geom.Key) []geom.Key {
	result := []geom.Key{origin}
	curr := origin
	for _, dir := range f.Coverage {
		curr = curr.Adjacent()[dir]
		result = append(result, curr)
	}

	return result
}

// StartOfPathFor figures out which hex should connect to a pathway for thie
// feature when it is placed at the passed origin.
func (f *Feature) StartOfPathFor(origin geom.Key) geom.Key {
	if len(f.PathConnect) == 0 {
		panic("There is no PathConnect for this feature")
	}
	curr := origin
	for _, dir := range f.PathConnect {
		curr = curr.Adjacent()[dir]
	}
	return curr
}

var houseFeatures = []Feature{
	Feature{
		Sprite: game.Sprite{
			Texture: "embark-tiles.png",

			X: 0, Y: 24,
			W: 50, H: 42,
			OffsetY: -15,
		},
		Coverage:    []geom.DirectionType{geom.N, geom.SW, geom.N},
		PathConnect: []geom.DirectionType{geom.NE},
	},
	Feature{
		Sprite: game.Sprite{
			Texture: "embark-tiles.png",

			X: 50, Y: 24,
			W: 50, H: 42,
			OffsetX: -15,
			OffsetY: -15,
		},
		Coverage:    []geom.DirectionType{geom.NW, geom.NE, geom.NW},
		PathConnect: []geom.DirectionType{geom.NW, geom.SW},
	},
	Feature{
		Sprite: game.Sprite{
			Texture: "embark-tiles.png",

			X: 100, Y: 24,
			W: 50, H: 42,
			OffsetX: 15,
			OffsetY: -15,
		},
		Coverage:    []geom.DirectionType{geom.NE, geom.NW, geom.NE},
		PathConnect: []geom.DirectionType{geom.NE, geom.SE},
	},
	Feature{
		Sprite: game.Sprite{
			Texture: "embark-tiles.png",

			X: 150, Y: 24,
			W: 50, H: 42,
			OffsetY: -15,
		},
		Coverage:    []geom.DirectionType{geom.N, geom.SE, geom.N},
		PathConnect: []geom.DirectionType{geom.NW},
	},
}

var faeGateFeature = Feature{
	Sprite: game.Sprite{
		Texture: "embark-tiles.png",

		X: 0, Y: 66,
		W: 50, H: 38,
		OffsetX: 0,
		OffsetY: -13,
	},
	Coverage: []geom.DirectionType{geom.NW, geom.N, geom.SE, geom.NE, geom.S},
}

var windmillFeature = Feature{
	Sprite: game.Sprite{
		Texture: "embark-tiles.png",

		X: 50, Y: 66,
		W: 50, H: 38,
		OffsetY: -13,
	},
	Coverage: []geom.DirectionType{geom.NE, geom.NW, geom.SW},
}

var flavorFeatures = []Feature{
	Feature{
		Sprite: game.Sprite{
			Texture: "embark-tiles.png",

			X: 0, Y: 203,
			W: 35, H: 29,
			OffsetX: 7,
			OffsetY: -9,
		},
		Coverage: []geom.DirectionType{geom.NE, geom.NW},
	},
	Feature{
		Sprite: game.Sprite{
			Texture: "embark-tiles.png",

			X: 35, Y: 203,
			W: 35, H: 29,
			OffsetX: -8,
			OffsetY: -9,
		},
		Coverage: []geom.DirectionType{geom.NW, geom.NE},
	},
	Feature{
		Sprite: game.Sprite{
			Texture: "embark-tiles.png",

			X: 0, Y: 232,
			W: 20, H: 12,
		},
	},
	Feature{
		Sprite: game.Sprite{
			Texture: "embark-tiles.png",

			X: 20, Y: 232,
			W: 20, H: 12,
		},
	},
	Feature{
		Sprite: game.Sprite{
			Texture: "embark-tiles.png",

			X: 70, Y: 203,
			W: 35, H: 29,
			OffsetX: 7,
			OffsetY: -9,
		},
		Coverage: []geom.DirectionType{geom.NE, geom.NW},
	},
	Feature{
		Sprite: game.Sprite{
			Texture: "embark-tiles.png",

			X: 105, Y: 203,
			W: 35, H: 29,
			OffsetX: -8,
			OffsetY: -9,
		},
		Coverage: []geom.DirectionType{geom.NW, geom.NE},
	},
}
