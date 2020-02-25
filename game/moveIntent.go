package game

// MoveIntent is a Component that indicates that this Entity should move.
type MoveIntent struct {
	X, Y float64
}

// Type of this Component.
func (MoveIntent) Type() string {
	return "MoveIntent"
}
