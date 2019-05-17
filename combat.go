package main

import (
	"time"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/geom"
)

// CombatState enumerates the States that a Combat could be in.
type CombatState int

const (
	// AwaitingInput is when the combat is waiting for the local, human player to make a move.
	AwaitingInputState CombatState = iota
	// ExecutingState is when a move or action is being played out by a character.
	ExecutingState
	// ThinkingState is when an AI-controller player is waiting to get command.
	ThinkingState
	// PreparingState is when no characters is prepared enough to make a move.
	PreparingState
)

// Combat is a game-mode. It processes turns-based Combat until one or the other
// team is knocked out.
type Combat struct {
	// Combat should own systems that are only relevant to Combat. A Turns coordinator, a preparation timer
	mgr    *ecs.World
	bus    *event.Bus
	field  *geom.Field
	nav    *game.Navigator
	camera *Camera
	state  CombatState

	x, y             int     // where the mouse last was in screen coordinates
	wx, wy           float64 // where the mouse last was in world coordinates
	screenW, screenH float64 // most recent dimensions of the window

	// actors ActorSystem...?
	cursors *game.CursorSystem
	intents *game.IntentSystem
}

// NewCombat should accept two opposing squads of characters, a list of
// M,N,Terrain, and a tileset, and a set of environmental effects.
// Or it could accept a *ecs.World, which would by convention contain two squads
// of characters, and some other object that contains terrain, environmental
// effects etc.
func NewCombat(mgr *ecs.World, camera *Camera /**/) *Combat {
	f, _ := geom.NewField(8, 24)
	bus := &event.Bus{}

	c := Combat{
		mgr:     mgr,
		bus:     bus,
		field:   f,
		nav:     game.NewNavigator(bus),
		camera:  camera,
		state:   AwaitingInputState,
		cursors: game.NewCursorSystem(mgr),
		intents: game.NewIntentSystem(mgr, bus, f),
	}
	c.bus.Subscribe(event.MovementConcluded, c.handleMovementConcluded)
	return &c
}

// Begin should be called at the start of an engagement to set up components
// necessary for the combat.
func (c *Combat) Begin() {
	/*
		At the start of Combat, we need to add a sprite and position component to
		every actor, because a Combat should be the thing responsible for deciding
		how to render an actor on the field.
	*/
	c.camera.Center(c.field.Width()/2, c.field.Height()/2)
}

// End should be called at the resolution of a combat encounter. It removes
// combat-specific Components.
func (c *Combat) End() {

}

// Run a frame of this Combat.
func (c *Combat) Run(elapsed time.Duration) {
	c.nav.Update(c.mgr, elapsed)
	c.intents.Update()
}

// Interaction is the way to notify the Combat that a mouse click or touch event
// occurred.
func (c *Combat) Interaction(x, y int) {
	if c.state == AwaitingInputState {
		actor := c.actorAwaitingInput()

		wx, wy := c.camera.ScreenToWorld(x, y)

		m := game.MoveIntent{X: wx, Y: wy}
		c.mgr.AddComponent(actor, &m)

		c.state = ExecutingState

		c.cursors.Clear()
	}
}

// MousePosition is the way to notify the Combat that the mouse has a new
// position.
func (c *Combat) MousePosition(x, y int) {
	wx, wy := c.camera.ScreenToWorld(x, y)
	if c.state == AwaitingInputState {
		c.cursors.Clear()

		a := c.mgr.Component(c.actorAwaitingInput(), "Actor").(*game.Actor)

		switch a.Size {
		case game.SMALL:
			h := c.field.At(int(wx), int(wy))
			if h == nil {
				break
			}
			c.cursors.Add(h.X(), h.Y(), a.Size)

		case game.MEDIUM:
			h := c.field.At4(int(wx), int(wy))
			if h == nil {
				break
			}
			c.cursors.Add(h.X(), h.Y(), a.Size)

		case game.LARGE:
			h := c.field.At7(int(wx), int(wy))
			if h == nil {
				break
			}
			c.cursors.Add(h.X(), h.Y(), a.Size)
		}
	}

	c.x = x
	c.y = y
	c.wx = wx
	c.wy = wy
}

func (c *Combat) actorAwaitingInput() ecs.Entity {
	entities := c.mgr.Get([]string{"Actor"})
	if len(entities) == 0 {
		// FIXME: this is a flow error - there should always be an entity
		return 0
	}
	return entities[0]
}

func (c *Combat) handleMovementConcluded(t event.Typer) {
	c.state = AwaitingInputState
	c.MousePosition(c.x, c.y)
}
