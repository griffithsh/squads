package data

import (
	"fmt"

	"github.com/griffithsh/squads/game"
)

var internalAppearances = map[AppearanceKey]*game.Appearance{
	// Any debug professions that are staticly compiled into the binary should
	// have their appearances configured here.
}

// AppearanceKey uniquely identifies a Profession and Sex.
type AppearanceKey struct {
	Sex        game.CharacterSex
	Profession string
	Hair       string
	Skin       string
}

// Appearance retrieves an appropriate Appearance object to use for a character in combat.
func (a *Archive) Appearance(profession string, sex game.CharacterSex, hair string, skin string) *game.Appearance {
	// FIXME: implementation
	key := AppearanceKey{
		Profession: profession,
		Sex:        sex,
		Hair:       hair,
		Skin:       skin,
	}
	v, ok := a.appearances[key]
	if !ok {
		panic(fmt.Sprintf("no appearance for %v", key))
	}
	return v
}

// HairVariations returns the list of available hair colors.
func (a *Archive) HairVariations() []string {
	return a.hairColors
}

// SkinVariations returns the list of available skin colors.
func (a *Archive) SkinVariations() []string {
	return a.skinColors
}
