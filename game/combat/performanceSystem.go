package combat

import (
	"fmt"
	"time"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/geom"
	"github.com/griffithsh/squads/res"
)

/*
Actors Perform Movement, Striking with a weapon, Spellcasting, Hurt, Victory,
and Dying, otherwise they are Idle.

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

	bus.Subscribe(game.CombatActorMoving{}.Type(), ps.handleActorMoving)

	return &ps
}

// Update the System.
func (ps *PerformanceSystem) Update(elapse time.Duration) {
	// For any actors without a sprite, apply their idling animation
	for _, e := range ps.mgr.Get([]string{"Actor", "Position"}) {
		if _, ok := ps.mgr.Component(e, "Sprite").(*game.Sprite); ok {
			continue
		}

		actor := ps.mgr.Component(e, "Actor").(*game.Actor)
		facer := ps.mgr.Component(e, "Facer").(*game.Facer)

		fa := get(animationId{actor.Profession, actor.Sex, game.PerformIdle, facer.Face})

		// Start at a random point of the Idle animation.
		ps.mgr.AddComponent(e, fa.Randomise())
	}
}

var missing = map[animationId]struct{}{}

func reportMissing(id animationId) {
	if _, ok := missing[id]; ok {
		return
	}
	fmt.Println("missing animation:", id)
	missing[id] = struct{}{}
}

func get(id animationId) game.FrameAnimation {
	if fa, ok := all[id]; ok {
		return fa
	}
	reportMissing(id)
	return notFound()
}

func notFound() game.FrameAnimation {
	return game.FrameAnimation{
		Frames: []game.Sprite{
			{Texture: "tranquility-plus-39-palette.png", W: 8, H: 8},
		},
		Timings: []time.Duration{time.Second},
	}
}

func (ps *PerformanceSystem) handleActorMoving(t event.Typer) {
	ev := t.(*game.CombatActorMoving)
	e := ev.Entity
	actor := ps.mgr.Component(e, "Actor").(*game.Actor)
	facer := ps.mgr.Component(e, "Facer").(*game.Facer)

	// If the facing has changed, then we need to edit the FrameAnimation.
	if ev.OldFacing != ev.NewFacing {
		fa := get(animationId{actor.Profession, actor.Sex, game.PerformMove, facer.Face})
		ps.mgr.AddComponent(e, &fa)
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
			fa := get(animationId{actor.Profession, actor.Sex, game.PerformMove, facer.Face})
			ps.mgr.AddComponent(e, &fa)
		}
	}
}

type animationId struct {
	Profession  game.ActorProfession
	Sex         game.ActorSex
	Performance game.ActorPerformance
	Facing      geom.DirectionType
}

// links is the declaration map between animationId and string keys present in
// res.All.
var links = map[animationId]string{
	animationId{game.Villager, game.Male, game.PerformIdle, geom.N}:  "Villager-Male-Idle",
	animationId{game.Villager, game.Male, game.PerformIdle, geom.S}:  "Villager-Male-Idle",
	animationId{game.Villager, game.Male, game.PerformIdle, geom.SE}: "Villager-Male-Idle",
	animationId{game.Villager, game.Male, game.PerformIdle, geom.SW}: "Villager-Male-Idle",
	animationId{game.Villager, game.Male, game.PerformIdle, geom.NE}: "Villager-Male-Idle",
	animationId{game.Villager, game.Male, game.PerformIdle, geom.NW}: "Villager-Male-Idle",

	animationId{game.Villager, game.Male, game.PerformMove, geom.N}:  "Villager-Male-Move",
	animationId{game.Villager, game.Male, game.PerformMove, geom.S}:  "Villager-Male-Move",
	animationId{game.Villager, game.Male, game.PerformMove, geom.SE}: "Villager-Male-Move",
	animationId{game.Villager, game.Male, game.PerformMove, geom.SW}: "Villager-Male-Move",
	animationId{game.Villager, game.Male, game.PerformMove, geom.NE}: "Villager-Male-Move",
	animationId{game.Villager, game.Male, game.PerformMove, geom.NW}: "Villager-Male-Move",

	animationId{game.Wolf, game.Male, game.PerformIdle, geom.N}:  "Wolf",
	animationId{game.Wolf, game.Male, game.PerformIdle, geom.S}:  "Wolf",
	animationId{game.Wolf, game.Male, game.PerformIdle, geom.SE}: "Wolf",
	animationId{game.Wolf, game.Male, game.PerformIdle, geom.SW}: "Wolf",
	animationId{game.Wolf, game.Male, game.PerformIdle, geom.NE}: "Wolf",
	animationId{game.Wolf, game.Male, game.PerformIdle, geom.NW}: "Wolf",

	animationId{game.Wolf, game.Male, game.PerformMove, geom.N}:  "Wolf",
	animationId{game.Wolf, game.Male, game.PerformMove, geom.S}:  "Wolf",
	animationId{game.Wolf, game.Male, game.PerformMove, geom.SE}: "Wolf",
	animationId{game.Wolf, game.Male, game.PerformMove, geom.SW}: "Wolf",
	animationId{game.Wolf, game.Male, game.PerformMove, geom.NE}: "Wolf",
	animationId{game.Wolf, game.Male, game.PerformMove, geom.NW}: "Wolf",

	animationId{game.Giant, game.Male, game.PerformIdle, geom.N}:  "Giant",
	animationId{game.Giant, game.Male, game.PerformIdle, geom.S}:  "Giant",
	animationId{game.Giant, game.Male, game.PerformIdle, geom.SE}: "Giant",
	animationId{game.Giant, game.Male, game.PerformIdle, geom.SW}: "Giant",
	animationId{game.Giant, game.Male, game.PerformIdle, geom.NE}: "Giant",
	animationId{game.Giant, game.Male, game.PerformIdle, geom.NW}: "Giant",

	animationId{game.Wolf, game.Male, game.PerformMove, geom.N}:  "Wolf",
	animationId{game.Wolf, game.Male, game.PerformMove, geom.S}:  "Wolf",
	animationId{game.Wolf, game.Male, game.PerformMove, geom.SE}: "Wolf",
	animationId{game.Wolf, game.Male, game.PerformMove, geom.SW}: "Wolf",
	animationId{game.Wolf, game.Male, game.PerformMove, geom.NE}: "Wolf",
	animationId{game.Wolf, game.Male, game.PerformMove, geom.NW}: "Wolf",

	animationId{game.Skeleton, game.Male, game.PerformIdle, geom.N}:  "Skeleton-Idle-N",
	animationId{game.Skeleton, game.Male, game.PerformIdle, geom.S}:  "Skeleton-Idle-S",
	animationId{game.Skeleton, game.Male, game.PerformIdle, geom.SE}: "Skeleton-Idle-SE",
	animationId{game.Skeleton, game.Male, game.PerformIdle, geom.SW}: "Skeleton-Idle-SW",
	animationId{game.Skeleton, game.Male, game.PerformIdle, geom.NE}: "Skeleton-Idle-NE",
	animationId{game.Skeleton, game.Male, game.PerformIdle, geom.NW}: "Skeleton-Idle-NW",

	animationId{game.Skeleton, game.Male, game.PerformMove, geom.N}:  "Skeleton-Move-N",
	animationId{game.Skeleton, game.Male, game.PerformMove, geom.S}:  "Skeleton-Move-S",
	animationId{game.Skeleton, game.Male, game.PerformMove, geom.SE}: "Skeleton-Move-SE",
	animationId{game.Skeleton, game.Male, game.PerformMove, geom.SW}: "Skeleton-Move-SW",
	animationId{game.Skeleton, game.Male, game.PerformMove, geom.NE}: "Skeleton-Move-NE",
	animationId{game.Skeleton, game.Male, game.PerformMove, geom.NW}: "Skeleton-Move-NW",
}

// final map between animationId and game.FrameAnimation.
var all = map[animationId]game.FrameAnimation{}

// init function that populates the "all" map with the keys from links and the
// values from res.All.
func init() {
	for k, name := range links {
		a, ok := res.All[name]
		if !ok {
			panic(fmt.Sprintf("links misconfigured: \"%s\" not found in res.All", name))
		}
		fa := game.FrameAnimation{}
		for _, frame := range a.Frames {
			fa.Frames = append(fa.Frames, game.Sprite{
				Texture: frame.Texture,
				X:       frame.X,
				Y:       frame.Y,
				W:       frame.W,
				H:       frame.H,
				OffsetX: frame.OffsetX,
				OffsetY: frame.OffsetY,
			})
			fa.Timings = append(fa.Timings, frame.Duration)
		}
		all[k] = fa
	}
}
