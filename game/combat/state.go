package combat

//go:generate stringer -type=State

// State enumerates the States that a Combat could be in.
type State int

const (
	// Uninitialised is that default state of the combat.Manager.
	Uninitialised State = iota
	// AwaitingInputState is when the combat is waiting for the local, human player to make a move.
	AwaitingInputState
	// SelectingTargetState is when the local player is picking a hex to use a skill on.
	SelectingTargetState
	// ConfirmingSelectedTargetState is when the player has clicked or tapped on
	// a hex, and must tap on it again to use the selected skill on that hex.
	ConfirmingSelectedTargetState
	// ExecutingState is when a move or action is being played out by a character.
	ExecutingState
	// ThinkingState is when an AI-controller player is waiting to get command.
	ThinkingState
	// PreparingState is when no characters is prepared enough to make a move.
	PreparingState
	// Celebration occurs when there is only one team left, and they are
	// celebrating their victory.
	Celebration
	// FadingIn is when the combat is first starting, or returning from a menu,
	// and the curtain that obscures the scene change is disappearing.
	FadingIn
	//FadingOut is when the combat is going to another scene, and the curtain
	//that obscures the scene change is appearing.
	FadingOut
)

// Value allows simple States without context to implement the StateContext
// interface.
func (s State) Value() State {
	return s
}

// StateContext is something that has a State value, but could also include
// additional context for that state.
type StateContext interface {
	Value() State
}
