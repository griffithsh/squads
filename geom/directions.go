package geom

//go:generate stringer -type=DirectionType

// DirectionType enumerates directions.
type DirectionType int

// DirectionTypes represent the 6 edges of a hexagon.
const (
	N DirectionType = iota
	NE
	SE
	S
	SW
	NW
)

// Opposite provides the 180 degree opposite to any cardinal DirectionType.
var Opposite = map[DirectionType]DirectionType{
	S:  N,
	SW: NE,
	NW: SE,
	N:  S,
	NE: SW,
	SE: NW,
}
