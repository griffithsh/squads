package combat

import (
	"fmt"
	"sort"
	"time"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/ui"
)

/*
There are two (pseudo-public) methods to control all UI groups.

ShowX
HideX

Another method is only called internally - by the Update method.

RepaintX

- Show must add all entities, their tags, and parent relationships and the
invalidated tag. Because it's adding entities, it needs to ensure that it
does not add duplicates. It's the only method that should call to
mgr.NewEntity().

- Hide must call DestroyEntity on all the Entities that compose the group.
Anything that is added by Show.

- Repaint must add/replace the Position, Sprite, Font etc Components,
and then remove the invalidated tag.

There is no hierarchy between groups, each set of these functions is an
island disconnected from the others, and is only responsible for its own
Entities.

HUD element groups:

AllActors
TimePassingIcon
CurrentActor
	Hover-er
	Portrait
	Name
CurrentActorStats
	Health
	Energy
	Action
	Preparation
Skills
	Icon,Interactive
TurnQueue[]
	Portrait
	Prep
*/

var (
	combatHUDTag            = "COMBAT_HUD"
	liveActorsTag           = combatHUDTag + ".LIVE_ACTORS"
	timePassingTag          = combatHUDTag + ".TIME_PASSING"
	currentActorTag         = combatHUDTag + ".CURRENT_ACTOR"
	currentActorHovererTag  = currentActorTag + ".HOVERER"
	currentActorPortraitTag = currentActorTag + ".PORTRAIT"
	currentActorNameTag     = currentActorTag + ".NAME"
	currentActorStatsTag    = currentActorTag + ".STATS"
	skillsTag               = currentActorTag + ".SKILLS"
	turnQueueTag            = combatHUDTag + ".TURN_QUEUE"
	endTurnButtonTag        = combatHUDTag + ".END_TURN_BUTTON"

	invalidatedTag = combatHUDTag + ".INVALIDATED"
)

// HUD is a heads up display for the combat game state.
type HUD struct {
	mgr              *ecs.World
	bus              *event.Bus
	scale            float64
	layer            int
	centerX, centerY float64 // center of the game's window, or half the width and height.
	lastCombatState  State
}

// NewHUD construct a HUD.
func NewHUD(mgr *ecs.World, bus *event.Bus, screenX int, screenY int) *HUD {
	hud := HUD{
		mgr:             mgr,
		bus:             bus,
		scale:           2,
		layer:           100,
		centerX:         float64(screenX) / 2,
		centerY:         float64(screenY) / 2,
		lastCombatState: AwaitingInputState,
	}

	bus.Subscribe(game.WindowSizeChanged{}.Type(), hud.handleWindowSizeChanged)
	bus.Subscribe(game.CombatBegan{}.Type(), hud.handleCombatBegan)
	bus.Subscribe(game.CombatStatModified{}.Type(), hud.handleCombatStatModified)
	bus.Subscribe(StateTransition{}.Type(), hud.handleCombatStateTransition)
	bus.Subscribe(DifferentHexSelected{}.Type(), hud.handleDifferentHexSelected)

	return &hud
}

func (hud *HUD) handleWindowSizeChanged(e event.Typer) {
	wsc := e.(*game.WindowSizeChanged)
	hud.centerX, hud.centerY = float64(wsc.NewW)/2, float64(wsc.NewH)/2

	if hud.mgr.AnyTagged(timePassingTag) == 0 {
		return
	}

	hud.showTimePassingIcon()
}

func (hud *HUD) handleCombatBegan(event.Typer) {
	e := hud.mgr.NewEntity()
	hud.mgr.Tag(e, combatHUDTag)

	hud.showTurnQueue()
	hud.showLiveActors()
}

func (hud *HUD) handleCombatStatModified(ev event.Typer) {
	csm := ev.(*game.CombatStatModified)

	switch csm.Stat {
	case game.PrepStat:
		hud.showTurnQueue()
	}
}

