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
	Name    string                   `json:"name"`
	Sexes   []CharacterSex           `json:"-"`
	Idle    PerformancesForDirection `json:"idle"`
	Move    PerformancesForDirection `json:"move"`
	Attack  PerformancesForDirection `json:"attack"`
	Spell   []Frame                  `json:"spell"`
	Death   []Frame                  `json:"death"`
	Rise    []Frame                  `json:"rise"`
	Victory []Frame                  `json:"victory"`

	MoveSpeed  time.Duration `json:"-"`
	AttackApex time.Duration `json:"-"`
	SpellApex  time.Duration `json:"-"`
}

// UnmarshalJSON exists to extract the string based Sex values into enum values.
func (ps *PerformanceSet) UnmarshalJSON(data []byte) error {
	type alias PerformanceSet

	var a struct {
		alias
		Sexes        []string `json:"sexes"`
		MoveSpeedMs  int      `json:"moveSpeed"`
		AttackApexMs int      `json:"attackApexMs"`
		SpellApexMs  int      `json:"spellApexMs"`
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

	ps.MoveSpeed = time.Duration(a.MoveSpeedMs) * time.Millisecond
	ps.AttackApex = time.Duration(a.AttackApexMs) * time.Millisecond
	ps.SpellApex = time.Duration(a.SpellApexMs) * time.Millisecond

	return nil
}

func (ps PerformanceSet) MarshalJSON() ([]byte, error) {
	type v PerformanceSet
	alias := struct {
		v
		Sexes        []string `json:"sexes"`
		MoveSpeed    int      `json:"moveSpeed"`
		AttackApexMs int      `json:"attackApexMs"`
		SpellApexMs  int      `json:"spellApexMs"`
	}{
		v:            v(ps),
		MoveSpeed:    int(ps.MoveSpeed / time.Millisecond),
		AttackApexMs: int(ps.AttackApex / time.Millisecond),
		SpellApexMs:  int(ps.SpellApex / time.Millisecond),
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
