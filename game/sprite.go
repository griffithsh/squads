package game

// Sprite is a renderable slice of a texture.
type Sprite struct {
	Texture    string
	X, Y, W, H int

	// This could include a color, but does not for now, as there are no uses for it.
	// Color *color.RGBA
}

// Type of this Component.
func (s *Sprite) Type() string {
	return "Sprite"
}