func (hud *HUD) handleCombatStateTransition(ev event.Typer) {
	cst := ev.(*StateTransition)
	hud.lastCombatState = cst.New

	if cst.New == AwaitingInputState || cst.New == SelectingTargetState {
		hud.showSkills()
		if cst.Old != AwaitingInputState && cst.Old != SelectingTargetState {
			hud.showCurrentActor()
			hud.showCurrentActorStats()
		}
	} else {
		hud.hideSkills()
		hud.hideCurrentActor()
		hud.hideCurrentActorStats()
	}

	if cst.New == PreparingState {
		hud.showTimePassingIcon()
	} else {
		hud.hideTimePassingIcon()
	}
}

func (hud *HUD) handleDifferentHexSelected(ev event.Typer) {

	// TODO!
	// invalidate
}

// Update the HUD. Synchronise the current game state to the Entities that compose it.
func (hud *HUD) Update(elapsed time.Duration) {
	var e ecs.Entity

	e = hud.mgr.AnyTagged(liveActorsTag)
	if e != 0 && hud.mgr.HasTag(e, invalidatedTag) {
		hud.repaintLiveActors()
	}
	e = hud.mgr.AnyTagged(timePassingTag)
	if e != 0 && hud.mgr.HasTag(e, invalidatedTag) {
		hud.repaintTimePassingIcon()
	}
	e = hud.mgr.AnyTagged(currentActorTag)
	if e != 0 && hud.mgr.HasTag(e, invalidatedTag) {
		hud.repaintCurrentActor()
	}
	e = hud.mgr.AnyTagged(currentActorStatsTag)
	if e != 0 && hud.mgr.HasTag(e, invalidatedTag) {
		hud.repaintCurrentActorStats()
	}
	e = hud.mgr.AnyTagged(skillsTag)
	if e != 0 && hud.mgr.HasTag(e, invalidatedTag) {
		hud.repaintSkills()
	}
	e = hud.mgr.AnyTagged(turnQueueTag)
	if e != 0 && hud.mgr.HasTag(e, invalidatedTag) {
		hud.repaintTurnQueue()
	}
}

const maxLiveActors int = 25

func (hud *HUD) showLiveActors() {
	for _, e := range hud.mgr.Tagged(liveActorsTag) {
		hud.mgr.DestroyEntity(e)
	}

	for i := 0; i < maxLiveActors; i++ {
		e := hud.mgr.NewEntity()

		hud.mgr.Tag(e, liveActorsTag)
		hud.mgr.Tag(e, invalidatedTag)
	}
}

func (hud *HUD) hideLiveActors() {
	for _, e := range hud.mgr.Tagged(liveActorsTag) {
		hud.mgr.DestroyEntity(e)
	}
}

func (hud *HUD) repaintLiveActors() {
	entities := hud.mgr.Get([]string{"Actor"})
	for i, slot := range hud.mgr.Tagged(liveActorsTag) {
		if i < len(entities) {
			spr := game.Sprite{
				Texture: "cursors.png",
			}
			actor := hud.mgr.Component(entities[i], "Actor").(*game.Actor)
			switch actor.Size {
			case game.SMALL:
				spr.X = 0
				spr.Y = 0
				spr.W = 24
				spr.H = 16
			case game.MEDIUM:
				spr.X = 0
				spr.Y = 32
				spr.W = 58
				spr.H = 32
			case game.LARGE:
				spr.X = 0
				spr.Y = 64
				spr.W = 58
				spr.H = 48
			}
			hud.mgr.AddComponent(slot, &spr)
			hud.mgr.AddComponent(slot, &game.Leash{
				Owner:       entities[i],
				LayerOffset: -1,
			})
		} else {
			// hide cursor
			hud.mgr.RemoveComponent(slot, &game.Sprite{})
			hud.mgr.RemoveComponent(slot, &game.Position{})
			hud.mgr.RemoveComponent(slot, &game.Leash{})
		}
	}
}

