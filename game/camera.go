package game

// Camera is a class that stores focus and zoom values
type Camera struct {
	focusX, focusY   float64
	zoom             float64
	screenW, screenH int
}

// NewCamera creates a new camera for a view of the requested width and height
func NewCamera(width, height int) *Camera {
	return &Camera{
		focusX:  0,
		focusY:  0,
		zoom:    3.0,
		screenW: width,
		screenH: height,
	}
}

// Center the view on a point.
func (c *Camera) Center(x, y float64) {
	c.focusX = x
	c.focusY = y
}

// GetX coordinate of the camera.
func (c *Camera) GetX() float64 {
	return c.focusX
}

// SetX coordinate of the camera.
func (c *Camera) SetX(x float64) {
	c.focusX = x
}

// GetY coordinate of the camera.
func (c *Camera) GetY() float64 {
	return c.focusY
}

// SetY coordinate of the camera.
func (c *Camera) SetY(y float64) {
	c.focusY = y
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

// ScreenToWorld translates screen coordinates to world coordinates based on the
// current zoom and focus values.
func (c *Camera) ScreenToWorld(sx, sy int) (wx, wy float64) {
	zoom := c.GetZoom()
	wx, wy = float64(sx)/zoom, float64(sy)/zoom

	// Correct for camera focus.
	wx, wy = wx+c.GetX(), wy+c.GetY()

	// Correct for size of screen (!?).
	wx, wy = wx-float64(c.screenW)/2/zoom, wy-float64(c.screenH)/2/zoom

	return wx, wy
}

// Modulo an x,y screen coordinate so that things that are beyond the limits of
// the screen are wrapped around. Use case is when you want to display something
// aligned to the bottom of the screen, you can set a negative y coordinate, and
// Modulo will convert that to screen height - y.
func (c *Camera) Modulo(x, y int) (int, int) {
	mx := x % c.screenW
	if mx < 0 {
		mx = mx + c.screenW
	}
	my := y % c.screenH
	if my < 0 {
		my = my + c.screenH
	}
	return mx, my
}
