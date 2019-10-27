package overworld

import (
	"time"

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

// Begin a Manager session.
func (m *Manager) Begin(d Data) {
	// Add new entities for the squad, overworld terrain, etc?
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

	for _, n := range d.Nodes {
		e := m.mgr.NewEntity()
		m.mgr.Tag(e, "overworld")
		m.mgr.AddComponent(e, &game.Sprite{
			Texture: "overworld-jumbo-hexes.png",

			X: 0, Y: 0,
			W: 144, H: 96,
		})
		h := geom.Hex{M: n.ID.M, N: n.ID.N}
		m.mgr.AddComponent(e, &game.Position{
			Center: game.Center{
				X: h.X() * 6, Y: h.Y() * 6,
			},
			Layer: 1,
		})

		add := func(x int) {
			e := m.mgr.NewEntity()
			m.mgr.Tag(e, "overworld")
			m.mgr.AddComponent(e, &game.Sprite{
				Texture: "overworld-jumbo-hexes.png",

				X: x, Y: 0,
				W: 144, H: 96,
			})
			h := geom.Hex{M: n.ID.M, N: n.ID.N}
			m.mgr.AddComponent(e, &game.Position{
				Center: game.Center{
					X: h.X() * 6, Y: h.Y() * 6,
				},
				Layer: 1,
			})
		}

		if _, ok := n.Directions[geom.S]; ok {
			add(1 * 144)
		}
		if _, ok := n.Directions[geom.SW]; ok {
			add(2 * 144)
		}
		if _, ok := n.Directions[geom.NW]; ok {
			add(3 * 144)
		}
		if _, ok := n.Directions[geom.N]; ok {
			add(4 * 144)
		}
		if _, ok := n.Directions[geom.NE]; ok {
			add(5 * 144)
		}
		if _, ok := n.Directions[geom.SE]; ok {
			add(6 * 144)
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
