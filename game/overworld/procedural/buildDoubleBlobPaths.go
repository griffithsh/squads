package procedural

import (
	"github.com/griffithsh/squads/geom"
)

// buildDoubleBlobPaths is intended to implement maps where there should be a
// single intersection point between two distinct mazes. Imagine a dark forest
// with a river running through the middle where there is only one bridge.
func buildDoubleBlobPaths(seed int64, level int) (Paths, error) {
	// prng := rand.New(rand.NewSource(seed))
	return map[geom.Key]Placement{}, nil
}