func (hud *HUD) showTimePassingIcon() {
	e := hud.mgr.AnyTagged(timePassingTag)
	if e != 0 {
		hud.mgr.Tag(e, invalidatedTag)
		return
	}
	e = hud.mgr.NewEntity()
	hud.mgr.Tag(e, timePassingTag)
	hud.mgr.Tag(e, invalidatedTag)
}

func (hud *HUD) hideTimePassingIcon() {
	e := hud.mgr.AnyTagged(timePassingTag)
	if e == 0 {
		return
	}
	hud.mgr.DestroyEntity(e)
}

func (hud *HUD) repaintTimePassingIcon() {
	e := hud.mgr.AnyTagged(timePassingTag)
	if e == 0 {
		return
	}
	hud.mgr.AddComponent(e, &game.Sprite{
		Texture: "hud.png",
		X:       16,
		Y:       0,
		W:       16,
		H:       24,
	})
	hud.mgr.AddComponent(e, &game.Scale{
		X: hud.scale,
		Y: hud.scale,
	})
	hud.mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: hud.centerX,
			Y: hud.centerY,
		},
		Layer:    hud.layer,
		Absolute: true,
	})
	hud.mgr.RemoveTag(e, invalidatedTag)
}

func (hud *HUD) showCurrentActor() {
	e := hud.mgr.AnyTagged(currentActorTag)
	if e != 0 {
		hud.mgr.Tag(e, invalidatedTag)
		return
	}
	parent := hud.mgr.NewEntity()
	hud.mgr.Tag(parent, currentActorTag)
	hud.mgr.Tag(parent, invalidatedTag)

	e = hud.mgr.NewEntity()
	hud.mgr.AddComponent(e, &ecs.Parent{Value: parent})
	hud.mgr.Tag(e, currentActorHovererTag)

	e = hud.mgr.NewEntity()
	hud.mgr.AddComponent(e, &ecs.Parent{Value: parent})
	hud.mgr.Tag(e, currentActorPortraitTag)

	e = hud.mgr.NewEntity()
	hud.mgr.AddComponent(e, &ecs.Parent{Value: parent})
	hud.mgr.Tag(e, currentActorNameTag)
}

func (hud *HUD) hideCurrentActor() {
	e := hud.mgr.AnyTagged(currentActorTag)
	if e == 0 {
		return
	}
	hud.mgr.DestroyEntity(e)
}

func (hud *HUD) repaintCurrentActor() {
	parent := hud.mgr.AnyTagged(currentActorHovererTag)
	if parent == 0 {
		return
	}

	// Repaint the hovering arrow that points to the current actor.
	e := hud.mgr.AnyTagged(currentActorHovererTag)
	ex, ok := hud.mgr.Single([]string{"Actor", "TurnToken"})
	if !ok {
		return
	}
	pos := hud.mgr.Component(ex, "Position").(*game.Position)

	hud.mgr.AddComponent(e, &game.Sprite{
		Texture: "hud.png",
		X:       32,
		Y:       0,
		W:       18,
		H:       5,
		OffsetY: -32,
	})
	hud.mgr.AddComponent(e, &game.Position{
		Center: pos.Center,
		Layer:  hud.layer,
	})
	hud.mgr.AddComponent(e, game.NewHoverAnimation())

	// Repaint the current Actor's portrait.
	e = hud.mgr.AnyTagged(currentActorPortraitTag)
	actor := hud.mgr.Component(ecs.Must(hud.mgr.Single([]string{"Actor", "TurnToken"})), "Actor").(*game.Actor)

	hud.mgr.AddComponent(e, &actor.BigIcon)

	hud.mgr.AddComponent(e, &game.Scale{
		X: hud.scale,
		Y: hud.scale,
	})
	hud.mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: 30 * hud.scale,
			Y: -30 * hud.scale,
		},
		Layer:    hud.layer,
		Absolute: true,
	})

	// Repaint the current actor's name.
	e = hud.mgr.AnyTagged(currentActorNameTag)
	hud.mgr.AddComponent(e, &game.Font{
		Text: actor.Name,
	})
	hud.mgr.AddComponent(e, &game.Scale{
		X: hud.scale,
		Y: hud.scale,
	})
	hud.mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: 58 * hud.scale,
			Y: -54 * hud.scale,
		},
		Layer:    hud.layer + 1,
		Absolute: true,
	})
	hud.mgr.RemoveTag(parent, invalidatedTag)
}

