package procedural

import (
	"fmt"
	"math/rand"

	"github.com/griffithsh/squads/geom"
)

// Paths represents a procedurally generated set of pathways.
// FIXME: this should also include a start point and end point.
type Paths map[geom.Key]Placement

type pathFunc func(prng *rand.Rand, level int) Paths

func (f *pathFunc) UnmarshalJSON(b []byte) error {
	switch string(b) {
	case `"Maze"`:
		*f = buildMazePaths
		return nil
	case `"Vine"`:
		*f = buildVinePaths
		return nil
	case `"Ring"`:
		*f = buildRingPaths
		return nil
	case `"DoubleBlob"`:
		*f = buildDoubleBlobPaths
		return nil
	}
	return fmt.Errorf("unknown value for pathFunc %q", string(b))
}
