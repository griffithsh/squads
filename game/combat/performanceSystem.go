package combat

import (
	"time"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/geom"
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

	// TODO:
	// each handler applies an animation. Some animations wear off, and leave
	// only a sprite. Others continue indefinitely. In Update, Participants
	// without an Animation receive the idle animation.
	bus.Subscribe(ParticipantMoving{}.Type(), ps.handleParticipantMoving)
	bus.Subscribe(CharacterCelebrating{}.Type(), ps.handleCharacterCelebrating)
	bus.Subscribe(UsingSkill{}.Type(), ps.handleUsingSkill)
	bus.Subscribe(ParticipantDied{}.Type(), ps.handleParticipantDied)

	return &ps
}

// Update the System.
func (ps *PerformanceSystem) Update(elapse time.Duration) {
	// For any Participants without a sprite, apply their idling animation
	for _, e := range ps.mgr.Get([]string{"Participant", "Position"}) {
		if _, ok := ps.mgr.Component(e, "Sprite").(*game.Sprite); ok {
			continue
		}

		facer := ps.mgr.Component(e, "Facer").(*game.Facer)

		performances := ps.getPerformances(e)
		frames := performances.Idle.ForDirection(facer.Face)
		fa := game.NewFrameAnimationFromFrames(frames)

		// Start at a random point of the Idle animation.
		ps.mgr.AddComponent(e, fa.Randomise())
	}
}

func (ps *PerformanceSystem) getPerformances(e ecs.Entity) *game.PerformanceSet {
	participant := ps.mgr.Component(e, "Participant").(*Participant)
	prof := participant.Profession
	sex := participant.Sex
	return ps.archive.Performances(prof, sex)
}

func (ps *PerformanceSystem) handleParticipantMoving(t event.Typer) {
	ev := t.(*ParticipantMoving)
	e := ev.Entity
	facer := ps.mgr.Component(e, "Facer").(*game.Facer)
	performances := ps.getPerformances(ev.Entity)
	frames := performances.Move.ForDirection(facer.Face)

	// If the facing has changed, then we need to edit the FrameAnimation.
	if ev.OldFacing != ev.NewFacing {
		fa := game.NewFrameAnimationFromFrames(frames)
		ps.mgr.AddComponent(e, fa)
	}

	if ev.NewSpeed == 0 {
		// If the entity has stopped moving, then we must delete the sprite so
		// that Update can add the Idle animation in.
		ps.mgr.RemoveComponent(e, &game.Sprite{})
		ps.mgr.RemoveComponent(e, &game.AnimationSpeed{})
	} else if ev.OldSpeed != ev.NewSpeed {
		// Otherwise the speed has changed ...
		ps.mgr.AddComponent(e, &game.AnimationSpeed{
			Speed: ev.NewSpeed,
		})
		if ev.OldSpeed == 0 {
			fa := game.NewFrameAnimationFromFrames(frames)
			ps.mgr.AddComponent(e, fa)
		}
	}
}

func (ps *PerformanceSystem) handleCharacterCelebrating(t event.Typer) {
	ev := t.(*CharacterCelebrating)
	e := ev.Entity

	performances := ps.getPerformances(e)
	fa := game.NewFrameAnimationFromFrames(performances.Victory)
	fa.EndBehavior = game.HoldLastFrame
	ps.mgr.AddComponent(e, fa)
}

func (ps *PerformanceSystem) handleUsingSkill(t event.Typer) {
	ev := t.(*UsingSkill)
	e := ev.User

	performances := ps.getPerformances(e)
	skill := ps.archive.Skill(ev.Skill)
	var frames []game.Frame
	frames = performances.Spell
	if skill.IsAttack() {
		facer := ps.mgr.Component(e, "Facer").(*game.Facer)

		switch facer.Face {
		case geom.N:
			frames = performances.Attack.N
		case geom.S:
			frames = performances.Attack.S
		case geom.SE:
			frames = performances.Attack.SE
		case geom.SW:
			frames = performances.Attack.SW
		case geom.NE:
			frames = performances.Attack.NE
		case geom.NW:
			frames = performances.Attack.NW
		}
	}
	fa := game.NewFrameAnimationFromFrames(frames)
	fa.EndBehavior = game.HoldLastFrame
	ps.mgr.AddComponent(e, fa)
}

func (ps *PerformanceSystem) handleParticipantDied(t event.Typer) {
	pde := t.(*ParticipantDied)

	performances := ps.getPerformances(pde.Entity)
	fa := game.NewFrameAnimationFromFrames(performances.Death)
	fa.EndBehavior = game.HoldLastFrame
	ps.mgr.AddComponent(pde.Entity, fa)
}
