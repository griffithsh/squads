package overworld

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/geom"
	"github.com/griffithsh/squads/ui"
)

type Archive interface {
	PedestalAppearances(sinister bool) []int
	GetRecipes() []*Recipe
	GetAnimation(name string) game.FrameAnimation
}

// Manager is a game state that allows the player to pick which path to take,
// and which combat to enter etc.
type Manager struct {
	mgr     *ecs.World
	bus     *event.Bus
	recipes []*Recipe
	archive Archive

	screenW, screenH int

	dormant bool
	state   State

	fogged map[geom.Key]ecs.Entity

	rng *rand.Rand
}

// NewManager creates a new overworld Manager.
func NewManager(mgr *ecs.World, bus *event.Bus, archive Archive) *Manager {
	m := Manager{
		mgr:     mgr,
		bus:     bus,
		recipes: archive.GetRecipes(),
		archive: archive,

		dormant: false,
		state:   Uninitialised,

		fogged: make(map[geom.Key]ecs.Entity),
	}

	bus.Subscribe(TokensCollided{}.Type(), m.handleTokensCollided)
	bus.Subscribe(game.WindowSizeChanged{}.Type(), m.handleWindowSizeChanged)

	return &m
}

func (m *Manager) handleSquadTokensCollided(e1, e2 ecs.Entity) {
	m.setState(FadingOut)
	m.mgr.AddComponent(m.mgr.NewEntity(), &game.DiagonalMatrixWipe{
		W: m.screenW, H: m.screenH,
		Obscuring: true,
		OnComplete: func() {
			squads := []ecs.Entity{}
			for _, e := range []ecs.Entity{e1, e2} {
				token := m.mgr.Component(e, "Token").(*Token)
				if m.mgr.Component(token.Presence, "Squad") != nil {
					squads = append(squads, token.Presence)
				}
			}
			m.bus.Publish(&CombatInitiated{
				Squads: squads,
				// TODO: info about terrain
			})
		},
	})
}

func (m *Manager) handleExitGateCollided() {
	m.setState(FadingOut)
	m.mgr.AddComponent(m.mgr.NewEntity(), &game.DiagonalMatrixWipe{
		W: m.screenW, H: m.screenH,
		Obscuring: true,
		OnComplete: func() {
			m.bus.Publish(&Complete{})
		},
	})
}

func (m *Manager) handleTokensCollided(t event.Typer) {
	ev := t.(*TokensCollided)

	p := m.mgr.Component(ev.E1, "Position").(*game.Position)
	m.bus.Publish(&game.SomethingInteresting{
		X: p.Center.X,
		Y: p.Center.Y,
	})

	var sum int
	for _, e := range []ecs.Entity{ev.E1, ev.E2} {
		sum += int(m.mgr.Component(e, "Token").(*Token).Category)
	}

	switch sum {
	case int(SquadToken + SquadToken):
		m.handleSquadTokensCollided(ev.E1, ev.E2)
	case int(SquadToken + GateToken):
		m.handleExitGateCollided()
	default:
		fmt.Println("unknown TokenCollision with sum:", sum)
	}
}

func (m *Manager) handleWindowSizeChanged(e event.Typer) {
	wsc := e.(*game.WindowSizeChanged)
	m.screenW, m.screenH = wsc.NewW, wsc.NewH
}

// randInHex generates a random point in an overworld hex.
func randInHex() (float64, float64) {
	rad := rand.Float64() * math.Pi * 2
	sin, cos := math.Sincos(rad)

	w, h := 144.0, 96.0
	factor := 0.2
	return w * factor * sin, h * factor * cos
}

// newNodeClickHandler creates a new click handler for Node n.
func (m *Manager) newNodeClickHandler(n *Node) func(x, y float64) {
	// Closure to capture value of n and provide a function that matches the
	// signature of ui.Interactive.Trigger.
	return func(x, y float64) {
		if m.state != AwaitingInputState {
			return
		}

		// We need to know if n (the node we clicked on) is connected to the
		// node the overworld token is on. If it's not, it's invalid to move
		// there.

		// Find the Token that belongs to the player's squad.
		var e ecs.Entity
		for _, maybe := range m.mgr.Get([]string{"Token", "Team"}) {
			team := m.mgr.Component(maybe, "Team").(*game.Team)
			if team.Control == game.LocalControl {
				e = maybe
				break
			}
		}

		t := m.mgr.Component(e, "Token").(*Token)
		var connected bool
		for _, neighbor := range n.Connected {
			if neighbor == t.Key {
				connected = true
				break
			}
		}
		if connected {
			m.setState(AnimatingState)

			refPos := m.mgr.Component(n.e, "Position").(*game.Position)
			m.mgr.AddComponent(e, &Traversal{
				Duration:    800 * time.Millisecond,
				Destination: refPos.Center,
				Complete: func() {
					m.bus.Publish(&TokenMoved{
						E:    e,
						From: t.Key,
						To:   n.ID,
					})
					t.Key = n.ID
					m.setState(AwaitingInputState)

					// We've arrived at a new node, so update what is fogged.
					for _, key := range m.playerVision() {
						e, ok := m.fogged[key]
						if !ok {
							continue
						}

						// Start the "fog-revealing" animation.
						fa := m.archive.GetAnimation("overworld-reveal-grass")

						fa.EndBehavior = game.DestroyEntity
						m.mgr.AddComponent(e, &fa)
						delete(m.fogged, key)
					}
				},
			})
		}
	}
}

