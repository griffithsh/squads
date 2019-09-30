package overworld

import (
	"time"

	"github.com/griffithsh/squads/game"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
)

// Manager is a game state that allows the player to pick which path to take,
// and which combat to enter etc.
type Manager struct {
	mgr *ecs.World
	bus *event.Bus

	paused bool
}

// NewManager creates a new overworld Manager.
func NewManager(mgr *ecs.World, bus *event.Bus) *Manager {
	return &Manager{
		mgr: mgr,
		bus: bus,

		paused: false,
	}
}

// Begin a Manager session.
func (m *Manager) Begin() {
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
		Layer: 2001,
	})
}

// Pause toggles the paused state of the overworld Manager. When paused, nothing
// is rendered, and player interactions are discarded.
func (m *Manager) Pause() {
	if m.paused {
		for _, e := range m.mgr.Tagged("overworld") {
			m.mgr.RemoveComponent(e, &game.Hidden{})
		}
	} else {
		for _, e := range m.mgr.Tagged("overworld") {
			m.mgr.AddComponent(e, &game.Hidden{})
		}
	}

	m.paused = !m.paused
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
	if m.paused {
		return
	}
	// accept input from hardware abstraction layer
}

// MousePosition handles a change in the mouse position from the player.
func (m *Manager) MousePosition(x, y int) {
	if m.paused {
		return
	}
	// accept input from hardware abstraction layer
}

// Run the Manager.
func (m *Manager) Run(elapsed time.Duration) {
	if m.paused {
		return
	}
	// todo
}
