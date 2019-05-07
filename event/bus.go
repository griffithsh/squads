package event

// Bus intemediates publishers and subscribers via an event interface.
type Bus struct {
	subs map[Type][]Subscriber
}

// Subscribe to all events of a type.
// FIXME: How could Unsubscribe be implemented?
func (b *Bus) Subscribe(t Type, f Subscriber) {
	b.ensure()
	b.subs[t] = append(b.subs[t], f)
}

// Publish an event to all Subscribers to that type.
func (b *Bus) Publish(t Typer) {
	b.ensure()
	subscriptions, ok := b.subs[t.Type()]

	if !ok {
		return
	}

	for _, f := range subscriptions {
		f(t)
	}
}

func (b *Bus) ensure() {
	if b.subs == nil {
		b.subs = map[Type][]Subscriber{}
	}
}
