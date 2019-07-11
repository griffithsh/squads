package combat

import "math/rand"

// Team is a Component that represents a grouping of Actors that work together,
// and share a source of control.
type Team struct {
	ID int64
	// Control string // local-control|network-control|ai
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
