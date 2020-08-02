package game

import (
	"time"
)

type Frame struct {
	DurationMs int    `json:"durationMs"`
	Sprite     Sprite `json:"sprite"`
}

// Appearance is how a Character appears in combat.
type Appearance struct {
	Participant Sprite
	Portrait    Sprite
}

// NewFrameAnimationFromFrames creates a new animation from a slice of Frames.
func NewFrameAnimationFromFrames(frames []Frame) *FrameAnimation {
	fa := FrameAnimation{}

	for _, frame := range frames {
		fa.Frames = append(fa.Frames, frame.Sprite)
		fa.Timings = append(fa.Timings, time.Duration(frame.DurationMs)*time.Millisecond)
	}
	return &fa
}
