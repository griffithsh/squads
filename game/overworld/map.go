package overworld

import (
	"sort"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/geom"
	"github.com/griffithsh/squads/squad"
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

	// Enemies stores rolled enemy squad locations and their types.
	Enemies map[geom.Key]squad.RecipeID

	// Start stores the rolled location for where the player should start in
	// this overworld map.
	Start geom.Key
}

// SortedNodeKeys provides the geom.Keys that appear in the Nodes map, sorted by
// M, then N.
func (m *Map) SortedNodeKeys() []geom.Key {
	result := []geom.Key{}
	for k := range m.Nodes {
		result = append(result, k)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].M != result[j].M {
			return result[i].M < result[j].M
		}
		return result[i].N < result[j].N
	})
	return result
}

type TileID int

const (
	Grass TileID = iota
	Stone
	Trees
)
