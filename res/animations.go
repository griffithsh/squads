package res

import "time"

type Frame struct {
	Texture                      string
	X, Y, W, H, OffsetX, OffsetY int
	Duration                     time.Duration
}

type Animation struct {
	Frames []Frame
}

var All = map[string]Animation{
	"Villager-Male-Idle": {[]Frame{
		{"figure.png", 0, 0, 24, 48, 0, -16, 10000 * time.Millisecond},
		{"figure.png", 48, 0, 24, 48, 0, -16, 2000 * time.Millisecond},
		{"figure.png", 0, 0, 24, 48, 0, -16, 750 * time.Millisecond},
		{"figure.png", 24, 0, 24, 48, 0, -16, 300 * time.Millisecond},
	}},
	"Villager-Male-Move": {[]Frame{
		{"figure.png", 0, 48, 24, 48, 0, -16, 150 * time.Millisecond},
		{"figure.png", 24, 48, 24, 48, 0, -16, 150 * time.Millisecond},
		{"figure.png", 48, 48, 24, 48, 0, -16, 150 * time.Millisecond},
	}},
	"Giant": {[]Frame{
		{"giant.png", 0, 0, 48, 96, 0, -32, 1 * time.Second},
	}},
	"Wolf": {[]Frame{
		{"wolf.png", 0, 0, 58, 48, 0, -4, 1 * time.Second},
	}},
}
