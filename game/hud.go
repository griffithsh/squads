package game

import (
	"fmt"
	"math"

	"github.com/faiface/pixel"
	"github.com/griffithsh/squads/pixelutil"
)

// HUD holds rendering state.
type HUD struct {
	w, h int

	zoom float64

	composed []*pixel.Batch
	dirty    bool
}

// NewHUD for a screen of width and height.
func NewHUD(width, height int) *HUD {
	hud := &HUD{
		w:    width,
		h:    height,
		zoom: 3,
	}
	if err := hud.compose(); err != nil {
		// FIXME: panic? log? something else?
		return nil
	}

	return hud
}

// compose the HUD information into pixel.Sprites from the game's state.
func (hud *HUD) compose() error {
	p, err := pixelutil.LoadPicture("hud.png")
	if err != nil {
		return fmt.Errorf("pixelutil.LoadPicture: %v", err)
	}

	batch := pixel.NewBatch(&pixel.TrianglesData{}, p)
	cam := pixel.IM.Scaled(pixel.ZV, hud.zoom)
	batch.SetMatrix(cam)

	// Compose some icons along the bottom of the screen
	for i := 0; i < 5; i++ {
		sprite := pixelutil.NewSprite(p, 0, 0, 16, 24)

		// Using faiface/pixel inverted Y coordinates for simple fixed
		// positioning along the bottom of the screen.
		move := pixel.Vec{X: 16 + float64(24*i), Y: 12 + 2}

		sprite.Draw(batch, pixel.IM.Moved(move))
	}

	hud.composed = []*pixel.Batch{
		batch,
	}

	return nil
}

// Render the HUD to the passed Target.
func (hud *HUD) Render(win pixel.Target) {
	if hud.dirty {
		hud.compose()
		hud.dirty = false
	}
	for _, batch := range hud.composed {
		batch.Draw(win)
	}
}

// SetZoom of the HUD.
func (hud *HUD) SetZoom(zoom float64) {
	if hud.zoom != math.Round(zoom) {
		hud.zoom = math.Round(zoom)
		hud.dirty = true
	}
}
