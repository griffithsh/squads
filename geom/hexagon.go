package geom

// Hex is a part of a Field.
type Hex struct {
	// N, NE, SE, S, SW, NW neighbors
	x, y float64
	key  Key
}

// Key of the Hexagon. Where it exists in the Field grid.
func (h Hex) Key() Key {
	return h.key
}

// Center of the Hexagon in world coordinates.
func (h Hex) Center() (float64, float64) {
	return h.x, h.y
}
