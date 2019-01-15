package game

// Facer is a Component that represents a direction to face in.
type Facer struct {
	Face DirectionType
}

// Type of this Component.
func (f *Facer) Type() string {
	return "Facer"
}
