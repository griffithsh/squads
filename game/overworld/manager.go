package overworld

import (
	"math"
	"math/rand"
	"time"

	"github.com/griffithsh/squads/ui"

	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/geom"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
)

// Manager is a game state that allows the player to pick which path to take,
// and which combat to enter etc.
type Manager struct {
	mgr *ecs.World
	bus *event.Bus

	dormant bool
	state   State
}

// NewManager creates a new overworld Manager.
func NewManager(mgr *ecs.World, bus *event.Bus) *Manager {
	m := Manager{
		mgr: mgr,
		bus: bus,

		dormant: false,
		state:   Uninitialised,
	}

	bus.Subscribe(TokensCollided{}.Type(), m.handleTokensCollided)

	return &m
}

func (m *Manager) handleTokensCollided(t event.Typer) {
	ev := t.(*TokensCollided)

	p := m.mgr.Component(ev.E1, "Position").(*game.Position)
	m.bus.Publish(&game.SomethingInteresting{
		X: p.Center.X,
		Y: p.Center.Y,
	})

	m.setState(FadingOut)
	m.mgr.AddComponent(m.mgr.NewEntity(), &game.DiagonalMatrixWipe{
		W: 1024, H: 768, // FIXME: need access to screen dimensions
		Obscuring: true,
		OnComplete: func() {
			squads := []ecs.Entity{}
			for _, e := range []ecs.Entity{ev.E1, ev.E2} {
				token := m.mgr.Component(e, "Token").(*Token)
				squads = append(squads, token.Squad)
			}
			m.bus.Publish(&CombatInitiated{
				Squads: squads,
				// info about terrain
			})
		},
	})
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

		// We need to know if n - the node we clicked on - is
		// connected to the node the overworld token is on.

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
		for _, neighbor := range n.Directions {
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

// Begin a Manager session.
func (m *Manager) Begin(d Data) {
	m.setState(AwaitingInputState)
	// Add new entities for the squad, overworld terrain, etc?
	// TODO

	// Show nodes/cities/halts.
	positions := map[geom.Key]game.Center{}
	for _, n := range d.Nodes {
		e := m.mgr.NewEntity()
		n.e = e
		m.mgr.Tag(e, "overworld")
		m.mgr.AddComponent(e, &game.Sprite{
			Texture: "overworld-nodes.png",

			X: 0, Y: 0,
			W: 24, H: 16,
		})
		m.mgr.AddComponent(e, &ui.Interactive{
			W: 32, H: 24,
			Trigger: m.newNodeClickHandler(n),
		})

		x, y := geom.XY(n.ID.M, n.ID.N, 144, 96)

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

		e = m.mgr.NewEntity()
		m.mgr.Tag(e, "overworld")
		m.mgr.AddComponent(e, &game.Sprite{
			Texture: "overworld-grass.png",

			X: 0, Y: 0,
			W: 144, H: 96,
		})

		x, y = geom.XY(n.ID.M, n.ID.N, 144, 96)

		m.mgr.AddComponent(e, &game.Position{
			Center: game.Center{
				X: x, Y: y,
			},
			Layer: 5,
		})
	}

	keys := make([]geom.Key, 0, len(d.Nodes))
	for k := range d.Nodes {
		keys = append(keys, k)
	}
	rand.Shuffle(len(keys), func(i, j int) { keys[i], keys[j] = keys[j], keys[i] })
	enemies := game.NewTeam()
	for i, k := range keys {
		start := d.Nodes[k]
		position := m.mgr.Component(start.e, "Position").(*game.Position)
		if i == 0 {
			// Add a Token to mark where the player's squad is.
			e := m.mgr.NewEntity()
			m.mgr.Tag(e, "overworld")
			m.mgr.Tag(e, "player")
			m.mgr.AddComponent(e, m.playerTeam())
			m.mgr.AddComponent(e, &game.Sprite{
				Texture: "figure.png",

				X: 0, Y: 0,
				W: 24, H: 48,
				OffsetY: -16,
			})
			m.mgr.AddComponent(e, &game.Position{
				Center: game.Center{
					X: position.Center.X, Y: position.Center.Y,
				},
				Layer: position.Layer + 1,
			})
			m.mgr.AddComponent(e, &Token{
				Key:   start.ID,
				Squad: m.playerSquad(),
			})

			// Publish a focus event for the camera.
			m.bus.Publish(&game.SomethingInteresting{
				X: position.Center.X,
				Y: position.Center.Y,
			})

			continue
		}

		roll := rand.Intn(6)
		if roll == 0 {
			// Luckily, no baddies on this node.
			continue
		}

		// Add a Squad, and visible Token to the overworld map.
		e := m.mgr.NewEntity()
		m.mgr.Tag(e, "overworld")
		enemies.Control = game.ComputerControl
		m.mgr.AddComponent(e, enemies)
		m.mgr.AddComponent(e, &game.Squad{})
		squad := m.mgr.Component(e, "Squad").(*game.Squad)
		m.mgr.AddComponent(e, &game.Sprite{
			Texture: "wolf.png",

			X: 0, Y: 0,
			W: 64, H: 64,
		})
		m.mgr.AddComponent(e, enemies)
		m.mgr.AddComponent(e, &game.Position{
			Center: game.Center{
				X: position.Center.X, Y: position.Center.Y,
			},
			Layer: position.Layer + 1,
		})
		m.mgr.AddComponent(e, &Token{
			Key:   start.ID,
			Squad: e,
		})

		// Add a baddy to this Squad.
		e = m.mgr.NewEntity()
		m.mgr.Tag(e, "overworld")
		m.mgr.Tag(e, "baddy")
		m.mgr.AddComponent(e, enemies)
		m.mgr.AddComponent(e, &game.Character{
			Name:                 "Dumble",
			Size:                 game.SMALL,
			Sex:                  game.Male,
			Profession:           game.Skeleton,
			PreparationThreshold: 1650,
			ActionPoints:         60,
			SmallIcon: game.Sprite{
				Texture: "hud.png",
				X:       0,
				Y:       154,
				W:       26,
				H:       26,
			},
			BigIcon: game.Sprite{
				Texture: "hud.png",
				X:       0,
				Y:       102,
				W:       52,
				H:       52,
			},
		})
		squad.Members = append(squad.Members, e)
	}

	// Now show connections between the nodes.
	type connectKey struct {
		M1, N1, M2, N2 int
	}
	connected := map[connectKey]struct{}{}
	for _, n := range d.Nodes {
		for _, other := range n.Directions {
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
					Texture: "overworld-nodes.png",

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

// Enable the overworld Manager, responding to input and rendering the state of
// the overworld.
func (m *Manager) Enable() {
	if m.dormant {
		m.mgr.AddComponent(m.mgr.NewEntity(), &game.DiagonalMatrixWipe{
			W: 1024, H: 768, // FIXME: need access to screen dimensions
		})
		for _, e := range m.mgr.Tagged("overworld") {
			m.mgr.RemoveComponent(e, &game.Hidden{})
		}
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
