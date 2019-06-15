package event

import "github.com/griffithsh/squads/ecs"

// ActorMovementConcluded occurs when an actor has finished their movement.
type ActorMovementConcluded struct {
	Entity ecs.Entity
}

// Type of the Event.
func (ActorMovementConcluded) Type() Type {
	return MovementConcluded
}

// CombatBegun occurs when an actor has finished their movement.
type CombatBegun struct{}

// Type of the Event.
func (CombatBegun) Type() Type {
	return CombatBegunType
}

// EndTurnRequested occurs when an actor has finished their movement.
type EndTurnRequested struct{}

// Type of the Event.
func (EndTurnRequested) Type() Type {
	return EndTurnRequestedType
}
