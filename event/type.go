package event

// Type enumerates directions.
type Type string

// Type implements the Typer interface, so that simple events without data can be Published.
func (ty Type) Type() Type {
	return ty
}

// Typer is an awkward thing that represents anything that provides its type.
type Typer interface {
	Type() Type
}

// Subscriber is anything that is subscribed to an event type.
type Subscriber func(Typer)
