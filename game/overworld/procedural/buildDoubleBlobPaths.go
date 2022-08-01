package procedural

import (
	"math/rand"

	"github.com/griffithsh/squads/geom"
)

// buildDoubleBlobPaths is intended to implement maps where there should be a
// single intersection point between two distinct mazes. Imagine a dark forest
// with a river running through the middle where there is only one bridge.
func buildDoubleBlobPaths(prng *rand.Rand, level int) Paths {
	return map[geom.Key]Placement{}
}
