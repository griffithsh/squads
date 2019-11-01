package ui

import (
	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
)

// Interactive tags Components that can be interacted with in some way.
type Interactive struct {
	Trigger func()
}

// Type of this Component.
func (*Interactive) Type() string {
	return "Interactive"
}

// Interactive2 tags Components that can be interacted with in some way.
type Interactive2 struct {
	W, H    float64
	Trigger func()
}

// Type of this Component.
func (*Interactive2) Type() string {
	return "Interactive2"
}

// InteractiveSystem pumps events through Interactives.
type InteractiveSystem struct {
	mgr *ecs.World
}

// Handle an Interact event.
func (is *InteractiveSystem) Handle(ev *Interact) {
	for _, e := range is.mgr.Get([]string{"Interactive2", "Position"}) {
		interactive := is.mgr.Component(e, "Interactive2").(*Interactive2)
		w, h := interactive.W, interactive.H
		position := is.mgr.Component(e, "Position").(*game.Position)
		if is.mgr.Component(e, "Scale") != nil {
			scale := is.mgr.Component(e, "Scale").(*game.Scale)
			w *= scale.X
			h *= scale.Y
		}

		if ev.Absolute != position.Absolute {
			continue
		}

		if ev.X < position.Center.X-w/2 {
			continue
		}

		if ev.X > position.Center.X+w/2 {
			continue
		}

		if ev.Y < position.Center.Y-h/2 {
			continue
		}

		if ev.Y > position.Center.Y+h/2 {
			continue
		}

		interactive.Trigger()
	}
}

// NewInteractiveSystem creates a new InteractiveSystem.
func NewInteractiveSystem(mgr *ecs.World, bus *event.Bus) *InteractiveSystem {
	is := InteractiveSystem{
		mgr: mgr,
	}
	bus.Subscribe(Interact{}.Type(), func(t event.Typer) {
		ev := t.(*Interact)
		is.Handle(ev)
	})
	return &is
}
