package targeting

type BrushType int

//go:generate stringer -type=BrushType

// BrushTypeFromString converts strings to SelectableTypes. Mostly useful
// for Unmarshaling.
func BrushTypeFromString(s string) *BrushType {
	for i := 0; i < len(_BrushType_index)-1; i++ {
		c := BrushType(i)

		if c.String() == s {
			return &c
		}
	}
	return nil
}

const (
	// SingleHex means that only the targeted Hex is highlighted.
	SingleHex BrushType = iota

	// WithinRangeOfTarget indicates that all hexes that are within the
	// configured ranges of the target are painted.
	WithinRangeOfTarget

	// WithinRangeOfOrigin indicates that all hexes that are within the
	// configured ranges of the origin are painted.
	WithinRangeOfOrigin

	// LinearBrush uses a direction and length to indicate which hexes are painted.
	// LinearBrush

	// Arc? Radius (probably 1, or maye 2?), Length (how many hexes to paint
	// around the circle, negative for anticlockwise), Begin (if you want the
	// selected adjacent hex, as well as the ones to each side, you can begin on
	// -1 and use length 3) .
	// Arc

	// None means no selected hexes. This could be used for ... what? Nothing?
	// None
)

type Brush struct {
	Type BrushType

	// Only for AoE BrushTypes.
	MinRange int
	MaxRange int

	// LinearLength    int
	// LinearDirection geom.DirectionType // should be a contextual direction - "Forward", "Back-Left"
}
