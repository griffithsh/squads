package main

type Vec struct {
	X, Y float64
}
type Rect struct {
	Min, Max Vec
}

// Camera is a class that stores focus and zoom values
type Camera struct {
	pos        Vec
	zoom       float64
	viewBounds Rect
}

// NewCamera creates a new camera for a view of the requested width and height
func NewCamera(width, height float64) *Camera {
	return &Camera{
		pos:        Vec{X: 0, Y: 0},
		zoom:       3.0,
		viewBounds: Rect{Max: Vec{X: width, Y: height}},
	}
}

// Center the view on a point.
func (c *Camera) Center(p Vec) {
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
