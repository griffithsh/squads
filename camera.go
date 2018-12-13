package main

import (
	"math"

	"github.com/faiface/pixel"
)

type Camera struct {
	pos        pixel.Vec
	zoom       float64
	viewBounds pixel.Rect
}

func NewCamera(w, h float64) *Camera {
	return &Camera{
		pos:        pixel.Vec{X: 0, Y: 0},
		zoom:       3.0,
		viewBounds: pixel.Rect{Max: pixel.Vec{X: w, Y: h}},
	}
}

func (c *Camera) View() pixel.Matrix {
	// faiface/pixel inverts the Y coordinate
	vFlip := c.pos
	vFlip.Y *= -1

	return pixel.IM.Scaled(vFlip, math.Round(c.zoom)).Moved(c.viewBounds.Center().Sub(vFlip))
}

// Center the view on a point.
func (c *Camera) Center(p pixel.Vec) {
	c.pos = p
}

func (c *Camera) GetX() float64 {
	return c.pos.X
}

func (c *Camera) SetX(x float64) {
	c.pos.X = x
}

func (c *Camera) GetY() float64 {
	return c.pos.Y
}

func (c *Camera) SetY(y float64) {
	c.pos.Y = y
}

func (c *Camera) GetZoom() float64 {
	return c.zoom
}
func (c *Camera) SetZoom(zoom float64) {
	// don't zoom out further than 1:1 pixel ratio
	if zoom < 1.0 {
		zoom = 1.0
	}

	c.zoom = zoom
}
