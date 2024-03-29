package combat

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/geom"
	"github.com/griffithsh/squads/skill"
)

/*
Skill ideas:

[*] Attack that deals damage.

[ ] Summon a skeleton from a dead character, preventing resurrection.

[ ] Banish a character to a place without time, preventing them from preparing,
being affected by DoTs, or taking their turn.

[ ] Equalise HP of the user and target.

[ ] Transform into Animal form on activation - changing prep and AP values.

[ ] Apply a status affect like bleeding, poison, weakness etc.

[ ] Summon a wall of brambles to create a temporary obstacle on the field.

[ ] Create a projectile that has varying flight-time to its target, and deals
damage when it lands.

[ ] Dispell all effects on the target.

Passive skills:

[ ] Summon an animal familiar at the start of combat. Cannot be activated.

[ ] Summon a twin at the start of combat. Cannot be activated.

[ ] An aura that effects all allies and gives them effective +1 to all masteries.

*/

type skillExecutor struct {
	mgr     *ecs.World
	bus     *event.Bus
	field   *geom.Field
	archive SkillArchive

	inPlay []*skillExecutionContext
}

func newSkillExecutor(mgr *ecs.World, bus *event.Bus, field *geom.Field, archive SkillArchive) *skillExecutor {
	se := skillExecutor{
		mgr:     mgr,
		bus:     bus,
		field:   field,
		archive: archive,
	}
	se.bus.Subscribe(UsingSkill{}.Type(), se.handleUsingSkill)

	return &se
}

type effect struct {
	when time.Duration
	what interface{}
}

type skillExecutionContext struct {
	ev       *UsingSkill
	age      time.Duration
	desc     *skill.Description
	effects  []effect
	affected []ecs.Entity
	targeted []geom.Key
	missCalc func() bool
}

// determineAffected collects the Entities that the usage of this skill affects,
// primarily based on TargetingBrush. It's perfectly normal for this to return
// no entites, as skills (AoEs, summons) can be cast on empty hexes. It also
// returns the hexes in the field that are targeted.
func (se *skillExecutor) determineAffected(ev *UsingSkill, s *skill.Description) ([]ecs.Entity, []geom.Key) {
	affected := []ecs.Entity{}

	user := se.mgr.Component(ev.User, "Obstacle").(*game.Obstacle)
	origin := geom.Key{M: user.M, N: user.N}
	_, painted := s.Targeting.Execute(ev.Selected.Key(), origin)

	for _, e := range se.mgr.Get([]string{"Participant"}) {
		// Defiled Participants do not have an Obstacle.
		o, exists := se.mgr.Component(e, "Obstacle").(*game.Obstacle)
		if !exists {
			continue
		}
		for _, k := range painted {
			if k.M == o.M && k.N == o.N {
				affected = append(affected, e)
				break
			}
		}
	}

	return affected, painted
}

// createRealiser creates a new timing point realiser for figuring out when
// effects with virtual timing points should be executed.
// FIXME: delete this?
func (se *skillExecutor) createRealiser(ev *UsingSkill, s *skill.Description) func(skill.Timing) time.Duration {
	m := map[skill.TimingPoint]time.Duration{
		skill.AttackApexTimingPoint: 0,
		skill.EndTimingPoint:        0,
	}
	return skill.NewTimingRealiser(m)
}

func (se *skillExecutor) handleUsingSkill(t event.Typer) {
	ev := t.(*UsingSkill)
	s := se.archive.Skill(ev.Skill)

	// if the skill is tagged Attack, then we need to apply the chance to hit
	// modifier to the base chance to hit. Spells do not miss and ignore the
	// chance to hit modifier.
	missCalc := func() bool {
		return false
	}
	if s.IsAttack() {
		usingParticipant := se.mgr.Component(ev.User, "Participant").(*Participant)
		baseChance := usingParticipant.chanceToHit()

		chance := baseChance + ((1.0 - baseChance) * s.AttackChanceToHitModifier)
		if chance > 1.0 {
			chance = 1.0
		} else if chance < 0 {
			chance = 0
		}
		missCalc = func() bool {
			roll := rand.Float64()
			return roll > chance
		}
	}

	realiser := se.createRealiser(ev, s)

	// Take a copy of effects so that we can remove them without mutating the
	// "reference" copy. Also use the realiser to convert to concrete timing
	// points.
	effects := make([]effect, len(s.Effects))
	for i := 0; i < len(s.Effects); i++ {
		effects[i] = effect{
			when: realiser(s.Effects[i].When),
			what: s.Effects[i].What,
		}
	}

	affected, targeted := se.determineAffected(ev, s)

	se.inPlay = append(se.inPlay, &skillExecutionContext{
		ev:       ev,
		desc:     s,
		effects:  effects,
		affected: affected,
		targeted: targeted,
		missCalc: missCalc,
	})

	// Apply costs of skill to user.
	participant := se.mgr.Component(ev.User, "Participant").(*Participant)
	for ty, amount := range s.Costs {
		switch ty {
		case skill.CostsActionPoints:
			participant.ActionPoints.Cur -= amount
		case skill.CostsMana:
			// TODO:
		default:
			panic(fmt.Sprintf("skillExector.handleUsingSkill: Cost Type %T not implemented", ty))
		}
	}
}

