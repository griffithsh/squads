package event

// Type enumerates directions.
type Type int

func (ty Type) Type() Type {
	return ty
}

// Types represent things that have happened.
const (
	AwaitingPlayerInputType Type = iota
	CombatBegunType
	CombatStateTransitionType
	CombatStatModifiedType
	EndTurnRequestedType
	MovementConcluded
)

// Typer is an awkward thing that represents anything that provides its type.
type Typer interface {
	Type() Type
}

// Subscriber is anything that is subscribed to an event type.
type Subscriber func(Typer)
