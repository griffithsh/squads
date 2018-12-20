package main

import (
	"math"

	"github.com/faiface/pixel"
)

// Camera is a class that composes View Matrices for rendering with faiface/pixel.
type Camera struct {
	pos        pixel.Vec
	zoom       float64
	viewBounds pixel.Rect
}

// NewCamera creates a new camera for a view of the requested width and height
func NewCamera(width, height float64) *Camera {
	return &Camera{
		pos:        pixel.Vec{X: 0, Y: 0},
		zoom:       3.0,
		viewBounds: pixel.Rect{Max: pixel.Vec{X: width, Y: height}},
	}
}

// View composes the camera's projection matrix.
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

// GetX coordinate of the camera.
func (c *Camera) GetX() float64 {
	return c.pos.X
}

// SetX coordinate of the camera.
func (c *Camera) SetX(x float64) {
	c.pos.X = x
}

// GetY coordinate of the camera.
func (c *Camera) GetY() float64 {
	return c.pos.Y
}

// SetY coordinate of the camera.
func (c *Camera) SetY(y float64) {
	c.pos.Y = y
}

// GetZoom of the camera.
func (c *Camera) GetZoom() float64 {
	return c.zoom
}

// SetZoom of the camera.
func (c *Camera) SetZoom(zoom float64) {
	// don't zoom out further than 1:1 pixel ratio
	if zoom < 1.0 {
		zoom = 1.0
	}

	c.zoom = zoom
}