func (hud *HUD) showCurrentActorStats() {
	e := hud.mgr.AnyTagged(currentActorStatsTag)
	if e != 0 {
		hud.mgr.Tag(e, invalidatedTag)
		return
	}
	parent := hud.mgr.NewEntity()
	hud.mgr.Tag(parent, currentActorStatsTag)
	hud.mgr.Tag(parent, invalidatedTag)

	// Create enough child entities for every stat and its label.
	children := make([]ecs.Entity, 8)
	for i := 0; i < 8; i++ {
		e := hud.mgr.NewEntity()
		hud.mgr.AddComponent(e, &ecs.Parent{
			Value: parent,
		})
		children[i] = e
	}
	hud.mgr.AddComponent(parent, &ecs.Children{
		Value: children,
	})
}

func (hud *HUD) hideCurrentActorStats() {
	e := hud.mgr.AnyTagged(currentActorStatsTag)
	if e == 0 {
		return
	}
	hud.mgr.DestroyEntity(e)
}

func (hud *HUD) repaintCurrentActorStats() {
	parent := hud.mgr.AnyTagged(currentActorStatsTag)
	children := hud.mgr.Component(parent, "Children").(*ecs.Children)

	e := ecs.Must(hud.mgr.Single([]string{"Actor", "CombatStats", "TurnToken"}))
	actor := hud.mgr.Component(e, "Actor").(*game.Actor)
	stats := hud.mgr.Component(e, "CombatStats").(*game.CombatStats)
	labels := []string{
		"Health:",
		"?/N", // TODO
		"Energy:",
		"?/N", // TODO
		"Action:",
		fmt.Sprintf("%d/%d", stats.ActionPoints, actor.ActionPoints),
		"Prep:",
		fmt.Sprintf("%d/%d", stats.CurrentPreparation, actor.PreparationThreshold),
	}
	for i, child := range children.Value {
		// Magical x,y coordinates for stats are:
		// x = 60 or 102
		// y = -43, -36, -29, or -22
		x, y := 60.0+float64(i%2)*42.0, -43.0+float64(i/2*7)
		hud.mgr.AddComponent(child, &game.Position{
			Center: game.Center{
				X: x * hud.scale,
				Y: y * hud.scale,
			},
			Layer:    hud.layer,
			Absolute: true,
		})

		hud.mgr.AddComponent(child, &game.Scale{
			X: hud.scale,
			Y: hud.scale,
		})

		hud.mgr.AddComponent(child, &game.Font{
			Text: labels[i],
			Size: "small",
		})
	}
	hud.mgr.RemoveTag(parent, invalidatedTag)
}

func (hud *HUD) showSkills() {
	// Create a parent entity tagged with invalidatedTag and skillsTag
	e := hud.mgr.AnyTagged(skillsTag)
	if e != 0 {
		hud.mgr.Tag(e, invalidatedTag)
		return
	}
	e = hud.mgr.NewEntity()
	hud.mgr.Tag(e, skillsTag)

	// Create 7x2 child Entities with parent Components
	children := make([]ecs.Entity, 14)
	for i := 0; i < 14; i++ {
		children[i] = hud.mgr.NewEntity()
		hud.mgr.AddComponent(children[i], &ecs.Parent{
			Value: e,
		})
	}

	// Append all child entities to parent
	hud.mgr.AddComponent(e, &ecs.Children{
		Value: children,
	})
	hud.mgr.Tag(e, invalidatedTag)
}

