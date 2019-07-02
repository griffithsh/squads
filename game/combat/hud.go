package combat

import (
	"fmt"
	"sort"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/ui"
)

/*
HUD element hierarchy design:

TimePassingIcon
CombatHUD
	CurrentActor
		Portrait
		Name
		Health
		Mana
		Action
		Preparation
		Skills
			[]
				Icon
	TurnQueue
		[]
			Portrait
			Prep
*/

var (
	combatHUDTag     = "COMBAT_HUD"
	timePassingTag   = combatHUDTag + ".TIME_PASSING"
	currentActorTag  = combatHUDTag + ".CURRENT_ACTOR"
	hpLabelTag       = currentActorTag + ".HP"
	energyLabelTag   = currentActorTag + ".ENERGY"
	actionLabelTag   = currentActorTag + ".ACTION"
	prepLabelTag     = currentActorTag + ".PREPARATION"
	skillsTag        = currentActorTag + ".SKILLS"
	turnQueueTag     = combatHUDTag + ".TURN_QUEUE"
	endTurnButtonTag = combatHUDTag + ".END_TURN_BUTTON"
)

// HUD is a heads up display for the combat game state.
type HUD struct {
	mgr              *ecs.World
	bus              *event.Bus
	scale            float64
	layer            int
	centerX, centerY float64 // center of the game's window, or half the width and height.
}

// NewHUD construct a HUD.
func NewHUD(mgr *ecs.World, bus *event.Bus, screenX int, screenY int) *HUD {
	hud := HUD{
		mgr:     mgr,
		bus:     bus,
		scale:   2,
		layer:   100,
		centerX: float64(screenX) / 2,
		centerY: float64(screenY) / 2,
	}

	bus.Subscribe(game.WindowSizeChanged{}.Type(), hud.handleWindowSizeChanged)
	bus.Subscribe(game.CombatBegan{}.Type(), hud.handleCombatBegan)
	bus.Subscribe(StateTransition{}.Type(), hud.handleCombatStateTransition)
	bus.Subscribe(game.CombatStatModified{}.Type(), hud.handleCombatStatModified)

	return &hud
}

func (hud *HUD) handleWindowSizeChanged(e event.Typer) {
	wsc := e.(*game.WindowSizeChanged)
	hud.centerX, hud.centerY = float64(wsc.NewW)/2, float64(wsc.NewH)/2

	if hud.mgr.AnyTagged(timePassingTag) == 0 {
		return
	}

	hud.destroyTimePassingIcon()
	hud.createTimePassingIcon()
}

func (hud *HUD) handleCombatBegan(event.Typer) {
	hud.create()
}

func (hud *HUD) handleCombatStateTransition(ev event.Typer) {
	cst := ev.(*StateTransition)

	// when we are awaiting input, then we should just create the current
	// actor, because destroy should already have happened.
	if cst.New == AwaitingInputState {
		hud.createCurrentActor(hud.mgr.AnyTagged(combatHUDTag))
	} else {
		hud.destroyCurrentActor()
	}

	if cst.New == PreparingState {
		hud.createTimePassingIcon()
	} else {
		hud.destroyTimePassingIcon()
	}
}

func (hud *HUD) handleCombatStatModified(ev event.Typer) {
	csm := ev.(*game.CombatStatModified)

	switch csm.Stat {
	case game.PrepStat:
		hud.destroyTurnQueue()
		hud.createTurnQueue(hud.mgr.AnyTagged(combatHUDTag))
	}

}

// create entire HUD hierarchy of ui elements.
func (hud *HUD) create() {
	e := hud.mgr.NewEntity()
	hud.mgr.Tag(e, combatHUDTag)

	hud.createCurrentActor(e)
	hud.createTurnQueue(e)
}

func (hud *HUD) createPortrait(parent ecs.Entity) {
	e := hud.mgr.NewEntity()
	hud.mgr.AddComponent(e, &ecs.Parent{
		Value: parent,
	})

	hud.mgr.AddComponent(e, &game.Sprite{
		Texture: "hud.png",
		X:       0,
		Y:       24,
		W:       52,
		H:       52,
	})
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
}

func (hud *HUD) createName(parent ecs.Entity) {
	e, ok := hud.mgr.Single([]string{"TurnToken", "Actor"})
	if !ok {
		return
	}
	actor := hud.mgr.Component(e, "Actor").(*game.Actor)
	e = hud.mgr.NewEntity()
	hud.mgr.AddComponent(e, &ecs.Parent{
		Value: parent,
	})

	hud.mgr.AddComponent(e, &game.Font{
		Text: actor.Name,
	})
	hud.mgr.AddComponent(e, &game.Scale{
		X: hud.scale,
		Y: hud.scale,
	})
	hud.mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: 12,
			Y: -32,
		},
		Layer:    hud.layer + 1,
		Absolute: true,
	})
}

