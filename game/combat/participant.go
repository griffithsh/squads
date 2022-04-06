package combat

import (
	"fmt"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/game/item"
	"github.com/griffithsh/squads/skill"
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

type injury struct {
	// Value is the duration of the injury in preparation for all Injuries
	// except Poison. For the poison injury it represents intensity.
	Value int

	// Remainder is (for Bleeding and Burning) the preparation that has been
	// removed from Value due to preparation being applied, that has not been
	// rounded up into a whole number unit of damage.
	Remainder int
}

// Participant is a transient aggregation of the stats of a Character for the purposes
// of combat. Participants are created at the beginning of combat and are destroyed at
// the end of combat.
type Participant struct {
	Character ecs.Entity
	// Hexes occupied? Do we merge with Obstacle?

	Level int

	Name               string
	Hair               string
	Skin               string
	SmallPortraitBG    game.Sprite
	BigPortraitBG      game.Sprite
	SmallPortraitFrame game.Sprite
	BigPortraitFrame   game.Sprite
	SmallIcon          game.Sprite // (26x26)
	BigIcon            game.Sprite // (52x52)

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

	// Injuries stores the current Injuries the Participant is suffering from.
	Injuries map[skill.InjuryType]*injury

	Disambiguator float64

	Masteries map[game.Mastery]int

	ItemStats map[item.Modifier]float64

	EquippedWeaponClass item.Class
	// Skills should not change while in combat.
	Skills []skill.ID
}

// Type of this Component.
func (*Participant) Type() string {
	return "Participant"
}

// baseDamage of a combat participant.
func (p *Participant) baseDamage() (baseMin, baseMax float64) {
	mult, ok := p.ItemStats[item.BaseDamageModifier]
	if !ok {
		mult = 1.0
	}
	// If these values are not present, then default zero is appropriate.
	min := p.ItemStats[item.BaseMinDamageModifier]
	max := p.ItemStats[item.BaseMaxDamageModifier]
	min *= mult
	max *= mult

	var strBonus, agiBonus float64
	switch p.EquippedWeaponClass {
	case item.UnarmedClass:
		strBonus = 0.75
	case item.SwordClass:
		strBonus = 0.65
		agiBonus = 0.35
	case item.AxeClass:
		strBonus = 0.85
		agiBonus = 0.15
	case item.ClubClass:
		strBonus = 1
	case item.DaggerClass:
		agiBonus = 1
	case item.SlingClass:
		strBonus = 0.1
		agiBonus = 0.9
	case item.BowClass:
		strBonus = 0.2
		agiBonus = 0.8
	case item.SpearClass:
		strBonus = 0.75
		agiBonus = 0.75
	case item.PolearmClass:
		strBonus = 1.0
		agiBonus = 0.5
	case item.StaffClass:
		strBonus = .6
		agiBonus = .6
	case item.WandClass:
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
