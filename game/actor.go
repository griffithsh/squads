package game

//go:generate stringer -output=./actor_string.go -type=ActorSize,ActorProfession,ActorSex,ActorPerformance

// ActorSize enumerates Sizes for Actors.
type ActorSize int

// ActorSize represents the sizes that an Actor can be. Small Actrors take only
// one Hex, Medium, take 4, and Large take 7.
const (
	SMALL ActorSize = iota
	MEDIUM
	LARGE
)

type ActorProfession int

const (
	Villager ActorProfession = iota
	Wolf
	Giant
	Skeleton
)

type ActorSex int

const (
	Male ActorSex = iota
	Female
)

type ActorPerformance int

const (
	PerformIdle ActorPerformance = iota
	PerformMove
	PerformSkill1
	PerformSkill2
	PerformSkill3
	PerformHurt
	PerformDying
	PerformVictory
)

// Actor is a component that can be commanded to do things. Or maybe it's just an animator?
type Actor struct {
	// Things that don't affect gameplay.
	Name      string
	SmallIcon Sprite // (26x26)
	BigIcon   Sprite // (52x52)

	// Intrinsic to the Actor
	Size        ActorSize
	Profession  ActorProfession
	Sex         ActorSex
	Performance ActorPerformance

	PreparationThreshold int // Preparation required to take a turn
	ActionPoints         int
}

// Type of this Component.
func (*Actor) Type() string {
	return "Actor"
}
