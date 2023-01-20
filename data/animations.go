package data

import (
	"time"

	"github.com/griffithsh/squads/game"
)

// animations hard-codes some animations into the base squads game.
var animations = map[string]game.FrameAnimation{
	"overworld-reveal-grass": {
		Frames: []game.Sprite{
			{Texture: "deprecated/overworld-grass.png", X: 0, Y: 96, W: 144, H: 96},
			{Texture: "deprecated/overworld-grass.png", X: 0, Y: 192, W: 144, H: 96},
			{Texture: "deprecated/overworld-grass.png", X: 0, Y: 288, W: 144, H: 96},
			{Texture: "deprecated/overworld-grass.png", X: 0, Y: 384, W: 144, H: 96},
		},
		Timings: []time.Duration{
			30 * time.Millisecond,
			30 * time.Millisecond,
			30 * time.Millisecond,
			30 * time.Millisecond,
		},
	},
	"overworld-hide-card": {Frames: []game.Sprite{
		{Texture: "overworld/cards.png", X: 256, Y: 0, W: 128, H: 192},
		{Texture: "overworld/cards.png", X: 384, Y: 0, W: 128, H: 192},
		{Texture: "overworld/cards.png", X: 0, Y: 192, W: 128, H: 192},
		{Texture: "overworld/cards.png", X: 128, Y: 192, W: 128, H: 192},
		{Texture: "overworld/cards.png", X: 256, Y: 192, W: 128, H: 192},
		{Texture: "overworld/cards.png", X: 384, Y: 192, W: 128, H: 192},
	},
		Timings: []time.Duration{
			30 * time.Millisecond,
			30 * time.Millisecond,
			40 * time.Millisecond,
			50 * time.Millisecond,
			60 * time.Millisecond,
			70 * time.Millisecond,
		},
	},
}

func (a *Archive) GetAnimation(name string) game.FrameAnimation {
	return animations[name]
}
