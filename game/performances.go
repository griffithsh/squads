package game

import (
	"encoding/json"
	"time"

	"github.com/griffithsh/squads/geom"
)

type Frame struct {
	DurationMs int    `json:"durationMs"`
	Sprite     Sprite `json:"sprite"`
}

type PerformancesForDirection struct {
	N  []Frame `json:"n"`
	S  []Frame `json:"s"`
	NE []Frame `json:"ne"`
	NW []Frame `json:"nw"`
	SE []Frame `json:"se"`
	SW []Frame `json:"sw"`
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

// PerformanceSet is a set of performances for a profession and sex(es).
type PerformanceSet struct {
	Name  string                   `json:"name"`
	Sexes []CharacterSex           `json:"-"`
	Idle  PerformancesForDirection `json:"idle"`
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

// MarshalJSON implements the json marshaling interface.
func (ps PerformanceSet) MarshalJSON() ([]byte, error) {
	type v PerformanceSet
	alias := struct {
		v
		Sexes []string `json:"sexes"`
	}{
		v: v(ps),
	}
	for _, sex := range ps.Sexes {
		if sex == Male {
			alias.Sexes = append(alias.Sexes, "XX")
		} else if sex == Female {
			alias.Sexes = append(alias.Sexes, "XY")
		}
	}

	return json.Marshal(&alias)
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
