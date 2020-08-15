package embark

import (
	"fmt"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/ui"
)

// Archive is what is required by embark of any archive data provider.
type Archive interface {
	Profession(profession string) *game.ProfessionDetails
	Names() map[string][]string
	Appearance(profession string, sex game.CharacterSex, hair string, skin string) *game.Appearance
	HairVariations() []string
	SkinVariations() []string
}

// Manager holds state and provides methods to control that state for an embark
// screen. This screen allows the player to configure the Characters they will
// start their run with.
type Manager struct {
	mgr     *ecs.World
	bus     *event.Bus
	archive Archive

	screenW, screenH int

	// prepared lists Characters who are about to embark.
	prepared []ecs.Entity
	// villagers lists Characters who could be selected for embarking.
	villagers []ecs.Entity
}

// NewManager creates a new Manager in a default state. You should call Begin to start the Manager.
func NewManager(mgr *ecs.World, bus *event.Bus, archive Archive) *Manager {
	em := Manager{
		mgr:     mgr,
		bus:     bus,
		archive: archive,
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
	const sheetW float64 = 148 // sheetW is the width of each character sheet
	for i, villager := range em.villagers {
		char := em.mgr.Component(villager, "Character").(*game.Character)
		equip, _ := em.mgr.Component(villager, "Equipment").(*game.Equipment)

		handlePrepare := func(i int, villager ecs.Entity) func(x, y float64) {
			return func(x, y float64) {
				em.villagers = append(em.villagers[:i], em.villagers[i+1:]...)
				em.prepared = append(em.prepared, villager)
				em.repaint()
			}
		}
		em.paintChar(char, equip, float64(i)*sheetW+lMargin, tMargin, handlePrepare(i, villager))
	}

	for i, villager := range em.prepared {
		char := em.mgr.Component(villager, "Character").(*game.Character)
		equip, _ := em.mgr.Component(villager, "Equipment").(*game.Equipment)
		em.paintChar(char, equip, float64(i)*sheetW+lMargin, tMargin+200, nil)
	}

	if len(em.prepared) > 0 {
		// You can embark!
		var e ecs.Entity

		e = ui.ButtonBackground(em.mgr, 30, 15, float64(em.screenW)/2-10, float64(em.screenH)-67, 100, true)
		em.mgr.Tag(e, "embark")
		em.mgr.Tag(e, "embark-hud")

		e = em.mgr.NewEntity()
		em.mgr.Tag(e, "embark")
		em.mgr.Tag(e, "embark-hud")

		em.mgr.AddComponent(e, &game.Font{
			Text: "Go",
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

	g := newGenerator(em.archive)
	for i := 0; i < 5; i++ {
		e := em.mgr.NewEntity()
		em.mgr.Tag(e, "embark")
		em.mgr.AddComponent(e, g.generateChar())
		em.mgr.AddComponent(e, &game.Equipment{
			Weapon: g.generateWeapon(),
		})
		em.villagers = append(em.villagers, e)
	}
}

func (em *Manager) paintChar(char *game.Character, equip *game.Equipment, left float64, top float64, handlePrepare func(x, y float64)) {
	container := em.mgr.NewEntity()
	em.mgr.Tag(container, "embark")
	em.mgr.Tag(container, "embark-characters")

	prof := em.archive.Profession(char.Profession)
	var e ecs.Entity

	// Panel
	e = ui.Panel(em.mgr, 144, 200, left-8, top-8, 90, false)
	em.mgr.Dependency(container, e)

	// Name
	e = em.mgr.NewEntity()
	em.mgr.Dependency(container, e)
	em.mgr.AddComponent(e, &game.Font{
		Text: char.Name,
	})
	em.mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: left,
			Y: top,
		},
		Layer: 100,
	})

	// Icon (BG, portrait, then frame)
	center := game.Center{
		X: left + 108,
		Y: top + 56/2,
	}
	e = em.mgr.NewEntity()
	em.mgr.Dependency(container, e)
	em.mgr.AddComponent(e, &game.Position{
		Center: center,
		Layer:  99,
	})
	em.mgr.AddComponent(e, &game.PortraitBGBig[char.PortraitBG])

	e = em.mgr.NewEntity()
	em.mgr.Dependency(container, e)
	em.mgr.AddComponent(e, &game.Position{
		Center: center,
		Layer:  100,
	})
	app := em.archive.Appearance(char.Profession, char.Sex, char.Hair, char.Skin)
	spr := app.BigIcon()
	em.mgr.AddComponent(e, &spr)

	e = em.mgr.NewEntity()
	em.mgr.Dependency(container, e)
	em.mgr.AddComponent(e, &game.Position{
		Center: center,
		Layer:  101,
	})
	em.mgr.AddComponent(e, &game.PortraitFrameBig[char.PortraitFrame])

	// Profession
	e = em.mgr.NewEntity()
	em.mgr.Dependency(container, e)
	em.mgr.AddComponent(e, &game.Font{
		Text: char.Profession,
		Size: "small",
	})
	em.mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: left,
			Y: top + 12,
		},
		Layer: 100,
	})

	// Level
	e = em.mgr.NewEntity()
	em.mgr.Dependency(container, e)
	em.mgr.AddComponent(e, &game.Font{
		Text: fmt.Sprintf("Lvl: %d", char.Level),
		Size: "small",
	})
	em.mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: left,
			Y: top + 12 + 8,
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
			X: left,
			Y: top + 12 + 16,
		},
		Layer: 100,
	})

	// Prep
	e = em.mgr.NewEntity()
	em.mgr.Dependency(container, e)
	em.mgr.AddComponent(e, &game.Font{
		Text: fmt.Sprintf("Prep: %d", char.InherantPreparation+prof.Preparation+equip.WeaponPreparation()),
		Size: "small",
	})
	em.mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: left,
			Y: top + 12 + 24,
		},
		Layer: 100,
	})

	// Action Points
	e = em.mgr.NewEntity()
	em.mgr.Dependency(container, e)
	em.mgr.AddComponent(e, &game.Font{
		Text: fmt.Sprintf("AP: %d", char.InherantActionPoints+prof.ActionPoints+equip.WeaponActionPoints()),
		Size: "small",
	})
	em.mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: left,
			Y: top + 12 + 32,
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
			X: left,
			Y: top + 12 + 40,
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
			X: left,
			Y: top + 12 + 48,
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
			X: left,
			Y: top + 12 + 56,
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
			X: left,
			Y: top + 12 + 64,
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
				X: left,
				Y: top + 88 + float64(used)*8,
			},
			Layer: 100,
		})

		used++
	}

	if handlePrepare == nil {
		// Don't add a button if there is no handler provided.
		return
	}

	// Prepare button (add villager to the preparing squad).
	e = ui.ButtonBackground(em.mgr, 48, 15, left, top+170, 90, false)
	em.mgr.Dependency(container, e)

	e = em.mgr.NewEntity()
	em.mgr.Dependency(container, e)
	em.mgr.AddComponent(e, &game.Font{
		Text: "Prepare",
	})
	em.mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: left + 3,
			Y: top + 173,
		},
		Layer: 100,
	})

	e = em.mgr.NewEntity()
	em.mgr.Dependency(container, e)
	em.mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: left + 48/2,
			Y: top + 177,
		},
		Layer: 90,
	})
	em.mgr.AddComponent(e, &ui.Interactive{
		W: 48, H: 15,
		Trigger: handlePrepare,
	})
}
