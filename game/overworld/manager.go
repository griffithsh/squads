package overworld

import (
	"fmt"
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
}

// NewManager creates a new overworld Manager.
func NewManager(mgr *ecs.World, bus *event.Bus) *Manager {
	return &Manager{
		mgr: mgr,
		bus: bus,

		dormant: false,
	}
}

// randInHex generates a random point in an overworld hex.
func randInHex() (float64, float64) {
	rad := rand.Float64() * math.Pi * 2
	sin, cos := math.Sincos(rad)

	w, h := 144.0, 96.0
	factor := 0.2
	return w * factor * sin, h * factor * cos
}

// Begin a Manager session.
func (m *Manager) Begin(d Data) {
	// Add new entities for the squad, overworld terrain, etc?
	// TODO

	// Add some random figure to show pan/zoom.
	e := m.mgr.NewEntity()
	m.mgr.Tag(e, "overworld")
	m.mgr.AddComponent(e, &game.Sprite{
		Texture: "figure.png",

		X: 0, Y: 0,
		W: 24, H: 48,
		OffsetY: -6,
	})
	m.mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: 12, Y: 48,
		},
		Layer: 101,
	})

	// Show nodes/cities/halts.
	positions := map[geom.Key]game.Center{}
	for _, n := range d.Nodes {
		e := m.mgr.NewEntity()
		m.mgr.Tag(e, "overworld")
		m.mgr.AddComponent(e, &game.Sprite{
			Texture: "overworld-nodes.png",

			X: 0, Y: 0,
			W: 24, H: 16,
		})
		f := func(n *Node) func() {
			return func() {
				fmt.Printf("click: %v\n", n.ID)
			}
		}
		m.mgr.AddComponent(e, &ui.Interactive2{
			W: 32, H: 24,
			Trigger: f(n),
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

// Interaction handles an interaction from the player at x,y.
func (m *Manager) Interaction(x, y int) {
	if m.dormant {
		return
	}
	// accept input from hardware abstraction layer
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
