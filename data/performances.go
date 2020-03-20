package data

import (
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
		return &game.PerformanceSet{
			// TODO: define a default performance set.
		}
	}
	return set
}
