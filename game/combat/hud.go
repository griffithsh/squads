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
*/

var (
	combatHUDTag                  = "COMBAT_HUD"
	timePassingTag                = combatHUDTag + ".TIME_PASSING"
	currentParticipantTag         = combatHUDTag + ".CURRENT_PARTICIPANT"
	currentParticipantHovererTag  = currentParticipantTag + ".HOVERER"
	currentParticipantPortraitTag = currentParticipantTag + ".PORTRAIT"
	currentParticipantNameTag     = currentParticipantTag + ".NAME"
	currentParticipantStatsTag    = currentParticipantTag + ".STATS"
	skillsTag                     = currentParticipantTag + ".SKILLS"
	turnQueueTag                  = combatHUDTag + ".TURN_QUEUE"
	endTurnButtonTag              = combatHUDTag + ".END_TURN_BUTTON"

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

	// Whose turn is it?
	turnToken ecs.Entity

	dormant bool
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
	bus.Subscribe(StatModified{}.Type(), hud.handleCombatStatModified)
	bus.Subscribe(StateTransition{}.Type(), hud.handleCombatStateTransition)
	bus.Subscribe(ParticipantTurnChanged{}.Type(), hud.handleParticipantTurnChanged)

	return &hud
}

// Enable is the opposite of Disable. It shows anything that should be shown based on the current state.
func (hud *HUD) Enable() {
	if hud.lastCombatState == AwaitingInputState || hud.lastCombatState == SelectingTargetState {
		hud.showSkills()
		hud.showCurrentParticipant()
		hud.showCurrentParticipantStats()
	}

	if hud.lastCombatState == PreparingState {
		hud.showTimePassingIcon()
	}

	hud.showTurnQueue()

	hud.dormant = false
}

// Disable the hud and everything in it, and ignore events that would show parts of
// the hud until Enable() is called.
func (hud *HUD) Disable() {
	hud.hideSkills()
	hud.hideCurrentParticipant()
	hud.hideCurrentParticipantStats()
	hud.hideTimePassingIcon()
	hud.hideTurnQueue()

	hud.dormant = true
}

func (hud *HUD) handleWindowSizeChanged(e event.Typer) {
	wsc := e.(*game.WindowSizeChanged)
	hud.centerX, hud.centerY = float64(wsc.NewW)/2, float64(wsc.NewH)/2

	if hud.dormant {
		return
	}

	if hud.mgr.AnyTagged(timePassingTag) == 0 {
		return
	}

	hud.showTimePassingIcon()
}

func (hud *HUD) handleCombatBegan(event.Typer) {
	e := hud.mgr.NewEntity()
	hud.mgr.Tag(e, combatHUDTag)

	if hud.dormant {
		return
	}

	hud.showTurnQueue()
}

func (hud *HUD) handleCombatStatModified(ev event.Typer) {
	if hud.dormant {
		return
	}

	csm := ev.(*StatModified)

	switch csm.Stat {
	case game.PrepStat:
		hud.showTurnQueue()
	}
}

func (hud *HUD) handleCombatStateTransition(ev event.Typer) {
	if hud.dormant {
		return
	}

	cst := ev.(*StateTransition)
	hud.lastCombatState = cst.New.Value()
	n, o := cst.New.Value(), cst.Old.Value()

	if n == AwaitingInputState || n == SelectingTargetState {
		hud.showSkills()
		if o != AwaitingInputState && o != SelectingTargetState {
			hud.showCurrentParticipant()
			hud.showCurrentParticipantStats()
		}
	} else {
		hud.hideSkills()
		hud.hideCurrentParticipant()
		hud.hideCurrentParticipantStats()
	}

	if cst.New == PreparingState {
		hud.showTimePassingIcon()
	} else {
		hud.hideTimePassingIcon()
	}
}

func (hud *HUD) handleParticipantTurnChanged(t event.Typer) {
	ev := t.(*ParticipantTurnChanged)
	hud.turnToken = ev.Entity
}

