package combat

import (
	"time"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
)

// PerformanceSystem sets appropriate Animations for Actors in a Combat based on
// what's happening.
type PerformanceSystem struct {
	mgr *ecs.World
}

// NewPerformanceSystem creates a new PerformanceSystem.
func NewPerformanceSystem(mgr *ecs.World, bus *event.Bus) *PerformanceSystem {
	ps := PerformanceSystem{
		mgr: mgr,
	}

	bus.Subscribe(game.CombatActorMovementCommenced{}.Type(), ps.handleActorMoving)
	bus.Subscribe(game.CombatActorMovementConcluded{}.Type(), ps.handleActorStopped)

	return &ps
}

// Update the System.
func (ps *PerformanceSystem) Update(mgr *ecs.World, elapse time.Duration) {
	// For any actors without a sprite, apply their idling animation
	for _, e := range mgr.Get([]string{"Actor", "Position"}) {
		if _, ok := mgr.Component(e, "Sprite").(*game.Sprite); ok {
			continue
		}

		// FIXME: need to switch on Actor profession, sex, facing, etc.
		// actor := mgr.Component(e, "Actor").(*game.Actor)
		mgr.AddComponent(e, &game.Sprite{})
		fa := game.FrameAnimation{
			Frames: []game.Sprite{
				game.Sprite{
					Texture: "figure.png",
					X:       0,
					Y:       0,
					W:       24,
					H:       48,
					OffsetY: -16,
				},
				game.Sprite{
					Texture: "figure.png",
					X:       48,
					Y:       0,
					W:       24,
					H:       48,
					OffsetY: -16,
				},
				game.Sprite{
					Texture: "figure.png",
					X:       24,
					Y:       0,
					W:       24,
					H:       48,
					OffsetY: -16,
				},
			},
			Timings: []time.Duration{1500 * time.Millisecond, 300 * time.Millisecond, 300 * time.Millisecond},
		}
		mgr.AddComponent(e, fa.Randomise())
	}
}

func (ps *PerformanceSystem) handleActorMoving(t event.Typer) {
	ev := t.(*game.CombatActorMovementCommenced)
	e := ev.Entity
	actor := ps.mgr.Component(e, "Actor").(*game.Actor)

	// FIXME: need to switch on Actor profession, sex, direction of movement (facing), rate of movement, etc.
	if actor.Size == game.SMALL {
		// set moving animation
		ps.mgr.AddComponent(e, &game.FrameAnimation{
			Frames: []game.Sprite{
				game.Sprite{
					Texture: "figure.png",
					X:       0,
					Y:       48,
					W:       24,
					H:       48,
					OffsetY: -16,
				},
				game.Sprite{
					Texture: "figure.png",
					X:       48,
					Y:       48,
					W:       24,
					H:       48,
					OffsetY: -16,
				},
				game.Sprite{
					Texture: "figure.png",
					X:       24,
					Y:       48,
					W:       24,
					H:       48,
					OffsetY: -16,
				},
			},
			Timings: []time.Duration{55 * time.Millisecond, 45 * time.Millisecond, 60 * time.Millisecond},
		})
	}

}

func (ps *PerformanceSystem) handleActorStopped(t event.Typer) {
	ev := t.(*game.CombatActorMovementConcluded)
	e := ev.Entity
	actor := ps.mgr.Component(e, "Actor").(*game.Actor)

	// FIXME: need to switch on Actor profession, sex, direction of movement, rate of movement, etc.
	if actor.Size == game.SMALL {
		// set idle animation
		ps.mgr.AddComponent(e, &game.FrameAnimation{
			Frames: []game.Sprite{
				game.Sprite{
					Texture: "figure.png",
					X:       0,
					Y:       0,
					W:       24,
					H:       48,
					OffsetY: -16,
				},
				game.Sprite{
					Texture: "figure.png",
					X:       48,
					Y:       0,
					W:       24,
					H:       48,
					OffsetY: -16,
				},
				game.Sprite{
					Texture: "figure.png",
					X:       24,
					Y:       0,
					W:       24,
					H:       48,
					OffsetY: -16,
				},
			},
			Timings: []time.Duration{1500 * time.Millisecond, 300 * time.Millisecond, 300 * time.Millisecond},
		})
	}
}
