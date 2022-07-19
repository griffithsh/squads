package procedural

import (
	"math/rand"
	"time"

	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/geom"
)

type Placement struct {
	Connections map[geom.DirectionType]struct{}
	// what sprites to use for connections?
}

type Generated struct {
	Paths       map[geom.Key]Placement
	BaseTerrain game.Sprite
	RoadSprites map[geom.DirectionType]game.Sprite
}

// Generate should take a recipe and output an overworld map.
func Generate() Generated {
	prng := rand.New(rand.NewSource(time.Now().Unix()))
	// switch on recipe strategy

	// build paths

	// apply terrain rules to finished paths

	// overwrite with doodads

	// return an object that contains info on how to render the overworld as
	// well as programmatic info on what hexes are navigable.
	paths := buildVanillaPaths(prng, 0)

	roads := map[geom.DirectionType]game.Sprite{
		geom.N: {
			Texture: "temporary.png",
			W:       68,
			H:       34,
			X:       0,
			Y:       68,
		},
		geom.NE: {
			Texture: "temporary.png",
			W:       68,
			H:       34,
			X:       0,
			Y:       0,
		},
		geom.SE: {
			Texture: "temporary.png",
			W:       68,
			H:       34,
			X:       0,
			Y:       34,
		},
		geom.S: {
			Texture: "temporary.png",
			W:       68,
			H:       34,
			X:       68,
			Y:       68,
		},
		geom.SW: {
			Texture: "temporary.png",
			W:       68,
			H:       34,
			X:       68,
			Y:       34,
		},
		geom.NW: {
			Texture: "temporary.png",
			W:       68,
			H:       34,
			X:       68,
			Y:       0,
		},
	}

	return Generated{
		Paths: paths,
		BaseTerrain: game.Sprite{
			Texture: "temporary.png",
			W:       68,
			H:       34,
			X:       136,
			Y:       0,
		},
		RoadSprites: roads,
	}
}
