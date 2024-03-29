package procedural

import (
	"fmt"

	"github.com/griffithsh/squads/geom"
)

// Paths represents a procedurally generated set of pathways.
type Paths struct {
	Algorithm string
	Seed      int64

	Start    geom.Key
	Goal     geom.Key
	Nodes    map[geom.Key]Placement
	Specials map[string][]geom.Key
}

func (paths Paths) Connect(a, b geom.Key) {
	to, ok := a.Neighbors()[b]
	if !ok {
		panic(fmt.Sprintf("cannot connect non-neighbors %v and %v", a, b))
	}
	if paths.Nodes == nil {
		paths.Nodes = map[geom.Key]Placement{}
	}
	if _, ok := paths.Nodes[a]; !ok {
		paths.Nodes[a] = Placement{
			Connections: map[geom.DirectionType]struct{}{},
		}
	}
	paths.Nodes[a].Connections[to] = struct{}{}

	back, ok := b.Neighbors()[a]
	if !ok {
		panic(fmt.Sprintf("cannot connect non-neighbors %v and %v", b, a))
	}
	if _, ok := paths.Nodes[b]; !ok {
		paths.Nodes[b] = Placement{
			Connections: map[geom.DirectionType]struct{}{},
		}
	}
	paths.Nodes[b].Connections[back] = struct{}{}
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
