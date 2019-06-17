package combat

import (
	"fmt"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/ui"
)

/*
HUD element hierarchy design:

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
	EndTurnButton
*/

var (
	combatHUDTag     = "COMBAT_HUD"
	currentActorTag  = combatHUDTag + ".CURRENT_ACTOR"
	hpLabelTag       = currentActorTag + ".HPLabelTag"
	energyLabelTag   = currentActorTag + ".HPLabelTag"
	actionLabelTag   = currentActorTag + ".HPLabelTag"
	prepLabelTag     = currentActorTag + ".HPLabelTag"
	turnQueueTag     = combatHUDTag + ".TURN_QUEUE"
	endTurnButtonTag = combatHUDTag + ".END_TURN_BUTTON"
)

// HUD is a heads up display for the combat game state.
type HUD struct {
	mgr   *ecs.World
	bus   *event.Bus
	scale float64
	layer int
}

// NewHUD construct a HUD.
func NewHUD(mgr *ecs.World, bus *event.Bus) *HUD {
	hud := HUD{
		mgr:   mgr,
		bus:   bus,
		scale: 2,
		layer: 100,
	}

	bus.Subscribe(event.CombatBegunType, hud.handleCombatBegan)
	bus.Subscribe(event.AwaitingPlayerInputType, hud.handleAwaitingInput)
	return &hud
}

func (hud *HUD) handleCombatBegan(event.Typer) {
	hud.create()
}

func (hud *HUD) handleAwaitingInput(event.Typer) {
	hud.destroyCurrentActor()
	hud.createCurrentActor(hud.mgr.AnyTagged(combatHUDTag))
}

// create entire HUD hierarchy of ui elements.
func (hud *HUD) create() {
	e := hud.mgr.NewEntity()
	hud.mgr.Tag(e, combatHUDTag)

	hud.createCurrentActor(e)
	// hud.createTurnQueue(e) // TODO
	hud.createEndTurnButton(e)
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
		W:       96,
		H:       96,
	})
	hud.mgr.AddComponent(e, &game.Scale{
		X: hud.scale,
		Y: hud.scale,
	})
	hud.mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: 112,
			Y: -112,
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
	label("Health:", 110, -92)
	label("Energy:", 110, -80)
	label("Action:", 110, -68)
	label("Prep:", 110, -56)

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

	labelRightTagged("10/110", 170, -92, hpLabelTag)
	labelRightTagged("70/110", 170, -80, energyLabelTag)

	actionLabel := fmt.Sprintf("%d/%d", stats.ActionPoints, actor.ActionPoints)
	labelRightTagged(actionLabel, 170, -68, actionLabelTag)

	prepLabel := fmt.Sprintf("%d/%d", stats.CurrentPreparation, actor.PreparationThreshold)
	labelRightTagged(prepLabel, 170, -56, prepLabelTag)

}

func (hud *HUD) createSkills(parent ecs.Entity) {}

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

func (hud *HUD) createEndTurnButton(parent ecs.Entity) {
	e := hud.mgr.NewEntity()
	hud.mgr.Tag(e, endTurnButtonTag)
	hud.mgr.AddComponent(e, &ecs.Parent{
		Value: parent,
	})

	hud.mgr.AddComponent(e, &game.Sprite{
		Texture: "hud.png",
		X:       16,
		Y:       0,
		W:       46,
		H:       14,
	})
	hud.mgr.AddComponent(e, &game.Scale{
		X: hud.scale,
		Y: hud.scale,
	})
	hud.mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: -80,
			Y: 32,
		},
		Layer:    hud.layer,
		Absolute: true,
	})

	hud.mgr.AddComponent(e, &ui.Interactive{
		Trigger: func() {
			hud.bus.Publish(&event.EndTurnRequested{})
		},
	})
}
