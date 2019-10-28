package overworld

import "github.com/griffithsh/squads/geom"

type Node struct {
	ID geom.Key
	// Connected neighbors by DirectionType.
	Directions map[geom.DirectionType]geom.Key
}

// Data that describes an overworld.
type Data struct {
	Nodes map[geom.Key]*Node
}
