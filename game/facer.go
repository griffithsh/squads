package game

import "github.com/griffithsh/squads/geom"

// Facer is a Component that represents a direction to face in.
type Facer struct {
	Face geom.DirectionType
}

// Type of this Component.
func (f *Facer) Type() string {
	return "Facer"
}
