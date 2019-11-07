package overworld

import (
	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
)

// CollisionSystem generates collision events when two Tokens occupy the same node.
type CollisionSystem struct {
	mgr *ecs.World
	bus *event.Bus
}

// NewCollisionSystem creates a new Collision System.
func NewCollisionSystem(mgr *ecs.World, bus *event.Bus) *CollisionSystem {
	cs := CollisionSystem{
		mgr: mgr,
		bus: bus,
	}
	bus.Subscribe(TokenMoved{}.Type(), cs.handleTokenMoved)
	return &cs
}

func (cs *CollisionSystem) handleTokenMoved(t event.Typer) {
	ev := t.(*TokenMoved)
	team := cs.mgr.Component(ev.E, "Team").(*game.Team)
	for _, e := range cs.mgr.Get([]string{"Token", "Team"}) {
		// Entities don't collide with themself.
		if ev.E == e {
			continue
		}

		// If they're on the same team, no collision.
		otherTeam := cs.mgr.Component(e, "Team").(*game.Team)
		if team.ID == otherTeam.ID {
			continue
		}

		otherToken := cs.mgr.Component(e, "Token").(*Token)
		if otherToken.Key == ev.To {
			cs.bus.Publish(&TokensCollided{
				E1: ev.E,
				E2: e,
				At: ev.To,
			})
		}
	}
}
