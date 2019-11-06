package overworld

//go:generate stringer -type=State

// State enumerates the States that an Overworld could be in.
type State int

const (
	// AnimatingState is when a squad is moving.
	AnimatingState State = iota

	// AwaitingInputState is when the Overworld is waiting for the local, human player to move their Squad.
	AwaitingInputState

	// Uninitialised is before the overworld Manager has been Begun().
	Uninitialised
)
