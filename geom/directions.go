package geom

import (
	"fmt"
)

//go:generate stringer -type=DirectionType

// DirectionType enumerates directions.
type DirectionType int

// DirectionTypes represent the 6 directions that a Mover can face.
const (
	N DirectionType = iota
	NE
	SE
	S
	SW
	NW

	CENTER
)

// Direction calculates which hexagonal direction the vector of x,y aligns with.
func Direction(x, y float64) (DirectionType, error) {
	switch {
	case x == 0 && y == 16:
		return S, nil
	case x == 0 && y == -16:
		return N, nil
	case x == 17 && y == 8:
		return SE, nil
	case x == -17 && y == 8:
		return SW, nil
	case x == 17 && y == -8:
		return NE, nil
	case x == -17 && y == -8:
		return NW, nil
	}
	return DirectionType(0), fmt.Errorf("unhandled: %f,%f", x, y)
}

// Opposite provides the 180 degree opposite to any cardinal DirectionType.
var Opposite = map[DirectionType]DirectionType{
	S:  N,
	SW: NE,
	NW: SE,
	N:  S,
	NE: SW,
	SE: NW,
}

// BestFacing returns the DirectionType that best represents the direction that
// dest is from origin.
func BestFacing(origin, dest Key) DirectionType {
	o := Hex{M: origin.M, N: origin.N}
	d := Hex{M: dest.M, N: dest.N}

	dx := d.X() - o.X()
	dy := d.Y() - o.Y()

	if o.Y() <= d.Y() {
		// Then we're heading S, SE, or SW.
		if o.X() <= d.X() {
			// Then it's either S or SE
			if dx/dy <= 30/56 {
				return S
			}
			return SE
		}
		// Then it's either S or SW
		if dx/dy >= -30/65 {
			return S
		}
		return SW
	}
	// Then we're heading N, NE, or NW
	if d.X() >= o.X() {
		// Then it's either N or NE
		if dx/dy >= 30/-56 {
			return N
		}
		return NE
	}
	// Then it's either N or NW
	if dx/dy < -30/-56 {
		return N
	}
	return NW
}
