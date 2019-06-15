package combat

import (
	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/ui"
)

var (
	combatHUDTag            = "COMBAT_HUD"
	currentActorControlsTag = combatHUDTag + ".CURRENT_ACTOR"
	turnQueueTag            = combatHUDTag + ".TURN_QUEUE"
	endTurnButtonTag        = combatHUDTag + ".END_TURN_BUTTON"
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
	return &hud
}

func (hud *HUD) handleCombatBegan(event.Typer) {
	hud.create()
}

// create entire HUD hierarchy of ui elements.
func (hud *HUD) create() {
	// TODO

	hud.createEndTurnButton()
}

func (hud *HUD) createEndTurnButton() {
	e := hud.mgr.NewEntity()
	hud.mgr.Tag(e, endTurnButtonTag)
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

/*
CombatHUD
	CurrentActorControls
		Portrait
		Name
		Stats
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

// 1. On construct, create everything

// 2. When *anything* changes recreate *everything*.

// 3. Start breaking down create functions into composable pieces, so not
// everything needs to be recreated all the time.
