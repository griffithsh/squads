package game

//go:generate stringer -type=ActorSize

// ActorSize enumerates Sizes for Actors.
type ActorSize int

// ActorSize represents the sizes that an Actor can be. Small Actrors take only
// one Hex, Medium, take 4, and Large take 7.
const (
	SMALL ActorSize = iota
	MEDIUM
	LARGE
)

// Actor is a component that can be commanded to do things. Or maybe it's just an animator?
type Actor struct {
	// Things that don't affect gameplay.
	Name      string
	SmallIcon Sprite // (26x26)
	BigIcon   Sprite // (52x52)

	// Intrinsic to the Actor
	Size ActorSize
	// Sex

	PreparationThreshold int // Preparation required to take a turn
	ActionPoints         int
}

// Type of this Component.
func (*Actor) Type() string {
	return "Actor"
}
