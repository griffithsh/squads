package procedural

import (
	"math/rand"

	"github.com/griffithsh/squads/geom"
)

// buildVinePaths generates long maps, as a single correct pathway with
// alternative dead-ending branches is generated trending in a single direction.
// It is intended to implement beaches.
func buildVinePaths(prng *rand.Rand, level int) Paths {
	// start from a seed, grow in a direction
	// allow short branches

	return map[geom.Key]Placement{}
}
