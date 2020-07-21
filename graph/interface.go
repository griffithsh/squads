package graph

// Vertex is a node in the graph. It's an empty interface so that consumers of
// this package are free to provide any type they please.
type Vertex interface{}

// CostFunc is a function that provides the cost of moving from one Vertex to another.
type CostFunc func(Vertex, Vertex) float64

// EdgeFunc is a function that provides the Vertices that are connected to a Vertex.
type EdgeFunc func(Vertex) []Vertex

// Heuristic is a function that provides a guess at the distance between two Vertices.
type Heuristic func(Vertex, Vertex) float64
