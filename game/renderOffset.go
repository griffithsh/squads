package game

// RenderOffset is a value that offsets the position the Component should be
// rendered at. It can be used for visual effects or animations that should not
// affect the real Position of the Entity, just where it appears to be.
type RenderOffset struct {
	X, Y int
}

// Type of this Component.
func (s *RenderOffset) Type() string {
	return "RenderOffset"
}
