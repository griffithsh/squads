package game

//go:generate stringer -output=./character_string.go -type=CharacterSize,CharacterProfession,CharacterSex,CharacterPerformance

// CharacterSize enumerates Sizes for Actors.
type CharacterSize int

// CharacterSize represents the sizes that an Actor can be. Small Actrors take only
// one Hex, Medium, take 4, and Large take 7.
const (
	SMALL CharacterSize = iota
	MEDIUM
	LARGE
)

type CharacterProfession int

const (
	Villager CharacterProfession = iota
	Wolf
	Giant
	Skeleton
)

type CharacterSex int

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
	Name string
	// SmallIcon Sprite // (26x26)
	// BigIcon   Sprite // (52x52)

	// Intrinsic to the Character
	Size       CharacterSize
	Profession CharacterProfession
	Sex        CharacterSex

	PreparationThreshold int // Preparation required to take a turn
	ActionPoints         int

	Level                uint
	StrengthPerLevel     float64
	DexterityPerLevel    float64
	IntelligencePerLevel float64
	VitalityPerLevel     float64
}

// Type of this Component.
func (*Character) Type() string {
	return "Character"
}
