package combat

import (
	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/skill"
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
	bus.Subscribe(InjuryApplied{}.Type(), result.handleInjuryApplied)

	return &result
}

// percentDamageOverTime calculate how much damage to deal for a percentage of
// maximum health dealt damage over time effect. The percent arg is literally
// the percent, passing 7 means 7% of max health. perPrep means how long does it
// take to deal `percent` of the target's max health. max is the target's max
// health. elapsedPrep is how long to calculate the effect over.
// The function returns the damage inflicted as well as the amount of
// elapsedPrep that should be consumed by this damage.
// Passing 10,1000,100,1000 would result in 10% of 100 being dealt, with all
// 1000 elapsedPrep being consumed.
func percentDamageOverTime(percent int, perPrep int, max int, elapsedPrep int) (int, int) {
	dmg := ((max * percent) / 100) * elapsedPrep / perPrep
	// consumed := elapsedPrep - ((max*percent)/100)/perPrep*dmg
	var consumed int

	for i := 0; ; i++ {
		guess := ((max * percent) / 100) * i / perPrep
		if guess == dmg {
			consumed = i
			break
		}
	}

	return dmg, consumed
}

// bleedingDamageOverTime applies the bleeding injury rules to percentDamageOverTime.
func bleedingDamageOverTime(max int, elapsedPrep int) (int, int) {
	return percentDamageOverTime(7, 1000, max, elapsedPrep)
}

func (ds *damageSystem) ProcessDamageOverTime(elapsedPreparation int) {
	// for every participant affected by bleeding, poisoned, burning ...
	for _, e := range ds.mgr.Get([]string{"Participant"}) {
		participant := ds.mgr.Component(e, "Participant").(*Participant)

		for ty, injury := range participant.Injuries {
			switch ty {
			case skill.BleedingInjury:
				injury.Value -= elapsedPreparation
				injury.Remainder += elapsedPreparation
				if injury.Value < 0 {
					// if elapsed preparation exceeds the value, then remove
					// that much from the remainder too.
					injury.Remainder += injury.Value
				}

				damage, consumed := bleedingDamageOverTime(participant.maxHealth(), injury.Remainder)

				if damage > 0 {
					// Remove from the remainder, what we have converted to damage.
					injury.Remainder -= consumed

					ds.bus.Publish(&DamageAccepted{
						Target:     e,
						Amount:     damage,
						Reduced:    0,
						DamageType: game.PhysicalDamage,
					})
				}

				if injury.Value <= 0 {
					delete(participant.Injuries, ty)
				}
			}
		}
	}
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

	if target.CurrentHealth < 0 {
		target.CurrentHealth = 0
		target.Status = KnockedDown
		ds.bus.Publish(&ParticipantDied{ev.Target})
	} else {
		ds.mgr.AddComponent(ev.Target, &game.TakeDamageAnimation{})
	}
}

func (ds *damageSystem) handleInjuryApplied(event event.Typer) {
	ev := event.(*InjuryApplied)
	participant := ds.mgr.Component(ev.Target, "Participant").(*Participant)

	if participant.Injuries == nil {
		participant.Injuries = map[skill.InjuryType]*injury{
			ev.InjuryType: {Value: ev.Value},
		}
	} else if existing, ok := participant.Injuries[ev.InjuryType]; ok {
		existing.Value += ev.Value
	} else {
		participant.Injuries[ev.InjuryType] = &injury{
			Value: ev.Value,
		}
	}
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
