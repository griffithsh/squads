package overworld

import (
	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/geom"
)

// Node is a stop on the overworld that might be occupied by the player, an
// encounter with an enemy squad, the escape portal, a merchant, etc, or nothing
// at all.
type Node struct {
	ID geom.Key
	e  ecs.Entity

	// Connected neighbors mapped by the direction they lie.
	Connected map[geom.DirectionType]geom.Key
}

// Map of an overworld.
type Map struct {
	// Terrain stores the visible tiles of an overworld.
	Terrain map[geom.Key]TileID

	// Nodes stores the stops on the overworld, and the traversable paths
	// between them.
	Nodes map[geom.Key]*Node

	// Start stores the rolled location for where the player should start in
	// this overworld map.
	Start geom.Key
}

type TileID int

const (
	Grass TileID = iota
	Stone
	Trees
)
