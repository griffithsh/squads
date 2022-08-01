package procedural

import (
	"math/rand"

	"github.com/griffithsh/squads/geom"
)

// Generator holds configuration data to construct overworld maps.
type Generator struct {
	MakePaths pathFunc       `json:"pathGeneration"`
	Terrain   TerrainBuilder `json:"terrainBuilder"`
}

// Placement holds info about what roads or paths have been placed on a Key.
type Placement struct {
	Connections map[geom.DirectionType]struct{}
}

type Generated struct {
	Paths Paths

	// PathExtents are only being exposed for debugging purposes.
	PathExtents map[geom.DirectionType]geom.Key
	Terrain     map[geom.Key]Code
}

// Generate should take a recipe and output an overworld map.
func (g Generator) Generate(seed int64, level int) Generated {
	prng := rand.New(rand.NewSource(seed))

	// build paths
	paths := g.MakePaths(prng, level)

	terrainCodes := g.Terrain.Build(prng, paths)

	// overwrite standard terrain with doodads
	// TODO: ...

	// TODO: baddies

	// TODO: misc other encounters

	// return an object that contains info on how to render the overworld as
	// well as programmatic info on what hexes are navigable.
	return Generated{
		Paths:       paths,
		PathExtents: extentsOf(keysOf(paths)),
		Terrain:     terrainCodes,
	}
}
