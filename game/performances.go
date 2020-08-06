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
	Participant Sprite `json:"participant"`
	FaceX       int    `json:"faceX"`
	FaceY       int    `json:"faceY"`
	FeetX       int    `json:"feetX"`
	FeetY       int    `json:"feetY"`
}

// BigIcon creates a new Sprite to use as a 52 by 52 portrait.
func (a *Appearance) BigIcon() Sprite {
	return Sprite{
		Texture: a.Participant.Texture,
		X:       a.Participant.X + a.FaceX - 26,
		Y:       a.Participant.Y + a.FaceY - 26,
		W:       52, H: 52,
	}
}

// SmallIcon creates a new Sprite to use as a 26 by 26 portrait.
func (a *Appearance) SmallIcon() Sprite {
	return Sprite{
		Texture: a.Participant.Texture,
		X:       a.Participant.X + a.FaceX - 13,
		Y:       a.Participant.Y + a.FaceY - 13,
		W:       26, H: 26,
	}
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
