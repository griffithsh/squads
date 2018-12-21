package game

import (
	"image/color"
)

// Sprite is a renderable slice of a texture.
type Sprite struct {
	Texture    string
	X, Y, W, H int

	Color *color.RGBA
}

// Type of this Component.
func (s *Sprite) Type() string {
	return "Sprite"
}
