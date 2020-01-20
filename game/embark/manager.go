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

	// prepared lists Characters who are about to embark.
	prepared []ecs.Entity
	// villagers lists Characters who could be selected for embarking.
	villagers []ecs.Entity
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
	em.rollVillagers()
	em.repaint()
}

// End an embark Manager, resetting it to a default state.
func (em *Manager) End() {
	for _, e := range em.mgr.Tagged("embark") {
		em.mgr.DestroyEntity(e)
	}
}

func (em *Manager) handleWindowSizeChanged(e event.Typer) {
	wsc := e.(*game.WindowSizeChanged)
	em.screenW, em.screenH = wsc.NewW, wsc.NewH
}

// repaint synchronises the renderable Components to the Characters in
// em.prepared and em.villagers. It should be called after a change is made to
// who will embark.
func (em *Manager) repaint() {
	// Destroy all existing Entities used to render this Character.
	for _, e := range em.mgr.Tagged("embark-characters") {
		em.mgr.DestroyEntity(e)
	}

	for _, e := range em.mgr.Tagged("embark-hud") {
		em.mgr.DestroyEntity(e)
	}

	lMargin := 64.0
	tMargin := 16.0
	sheetW := 144.0 // sheetW is the width of each character sheet
	for i, villager := range em.villagers {
		char := em.mgr.Component(villager, "Character").(*game.Character)

		container := em.mgr.NewEntity()
		em.mgr.Tag(container, "embark")
		em.mgr.Tag(container, "embark-characters")

		var e ecs.Entity

		// Name
		e = em.mgr.NewEntity()
		em.mgr.Dependency(container, e)
		em.mgr.AddComponent(e, &game.Font{
			Text: char.Name,
		})
		em.mgr.AddComponent(e, &game.Position{
			Center: game.Center{
				X: float64(i)*sheetW + lMargin,
				Y: tMargin,
			},
			Layer: 100,
		})

		// Icon
		e = em.mgr.NewEntity()
		em.mgr.Dependency(container, e)
		em.mgr.AddComponent(e, &char.BigIcon)
		em.mgr.AddComponent(e, &game.Position{
			Center: game.Center{
				X: float64(i)*sheetW + lMargin + sheetW - 56/2 - 16,
				Y: tMargin + 56/2,
			},
			Layer: 100,
		})

		// Profession
		e = em.mgr.NewEntity()
		em.mgr.Dependency(container, e)
		em.mgr.AddComponent(e, &game.Font{
			Text: char.Profession.String(),
			Size: "small",
		})
		em.mgr.AddComponent(e, &game.Position{
			Center: game.Center{
				X: float64(i)*sheetW + lMargin,
				Y: tMargin + 12,
			},
			Layer: 100,
		})

		// Level
		e = em.mgr.NewEntity()
		em.mgr.Dependency(container, e)
		em.mgr.AddComponent(e, &game.Font{
			Text: fmt.Sprintf("Level: %d", char.Level),
			Size: "small",
		})
		em.mgr.AddComponent(e, &game.Position{
			Center: game.Center{
				X: float64(i)*sheetW + lMargin,
				Y: tMargin + 12 + 8,
			},
			Layer: 100,
		})

		// Sex
		e = em.mgr.NewEntity()
		em.mgr.Dependency(container, e)
		em.mgr.AddComponent(e, &game.Font{
			Text: fmt.Sprintf("Sex: %s", char.Sex),
			Size: "small",
		})
		em.mgr.AddComponent(e, &game.Position{
			Center: game.Center{
				X: float64(i)*sheetW + lMargin,
				Y: tMargin + 12 + 16,
			},
			Layer: 100,
		})

		// Prep
		e = em.mgr.NewEntity()
		em.mgr.Dependency(container, e)
		em.mgr.AddComponent(e, &game.Font{
			Text: fmt.Sprintf("Preparation: %d", char.PreparationThreshold),
			Size: "small",
		})
		em.mgr.AddComponent(e, &game.Position{
			Center: game.Center{
				X: float64(i)*sheetW + lMargin,
				Y: tMargin + 12 + 24,
			},
			Layer: 100,
		})

		// Str per level
		e = em.mgr.NewEntity()
		em.mgr.Dependency(container, e)
		em.mgr.AddComponent(e, &game.Font{
			Text: fmt.Sprintf("Str/Lvl: %.2f", char.StrengthPerLevel),
			Size: "small",
		})
		em.mgr.AddComponent(e, &game.Position{
			Center: game.Center{
				X: float64(i)*sheetW + lMargin,
				Y: tMargin + 12 + 32,
			},
			Layer: 100,
		})
		// Agi per level
		e = em.mgr.NewEntity()
		em.mgr.Dependency(container, e)
		em.mgr.AddComponent(e, &game.Font{
			Text: fmt.Sprintf("Agi/Lvl: %.2f", char.AgilityPerLevel),
			Size: "small",
		})
		em.mgr.AddComponent(e, &game.Position{
			Center: game.Center{
				X: float64(i)*sheetW + lMargin,
				Y: tMargin + 12 + 40,
			},
			Layer: 100,
		})

		// Int per level
		e = em.mgr.NewEntity()
		em.mgr.Dependency(container, e)
		em.mgr.AddComponent(e, &game.Font{
			Text: fmt.Sprintf("Int/Lvl: %.2f", char.IntelligencePerLevel),
			Size: "small",
		})
		em.mgr.AddComponent(e, &game.Position{
			Center: game.Center{
				X: float64(i)*sheetW + lMargin,
				Y: tMargin + 12 + 48,
			},
			Layer: 100,
		})

		// Vit per level
		e = em.mgr.NewEntity()
		em.mgr.Dependency(container, e)
		em.mgr.AddComponent(e, &game.Font{
			Text: fmt.Sprintf("Vit/Lvl: %.2f", char.VitalityPerLevel),
			Size: "small",
		})
		em.mgr.AddComponent(e, &game.Position{
			Center: game.Center{
				X: float64(i)*sheetW + lMargin,
				Y: tMargin + 12 + 56,
			},
			Layer: 100,
		})

		// Masteries
		used := 0
		for j := 0; j < 11; j++ {
			m := game.Mastery(j)
			mastery := char.Masteries[m]
			if mastery == 0 {
				continue
			}

			e = em.mgr.NewEntity()
			em.mgr.Dependency(container, e)
			em.mgr.AddComponent(e, &game.Font{
				Text: fmt.Sprintf("%s: %d", m.String(), mastery),
				Size: "small",
			})
			em.mgr.AddComponent(e, &game.Position{
				Center: game.Center{
					X: float64(i)*sheetW + lMargin,
					Y: tMargin + 80 + float64(used)*8,
				},
				Layer: 100,
			})

			used++
		}

		// embark button
		e = em.mgr.NewEntity()
		em.mgr.Dependency(container, e)
		em.mgr.AddComponent(e, &game.Sprite{
			Texture: "embark-button.png",

			X: 0, Y: 0,
			W: 64, H: 64,
		})
		em.mgr.AddComponent(e, &game.Position{
			Center: game.Center{
				X: float64(i)*sheetW + lMargin + 32,
				Y: tMargin + 80 + float64(used)*8 + 32,
			},
			Layer: 90,
		})

		handlePrepare := func(i int, villager ecs.Entity) func(x, y float64) {
			return func(x, y float64) {
				em.villagers = append(em.villagers[:i], em.villagers[i+1:]...)
				em.prepared = append(em.prepared, villager)
				em.repaint()
			}
		}
		em.mgr.AddComponent(e, &ui.Interactive{
			W: 48, H: 48,
			Trigger: handlePrepare(i, villager),
		})
	}

	// for _, e := range em.prepared {
	// 	// TODO
	// }

	if len(em.prepared) > 0 {
		// You can embark!
		e := em.mgr.NewEntity()
		em.mgr.Tag(e, "embark")
		em.mgr.Tag(e, "embark-hud")
		em.mgr.AddComponent(e, &game.Sprite{
			Texture: "embark-button.png",

			X: 0, Y: 0,
			W: 64, H: 64,
		})
		em.mgr.AddComponent(e, &game.Scale{
			X: 2,
			Y: 2,
		})
		em.mgr.AddComponent(e, &game.Position{
			Center: game.Center{
				X: float64(em.screenW) / 2,
				Y: float64(em.screenH) - 64,
			},
			Layer:    100,
			Absolute: true,
		})
		em.mgr.AddComponent(e, &ui.Interactive{
			W: 48, H: 48,
			Trigger: func(x, y float64) {
				em.bus.Publish(&SquadSelected{})

				e := em.mgr.NewEntity()
				em.mgr.Tag(e, "player")
				em.mgr.AddComponent(e, &game.Squad{})
				squad := em.mgr.Component(e, "Squad").(*game.Squad)
				players := game.NewTeam()
				em.mgr.AddComponent(e, players)

				// Add prepared villagers to the team and squad
				for _, e := range em.prepared {
					em.mgr.AddComponent(e, players)
					squad.Members = append(squad.Members, e)
					em.mgr.RemoveTag(e, "embark")
				}

				e = em.mgr.NewEntity()
				em.mgr.AddComponent(e, &game.DiagonalMatrixWipe{
					W: em.screenW, H: em.screenH,
					Obscuring: true,
					OnComplete: func() {
						em.bus.Publish(&Embarked{})
					},
				})
			},
		})
	}
}

// rollVillagers removes any Characters in em.villagers, and generates new ones.
func (em *Manager) rollVillagers() {
	for _, e := range em.villagers {
		em.mgr.DestroyEntity(e)
	}

	// Empty villagers slice while preserving capacity.
	em.villagers = em.villagers[:0]

	g := newGenerator()
	for i := 0; i < 5; i++ {
		e := em.mgr.NewEntity()
		em.mgr.Tag(e, "embark")
		em.mgr.AddComponent(e, g.generateChar())
		em.villagers = append(em.villagers, e)
	}
}