// Update the HUD. Synchronise the current game state to the Entities that compose it.
func (hud *HUD) Update(elapsed time.Duration) {
	var e ecs.Entity

	e = hud.mgr.AnyTagged(timePassingTag)
	if e != 0 && hud.mgr.HasTag(e, invalidatedTag) {
		hud.repaintTimePassingIcon()
	}
	e = hud.mgr.AnyTagged(currentParticipantTag)
	if e != 0 && hud.mgr.HasTag(e, invalidatedTag) {
		hud.repaintCurrentParticipant()
	}
	e = hud.mgr.AnyTagged(currentParticipantStatsTag)
	if e != 0 && hud.mgr.HasTag(e, invalidatedTag) {
		hud.repaintCurrentParticipantStats()
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

func (hud *HUD) showCurrentParticipant() {
	e := hud.mgr.AnyTagged(currentParticipantTag)
	if e != 0 {
		hud.mgr.Tag(e, invalidatedTag)
		return
	}
	parent := hud.mgr.NewEntity()
	hud.mgr.Tag(parent, currentParticipantTag)
	hud.mgr.Tag(parent, invalidatedTag)

	e = hud.mgr.NewEntity()
	hud.mgr.AddComponent(e, &ecs.Parent{Value: parent})
	hud.mgr.Tag(e, currentParticipantHovererTag)

	e = hud.mgr.NewEntity()
	hud.mgr.AddComponent(e, &ecs.Parent{Value: parent})
	hud.mgr.Tag(e, currentParticipantPortraitTag)

	e = hud.mgr.NewEntity()
	hud.mgr.AddComponent(e, &ecs.Parent{Value: parent})
	hud.mgr.Tag(e, currentParticipantNameTag)
}

func (hud *HUD) hideCurrentParticipant() {
	e := hud.mgr.AnyTagged(currentParticipantTag)
	if e == 0 {
		return
	}
	hud.mgr.DestroyEntity(e)
}

func (hud *HUD) repaintCurrentParticipant() {
	parent := hud.mgr.AnyTagged(currentParticipantTag)
	if parent == 0 {
		return
	}

	// Repaint the hovering arrow that points to the current Character.
	ex := hud.turnToken
	if ex == 0 || hud.mgr.Component(ex, "Position") == nil {
		return
	}

	pos := hud.mgr.Component(ex, "Position").(*game.Position)

	e := hud.mgr.AnyTagged(currentParticipantHovererTag)
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

	// Repaint the current Participant's portrait.
	e = hud.mgr.AnyTagged(currentParticipantPortraitTag)
	participant := hud.mgr.Component(hud.turnToken, "Participant").(*Participant)

	hud.mgr.AddComponent(e, &participant.BigIcon)

	hud.mgr.AddComponent(e, &game.Scale{
		X: hud.scale,
		Y: hud.scale,
	})
	hud.mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: 30 * hud.scale,
			Y: (hud.centerY * 2) - 30*hud.scale,
		},
		Layer:    hud.layer,
		Absolute: true,
	})

	// Repaint the current Participant's name.
	e = hud.mgr.AnyTagged(currentParticipantNameTag)
	hud.mgr.AddComponent(e, &game.Font{
		Text: participant.Name,
	})
	hud.mgr.AddComponent(e, &game.Scale{
		X: hud.scale,
		Y: hud.scale,
	})
	hud.mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: 58 * hud.scale,
			Y: (hud.centerY * 2) - 54*hud.scale,
		},
		Layer:    hud.layer + 1,
		Absolute: true,
	})
	hud.mgr.RemoveTag(parent, invalidatedTag)
}

