package data

import (
	"fmt"

	"github.com/griffithsh/squads/game"
)

var internalPerformances = map[PerformanceKey]*game.PerformanceSet{
	// Any debug professions that are staticly compiled into the binary should
	// have their performance animations configured here.
}

// PerformanceKey uniquely identifies a PerformanceSet for a Sex.
type PerformanceKey struct {
	Sex        game.CharacterSex
	Profession string
}

// Performances for Profession.
func (a *Archive) Performances(profession string, sex game.CharacterSex) *game.PerformanceSet {
	set, found := a.performances[PerformanceKey{Profession: profession, Sex: sex}]
	if !found {
		fmt.Printf("no performances for %s %s\n", sex.String(), profession)
		p := defaultPerformanceSet()
		a.performances[PerformanceKey{Profession: profession, Sex: sex}] = p
		return p
	}
	return set
}

func defaultPerformanceSet() *game.PerformanceSet {
	s := game.Sprite{
		Texture: "tranquility-plus-39-palette.png",
		W:       8, H: 8,
	}
	return &game.PerformanceSet{
		Sexes: []game.CharacterSex{game.Male, game.Female},
		Idle: game.PerformancesForDirection{
			N:  []game.Frame{{DurationMs: 10, Sprite: s}},
			S:  []game.Frame{{DurationMs: 10, Sprite: s}},
			NE: []game.Frame{{DurationMs: 10, Sprite: s}},
			NW: []game.Frame{{DurationMs: 10, Sprite: s}},
			SE: []game.Frame{{DurationMs: 10, Sprite: s}},
			SW: []game.Frame{{DurationMs: 10, Sprite: s}},
		},
		Move: game.PerformancesForDirection{
			N:  []game.Frame{{DurationMs: 10, Sprite: s}},
			S:  []game.Frame{{DurationMs: 10, Sprite: s}},
			NE: []game.Frame{{DurationMs: 10, Sprite: s}},
			NW: []game.Frame{{DurationMs: 10, Sprite: s}},
			SE: []game.Frame{{DurationMs: 10, Sprite: s}},
			SW: []game.Frame{{DurationMs: 10, Sprite: s}},
		},
		Attack: game.PerformancesForDirection{
			N:  []game.Frame{{DurationMs: 10, Sprite: s}},
			S:  []game.Frame{{DurationMs: 10, Sprite: s}},
			NE: []game.Frame{{DurationMs: 10, Sprite: s}},
			NW: []game.Frame{{DurationMs: 10, Sprite: s}},
			SE: []game.Frame{{DurationMs: 10, Sprite: s}},
			SW: []game.Frame{{DurationMs: 10, Sprite: s}},
		},
		Spell:   []game.Frame{{DurationMs: 10, Sprite: s}},
		Death:   []game.Frame{{DurationMs: 10, Sprite: s}},
		Victory: []game.Frame{{DurationMs: 10, Sprite: s}},
	}
}
