package game

import (
	"encoding/json"
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

type PerformanceSet struct {
	Sexes   []CharacterSex `json:"-"`
	Idle    PerformancesForDirection
	Move    PerformancesForDirection
	Attack  PerformancesForDirection
	Spell   []Frame
	Death   []Frame
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
