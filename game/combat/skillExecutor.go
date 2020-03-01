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

Attack that deals damage.

Summon a skeleton from a dead character, preventing resurrection.

Banish a character to a place without time, preventing them from preparing,
being affected by DoTs, or taking their turn.

Equalise HP of the user and target.

Transform into Animal form on activation - changing prep and AP values.

Apply a status affect like bleeding, poison, weakness etc.

Summon a wall of brambles to create a temporary obstacle on the field.

Create a projectile that has varying flight-time to its target, and deals
damage when it lands.

Dispell all effects on the target.

Passive skills:

Summon an animal familiar at the start of combat. Cannot be activated.

Summon a twin at the start of combat. Cannot be activated.

An aura that effects all allies and gives them effective +1 to all masteries.

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

type skillExecutionContext struct {
	ev      *UsingSkill
	age     time.Duration
	desc    *skill.Description
	effects []skill.Effect
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

func (se *skillExecutor) handleUsingSkill(t event.Typer) {
	ev := t.(*UsingSkill)
	s := se.archive.Skill(ev.Skill)

	// TODO: Complete implementation ...
	fmt.Printf("skillExecutor: Entity %d used %s on %v\n", ev.User, s.Name, ev.Selected)

	// Take a copy of effects so that we can remove them without mutating the
	// "reference" copy.
	effects := make([]skill.Effect, len(s.Effects))
	for i := 0; i < len(s.Effects); i++ {
		effects[i] = s.Effects[i]
	}

	se.inPlay = append(se.inPlay, &skillExecutionContext{
		ev:      ev,
		desc:    s,
		effects: effects,
	})

func (se *skillExecutor) executeEffect(effect skill.Effect, inPlay *skillExecutionContext) error {
	switch e := effect.(type) {
	case skill.DamageEffect:
		dereference := se.dereferencer(inPlay.ev.User)
		min := e.Min.Calculate(dereference)
		max := e.Max.Calculate(dereference)

		// Roll for damage between min and max.
		dmg := min
		if max != min {
			dmg += rand.Intn((max - min) + 1)
		}

		fmt.Printf("TODO: skillExecutor: DamageEffect: apply %d damage\n", dmg)
	default:
		return fmt.Errorf("unhandled skill effect type %T", e)
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
			target := inPlay.effects[0].Schedule()

			// When the first effect is not ready to fire yet, break out of this
			// skill, because none of the other effects should be ready either.
			if inPlay.age < target {
				break
			}

			// Apply effects on target (may or may not apply to characters).
			se.executeEffect(inPlay.effects[0], inPlay)

			// Because this effect has been executed, it should be removed.
			inPlay.effects = inPlay.effects[1:]
		}
	}
}
