package ui

import (
	"sort"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
)

// Interactive tags Components that can be interacted with in some way. Trigger
// is a function that is provided the x,y coordinate of the click either in
// world or screen coordinates, depending on the Position of the Interactive.
type Interactive struct {
	W, H    float64
	Trigger func(x, y float64)
}

// Type of this Component.
func (*Interactive) Type() string {
	return "Interactive"
}

// InteractiveSystem pumps events through Interactives.
type InteractiveSystem struct {
	mgr *ecs.World
}

// Handle an Interact event.
func (is *InteractiveSystem) Handle(ev *Interact) {
	type tmp struct {
		e           ecs.Entity
		interactive *Interactive
		position    *game.Position
		scale       *game.Scale
	}
	tmps := []tmp{}
	for _, e := range is.mgr.Get([]string{"Interactive", "Position"}) {
		t := tmp{
			e:           e,
			interactive: is.mgr.Component(e, "Interactive").(*Interactive),
			position:    is.mgr.Component(e, "Position").(*game.Position),
		}
		if is.mgr.Component(e, "Scale") != nil {
			t.scale = is.mgr.Component(e, "Scale").(*game.Scale)
		}
		tmps = append(tmps, t)
	}
	sort.Slice(tmps, func(i, j int) bool {
		return tmps[i].position.Layer > tmps[j].position.Layer
	})

	for _, t := range tmps {
		interactive := t.interactive
		w, h := interactive.W, interactive.H
		position := t.position
		if t.scale != nil {
			w *= t.scale.X
			h *= t.scale.Y
		}

		// Use the right coordinate framework (world or screen) for this
		// Interactive.
		x, y := ev.X, ev.Y
		if position.Absolute {
			x, y = ev.AbsoluteX, ev.AbsoluteY
		}

		if x < position.Center.X-w/2 {
			continue
		}
		if x > position.Center.X+w/2 {
			continue
		}

		if y < position.Center.Y-h/2 {
			continue
		}
		if y > position.Center.Y+h/2 {
			continue
		}

		// Only one Interactive should handle this, controlled by the Layer of
		// their position.
		interactive.Trigger(x, y)
		return
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
