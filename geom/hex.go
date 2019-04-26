package geom

import "fmt"

// Hex is a hexagon tile that the play field is composed from.
type Hex struct {
	M, N      int
	neighbors []*Hex
}

// String to implement Stringer.
func (h *Hex) String() string {
	return fmt.Sprintf("%d,%d (%f,%f)", h.M, h.N, h.X(), h.Y())
}

// X coordinate of the center of this hexagon.
func (h *Hex) X() float64 {
	oddXOffset := xStride / 2
	return (hexWidth / 2) + float64((xStride*h.M)+(h.N%2*oddXOffset))
}

// Y coordinate of the center of this hexagon.
func (h *Hex) Y() float64 {
	return yStride + float64(yStride*h.N)
}

// Neighbors of this Hex.
func (h *Hex) Neighbors() []*Hex {
	return h.neighbors
}