func (hud *HUD) hideSkills() {
	e := hud.mgr.AnyTagged(skillsTag)
	if e == 0 {
		return
	}
	hud.mgr.DestroyEntity(e)
}

func (hud *HUD) repaintSkills() {
	parent := hud.mgr.AnyTagged(skillsTag)
	children := hud.mgr.Component(parent, "Children").(*ecs.Children)

	type skill struct {
		sprite      game.Sprite
		interactive *ui.Interactive
	}

	skills := map[int]skill{
		// Cancel button
		0: skill{
			sprite: game.Sprite{
				Texture: "hud.png",
				X:       208,
				Y:       0,
				W:       24,
				H:       24,
			},
			interactive: &ui.Interactive{
				Trigger: func() {
					hud.bus.Publish(&game.CancelSkillRequested{})
				},
			},
		},
	}

	if hud.lastCombatState == AwaitingInputState {
		skills = map[int]skill{
			// Move
			0: skill{
				sprite: game.Sprite{
					Texture: "hud.png",
					X:       232,
					Y:       24,
					W:       24,
					H:       24,
				},
				interactive: &ui.Interactive{
					Trigger: func() {
						hud.bus.Publish(&game.MoveModeRequested{})
					},
				},
			},

			// Consumables
			7: skill{
				sprite: game.Sprite{
					Texture: "hud.png",
					X:       232,
					Y:       0,
					W:       24,
					H:       24,
				},
			},

			// Flee
			6: skill{
				sprite: game.Sprite{
					Texture: "hud.png",
					X:       184,
					Y:       24,
					W:       24,
					H:       24,
				},
			},

			// End turn
			13: skill{
				sprite: game.Sprite{
					Texture: "hud.png",
					X:       208,
					Y:       24,
					W:       24,
					H:       24,
				},
				interactive: &ui.Interactive{
					Trigger: func() {
						hud.bus.Publish(&game.EndTurnRequested{})
					},
				},
			},
		}
	}

	for i, child := range children.Value {
		x, y := i%7, i/7

		hud.mgr.AddComponent(child, &game.Scale{
			X: hud.scale,
			Y: hud.scale,
		})
		hud.mgr.AddComponent(child, &game.Position{
			Center: game.Center{
				X: (150 + float64(26*x)) * hud.scale,
				Y: (-42 + float64(26*y)) * hud.scale,
			},
			Layer:    hud.layer,
			Absolute: true,
		})

		s, ok := skills[i]
		if !ok {
			// Add empty skill slot sprite component
			hud.mgr.AddComponent(child, &game.Sprite{
				Texture: "hud.png",
				X:       184,
				Y:       0,
				W:       24,
				H:       24,
			})
			continue
		}
		hud.mgr.AddComponent(child, &s.sprite)
		if s.interactive != nil {
			hud.mgr.AddComponent(child, s.interactive)
		}
	}
	hud.mgr.RemoveTag(parent, invalidatedTag)
}

const turnQueueSlots int = 8
const entitiesPerTurnQueueSlot int = 3

func (hud *HUD) showTurnQueue() {
	e := hud.mgr.AnyTagged(turnQueueTag)
	if e != 0 {
		hud.mgr.Tag(e, invalidatedTag)
		return
	}
	e = hud.mgr.NewEntity()
	hud.mgr.Tag(e, turnQueueTag)
	hud.mgr.Tag(e, invalidatedTag)

	children := make([]ecs.Entity, turnQueueSlots*entitiesPerTurnQueueSlot)
	for i := 0; i < turnQueueSlots*entitiesPerTurnQueueSlot; i++ {
		child := hud.mgr.NewEntity()
		hud.mgr.AddComponent(child, &ecs.Parent{
			Value: e,
		})
		children[i] = child
	}
	hud.mgr.AddComponent(e, &ecs.Children{
		Value: children,
	})
}

