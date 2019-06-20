package event

import (
	"testing"
)

func TestBus(t *testing.T) {
	bus := &Bus{}

	subs := 0
	var ty Type
	unsub1x := bus.Subscribe(ty, func(Typer) {
		subs++
	})
	unsub10x := bus.Subscribe(ty, func(Typer) {
		subs += 10
	})

	bus.Publish(ty)

	if subs != 11 {
		t.Errorf("want 11 got %d", subs)
	}

	unsub1x()

	bus.Publish(ty)

	if subs != 21 {
		t.Errorf("want 21 got %d", subs)
	}
	unsub10x()

	bus.Publish(ty)

	if subs != 21 {
		t.Errorf("want 21 got %d", subs)
	}
}
