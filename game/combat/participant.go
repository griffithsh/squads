package combat

import (
	"fmt"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/game"
)

// EngagementStatus represents how fit for combat a Character is.
type EngagementStatus int

const (
	// Alive Characters can fight normally
	Alive EngagementStatus = iota

	// KnockedDown is the status of a Participant that lost all its Health
	// Points. KnockedDown Participants cannot perform any actions, do not
	// prepare, and are not affected by healing, but are still obstacles to
	// movement. Their body may be reanimated by a necromancer, or revived by a
	// skill with the resurrection property.
	KnockedDown

	// Defiled means that the Participant is permanently dead, and cannot be
	// revived or resurrected; Their corpse is gone, and no longer poses an
	// obstacle to movement.
	Defiled

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
	Level     int
	SmallIcon game.Sprite // (26x26)
	BigIcon   game.Sprite // (52x52)

	Profession string
	Sex        game.CharacterSex

	PreparationThreshold CurMax
	ActionPoints         CurMax
	BaseHealth           int
	CurrentHealth        int

	Strength     int
	Agility      int
	Intelligence int
	Vitality     int

	Status EngagementStatus // Alive, Knocked down, or Escaped

	Disambiguator float64

	Masteries map[game.Mastery]int

	// EquippedWeaponClass should not change while in combat.
	EquippedWeaponClass game.ItemClass
	ItemStats           map[game.Modifier]float64
}

// Type of this Component.
func (*Participant) Type() string {
	return "Participant"
}

// baseDamage of a combat participant.
func (p *Participant) baseDamage() (baseMin, baseMax float64) {
	mult, ok := p.ItemStats[game.BaseDamageModifier]
	if !ok {
		mult = 1.0
	}
	// If these values are not present, then default zero is appropriate.
	min, _ := p.ItemStats[game.BaseMinDamageModifier]
	max, _ := p.ItemStats[game.BaseMaxDamageModifier]
	min *= mult
	max *= mult

	var strBonus, agiBonus float64
	switch p.EquippedWeaponClass {
	case game.UnarmedClass:
		strBonus = 0.75
	case game.SwordClass:
		strBonus = 0.65
		agiBonus = 0.35
	case game.AxeClass:
		strBonus = 0.85
		agiBonus = 0.15
	case game.ClubClass:
		strBonus = 1
	case game.DaggerClass:
		agiBonus = 1
	case game.SlingClass:
		strBonus = 0.1
		agiBonus = 0.9
	case game.BowClass:
		strBonus = 0.2
		agiBonus = 0.8
	case game.SpearClass:
		strBonus = 0.75
		agiBonus = 0.75
	case game.PolearmClass:
		strBonus = 1.0
		agiBonus = 0.5
	case game.StaffClass:
		strBonus = .6
		agiBonus = .6
	case game.WandClass:
		strBonus = 0.1
		agiBonus = 0.1
	default:
		panic(fmt.Sprintf("unknown equipped weapon %v", p.EquippedWeaponClass))
	}

	weapBonus := float64(p.Strength)*strBonus + float64(p.Agility)*agiBonus

	min = min + min*weapBonus*0.15
	max = max + max*weapBonus*0.15
	return min, max
}

func (p *Participant) maxHealth() int {
	// TODO: Also include vitality and raw health from equipped items.
	return game.MaxHealth(p.BaseHealth, p.Vitality)
}
