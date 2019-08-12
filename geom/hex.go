package geom

import "fmt"

// Hex is a hexagon tile that the play field is composed from.
type Hex struct {
	M, N                 int
	neighbors            []*Hex
	neighborsByDirection map[DirectionType]*Hex
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

// Key returns the M,N coordinates of this Hex.
func (h *Hex) Key() Key {
	return Key{h.M, h.N}
}

// Hexes returns the Hexes that compose this Hex. In all cases this is simply
// the hex itself. This method is implemented so that it can share a common
// interface with the composite Hex types Hex4 and Hex7.
func (h *Hex) Hexes() []*Hex {
	return []*Hex{h}
}