func (hud *HUD) createStats(parent ecs.Entity) {
	label := func(text string, x, y float64) ecs.Entity {
		e := hud.mgr.NewEntity()
		hud.mgr.AddComponent(e, &ecs.Parent{
			Value: parent,
		})

		hud.mgr.AddComponent(e, &game.Font{
			Text: text,
		})
		hud.mgr.AddComponent(e, &game.Scale{
			X: hud.scale,
			Y: hud.scale,
		})
		hud.mgr.AddComponent(e, &game.Position{
			Center: game.Center{
				X: x * hud.scale,
				Y: y * hud.scale,
			},
			Layer:    hud.layer,
			Absolute: true,
		})
		return e
	}
	label("Health:", 58, -52)
	label("Energy:", 58, -40)
	label("Action:", 58, -28)
	label("Prep:", 58, -16)

	labelRightTagged := func(text string, x, y float64, tag string) {
		e := label("", x, y) // empty string because we're stomping it below
		hud.mgr.Tag(e, tag)
		hud.mgr.AddComponent(e, &game.Font{
			Text: text,
			// TODO: Style should be right aligned!
			// TODO: Some width? 32?
		})
	}

	e, ok := hud.mgr.Single([]string{"TurnToken", "Actor", "CombatStats"})
	if !ok {
		return
	}
	actor := hud.mgr.Component(e, "Actor").(*game.Actor)
	stats := hud.mgr.Component(e, "CombatStats").(*game.CombatStats)

	labelRightTagged("?/N", 100, -52, hpLabelTag)
	labelRightTagged("?/N", 100, -40, energyLabelTag)

	actionLabel := fmt.Sprintf("%d/%d", stats.ActionPoints, actor.ActionPoints)
	labelRightTagged(actionLabel, 100, -28, actionLabelTag)

	prepLabel := fmt.Sprintf("%d/%d", stats.CurrentPreparation, actor.PreparationThreshold)
	labelRightTagged(prepLabel, 100, -16, prepLabelTag)
}

var skillsOffset = struct {
	X, Y float64
}{
	150, -42,
}

func (hud *HUD) createSkillsTargetSelectionMode(parent ecs.Entity) {
	// Cancel button always appears in the first skill slot.
	x, y := 0, 0

	e := hud.mgr.NewEntity()
	hud.mgr.Tag(e, skillsTag)
	hud.mgr.AddComponent(e, &ecs.Parent{
		Value: parent,
	})
	hud.mgr.AddComponent(e, &game.Sprite{
		Texture: "hud.png",
		X:       208,
		Y:       0,
		W:       24,
		H:       24,
	})
	hud.mgr.AddComponent(e, &game.Scale{
		X: hud.scale,
		Y: hud.scale,
	})
	hud.mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: (skillsOffset.X + float64(26*x)) * hud.scale,
			Y: (skillsOffset.Y + float64(26*y)) * hud.scale,
		},
		Layer:    hud.layer,
		Absolute: true,
	})

	hud.mgr.AddComponent(e, &ui.Interactive{
		Trigger: func() {
			// hud.bus.Publish(&event.CancelTargetSelection{})
		},
	})
}