func (hud *HUD) showCurrentParticipantStats() {
	e := hud.mgr.AnyTagged(currentParticipantStatsTag)
	if e != 0 {
		hud.mgr.Tag(e, invalidatedTag)
		return
	}
	parent := hud.mgr.NewEntity()
	hud.mgr.Tag(parent, currentParticipantStatsTag)
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

func (hud *HUD) hideCurrentParticipantStats() {
	e := hud.mgr.AnyTagged(currentParticipantStatsTag)
	if e == 0 {
		return
	}
	hud.mgr.DestroyEntity(e)
}

func (hud *HUD) repaintCurrentParticipantStats() {
	parent := hud.mgr.AnyTagged(currentParticipantStatsTag)
	children := hud.mgr.Component(parent, "Children").(*ecs.Children)

	e := hud.turnToken
	participant := hud.mgr.Component(e, "Participant").(*Participant)
	labels := []string{
		"Health:",
		"?/N", // TODO
		"Energy:",
		"?/N", // TODO
		"Action:",
		fmt.Sprintf("%d/%d", participant.ActionPoints.Cur, participant.ActionPoints.Max),
		"Prep:",
		fmt.Sprintf("%d/%d", participant.PreparationThreshold.Cur, participant.PreparationThreshold.Max),
	}
	for i, child := range children.Value {
		// Magical x,y coordinates for stats are:
		// x = 60 or 102
		// y = -43, -36, -29, or -22
		x, y := 60.0+float64(i%2)*42.0, -43.0+float64(i/2*7)
		hud.mgr.AddComponent(child, &game.Position{
			Center: game.Center{
				X: x * hud.scale,
				Y: (hud.centerY * 2) + y*hud.scale,
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
				W: 24, H: 24,
				Trigger: func(x, y float64) {
					hud.bus.Publish(&CancelSkillRequested{})
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
					W: 24, H: 24,
					Trigger: func(x, y float64) {
						hud.bus.Publish(&SkillRequested{
							Code: game.BasicMovement,
						})
					},
				},
			},

			// weapon 1
			1: skill{
				sprite: game.Sprite{
					Texture: "hud.png",
					X:       160,
					Y:       0,
					W:       24,
					H:       24,
				},
				interactive: &ui.Interactive{
					W: 24, H: 24,
					Trigger: func(x, y float64) {
						hud.bus.Publish(&SkillRequested{
							Code: game.BasicAttack,
						})
					},
				},
			},
			// 2 - weapon 2
			// 8 - weapon 3
			// 9 - weapon 4

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
				interactive: &ui.Interactive{
					W: 24, H: 24,
					Trigger: func(x, y float64) {
						hud.bus.Publish(&AttemptingEscape{Entity: hud.turnToken})
					},
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
					W: 24, H: 24,
					Trigger: func(x, y float64) {
						hud.bus.Publish(&EndTurnRequested{})
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
				Y: hud.centerY*2 + (-42+float64(26*y))*hud.scale,
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
	parents := hud.mgr.Tagged(turnQueueTag)
	if len(parents) > 1 {
		panic("incorrect use of " + turnQueueTag + ": multiple found")
	}
	parent := parents[0]
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
	for _, e := range hud.mgr.Get([]string{"Participant"}) {
		participant := hud.mgr.Component(e, "Participant").(*Participant)

		// Knocked Down and Escaped Characters cannot take turns.
		if participant.Status != Alive {
			continue
		}

		// An awkward way of not including the Character with the TurnToken.
		if participant.PreparationThreshold.Cur == 0 {
			continue
		}

		q = append(q, v{
			e:         e,
			remaining: participant.PreparationThreshold.Max - participant.PreparationThreshold.Cur,
			current:   participant.PreparationThreshold.Cur,
			max:       participant.PreparationThreshold.Max,
			icon:      &participant.SmallIcon,
		})
	}
	sort.Slice(q, func(i, j int) bool {
		return q[i].remaining < q[j].remaining
	})

	children := hud.mgr.Component(parent, "Children").(*ecs.Children)
	x, y := 10, 10+13
	stride := 42
	for i := 0; i < turnQueueSlots; i++ {
		// participant's portrait icon
		child := children.Value[i*3]

		// If we have no more Characters for this slot, then hide it.
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
