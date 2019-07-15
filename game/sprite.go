package game

// Sprite is a renderable slice of a texture.
type Sprite struct {
	Texture          string
	X, Y, W, H       int
	OffsetX, OffsetY int

	// This could include a color, but does not for now, as there are no uses for it.
	// Color *color.RGBA
	// NB It might make more sense to keep Color/Tint as a separate Component
	// so that it can be applied to non-sprite renderable things, like solid-
	// color or bordered shapes.
}

// Type of this Component.
func (s *Sprite) Type() string {
	return "Sprite"
}
