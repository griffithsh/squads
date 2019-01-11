package game

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

// direction calculates the compass direction given an origin (om,on) and a
// destination (dm,dn).
func direction(om, on, dm, dn int) DirectionType {
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