func (hud *HUD) createSkills(parent ecs.Entity) {
	// Two rows of skills, showing a mix of default and personal skills of the actor:
	// 0,0. Move.
	// 0,1. Consumables - pops a selection modal
	// 1-5,0-1. Configurable skills
	// 6,0. Flee.
	// 6,1. End Turn.

	for y := 0; y < 2; y++ {
		for x := 0; x < 7; x++ {
			e := hud.mgr.NewEntity()
			hud.mgr.Tag(e, skillsTag)
			hud.mgr.AddComponent(e, &ecs.Parent{
				Value: parent,
			})
			hud.mgr.AddComponent(e, &game.Scale{
				X: hud.scale,
				Y: hud.scale,
			})
			hud.mgr.AddComponent(e, &game.Position{
				Center: game.Center{
					X: (skillsOffset.X + float64(26*x)) * hud.scale,
					Y: (skillsOffset.Y + float64(26*y)) * hud.scale,
				},
				Layer:    hud.layer,
				Absolute: true,
			})

			var spr = game.Sprite{
				Texture: "hud.png",
				X:       184,
				Y:       0,
				W:       24,
				H:       24,
			}
			var trigger = ui.Interactive{
				Trigger: func() {
					// TODO Publish ActorRequestedSkill{ skillid:?}
				},
			}
			if x == 0 && y == 0 {
				// Move
				spr = game.Sprite{
					Texture: "hud.png",
					X:       232,
					Y:       24,
					W:       24,
					H:       24,
				}
				// trigger = ui.Interactive{
				// 	Trigger: func() {
				// 		hud.bus.Publish(&event.MoveModeRequested{})
				// 	},
				// }
			} else if x == 0 && y == 1 {
				// Consumables
				spr = game.Sprite{
					Texture: "hud.png",
					X:       232,
					Y:       0,
					W:       24,
					H:       24,
				}
				// trigger = ui.Interactive{
				// 	Trigger: func() {
				// 		hud.bus.Publish(&event.ViewConsumablesRequested{})
				// 	},
				// }
			} else if x == 6 && y == 0 {
				// Flee
				spr = game.Sprite{
					Texture: "hud.png",
					X:       184,
					Y:       24,
					W:       24,
					H:       24,
				}
				// trigger = ui.Interactive{
				// 	Trigger: func() {
				// 		hud.bus.Publish(&event.FleeRequested{})
				// 	},
				// }
			} else if x == 6 && y == 1 {
				// End Turn
				spr = game.Sprite{
					Texture: "hud.png",
					X:       208,
					Y:       24,
					W:       24,
					H:       24,
				}
				trigger = ui.Interactive{
					Trigger: func() {
						hud.bus.Publish(&game.EndTurnRequested{})
					},
				}
			}
			hud.mgr.AddComponent(e, &spr)
			hud.mgr.AddComponent(e, &trigger)
		}
	}
}

func (hud *HUD) destroyCurrentActor() {
	e := hud.mgr.AnyTagged(currentActorTag)
	hud.mgr.DestroyEntity(e)
}

func (hud *HUD) createCurrentActor(parent ecs.Entity) {
	if _, ok := hud.mgr.Single([]string{"Actor", "TurnToken"}); !ok {
		return
	}

	e := hud.mgr.NewEntity()
	hud.mgr.Tag(e, currentActorTag)
	hud.mgr.AddComponent(e, &ecs.Parent{
		Value: parent,
	})

	hud.createPortrait(e)
	hud.createName(e)
	hud.createStats(e)
	hud.createSkills(e)
}

func (hud *HUD) createTurnQueue(parent ecs.Entity) {
	e := hud.mgr.NewEntity()
	hud.mgr.Tag(e, turnQueueTag)
	hud.mgr.AddComponent(e, &ecs.Parent{
		Value: parent,
	})

	type v struct {
		e            ecs.Entity
		remaining    int
		textureY     int
		current, max int
	}

	var q []v
	for _, e := range hud.mgr.Get([]string{"Actor", "CombatStats"}) {
		actor := hud.mgr.Component(e, "Actor").(*game.Actor)
		stats := hud.mgr.Component(e, "CombatStats").(*game.CombatStats)
		if stats.CurrentPreparation == 0 {
			continue
		}
		textureY := 0
		if actor.Size == game.MEDIUM {
			textureY = 52
		} else if actor.Size == game.LARGE {
			textureY = 104
		}
		q = append(q, v{
			e:         e,
			remaining: actor.PreparationThreshold - stats.CurrentPreparation,
			textureY:  textureY,
			current:   stats.CurrentPreparation,
			max:       actor.PreparationThreshold,
		})
	}
	sort.Slice(q, func(i, j int) bool {
		return q[i].remaining < q[j].remaining
	})

	x, y := 10, 10+13
	stride := 42
	for i, v := range q {
		// actor's small icon
		child := hud.mgr.NewEntity()
		hud.mgr.AddComponent(child, &ecs.Parent{
			Value: e,
		})
		hud.mgr.AddComponent(child, &game.Sprite{
			Texture: "hud.png",
			X:       v.textureY,
			Y:       76,
			W:       26,
			H:       26,
		})
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

		// current prep
		prepPerc := float64(v.current) / float64(v.max)
		child = hud.mgr.NewEntity()
		hud.mgr.AddComponent(child, &ecs.Parent{
			Value: e,
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

		// current/max prep text
		child = hud.mgr.NewEntity()
		hud.mgr.AddComponent(child, &ecs.Parent{
			Value: e,
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
}

func (hud *HUD) destroyTurnQueue() {
	hud.mgr.DestroyEntity(hud.mgr.AnyTagged(turnQueueTag))
}

func (hud *HUD) destroyTimePassingIcon() {
	e := hud.mgr.AnyTagged(timePassingTag)
	hud.mgr.DestroyEntity(e)
}

func (hud *HUD) createTimePassingIcon() {
	e := hud.mgr.NewEntity()
	hud.mgr.Tag(e, timePassingTag)
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

}
