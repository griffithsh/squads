package targeting

import "github.com/griffithsh/squads/geom"

type BrushType int

//go:generate go run github.com/dmarkham/enumer -type=BrushType -json

const (
	// SingleHex means that only the targeted Hex is highlighted.
	SingleHex BrushType = iota

	// WithinRangeOfTarget indicates that all hexes that are within the
	// configured ranges of the target are painted.
	WithinRangeOfTarget

	// WithinRangeOfOrigin indicates that all hexes that are within the
	// configured ranges of the origin are painted.
	WithinRangeOfOrigin

	// LinearFromOrigin uses a direction and length to indicate which hexes are painted.
	LinearFromOrigin

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

	LinearExtent    int
	LinearDirection geom.RelativeDirection
}
