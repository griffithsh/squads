package combat

import (
	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
)

type damageSystem struct {
	mgr *ecs.World
	bus *event.Bus
}

func newDamageSystem(mgr *ecs.World, bus *event.Bus) *damageSystem {
	result := damageSystem{
		mgr: mgr,
		bus: bus,
	}
	bus.Subscribe(DamageApplied{}.Type(), result.handleDamageApplied)

	return &result
}

func (ds *damageSystem) handleDamageApplied(event event.Typer) {
	ev := event.(*DamageApplied)

	if failed := ds.failure(ev); failed != "" {
		ds.bus.Publish(&DamageFailed{
			Target: ev.Target,
			Reason: failed,
		})
		return
	}

	accepted, ty, reduced := ds.reduce(ev)

	target := ds.mgr.Component(ev.Target, "Participant").(*Participant)

	target.CurrentHealth -= accepted

	// Let the UI know about this.
	ds.bus.Publish(&DamageAccepted{
		Target:     ev.Target,
		Amount:     accepted,
		Reduced:    reduced,
		DamageType: ty,
	})

	ds.mgr.AddComponent(ev.Target, &game.TakeDamageAnimation{})
}

// failure calculates whether the applied damage has failed to be applied or
// not.
func (ds *damageSystem) failure(ev *DamageApplied) string {
	// TODO: calculate chance to fail
	return ""
}

// reduce calculates what damage is accepted by an application, and what type of
// damage it was. (Because some targets may convert the type of damage
// received). It also returns how much of the original amount was reduced.
func (ds *damageSystem) reduce(ev *DamageApplied) (accepted int, ty game.DamageType, reduced int) {
	// TODO: calculate chance of negate/dodge
	return ev.Amount, ev.DamageType, 0
}
