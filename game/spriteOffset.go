package game

// SpriteOffset is a correction value that offsets the position the Sprite hould
// be rendered at.
type SpriteOffset struct {
	X, Y int
}

// Type of this Component.
func (s *SpriteOffset) Type() string {
	return "SpriteOffset"
}
