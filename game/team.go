package game

import "math/rand"

//go:generate stringer -type=TeamControl

// TeamControl marks the source of control.
type TeamControl int

const (
	// LocalControl means that the player is controlling this team.
	LocalControl TeamControl = iota

	// NoControl is for things that do not move
	NoControl

	// ComputerControl means that an AI script is controlling this team.
	ComputerControl
)

// Team is a Component that represents a grouping of Characters that work
// together, and share a source of control.
type Team struct {
	ID                 int64
	Control            TeamControl
	PedestalAppearance int
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
