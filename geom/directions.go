package geom

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

// Direction calculates the compass direction given an origin (om,on) and a
// destination (dm,dn). This function only works with adjacent Hexes.
func Direction(om, on, dm, dn int) DirectionType {
	switch {
	case on+2 == dn:
		return S
	case on-2 == dn:
		return N
	}
	if on%2 == 0 {
		switch {
		case om-1 == dm && on-1 == dn:
			return NW
		case om-1 == dm && on+1 == dn:
			return SW
		case om == dm && on-1 == dn:
			return NE
		case om == dm && on+1 == dn:
			return SE
		}
	} else {
		switch {
		case om == dm && on-1 == dn:
			return NW
		case om == dm && on+1 == dn:
			return SW
		case om+1 == dm && on-1 == dn:
			return NE
		case om+1 == dm && on+1 == dn:
			return SE
		}
	}

	return N
}
