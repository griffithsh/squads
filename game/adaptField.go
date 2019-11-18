package game

import (
	"github.com/griffithsh/squads/geom"
)

// AdaptField takes a geom.Field and adapts it to provide At() and Get() for the
// specified CharacterSize.
func AdaptField(f *geom.Field, sz CharacterSize) geom.LogicalField {
	switch sz {
	case MEDIUM:
		return geom.NewField4(f)
	case LARGE:
		return geom.NewField7(f)
	default:
		return geom.NewField1(f)

	}
}

func AdaptFieldObstacle(f *geom.Field, sz ObstacleType) geom.LogicalField {
	switch sz {
	case MediumCharacter:
		return geom.NewField4(f)
	case LargeCharacter:
		return geom.NewField7(f)
	default:
		return geom.NewField1(f)
	}
}
