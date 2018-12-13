package game

import (
	"fmt"
	"image/color"

	"github.com/griffithsh/squads/ecs"
)

const (
	xStride = 34
	yStride = 8
)

// Hex is a hexagon tile that an Actor can occupy.
type Hex struct {
	x, y int
}

func (h *Hex) Type() string {
	return "Hex"
}

// Board to play out encounters on. Collection of Hexes.
type Board struct {
	stride int // how many hexes are in a row
	hexes  []Hex
}

func NewBoard(mgr *ecs.World, w, h int) (*Board, error) {
	top, left := 0, 0

	oddXOffset := 17

	arr := make([]Hex, w*h)
	for i := 0; i < w*h; i++ {
		e := mgr.NewEntity()
		arr[i] = Hex{
			x: i % w,
			y: i / w,
		}
		mgr.AddComponent(e, &arr[i])
		mgr.AddComponent(e, &Sprite{
			Texture: "texture.png",
			X:       0,
			Y:       0,
			W:       24,
			H:       16,
			Color:   &color.RGBA{150, 150, 150, 63},
		})

		mgr.AddComponent(e, &Position{
			X:     float64(left + (xStride * arr[i].x) + (arr[i].y % 2 * oddXOffset)),
			Y:     float64(top+(yStride*arr[i].y)) - (16 / 2),
			Layer: 0,
		})

		// if i == 27 || i == 93 || i == 42 || i == 111 || i == 155 || i == 2 {
		if i%11 == 1 || i%17 == 1 || i%13 == 1 {
			e := mgr.NewEntity()
			mgr.AddComponent(e, &Sprite{
				Texture: "Untitled.png",
				X:       0,
				Y:       0,
				W:       24,
				H:       48,
			})
			mgr.AddComponent(e, &Position{
				X:     float64(left + (xStride * arr[i].x) + (arr[i].y % 2 * oddXOffset)),
				Y:     float64(top+(yStride*arr[i].y)) - (48 / 2),
				Layer: 1,
			})
		} else if i%19 == 0 {
			e := mgr.NewEntity()
			mgr.AddComponent(e, &Sprite{
				Texture: "Untitled.png",
				X:       24,
				Y:       0,
				W:       24,
				H:       48,
			})
			mgr.AddComponent(e, &Position{
				X:     float64(left + (xStride * arr[i].x) + (arr[i].y % 2 * oddXOffset)),
				Y:     float64(top+(yStride*arr[i].y)) - (48 / 2),
				Layer: 1,
			})
		}
	}

	return &Board{
		stride: w,
		hexes:  arr,
	}, nil
}

func (b *Board) At(x, y int) (*Hex, error) {
	ind := b.stride*y + x
	if ind >= len(b.hexes) {
		return nil, fmt.Errorf("not found: (%d,%d)", x, y)
	}
	return &b.hexes[b.stride*y+x], nil
}

func (b *Board) Width() float64 {
	// TODO!
	return float64(b.stride * xStride)

}

func (b *Board) Height() float64 {
	// TODO!
	return float64(len(b.hexes) / b.stride * yStride)
}
