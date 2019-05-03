package geom

import "fmt"

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
