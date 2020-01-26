package game

// Sprite is a renderable slice of a texture.
type Sprite struct {
	Texture    string
	X, Y, W, H int

	// OffsetX and OffsetY exist as they can be applied to individual frames of
	// an animation. They do not duplicate the functionality of RenderOffset, as
	// that applies at the Entity level, and could effect the rendering position
	// of things that are not Sprites. (i.e shape primitives)
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

// SpriteRepeat is a component that repeats (or "tiles") a sprite across a width
// and height.
type SpriteRepeat struct {
	W, H int
}

// Type of this Component.
func (*SpriteRepeat) Type() string {
	return "SpriteRepeat"
}
