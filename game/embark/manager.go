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
	// Create a button to press to embark
	e := em.mgr.NewEntity()
	em.mgr.Tag(e, "embark")

	// TODO: Generate 5 Characters, render "stat sheet" type things next to a
	// little villager avatar for each.
	lMargin := 64.0
	tMargin := 16.0
	sheetW := 144.0
	g := newGenerator()
	for i := 0; i < 5; i++ {
		char := g.generateChar()

		container := em.mgr.NewEntity()
		em.mgr.Tag(container, "embark")

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

		// Prep?
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

		// Embark button?
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

				// Create Characters to Populate the player's Squad.
				e = em.mgr.NewEntity()
				em.mgr.AddComponent(e, players)
				squad.Members = append(squad.Members, e)
				em.mgr.AddComponent(e, char)

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
