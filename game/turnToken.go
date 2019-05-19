package game

// TurnToken is a Component that is assigned to the Actor whose turn it is right
// now. There should only be one of these in existence at any one time.
type TurnToken struct{}

// Type of this Component.
func (*TurnToken) Type() string {
	return "TurnToken"
}
