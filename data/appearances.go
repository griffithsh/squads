package data

import (
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
	// Hair color
	// Skin color
}

// Appearance retrieves an appropriate Appearance object to use for a character in combat.
func (a *Archive) Appearance(profession string, sex game.CharacterSex) *game.Appearance {
	// FIXME: implementation
	return &game.Appearance{
		Participant: game.Sprite{
			Texture: "figure.png",

			X: 0, Y: 0,
			W: 24, H: 48,
		},
		Icon: game.Sprite{},
	}
}