func (m *Manager) setState(new State) {
	m.state = new
}

func (m *Manager) playerSquad() ecs.Entity {
	for _, e := range m.mgr.Tagged("player") {
		if m.mgr.Component(e, "Squad") != nil {
			return e
		}
	}
	return 0
}

// playerVision returns the Keys that the player can currently see, but not
// necessarily all of the hexes currently revealed in this overworld.
func (m *Manager) playerVision() []geom.Key {
	for _, e := range m.mgr.Tagged("player") {
		token, ok := m.mgr.Component(e, "Token").(*Token)
		if !ok {
			continue
		}
		// First entity tagged with "player" that has a Token Component is selected.
		result := []geom.Key{token.Key}
		for key := range token.Key.Neighbors() {
			result = append(result, key)
		}
		return result
	}
	return []geom.Key{}
}

func (m *Manager) playerTeam() *game.Team {
	for _, e := range m.mgr.Get([]string{"Squad"}) {
		if !m.mgr.HasTag(e, "player") {
			continue
		}
		// Found the player's squad.
		team := m.mgr.Component(e, "Team").(*game.Team)
		return team
	}
	return nil
}

func (m *Manager) boot(d Map) {
	f := geom.NewField(50, 47, 96)
	// Add a Sprite for every Node.
	positions := map[geom.Key]game.Center{}
	for _, n := range d.Nodes {
		e := m.mgr.NewEntity()
		n.e = e
		m.mgr.Tag(e, "overworld")
		m.mgr.AddComponent(e, &game.Sprite{
			Texture: "overworld/nodes.png",

			X: 0, Y: 0,
			W: 24, H: 16,
		})
		m.mgr.AddComponent(e, &ui.Interactive{
			W: 32, H: 24,
			Trigger: m.newNodeClickHandler(n),
		})

		x, y := f.Ktow(n.ID)

		rx, ry := randInHex()
		x = rx + x
		y = ry + y
		m.mgr.AddComponent(e, &game.Position{
			Center: game.Center{
				X: x, Y: y,
			},
			Layer: 10,
		})
		positions[n.ID] = game.Center{X: x, Y: y}
	}

	// Add terrain tiles.
	for k, tile := range d.Terrain {
		e := m.mgr.NewEntity()
		m.mgr.Tag(e, "overworld")
		switch tile {
		case Grass:
			m.mgr.AddComponent(e, &game.Sprite{
				Texture: "deprecated/overworld-grass.png",

				X: 144, Y: 0,
				W: 144, H: 96,
			})
		case Stone:
			m.mgr.AddComponent(e, &game.Sprite{
				Texture: "deprecated/overworld-grass.png",

				X: 144, Y: 96,
				W: 144, H: 96,
			})
		case Trees:
			m.mgr.AddComponent(e, &game.Sprite{
				Texture: "deprecated/overworld-grass.png",

				X: 288, Y: 0,
				W: 144, H: 96,
			})
		default:
			fmt.Printf("unknown tile %d at %v\n", tile, k)
		}

		x, y := f.Ktow(k)

		m.mgr.AddComponent(e, &game.Position{
			Center: game.Center{
				X: x, Y: y,
			},
			Layer: 5,
		})
	}

	// Add a Token to mark where the player's squad is.
	e := m.mgr.NewEntity()
	m.mgr.Tag(e, "overworld")
	m.mgr.Tag(e, "player")
	m.mgr.AddComponent(e, m.playerTeam())
	m.mgr.AddComponent(e, &game.Sprite{
		Texture: "deprecated/figure.png",

		X: 0, Y: 0,
		W: 24, H: 48,
		OffsetY: -16,
	})
	position := m.mgr.Component(d.Nodes[d.Start].e, "Position").(*game.Position)
	m.mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: position.Center.X, Y: position.Center.Y,
		},
		Layer: position.Layer + 1,
	})
	m.mgr.AddComponent(e, &Token{
		Key:      d.Start,
		Presence: m.playerSquad(),
		Category: SquadToken,
	})
	// Publish a focus event for the camera.
	m.bus.Publish(&game.SomethingInteresting{
		X: position.Center.X,
		Y: position.Center.Y,
	})

	// Add a token for the exit gate!
	e = m.mgr.NewEntity()
	m.mgr.Tag(e, "overworld")
	m.mgr.Tag(e, "gate")
	m.mgr.AddComponent(e, &game.Sprite{
		Texture: "overworld/tokens.png",

		X: 0, Y: 0,
		W: 64, H: 64,
		OffsetY: -8,
	})
	position = m.mgr.Component(d.Nodes[d.Gate].e, "Position").(*game.Position)
	m.mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: position.Center.X, Y: position.Center.Y,
		},
		Layer: position.Layer + 1,
	})
	m.mgr.AddComponent(e, &Token{
		Key:      d.Gate,
		Presence: e, // oh, it refs itself !?
		Category: GateToken,
	})
	npcs := game.NewTeam()
	npcs.Control = game.NoControl

	apps := m.archive.PedestalAppearances(true)
	npcs.PedestalAppearance = apps[rand.Intn(len(apps))]
	m.mgr.AddComponent(e, npcs)

	// Add a Token for every enemy Squad.
	for key, squadMembers := range d.Enemies {
		position := m.mgr.Component(d.Nodes[key].e, "Position").(*game.Position)
		// Add a Squad, and visible Token to the overworld map.
		e := m.mgr.NewEntity()
		m.mgr.Tag(e, "overworld")
		enemyTeam := game.NewTeam()
		enemyTeam.PedestalAppearance = apps[rand.Intn(len(apps))]
		enemyTeam.Control = game.ComputerControl
		m.mgr.AddComponent(e, enemyTeam)
		m.mgr.AddComponent(e, &game.Squad{})
		m.mgr.AddComponent(e, &game.Position{
			Center: game.Center{
				X: position.Center.X, Y: position.Center.Y,
			},
			Layer: position.Layer + 1,
		})
		m.mgr.AddComponent(e, &Token{
			Key:      key,
			Presence: e,
			Category: SquadToken,
		})

		m.mgr.AddComponent(e, &game.Sprite{
			Texture: "overworld/tokens.png",

			X: 64, Y: 32,
			W: 32, H: 32,
		})

		squad := m.mgr.Component(e, "Squad").(*game.Squad)
		for _, character := range squadMembers {
			e = m.mgr.NewEntity()
			m.mgr.Tag(e, "overworld")
			m.mgr.Tag(e, "baddy")
			m.mgr.AddComponent(e, character)
			m.mgr.AddComponent(e, enemyTeam)
			squad.Members = append(squad.Members, e)
		}
	}

	// Add fog over all the Terrain in the map.
	for k := range d.Terrain {
		if _, ok := m.fogged[k]; !ok {
			e := m.mgr.NewEntity()
			m.mgr.Tag(e, "overworld")
			m.fogged[k] = e
		}
	}

	// Remove anything the player can currently see from the fogged list.
	for _, key := range m.playerVision() {
		m.mgr.DestroyEntity(m.fogged[key])
		delete(m.fogged, key)
	}

	// Now add fog-of-war sprites over the fogged nodes.
	for k, e := range m.fogged {
		x, y := f.Ktow(k)

		m.mgr.AddComponent(e, &game.Sprite{
			Texture: "deprecated/overworld-grass.png",

			X: 0, Y: 0,
			W: 144, H: 96,
		})
		m.mgr.AddComponent(e, &game.Position{
			Center: game.Center{
				X: x, Y: y,
			},
			Layer: 100,
		})
	}

	// Now show the visible paths between the nodes.
	type connectKey struct {
		M1, N1, M2, N2 int
	}
	connected := map[connectKey]struct{}{}
	for _, n := range d.Nodes {
		for _, other := range n.Connected {
			conn := connectKey{other.M, other.N, n.ID.M, n.ID.N}
			if _, ok := connected[conn]; ok {
				continue
			}
			connected[connectKey{n.ID.M, n.ID.N, other.M, other.N}] = struct{}{}

			a := positions[other].X - positions[n.ID].X
			b := positions[other].Y - positions[n.ID].Y
			hypotenuse := math.Sqrt(a*a + b*b)
			steps := int(math.Round(hypotenuse / 24))
			if steps <= 1 {
				steps = 2
			}
			for i := 0; i < steps; i++ {
				if i == 0 {
					continue
				}
				e := m.mgr.NewEntity()
				m.mgr.Tag(e, "overworld")
				m.mgr.AddComponent(e, &game.Sprite{
					Texture: "overworld/nodes.png",

					X: 24, Y: 0,
					W: 8, H: 6,
				})
				x := positions[n.ID].X + float64(i)*a/float64(steps)
				y := positions[n.ID].Y + float64(i)*b/float64(steps)
				m.mgr.AddComponent(e, &game.Position{
					Center: game.Center{
						X: x,
						Y: y,
					},
					Layer: 10,
				})
			}
		}
	}
}

