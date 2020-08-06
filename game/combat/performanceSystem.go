package combat

import (
	"time"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
)

/*
Participants Perform Movement, Striking with a weapon, Spellcasting, Hurt,
Victory, and Dying, otherwise they are Idle.

Striking with a weapon and spellcasting are probably the same thing - "using
a skill" but different professions could have different numbers of skill
animations, and their skills could map to arbitrary animations.

The mapping between SkillA and "appropriate animation for SkillA when used by
$profession and $sex when $facing" could be owned either by the same thing
that owns the rest of the mappings, or the animation (i.e SKILL_A, SKILL_B)
used for $skill could be configured on the skills.

I think mapping skill to animation class on the skill is appropriate. This
would mean that the "UsedSkill" event would include a reference to the
animation "class" that best represents this skill.
*/

// PerformanceSystem sets appropriate Animations for Participants in a Combat based on
// what's happening.
type PerformanceSystem struct {
	mgr     *ecs.World
	archive SkillArchive
}

// NewPerformanceSystem creates a new PerformanceSystem.
func NewPerformanceSystem(mgr *ecs.World, bus *event.Bus, archive SkillArchive) *PerformanceSystem {
	ps := PerformanceSystem{
		mgr:     mgr,
		archive: archive,
	}

	bus.Subscribe(ParticipantMoving{}.Type(), ps.handleParticipantMoving)
	bus.Subscribe(CharacterCelebrating{}.Type(), ps.handleCharacterCelebrating)
	bus.Subscribe(UsingSkill{}.Type(), ps.handleUsingSkill)
	bus.Subscribe(ParticipantDied{}.Type(), ps.handleParticipantDied)
	bus.Subscribe(ParticipantRevived{}.Type(), ps.handleParticipantRevived)
	bus.Subscribe(ParticipantDefiled{}.Type(), ps.handleParticipantDefiled)

	return &ps
}

// Update the System, applying a Sprite Appearance to any participant without
// one.
func (ps *PerformanceSystem) Update(elapse time.Duration) {
	// For every Participant in the combat ...
	for _, e := range ps.mgr.Get([]string{"Participant"}) {
		// If the Entity has a FrameAnimation already, then they are animating
		// some action (or they might be idling already?), so don't change
		// anything.
		if _, ok := ps.mgr.Component(e, "Sprite").(*game.Sprite); ok {
			continue
		}

		participant := ps.mgr.Component(e, "Participant").(*Participant)
		prof := participant.Profession
		sex := participant.Sex
		appearance := ps.archive.Appearance(prof, sex, participant.Hair, participant.Skin)
		// copy the archive sprite
		spr := appearance.Participant

		spr.OffsetX = spr.W/2 - appearance.FeetX // almost always going to be zero
		spr.OffsetY = spr.H/2 - appearance.FeetY // should be negative number

		ps.mgr.AddComponent(e, &spr)
	}
}

func (ps *PerformanceSystem) handleParticipantMoving(t event.Typer) {
	ev := t.(*ParticipantMoving)

	// If the facing has changed, then we need to edit the FrameAnimation.
	if ev.OldFacing != ev.NewFacing {
	}

	if ev.NewSpeed == 0 {
		// If the entity has stopped moving, then we must delete the sprite so
		// that Update can add the Idle animation in.
	} else if ev.OldSpeed != ev.NewSpeed {
		// Otherwise the speed has changed ...
	}
}

func (ps *PerformanceSystem) handleCharacterCelebrating(t event.Typer) {
	// ev := t.(*CharacterCelebrating)

	// TODO: something?
}

func (ps *PerformanceSystem) handleUsingSkill(t event.Typer) {
	// ev := t.(*UsingSkill)

	// TODO!
}

func (ps *PerformanceSystem) handleParticipantDied(t event.Typer) {
	// ev := t.(*ParticipantDied)

	// TODO!
}

func (ps *PerformanceSystem) handleParticipantRevived(t event.Typer) {
	// ev := t.(*ParticipantRevived)
	// TODO: something?
}

func (ps *PerformanceSystem) handleParticipantDefiled(t event.Typer) {
	pde := t.(*ParticipantDefiled)

	ps.mgr.RemoveComponent(pde.Entity, &game.Sprite{})
	ps.mgr.RemoveComponent(pde.Entity, &game.FrameAnimation{})
}
