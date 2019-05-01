package game

import (
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// HUD holds rendering state.
type HUD struct {
}

// NewHUD for a screen of width and height.
func NewHUD(width, height int) *HUD {
	return &HUD{}
}

// Render the HUD to the passed screen.
func (hud *HUD) Render(screen *ebiten.Image, zoom, w, h float64) {
	img, _, err := ebitenutil.NewImageFromFile("hud.png", ebiten.FilterNearest)
	if err != nil {
		log.Fatalf("load texture: %v", err)
	}

	for i := 0; i < 5; i++ {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(16+float64(24*i), -28)
		op.GeoM.Scale(zoom, zoom)
		op.GeoM.Translate(0, h)
		screen.DrawImage(img, op)
	}
}
