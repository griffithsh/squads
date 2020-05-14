package game

import (
	"encoding/json"
	"time"

	"github.com/griffithsh/squads/geom"
)

type Frame struct {
	DurationMs int
	Sprite     Sprite
}

type PerformancesForDirection struct {
	N  []Frame
	S  []Frame
	NE []Frame
	NW []Frame
	SE []Frame
	SW []Frame
}

// ForDirection gets the value of []Frames for a direction.
func (pfd *PerformancesForDirection) ForDirection(dir geom.DirectionType) []Frame {
	switch dir {
	case geom.N:
		return pfd.N
	case geom.S:
		return pfd.S
	case geom.NE:
		return pfd.NE
	case geom.NW:
		return pfd.NW
	case geom.SE:
		return pfd.SE
	case geom.SW:
		return pfd.SW
	}
	return []Frame{}
}

// PerformanceSet is a set of performances for a profession and sex.
type PerformanceSet struct {
	Name    string
	Sexes   []CharacterSex `json:"-"`
	Idle    PerformancesForDirection
	Move    PerformancesForDirection
	Attack  PerformancesForDirection
	Spell   []Frame
	Death   []Frame
	Rise    []Frame
	Victory []Frame

	AttackApexMs int
	SpellApexMs  int
}

// UnmarshalJSON exists to extract the string based Sex values into enum values.
func (ps *PerformanceSet) UnmarshalJSON(data []byte) error {
	type alias PerformanceSet

	var a struct {
		alias
		Sexes []string `json:"sexes"`
	}

	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}

	*ps = PerformanceSet(a.alias)

	for _, sex := range a.Sexes {
		switch sex {
		case "XX":
			ps.Sexes = append(ps.Sexes, Female)
		case "XY":
			ps.Sexes = append(ps.Sexes, Male)
		}
	}

	return nil
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
