package procedural

import (
	"fmt"

	"github.com/griffithsh/squads/geom"
)

// Paths represents a procedurally generated set of pathways.
// FIXME: this should also include a start point and end point.
type Paths map[geom.Key]Placement

func (paths Paths) Connect(a, b geom.Key) {
	to, ok := a.Neighbors()[b]
	if !ok {
		panic(fmt.Sprintf("cannot connect non-neighbors %v and %v", a, b))
	}
	if _, ok := paths[a]; !ok {
		paths[a] = Placement{
			Connections: map[geom.DirectionType]struct{}{},
		}
	}
	paths[a].Connections[to] = struct{}{}

	back, ok := b.Neighbors()[a]
	if !ok {
		panic(fmt.Sprintf("cannot connect non-neighbors %v and %v", b, a))
	}
	if _, ok := paths[b]; !ok {
		paths[b] = Placement{
			Connections: map[geom.DirectionType]struct{}{},
		}
	}
	paths[b].Connections[back] = struct{}{}
}

type pathFunc func(seed int64, level int) (Paths, error)

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
