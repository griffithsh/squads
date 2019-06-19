package event

// Type enumerates directions.
type Type int

// Types represent things that have happened.
const (
	AwaitingPlayerInputType Type = iota
	CombatBegunType
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
