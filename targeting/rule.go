package targeting

import (
	"fmt"

	"github.com/griffithsh/squads/geom"
)

type Rule struct {
	Selectable Selectable
	Brush      Brush
}

// Execute is meant to give you the Keys that should be painted by this rule as
// executed on the given selection and origin.
func (r *Rule) Execute(selected, origin geom.Key) (selectable bool, paints []geom.Key) {
	// You might need to know the Keys that are permissable selections.
	// You might need to know the Keys that would be painted by the brush.
	// It might be impractical to calculate permissable selections in all contexts.
	// Return whether the selected Key is permissable or not, and the Keys that would be painted by the Brush.
	switch r.Selectable.Type {
	case SelectAnywhere: // cannot know permissable selections without access to the geom.Field.
		selectable = true
	case SelectWithin:
		distance := origin.HexesFrom(selected)
		if r.Selectable.MinRange <= distance && r.Selectable.MaxRange >= distance {
			selectable = true
		}
	case Untargeted:
		selectable = true
	default:
		panic(fmt.Sprintf("unhandled SelectableType %s", r.Selectable.Type))
	}

	switch r.Brush.Type {
	case SingleHex:
		paints = []geom.Key{selected}
	case WithinRangeOfTarget:
	case WithinRangeOfOrigin:
	// case LinearBrush:
	default:
		panic(fmt.Sprintf("unhandled BrushType %s", r.Brush.Type))
	}
	return selectable, paints
}
