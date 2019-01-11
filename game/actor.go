package game

// Actor is a component that can be commanded to do things. Or maybe it's just an animator?
type Actor struct {
	M, N int
	Busy bool
}

// Type of this Component.
func (a *Actor) Type() string {
	return "Actor"
}
