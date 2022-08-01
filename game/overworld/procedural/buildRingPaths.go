package procedural

import (
	"math/rand"

	"github.com/griffithsh/squads/geom"
)

// buildRingPaths might implement lakes, islands with impassable interiors, or atolls.
func buildRingPaths(prng *rand.Rand, level int) Paths {
	// pick N points on a ring, connect each to the next in turn, connecting the last to the first.
	return map[geom.Key]Placement{}
}
