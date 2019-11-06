package game

import "math/rand"

//go:generate stringer -type=TeamControl

// TeamControl marks the source of control.
type TeamControl int

const (
	// LocalControl means that the player is controlling this team.
	LocalControl TeamControl = iota

	// ComputerControl means that an AI script is controlling this team.
	ComputerControl

	// NetworkControl ?
)

// Team is a Component that represents a grouping of Actors that work together,
// and share a source of control.
type Team struct {
	ID      int64
	Control TeamControl
}

// NewTeam creates a new team.
func NewTeam( /*Control*/ ) *Team {
	return &Team{
		ID: rand.Int63(),
	}
}

// Type of this Component.
func (*Team) Type() string {
	return "Team"
}