func (se *skillExecutor) dereferencer(e ecs.Entity) func(s string) float64 {
	participant := se.mgr.Component(e, "Participant").(*Participant)
	min, max := participant.baseDamage()
	return func(s string) float64 {
		if strings.HasPrefix(s, "$") {
			switch s {
			// Core Stats
			case "$INT":
				return float64(participant.Intelligence)

			// Base Damage.
			case "$DMG-MIN":
				return min
			case "$DMG-MAX":
				return max

			// Masteries.
			case "$SHORT-RANGE-MELEE":
				return float64(participant.Masteries[game.ShortRangeMeleeMastery])
			case "$LONG-RANGE-MELEE":
				return float64(participant.Masteries[game.LongRangeMeleeMastery])
			case "$RANGED-COMBAT":
				return float64(participant.Masteries[game.RangedCombatMastery])
			case "$CRAFTS":
				return float64(participant.Masteries[game.CraftsmanshipMastery])
			case "$FIRE":
				return float64(participant.Masteries[game.FireMastery])
			case "$WATER":
				return float64(participant.Masteries[game.WaterMastery])
			case "$EARTH":
				return float64(participant.Masteries[game.EarthMastery])
			case "$AIR":
				return float64(participant.Masteries[game.AirMastery])
			case "$LIGHTNING":
				return float64(participant.Masteries[game.LightningMastery])
			case "$DARK":
				return float64(participant.Masteries[game.DarkMastery])
			case "$LIGHT":
				return float64(participant.Masteries[game.LightMastery])

			// Code is wrong.
			default:
				panic(fmt.Sprintf("dereference: unsupported variable \"%s\"", s))
			}
		}
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			panic(fmt.Sprintf("dereference \"%s\": %v", s, err))
		}
		return f
	}
}

func (se *skillExecutor) executeEffect(effect effect, inPlay *skillExecutionContext) error {
	whats := effect.what.([]interface{})
	for _, what := range whats {
		switch ef := what.(type) {
		case skill.DamageEffect:
			dereference := se.dereferencer(inPlay.ev.User)
			min := ef.Min.Calculate(dereference)
			max := ef.Max.Calculate(dereference)

			// Roll for damage between min and max.
			dmg := min
			if max != min {
				dmg += rand.Intn((max - min) + 1)
			}

			for _, affected := range inPlay.affected {
				if inPlay.missCalc() {
					se.bus.Publish(&DamageFailed{
						Target: affected,
						Reason: "Miss",
					})
					continue
				}
				se.bus.Publish(&DamageApplied{
					Amount:     dmg,
					Target:     affected,
					DamageType: game.PhysicalDamage,
					SkillType:  ef.Classification,
				})
			}
		case skill.ReviveEffect:
			for _, e := range inPlay.affected {
				participant := se.mgr.Component(e, "Participant").(*Participant)

				if participant.Status != KnockedDown {
					continue
				}

				participant.Status = Alive
				participant.CurrentHealth = 1
				se.bus.Publish(&ParticipantRevived{Entity: e})
			}
		case skill.HealEffect:
			for _, e := range inPlay.affected {
				participant := se.mgr.Component(e, "Participant").(*Participant)

				var heal int
				if ef.IsPercentage {
					heal = int(float64(participant.maxHealth()) * ef.Amount)
				} else {
					heal = int(ef.Amount)
				}
				participant.CurrentHealth = heal
			}
		case skill.DefileEffect:
			for _, e := range inPlay.affected {
				participant := se.mgr.Component(e, "Participant").(*Participant)
				if participant.Status == KnockedDown {
					participant.Status = Defiled
					se.bus.Publish(&ParticipantDefiled{Entity: e})
				}
			}
		case skill.SpawnParticipantEffect:
			dereference := se.dereferencer(inPlay.ev.User)
			for _, key := range inPlay.targeted {
				// FIXME: When executing a SpawnParticipantEffect, it is assumed
				// that the new participant is on the same team as the User of the
				// skill. This assumption might not always hold though.
				team := se.mgr.Component(inPlay.ev.User, "Team").(*game.Team)

				se.bus.Publish(&CharacterEnteredCombat{
					Level:      ef.Level.Calculate(dereference),
					Profession: ef.Profession,
					Team:       team,
					At:         key,
				})
			}
		case skill.InjuryEffect:
			for _, affected := range inPlay.affected {
				se.bus.Publish(&InjuryApplied{
					Target:     affected,
					InjuryType: ef.Type,
					Value:      ef.Value,
				})
			}
		default:
			return fmt.Errorf("unhandled skill effect type %T", ef)
		}
	}
	return nil
}

func (se *skillExecutor) Update(elapsed time.Duration) {
	if len(se.inPlay) == 0 {
		return
	}

	for i := 0; i < len(se.inPlay); i++ {
		inPlay := se.inPlay[i]
		inPlay.age += elapsed
		for {
			if len(inPlay.effects) == 0 {
				// End of skill; remove this skill execution context from se.inPlay.
				se.inPlay = append(se.inPlay[:i], se.inPlay[i+1:]...)
				i--

				if len(se.inPlay) == 0 {
					se.bus.Publish(&SkillUseConcluded{
						// FIXME: We cannot provide these values if we are using multiple
						// skills in play. Do we even need these values?
						0, "", nil,
					})
				}
				break
			}

			// When the first effect is not ready to fire yet, break out of this
			// skill, because none of the other effects should be ready either.
			if inPlay.age < inPlay.effects[0].when {
				break
			}

			// Apply effects on target (may or may not apply to characters).
			if err := se.executeEffect(inPlay.effects[0], inPlay); err != nil {
				panic(err)
			}

			// Because this effect has been executed, it should be removed.
			inPlay.effects = inPlay.effects[1:]
		}
	}
}
