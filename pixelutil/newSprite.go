package pixelutil

import "github.com/faiface/pixel"

// NewSprite wraps the pixel.NewSprite method, correcting faiface/pixel's
// inverted Y coordinates, so that slices of a texture can be created with
// texture coordinates.
func NewSprite(pic pixel.Picture, x, y, w, h int) *pixel.Sprite {
	height := pic.Bounds().H()
	r := pixel.R(float64(x), height-float64(y), float64(x+w), height-float64(h+y))
	return pixel.NewSprite(pic, r)
}