func (hud *HUD) hideTurnQueue() {
	e := hud.mgr.AnyTagged(turnQueueTag)
	if e == 0 {
		return
	}
	hud.mgr.DestroyEntity(e)
}

func (hud *HUD) repaintTurnQueue() {
	parent := hud.mgr.AnyTagged(turnQueueTag)
	if parent == 0 {
		return
	}
	type v struct {
		e            ecs.Entity
		remaining    int
		textureY     int
		current, max int
		icon         *game.Sprite
	}

	var q []v
	for _, e := range hud.mgr.Get([]string{"Actor", "CombatStats"}) {
		actor := hud.mgr.Component(e, "Actor").(*game.Actor)
		stats := hud.mgr.Component(e, "CombatStats").(*game.CombatStats)

		// An Awkward way of not including the Actor with the TurnToken.
		if stats.CurrentPreparation == 0 {
			continue
		}

		q = append(q, v{
			e:         e,
			remaining: actor.PreparationThreshold - stats.CurrentPreparation,
			current:   stats.CurrentPreparation,
			max:       actor.PreparationThreshold,
			icon:      &actor.SmallIcon,
		})
	}
	sort.Slice(q, func(i, j int) bool {
		return q[i].remaining < q[j].remaining
	})

	children := hud.mgr.Component(parent, "Children").(*ecs.Children)
	x, y := 10, 10+13
	stride := 42
	for i := 0; i < turnQueueSlots; i++ {
		// actor's portrait icon
		child := children.Value[i*3]

		// If we have no more Actors for this slot, then hide it.
		if i >= len(q) {
			hud.mgr.RemoveComponent(children.Value[i*3+0], &game.Sprite{})
			hud.mgr.RemoveComponent(children.Value[i*3+1], &game.Sprite{})
			hud.mgr.AddComponent(children.Value[i*3+2], &game.Font{})
			continue
		}

		v := q[i]
		hud.mgr.AddComponent(child, v.icon)
		hud.mgr.AddComponent(child, &game.Position{
			Center: game.Center{
				X: float64(13+x+i*stride) * hud.scale,
				Y: float64(y) * hud.scale,
			},
			Layer:    hud.layer,
			Absolute: true,
		})
		hud.mgr.AddComponent(child, &game.Scale{
			X: hud.scale,
			Y: hud.scale,
		})

		// current preparation progressbar
		prepPerc := float64(v.current) / float64(v.max)
		child = children.Value[i*3+1]
		hud.mgr.AddComponent(child, &ecs.Parent{
			Value: parent,
		})
		hud.mgr.AddComponent(child, &game.Sprite{
			Texture: "tranquility-plus-39-palette.png",
			X:       1,
			Y:       2,
			W:       1,
			H:       1,
		})
		hud.mgr.AddComponent(child, &game.Position{
			Center: game.Center{
				X: (13*prepPerc + float64(x+i*stride)) * hud.scale,
				Y: float64(y+13+2) * hud.scale,
			},
			Layer:    hud.layer + 1,
			Absolute: true,
		})
		hud.mgr.AddComponent(child, &game.Scale{
			X: hud.scale * 26 * prepPerc,
			Y: hud.scale * 4,
		})

		// current/max preparation text
		child = children.Value[i*3+2]
		hud.mgr.AddComponent(child, &ecs.Parent{
			Value: parent,
		})
		hud.mgr.AddComponent(child, &game.Font{
			Text: fmt.Sprintf("%d/%d", v.current, v.max),
			Size: "small",
		})
		hud.mgr.AddComponent(child, &game.Position{
			Center: game.Center{
				X: float64(x+i*stride) * hud.scale,
				Y: float64(y+13+2+4+2) * hud.scale,
			},
			Layer:    hud.layer,
			Absolute: true,
		})
		hud.mgr.AddComponent(child, &game.Scale{
			X: hud.scale,
			Y: hud.scale,
		})
	}
	hud.mgr.RemoveTag(parent, invalidatedTag)
}
