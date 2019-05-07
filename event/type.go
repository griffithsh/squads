package event

// Type enumerates directions.
type Type int

// Types represent things that have happened.
const (
	MovementConcluded Type = iota
)

type Typer interface {
	Type() Type
}

type Subscriber func(Typer)
