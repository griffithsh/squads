package game

import (
	"image/color"
)

type Sprite struct {
	Texture    string
	X, Y, W, H int

	Color *color.RGBA
}

func (s *Sprite) Type() string {
	return "Sprite"
}
