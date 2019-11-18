package overworld

//go:generate stringer -type=State

// State enumerates the States that an Overworld could be in.
type State int

const (
	// Uninitialised is before the overworld Manager has been Begun().
	Uninitialised State = iota

	// AnimatingState is when a squad is moving.
	AnimatingState

	// AwaitingInputState is when the Overworld is waiting for the local, human player to move their Squad.
	AwaitingInputState

	// FadingIn is when the combat is first starting, or returning from a menu,
	// and the curtain that obscures the scene change is disappearing.
	FadingIn
	//FadingOut is when the combat is going to another scene, and the curtain
	//that obscures the scene change is appearing.
	FadingOut
)
