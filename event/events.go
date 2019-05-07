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