func (m *Manager) handleCardSelected(e ecs.Entity, others []ecs.Entity, recipe *Recipe, lvl int) func(x, y float64) {
	return func(float64, float64) {
		// Remove the interactive we just clicked on so that we cannot
		// double-click by access.
		m.mgr.RemoveComponent(e, &ui.Interactive{})

		// Change this card to the selected card sprite.
		m.mgr.AddComponent(e, &game.Sprite{
			Texture: "overworld/cards.png",

			X: 128, Y: 0,
			W: 128, H: 192,
		})

		for _, e := range others {
			// remove interactive component from others
			m.mgr.RemoveComponent(e, &ui.Interactive{})

			// obscure others with fadeout animation
			anim := m.archive.GetAnimation("overworld-hide-card")
			anim.EndBehavior = game.HoldLastFrame

			obscure := m.mgr.NewEntity()
			m.mgr.Tag(obscure, "overworld")
			m.mgr.Tag(obscure, "overworld-pick-path")
			pos := m.mgr.Component(e, "Position").(*game.Position)
			posCopy := *pos
			posCopy.Layer += 10
			m.mgr.AddComponent(obscure, &posCopy)
			m.mgr.AddComponent(obscure, m.mgr.Component(e, "Scale"))
			m.mgr.AddComponent(obscure, &anim)

		}

		// start a screen wipe with an OnComplete that will ...
		m.setState(FadingOut)
		m.mgr.AddComponent(m.mgr.NewEntity(), &game.DiagonalMatrixWipe{
			W: m.screenW, H: m.screenH,
			Obscuring: true,
			OnComplete: func() {
				// ... destroy all components tagged "overworld-pick-path",
				for _, e := range m.mgr.Tagged("overworld-pick-path") {
					m.mgr.DestroyEntity(e)
				}

				// unwipe the screen,
				m.setState(FadingIn)
				m.mgr.AddComponent(m.mgr.NewEntity(), &game.DiagonalMatrixWipe{
					W: m.screenW, H: m.screenH,
					OnComplete: func() {
						m.setState(AwaitingInputState)
					},
					OnInitialised: func() {
						// and boot from the recipe.
						d := generate(m.rng, recipe, lvl)
						m.boot(d)
					},
				})
			},
		})
	}
}

