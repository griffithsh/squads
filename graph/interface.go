package graph

// Vertex is a node in the graph.
type Vertex interface{}

// CostFunc is a function that provides the cost of moving from one Vertex to another.
type CostFunc func(Vertex, Vertex) float64

// EdgeFunc is a function that provides connected Vertices.
type EdgeFunc func(Vertex) []Vertex

// Heuristic is a function that provides a guess at the distance between two Vertices.
type Heuristic func(Vertex, Vertex) float64
