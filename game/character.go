package game

//go:generate stringer -output=./character_string.go -type=CharacterProfession,CharacterSex,CharacterPerformance

type CharacterProfession int

const (
	Villager CharacterProfession = iota
	Wolf
	Giant
	Skeleton
)

// Preparation that this CharacterProfession provides.
func (p CharacterProfession) Preparation() int {
	switch p {
	case Villager:
		return 200
	case Wolf:
		return 400
	case Giant:
		return 800
	case Skeleton:
		return 900
	default:
		panic("unconfigured Preparation value for CharacterProfession " + p.String())
	}
}

// ActionPoints that this CharacterProfession provides.
func (p CharacterProfession) ActionPoints() int {
	switch p {
	case Villager:
		return 60
	case Wolf:
		return 40
	case Giant:
		return 40
	case Skeleton:
		return 40
	default:
		panic("unconfigured ActionPoints value for CharacterProfession " + p.String())
	}
}

// CharacterSex is not CharacterGender, NB.
type CharacterSex int

// CharacterSexes only have two values. They represent XY and XX chromosomes.
const (
	Male CharacterSex = iota
	Female
)

type CharacterPerformance int

const (
	PerformIdle CharacterPerformance = iota
	PerformMove
	PerformSkill1
	PerformSkill2
	PerformSkill3
	PerformHurt
	PerformDying
	PerformVictory
)

// Character represents a character at the run level, and persists until the run
// is over.
type Character struct {
	// Things that don't affect gameplay.
	Name      string
	SmallIcon Sprite // (26x26)
	BigIcon   Sprite // (52x52)

	// Intrinsic to the Character
	Profession CharacterProfession
	Sex        CharacterSex

	// InherantPreparation is the preparation value that this Character has as a
	// base, before values from their profession and equipped weapon are
	// applied.
	InherantPreparation int

	// InherantActionPoints is the base ActionPoints that this character has
	// before values from their profession and equipped weapon are applied.
	InherantActionPoints int

	CurrentHealth int
	MaxHealth     int

	Level                uint
	StrengthPerLevel     float64
	AgilityPerLevel      float64
	IntelligencePerLevel float64
	VitalityPerLevel     float64

	Disambiguator float64 // random number to order Characters when their Preparation etc collide.

	// Masteries indexed by the enum value.
	Masteries map[Mastery]int
}

// Type of this Component.
func (*Character) Type() string {
	return "Character"
}
