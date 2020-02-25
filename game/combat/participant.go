package combat

import (
	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/game"
)

// EngagementStatus represents how fit for combat a Character is.
type EngagementStatus int

const (
	// Alive Characters can fight normally
	Alive EngagementStatus = iota

	// KnockedDown Characters cannot perform any actions, do not prepare, and
	// are not affected by healing. Their body may be reanimated by a
	// necromancer, or revived by a skill with the resurrection property.
	KnockedDown

	// Escaped Characters have left combat, and cannot be affected by anything
	// that happens on the field.
	Escaped
)

//go:generate stringer -type=EngagementStatus

// CurMax represents a value which has both a Current and Maximum value.
type CurMax struct {
	Cur int
	Max int
}

// Participant is a transient aggregation of the stats of a Character for the purposes
// of combat. Participants are created at the beginning of combat and are destroyed at
// the end of combat.
type Participant struct {
	Character ecs.Entity
	// Hexes occupied? Do we merge with Obstacle?

	Name      string
	Level     uint
	SmallIcon game.Sprite // (26x26)
	BigIcon   game.Sprite // (52x52)

	Profession game.CharacterProfession
	Sex        game.CharacterSex

	PreparationThreshold CurMax
	ActionPoints         CurMax
	Health               CurMax

	Strength     int
	Agility      int
	Intelligence int
	Vitality     int

	Status EngagementStatus // Alive, Knocked down, or Escaped

	Disambiguator float64

	// EquippedWeaponClass should not change while in combat.
	EquippedWeaponClass game.ItemClass
	ItemStats           map[game.Modifier]float64
}

// Type of this Component.
func (*Participant) Type() string {
	return "Participant"
}
