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

// Key returns the M,N coordinates of this Hex4.
func (h *Hex4) Key() Key {
	return Key{h.M, h.N}
}

// Hexes returns the Hexes that compose this Hex4.
func (h *Hex4) Hexes() []*Hex {
	var result []*Hex
	for _, hex := range h.hexes {
		result = append(result, hex)
	}
	return result
}

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
