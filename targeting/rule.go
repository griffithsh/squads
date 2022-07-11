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
		paints = selected.ExpandBy(r.Selectable.MinRange, r.Selectable.MaxRange)
	case WithinRangeOfOrigin:
		paints = origin.ExpandBy(r.Selectable.MinRange, r.Selectable.MaxRange)
	case LinearFromOrigin:
		orientation := geom.FindDirection(origin, selected)

		direction := geom.Actualize(orientation, r.Brush.LinearDirection)
		mover := func(k geom.Key) geom.Key {
			switch direction {
			case geom.N:
				return k.ToN()
			case geom.S:
				return k.ToS()
			case geom.NE:
				return k.ToNE()
			case geom.NW:
				return k.ToNW()
			case geom.SE:
				return k.ToSE()
			case geom.SW:
				return k.ToSW()
			default:
				return k
			}
		}

		k := origin
		for i := 0; i < r.Brush.LinearExtent; i++ {
			k = mover(k)
			paints = append(paints, k)
		}

	default:
		panic(fmt.Sprintf("unhandled BrushType %s", r.Brush.Type))
	}
	return selectable, paints
}
