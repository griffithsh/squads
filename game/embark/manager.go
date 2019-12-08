package embark

import (
	"fmt"
	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/ui"
)

// Manager holds state and provides methods to control that state for an embark
// screen. This screen allows the player to configure the Characters they will
// start their run with.
type Manager struct {
	mgr *ecs.World
	bus *event.Bus

	screenW, screenH int
}

// NewManager creates a new Manager in a default state. You should call Begin to start the Manager.
func NewManager(mgr *ecs.World, bus *event.Bus) *Manager {
	em := Manager{
		mgr: mgr,
		bus: bus,
	}

	bus.Subscribe(game.WindowSizeChanged{}.Type(), em.handleWindowSizeChanged)
	return &em
}

// Begin an embark Manager, setting up Entities required to display and interact
// with the embark screen.
func (em *Manager) Begin() {
	// TODO:

	// Create a button to press to embark
	e := em.mgr.NewEntity()
	em.mgr.Tag(e, "embark")

	em.mgr.AddComponent(e, &game.Sprite{
		Texture: "embark-button.png",

		X: 0, Y: 0,
		W: 64, H: 64,
	})
	fmt.Println("embark.Manager::Begin", em)
	em.mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: 64 + 12,
			Y: 64 + 12,
		},
		Layer:    10,
		Absolute: true,
	})
	em.mgr.AddComponent(e, &ui.Interactive{
		W: 48, H: 48,
		Trigger: func(x, y float64) {
			fmt.Println("embarking!")
			em.setupSquad()
			em.bus.Publish(&SquadSelected{})
		},
	})
}

// End an embark Manager, resetting it to a default state.
func (em *Manager) End() {
	for _, e := range em.mgr.Tagged("embark") {
		em.mgr.DestroyEntity(e)
	}
}

func (em *Manager) handleWindowSizeChanged(e event.Typer) {
	fmt.Println("handleWindowSizeChanged", e)
	wsc := e.(*game.WindowSizeChanged)
	em.screenW, em.screenH = wsc.NewW, wsc.NewH
}

func (em *Manager) setupSquad() {
	// Create a Squad Entity.
	e := em.mgr.NewEntity()
	em.mgr.Tag(e, "player")
	em.mgr.AddComponent(e, &game.Squad{})
	squad := em.mgr.Component(e, "Squad").(*game.Squad)
	players := game.NewTeam()
	em.mgr.AddComponent(e, players)

	// Create Characters to Populate the player's Squad.
	e = em.mgr.NewEntity()
	em.mgr.AddComponent(e, players)
	squad.Members = append(squad.Members, e)
	em.mgr.AddComponent(e, &game.Character{
		Name:                 "Samithee",
		Size:                 game.SMALL,
		Sex:                  game.Male,
		Profession:           game.Villager,
		PreparationThreshold: 701,
		ActionPoints:         100,
		SmallIcon: game.Sprite{
			Texture: "hud.png",
			X:       0,
			Y:       76,
			W:       26,
			H:       26,
		},
		BigIcon: game.Sprite{
			Texture: "hud.png",
			X:       0,
			Y:       24,
			W:       52,
			H:       52,
		},
	})

	e = em.mgr.NewEntity()
	em.mgr.AddComponent(e, players)
	squad.Members = append(squad.Members, e)
	em.mgr.AddComponent(e, &game.Character{
		Name:                 "Timjamen",
		Size:                 game.SMALL,
		Sex:                  game.Male,
		Profession:           game.Villager,
		PreparationThreshold: 699,
		ActionPoints:         100,

		SmallIcon: game.Sprite{
			Texture: "portraits.png",
			X:       178,
			Y:       230,
			W:       26,
			H:       26,
		},
		BigIcon: game.Sprite{
			Texture: "portraits.png",
			X:       204,
			Y:       204,
			W:       52,
			H:       52,
		},
	})
}
