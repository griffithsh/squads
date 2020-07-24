package game

import (
	"math"
	"time"

	"github.com/griffithsh/squads/event"
)

// Camera is a class that stores focus and zoom values
type Camera struct {
	bus              event.Bus
	focusX, focusY   float64
	zoom             float64
	screenW, screenH int

	panning          bool
	targetX, targetY float64
}

// NewCamera creates a new camera for a view of the requested width and height
func NewCamera(width, height int, bus *event.Bus) *Camera {
	c := Camera{
		focusX:  0,
		focusY:  0,
		zoom:    2.0,
		screenW: width,
		screenH: height,
	}
	bus.Subscribe(WindowSizeChanged{}.Type(), c.handleWindowSizeChanged)
	bus.Subscribe(SomethingInteresting{}.Type(), c.handleSomethingInteresting)
	return &c
}

// Update the camera so that interesting things can be panned to.
func (c *Camera) Update(elapsed time.Duration) {
	if !c.panning {
		return
	}

	a, b := c.targetX-c.focusX, c.targetY-c.focusY
	hypo := math.Sqrt(a*a + b*b)
	if hypo <= 1.0 {
		// end
		c.panning = false
		c.Center(c.targetX, c.targetY)
		return
	}
	c.Center(a/8+c.focusX, b/6+c.focusY)
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

// GetW returns the width of the game's window.
func (c *Camera) GetW() int {
	return c.screenW
}

// GetH returns the height of the game's window.
func (c *Camera) GetH() int {
	return c.screenH
}

func (c *Camera) handleWindowSizeChanged(e event.Typer) {
	wsc := e.(*WindowSizeChanged)
	c.screenW, c.screenH = wsc.NewW, wsc.NewH
}

func (c *Camera) handleSomethingInteresting(t event.Typer) {
	ev := t.(*SomethingInteresting)
	c.panning = true
	c.targetX, c.targetY = ev.X, ev.Y
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
