package game

type MoveIntent struct {
	X, Y float64
}

// Type of this Component.
func (a *MoveIntent) Type() string {
	return "MoveIntent"
}