// randomRecipes returns three randomly selected recipes.
func (m *Manager) randomRecipes() []*Recipe {
	max := len(m.recipes)

	result := make([]*Recipe, 3)
	for i := 0; i < 3; i++ {
		ri := m.rng.Intn(max)
		result[i] = m.recipes[ri]
	}
	return result
}

// Begin a Manager session.
func (m *Manager) Begin(seed int64) {
	m.rng = rand.New(rand.NewSource(seed))

	m.setState(FadingIn)
	m.mgr.AddComponent(m.mgr.NewEntity(), &game.DiagonalMatrixWipe{
		W: m.screenW, H: m.screenH,
		OnInitialised: func() {
			// TODO: with the new prng, we should roll for three path options that are
			// presented as cards to the player. I *think* this means three recipes?
			// Sometimes this might be a special recipe that takes the player home and
			// ends their run though. We should also roll for a lvl for the
			// opponents to be, so that the player can elect to take a more
			// difficult path.
			cards := []ecs.Entity{m.mgr.NewEntity(), m.mgr.NewEntity(), m.mgr.NewEntity()}
			recipes := m.randomRecipes()
			for i, e := range cards {
				// lvl captures the difficulty of selecting this option.
				// TODO: it should be based on the level of the player's squad.
				lvl := m.rng.Intn(6) + 1
				recipe := recipes[i]

				others := []ecs.Entity{}
				switch i {
				case 0:
					others = append(others, cards[1], cards[2])
				case 1:
					others = append(others, cards[0], cards[2])
				case 2:
					others = append(others, cards[0], cards[1])
				}

				m.mgr.Tag(e, "overworld")
				m.mgr.Tag(e, "overworld-pick-path")

				diff := float64(128+12) * 2
				m.mgr.AddComponent(e, &game.Position{
					Center: game.Center{
						X: float64(m.screenW/2) - float64(diff) + diff*float64(i),
						Y: float64(m.screenH / 2),
					},

					Layer:    100,
					Absolute: true,
				})
				m.mgr.AddComponent(e, &game.Scale{
					X: 2,
					Y: 2,
				})
				m.mgr.AddComponent(e, &game.Sprite{
					Texture: "overworld/cards.png",

					X: 0, Y: 0,
					W: 128, H: 192,
				})
				m.mgr.AddComponent(e, &ui.Interactive{
					W: 128, H: 192,
					Trigger: m.handleCardSelected(e, others, recipe, lvl),
				})

				// Some text on the cards.
				e = m.mgr.NewEntity()
				m.mgr.Tag(e, "overworld")
				m.mgr.Tag(e, "overworld-pick-path")

				m.mgr.AddComponent(e, &game.Font{
					Text: recipe.Label,
				})
				m.mgr.AddComponent(e, &game.Position{
					Center: game.Center{
						X: float64(m.screenW/2) - float64(diff) + diff*float64(i) - 96,
						Y: float64(m.screenH/2) - 176,
					},

					Layer:    101,
					Absolute: true,
				})
				m.mgr.AddComponent(e, &game.Scale{
					X: 2,
					Y: 2,
				})

				// A preview image of the destination.
				e = m.mgr.NewEntity()
				m.mgr.Tag(e, "overworld")
				m.mgr.Tag(e, "overworld-pick-path")

				m.mgr.AddComponent(e, &game.Sprite{
					Texture: "overworld/cards.png",

					X: 0, Y: 384,
					W: 96, H: 72,
				})
				m.mgr.AddComponent(e, &game.Position{
					Center: game.Center{
						X: float64(m.screenW/2) - float64(diff) + diff*float64(i),
						Y: float64(m.screenH/2) - 72,
					},

					Layer:    101,
					Absolute: true,
				})
				m.mgr.AddComponent(e, &game.Scale{
					X: 2,
					Y: 2,
				})

				// How hard is this path?
				e = m.mgr.NewEntity()
				m.mgr.Tag(e, "overworld")
				m.mgr.Tag(e, "overworld-pick-path")

				m.mgr.AddComponent(e, &game.Font{
					Text: "LVL: " + strconv.Itoa(lvl),
					Size: "small",
				})
				m.mgr.AddComponent(e, &game.Position{
					Center: game.Center{
						X: float64(m.screenW/2) - float64(diff) + diff*float64(i) - 96,
						Y: float64(m.screenH/2) + 12,
					},

					Layer:    101,
					Absolute: true,
				})
				m.mgr.AddComponent(e, &game.Scale{
					X: 2,
					Y: 2,
				})
			}
		},
	})
}

