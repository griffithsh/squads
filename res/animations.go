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
		{"giant.png", 0, 0, 48, 96, 0, -16, 1 * time.Second},
	}},
	"Wolf": {[]Frame{
		{"wolf.png", 0, 0, 58, 48, 0, -4, 1 * time.Second},
	}},
	"Skeleton-Idle-S": {[]Frame{
		{"skeleton.png", 0, 0, 24, 48, 0, -16, 800 * time.Millisecond},
		{"skeleton.png", 24, 0, 24, 48, 0, -16, 450 * time.Millisecond},
		{"skeleton.png", 48, 0, 24, 48, 0, -16, 800 * time.Millisecond},
		{"skeleton.png", 72, 0, 24, 48, 0, -16, 450 * time.Millisecond},
	}},
	"Skeleton-Idle-N": {[]Frame{
		{"skeleton.png", 96, 0, 24, 48, 0, -16, 800 * time.Millisecond},
		{"skeleton.png", 120, 0, 24, 48, 0, -16, 450 * time.Millisecond},
		{"skeleton.png", 144, 0, 24, 48, 0, -16, 800 * time.Millisecond},
		{"skeleton.png", 168, 0, 24, 48, 0, -16, 450 * time.Millisecond},
	}},
	"Skeleton-Idle-SE": {[]Frame{
		{"skeleton.png", 192, 0, 24, 48, 0, -16, 800 * time.Millisecond},
		{"skeleton.png", 216, 0, 24, 48, 0, -16, 450 * time.Millisecond},
		{"skeleton.png", 240, 0, 24, 48, 0, -16, 800 * time.Millisecond},
		{"skeleton.png", 264, 0, 24, 48, 0, -16, 450 * time.Millisecond},
	}},
	"Skeleton-Idle-SW": {[]Frame{
		{"skeleton.png", 288, 0, 24, 48, 0, -16, 800 * time.Millisecond},
		{"skeleton.png", 312, 0, 24, 48, 0, -16, 450 * time.Millisecond},
		{"skeleton.png", 336, 0, 24, 48, 0, -16, 800 * time.Millisecond},
		{"skeleton.png", 360, 0, 24, 48, 0, -16, 450 * time.Millisecond},
	}},
	"Skeleton-Idle-NW": {[]Frame{
		{"skeleton.png", 192, 48, 24, 48, 0, -16, 800 * time.Millisecond},
		{"skeleton.png", 216, 48, 24, 48, 0, -16, 450 * time.Millisecond},
		{"skeleton.png", 240, 48, 24, 48, 0, -16, 800 * time.Millisecond},
		{"skeleton.png", 264, 48, 24, 48, 0, -16, 450 * time.Millisecond},
	}},
	"Skeleton-Idle-NE": {[]Frame{
		{"skeleton.png", 288, 48, 24, 48, 0, -16, 800 * time.Millisecond},
		{"skeleton.png", 312, 48, 24, 48, 0, -16, 450 * time.Millisecond},
		{"skeleton.png", 336, 48, 24, 48, 0, -16, 800 * time.Millisecond},
		{"skeleton.png", 360, 48, 24, 48, 0, -16, 450 * time.Millisecond},
	}},
	"Skeleton-Move-S": {[]Frame{
		{"skeleton.png", 0, 240, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 24, 240, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 48, 240, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 72, 240, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 96, 240, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 120, 240, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 144, 240, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 168, 240, 24, 48, 0, -16, 90 * time.Millisecond},
	}},
	"Skeleton-Move-N": {[]Frame{
		{"skeleton.png", 0, 288, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 24, 288, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 48, 288, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 72, 288, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 96, 288, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 120, 288, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 144, 288, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 168, 288, 24, 48, 0, -16, 90 * time.Millisecond},
	}},
	"Skeleton-Move-SE": {[]Frame{
		{"skeleton.png", 0, 48, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 24, 48, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 48, 48, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 72, 48, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 96, 48, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 120, 48, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 144, 48, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 168, 48, 24, 48, 0, -16, 90 * time.Millisecond},
	}},
	"Skeleton-Move-SW": {[]Frame{
		{"skeleton.png", 0, 96, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 24, 96, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 48, 96, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 72, 96, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 96, 96, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 120, 96, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 144, 96, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 168, 96, 24, 48, 0, -16, 90 * time.Millisecond},
	}},
	"Skeleton-Move-NW": {[]Frame{
		{"skeleton.png", 0, 144, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 24, 144, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 48, 144, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 72, 144, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 96, 144, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 120, 144, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 144, 144, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 168, 144, 24, 48, 0, -16, 90 * time.Millisecond},
	}},
	"Skeleton-Move-NE": {[]Frame{
		{"skeleton.png", 0, 192, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 24, 192, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 48, 192, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 72, 192, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 96, 192, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 120, 192, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 144, 192, 24, 48, 0, -16, 90 * time.Millisecond},
		{"skeleton.png", 168, 192, 24, 48, 0, -16, 90 * time.Millisecond},
	}},
}
