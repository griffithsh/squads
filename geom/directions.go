package geom

//go:generate go run github.com/dmarkham/enumer -type=DirectionType,RelativeDirection -json

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

// FindDirection returns the direction one would have to travel in from origin to
// arrive at projection. Returned values are rounded when the precise value is
// not one of the 6 hexagonal directions.
func FindDirection(origin, projection Key) DirectionType {
	if origin.M == projection.M {
		// we know it's N or S
		if origin.N > projection.N {
			return N
		}
		return S
	} else if origin.M < projection.M {
		// We know it's NE or SE
		if origin.M%2 == 0 {
			if origin.N > projection.N {
				return NE
			}
			return SE
		}
		if origin.N >= projection.N {
			return NE
		}
		return SE
	} else {
		// We know it's NW or SW
		if origin.M%2 == 0 {
			if origin.N > projection.N {
				return NW
			}
			return SW
		} else {
			if origin.N >= projection.N {
				return NW
			}
			return SW
		}
	}
}

type RelativeDirection int

const (
	Forward RelativeDirection = iota
	Behind
	ForwardLeft
	ForwardRight
	BackLeft
	BackRight
)

// Actualize a relative direction and a concrete direction together to form a
// new concrete direction. N and Behind make S. S and Forward make S etc.
func Actualize(facing DirectionType, relative RelativeDirection) DirectionType {
	switch relative {
	case Behind:
		return map[DirectionType]DirectionType{
			N:  S,
			S:  N,
			NE: SW,
			SW: NE,
			NW: SE,
			SE: NW,
		}[facing]
	case ForwardLeft:
		return map[DirectionType]DirectionType{
			N:  NW,
			S:  SE,
			NE: N,
			SW: S,
			NW: SW,
			SE: NE,
		}[facing]
	case ForwardRight:
		return map[DirectionType]DirectionType{
			N:  NE,
			S:  SW,
			NE: SE,
			SW: NW,
			NW: N,
			SE: S,
		}[facing]
	case BackLeft:
		return map[DirectionType]DirectionType{
			N:  SW,
			S:  NE,
			NE: NW,
			SW: SE,
			NW: S,
			SE: N,
		}[facing]
	case BackRight:
		return map[DirectionType]DirectionType{
			N:  SE,
			S:  NW,
			NE: S,
			SW: N,
			NW: NE,
			SE: SW,
		}[facing]
	case Forward:
		fallthrough
	default:
		return facing
	}
}
