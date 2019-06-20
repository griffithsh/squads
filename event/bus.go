package event

import "math/rand"

// Bus intemediates publishers and subscribers via an event interface.
type Bus struct {
	subs map[Type]map[int64]Subscriber
}

// ensure that the internal map is initialised.
func (b *Bus) ensure() {
	if b.subs == nil {
		b.subs = map[Type]map[int64]Subscriber{}
	}
}

func (b *Bus) unsubscribeKey() int64 {
	return rand.Int63()
}

// Subscribe to all events of a type. Returns an Unsubscribe function.
func (b *Bus) Subscribe(t Type, f Subscriber) func() {
	b.ensure()
	m, ok := b.subs[t]
	if !ok {
		b.subs[t] = map[int64]Subscriber{}
		m = b.subs[t]
	}
	var key int64
	for {
		_, ok := m[key]
		if ok {
			continue
		}
		key = b.unsubscribeKey()
		break
	}
	m[key] = f
	return func() {
		delete(m, key)
	}
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
