package game

//go:generate stringer -type=ActorSize

// ActorSize enumerates directions.
type ActorSize int

// ActorSize represent the 6 directions that a Mover can face.
const (
	SMALL ActorSize = iota
	MEDIUM
	LARGE
)

// Actor is a component that can be commanded to do things. Or maybe it's just an animator?
type Actor struct {
	// M, N int
	// Busy bool
	Size ActorSize
}

// Type of this Component.
func (a *Actor) Type() string {
	return "Actor"
}
