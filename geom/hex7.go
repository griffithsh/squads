package geom

// Hex7 represents a Hexagon and the surrounding 6 hexagons.
type Hex7 struct {
	M, N      int // corresponds to the center Hex.
	hexes     map[DirectionType]*Hex
	neighbors []*Hex7
}

// X coordinate of the center of this hexagon.
func (h *Hex7) X() float64 {
	oddXOffset := xStride / 2
	return (hexWidth / 2) + float64((xStride*h.M)+(h.N%2*oddXOffset))
}

// Y coordinate of the center of this hexagon.
func (h *Hex7) Y() float64 {
	return yStride + float64(yStride*h.N)
}

// Key returns the M,N coordinates of this Hex7.
func (h *Hex7) Key() Key {
	return Key{h.M, h.N}
}

// Hexes returns the Hexes that compose this Hex7.
func (h *Hex7) Hexes() []*Hex {
	var result []*Hex
	for _, hex := range h.hexes {
		result = append(result, hex)
	}
	return result
}
