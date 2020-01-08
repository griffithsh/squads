package game

//go:generate stringer -output=./character_string.go -type=CharacterProfession,CharacterSex,CharacterPerformance

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
	Name      string
	SmallIcon Sprite // (26x26)
	BigIcon   Sprite // (52x52)

	// Intrinsic to the Character
	Profession CharacterProfession
	Sex        CharacterSex

	PreparationThreshold int // Preparation required to take a turn
	ActionPoints         int
	CurrentHealth        int
	MaxHealth            int

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
