package geom

// Hex4 represents four hexes that can be occupied by a medium-sized unit.
type Hex4 struct {
	M, N      int // corresponds to the northern Hex.
	hexes     map[DirectionType]*Hex
	neighbors []*Hex4
}

// X coordinate of the center of this hexagon.
func (h *Hex4) X() float64 {
	oddXOffset := xStride / 2
	return (hexWidth / 2) + float64((xStride*h.M)+(h.N%2*oddXOffset))
}

// Y coordinate of the center of this hexagon.
func (h *Hex4) Y() float64 {
	return yStride + float64(yStride*h.N) + yStride
}

// Neighbors of this Hex.
func (h *Hex4) Neighbors() []*Hex4 {
	return h.neighbors
}

// Hexes returns the Hexes that compose this Hex4.
func (h *Hex4) Hexes() []*Hex {
	var result []*Hex
	for _, hex := range h.hexes {
		result = append(result, hex)
	}
	return result
}
