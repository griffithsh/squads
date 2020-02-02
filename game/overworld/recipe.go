package overworld

import "github.com/griffithsh/squads/geom"

// Recipe stores configuration of how to roll a map.
type Recipe struct {
	// Terrain stores the visible tiles of an overworld.
	Terrain map[geom.Key]TileID

	// Interesting stores locations in the map that should always be included as
	// part of the generated Nodes.
	Interesting []int
}

// Recipes stores all possible recipes of overworld maps.
var Recipes = []Recipe{
	{},
}