// Enable the overworld Manager, responding to input and rendering the state of
// the overworld.
func (m *Manager) Enable() {
	if m.dormant {
		m.mgr.AddComponent(m.mgr.NewEntity(), &game.DiagonalMatrixWipe{
			W: m.screenW, H: m.screenH,
			OnInitialised: func() {
				for _, e := range m.mgr.Tagged("overworld") {
					m.mgr.RemoveComponent(e, &game.Hidden{})
				}
			},
		})
		m.dormant = false
	}
}

// Disable the overworld Manager, ignoring input and not rendering the state of
// the overworld.
func (m *Manager) Disable() {
	if !m.dormant {
		for _, e := range m.mgr.Tagged("overworld") {
			m.mgr.AddComponent(e, &game.Hidden{})
		}
		m.dormant = true
	}
}

// End should be called when the current overworld map is complete, and the
// player is selecting another map to go to.
func (m *Manager) End() {
	// destroy or hide player entity, overworld components
	for _, e := range m.mgr.Tagged("overworld") {
		m.mgr.DestroyEntity(e)
	}
	m.fogged = make(map[geom.Key]ecs.Entity)
}

// MousePosition handles a change in the mouse position from the player.
func (m *Manager) MousePosition(x, y int) {
	if m.dormant {
		return
	}
	// accept input from hardware abstraction layer
}

// Run the Manager.
func (m *Manager) Run(elapsed time.Duration) {
	if m.dormant {
		return
	}
	// todo
}
