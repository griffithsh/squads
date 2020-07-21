package geom

import "math"

// DistanceSquared returns the distance between two Hexagons multiplied by
// itself. This is useful to use as a pathfinding heuristic when these values
// only need to be compared to other outputs of this function, and the cost
// of finding the square root adds nothing.
func DistanceSquared(h1, h2 *Hex) float64 {
	x1, y1 := h1.Center()
	x2, y2 := h2.Center()
	return (x1-x2)*(x1-x2) + (y1-y2)*(y1-y2)
}

// Distance between two Hexagons.
func Distance(h1, h2 *Hex) float64 {
	return math.Sqrt(DistanceSquared(h1, h2))
}
