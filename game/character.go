package game

//go:generate stringer -output=./character_string.go -type=CharacterSex,CharacterPerformance

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
	Profession string
	Sex        CharacterSex

	// InherantPreparation is the preparation value that this Character has as a
	// base, before values from their profession and equipped weapon are
	// applied.
	InherantPreparation int

	// InherantActionPoints is the base ActionPoints that this character has
	// before values from their profession and equipped weapon are applied.
	InherantActionPoints int

	CurrentHealth int
	BaseHealth    int

	Level                int
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
