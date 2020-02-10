package res

import "time"

// Frame is one frame of an Animation.
type Frame struct {
	Texture                      string
	X, Y, W, H, OffsetX, OffsetY int
	Duration                     time.Duration
}

// Animation is a collection of Frames.
type Animation struct {
	Frames []Frame
}

// Animations contains generic animations. Use one by calling game.NewFrameAnimation(res.Animations["..."]).
var Animations = map[string]Animation{
	"overworld-reveal-grass": {[]Frame{
		{"overworld-grass.png", 0, 96, 144, 96, 0, 0, 30 * time.Millisecond},
		{"overworld-grass.png", 0, 192, 144, 96, 0, 0, 30 * time.Millisecond},
		{"overworld-grass.png", 0, 288, 144, 96, 0, 0, 30 * time.Millisecond},
		{"overworld-grass.png", 0, 384, 144, 96, 0, 0, 30 * time.Millisecond},
	}},
	"overworld-hide-card": {[]Frame{
		{"overworld-cards.png", 256, 0, 128, 192, 0, 0, 30 * time.Millisecond},
		{"overworld-cards.png", 384, 0, 128, 192, 0, 0, 30 * time.Millisecond},
		{"overworld-cards.png", 0, 192, 128, 192, 0, 0, 40 * time.Millisecond},
		{"overworld-cards.png", 128, 192, 128, 192, 0, 0, 50 * time.Millisecond},
		{"overworld-cards.png", 256, 192, 128, 192, 0, 0, 60 * time.Millisecond},
		{"overworld-cards.png", 384, 192, 128, 192, 0, 0, 70 * time.Millisecond},
	}},
}
