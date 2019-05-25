package main

import (
	"fmt"
	"time"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/geom"
	"github.com/griffithsh/squads/ui"
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
		state:   PreparingState,
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

	c.addGrass()
	c.addTrees()
	c.constructHUD()

	// Upgrade all actors with components for visibility.
	entities := c.mgr.Get([]string{"Actor"})
	for _, e := range entities {
		actor := c.mgr.Component(e, "Actor").(*game.Actor)

		if actor.Size == game.SMALL {
			c.mgr.AddComponent(e, &game.Sprite{
				Texture: "figure.png",
				X:       0,
				Y:       0,
				W:       24,
				H:       48,
			})
			c.mgr.AddComponent(e, &game.SpriteOffset{
				Y: -16,
			})

			start := c.field.Get(0, 0)
			c.mgr.AddComponent(e, &game.Position{
				Center: game.Center{
					X: start.X(),
					Y: start.Y(),
				},
				Layer: 10,
			})

			c.mgr.AddComponent(e, &game.Obstacle{
				M:            0,
				N:            0,
				ObstacleType: game.SmallActor,
			})
		} else if actor.Size == game.MEDIUM {
			c.mgr.AddComponent(e, &game.Sprite{
				Texture: "wolf.png",
				X:       0,
				Y:       0,
				W:       58,
				H:       48,
			})
			c.mgr.AddComponent(e, &game.SpriteOffset{
				Y: -4,
			})

			start := c.field.Get4(0, 7)
			c.mgr.AddComponent(e, &game.Position{
				Center: game.Center{
					X: start.X(),
					Y: start.Y(),
				},
				Layer: 10,
			})

			c.mgr.AddComponent(e, &game.Obstacle{
				M:            0,
				N:            7,
				ObstacleType: game.MediumActor,
			})
		} else if actor.Size == game.LARGE {
			c.mgr.AddComponent(e, &game.Sprite{
				Texture: "figure.png",
				X:       0,
				Y:       0,
				W:       24,
				H:       48,
			})
			c.mgr.AddComponent(e, &game.SpriteOffset{
				Y: -32,
			})
			c.mgr.AddComponent(e, &game.Scale{X: 2, Y: 2})

			start := c.field.Get7(3, 8)
			c.mgr.AddComponent(e, &game.Position{
				Center: game.Center{
					X: start.X(),
					Y: start.Y(),
				},
				Layer: 10,
			})

			c.mgr.AddComponent(e, &game.Obstacle{
				M:            3,
				N:            8,
				ObstacleType: game.LargeActor,
			})
		}

		c.mgr.AddComponent(e, &game.CombatStats{
			CurrentPreparation: 0,
		})
		c.mgr.AddComponent(e, &game.Facer{Face: geom.S})
	}
}

// End should be called at the resolution of a combat encounter. It removes
// combat-specific Components.
func (c *Combat) End() {
	removals := []string{
		"Sprite",
		"SpriteOffset",
		"Scale",
		"Position",
		"Obstacle",
		"CombatStats",
		"Facer",
		"TurnToken",
	}
	for _, e := range c.mgr.Get([]string{"Actor"}) {
		for _, comp := range removals {
			c.mgr.RemoveType(e, comp)
		}
	}
}

// Run a frame of this Combat.
func (c *Combat) Run(elapsed time.Duration) {
	switch c.state {
	case PreparingState:
		// Use the elapsed time as a base for the preparation increment.
		increment := int(elapsed.Seconds() * 500)

		// But if any Actor requires less than that, then only use that amount
		// instead, so that no actor overshoots its PreparationThreshold.
		for _, e := range c.mgr.Get([]string{"Actor", "CombatStats"}) {
			s := c.mgr.Component(e, "CombatStats").(*game.CombatStats)
			actor := c.mgr.Component(e, "Actor").(*game.Actor)

			if actor.PreparationThreshold-s.CurrentPreparation < increment {
				increment = actor.PreparationThreshold - s.CurrentPreparation
			}
		}

		// prepared Entities are the ones
		prepared := []ecs.Entity{}

		// Now that we know the increment, we can apply it with confidence that
		// we will not over-prepare.
		for _, e := range c.mgr.Get([]string{"Actor", "CombatStats"}) {
			s := c.mgr.Component(e, "CombatStats").(*game.CombatStats)
			actor := c.mgr.Component(e, "Actor").(*game.Actor)

			s.CurrentPreparation += increment

			if s.CurrentPreparation >= actor.PreparationThreshold {
				prepared = append(prepared, e)
			}
		}

		if len(prepared) > 0 {
			e := prepared[0]
			s := c.mgr.Component(e, "CombatStats").(*game.CombatStats)
			actor := c.mgr.Component(e, "Actor").(*game.Actor)

			fmt.Printf("Awaiting input from %v (%d exceeds %d)\n", e, s.CurrentPreparation, actor.PreparationThreshold)

			s.CurrentPreparation = 0
			c.state = AwaitingInputState
			c.mgr.AddComponent(e, &game.TurnToken{})
		}

	case ExecutingState:
		c.nav.Update(c.mgr, elapsed)
	}

	c.intents.Update()
}

func (c *Combat) checkHUD(x, y int) bool {
	for _, e := range c.mgr.Get([]string{"Interactive", "Position", "Sprite"}) {
		position := c.mgr.Component(e, "Position").(*game.Position)
		// Only going to handle Absolute Components for now I think
		if !position.Absolute {
			continue
		}
		interactive := c.mgr.Component(e, "Interactive").(*ui.Interactive)
		sprite := c.mgr.Component(e, "Sprite").(*game.Sprite)
		scale := c.mgr.Component(e, "Scale").(*game.Scale)

		// Because Absolutely positioned components might have negative
		// position, we need to modulo them.
		px := int(position.Center.X) % c.camera.screenW
		if px < 0 {
			px = px + c.camera.screenW
		}
		py := int(position.Center.Y) % c.camera.screenH
		if py < 0 {
			py = py + c.camera.screenH
		}

		// Is the x,y of the interaction without the bounds of the
		// Interactive?
		minX := px - int(scale.X*float64(sprite.W)*0.5)
		if x < minX {
			continue
		}
		maxX := minX + int(float64(sprite.W)*scale.X)
		if x > maxX {
			continue
		}
		minY := py - int(scale.Y*float64(sprite.H)*0.5)
		if y < minY {
			continue
		}
		maxY := minY + int(float64(sprite.H)*scale.Y)
		if y > maxY {
			continue
		}

		// Trigger the Interactive and return to prevent other interactions from occurring.
		interactive.Trigger()
		return true
	}
	return false
}

// Interaction is the way to notify the Combat that a mouse click or touch event
// occurred.
func (c *Combat) Interaction(x, y int) {
	if c.state == AwaitingInputState {
		if handled := c.checkHUD(x, y); handled {
			return
		}

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

	c.x = x
	c.y = y
	c.wx = wx
	c.wy = wy

	if c.state == AwaitingInputState {
		c.cursors.Clear()
		e := c.actorAwaitingInput()
		if e == 0 {
			return
		}
		a := c.mgr.Component(e, "Actor").(*game.Actor)
		position := c.mgr.Component(e, "Position").(*game.Position)

		switch a.Size {
		case game.SMALL:
			h := c.field.At(int(wx), int(wy))
			if h == nil {
				break
			}
			c.cursors.Add(h.X(), h.Y(), a.Size)
			c.cursors.Add(position.Center.X, position.Center.Y, a.Size)

		case game.MEDIUM:
			h := c.field.At4(int(wx), int(wy))
			if h == nil {
				break
			}
			c.cursors.Add(h.X(), h.Y(), a.Size)
			c.cursors.Add(position.Center.X, position.Center.Y, a.Size)

		case game.LARGE:
			h := c.field.At7(int(wx), int(wy))
			if h == nil {
				break
			}
			c.cursors.Add(h.X(), h.Y(), a.Size)
			c.cursors.Add(position.Center.X, position.Center.Y, a.Size)
		}
	}

}

func (c *Combat) actorAwaitingInput() ecs.Entity {
	entities := c.mgr.Get([]string{"Actor", "TurnToken"})
	if len(entities) == 0 {
		// FIXME: this is a flow error - there should always be an entity
		return 0
	}
	return entities[0]
}

// syncActorObstacle updates the an Actor's Obstacle to be synchronised with its
// position. It should be called when an Actor has completed a move.
func (c *Combat) syncActorObstacle(evt event.ActorMovementConcluded) {
	actor := c.mgr.Component(evt.Entity, "Actor").(*game.Actor)
	obstacle := c.mgr.Component(evt.Entity, "Obstacle").(*game.Obstacle)
	position := c.mgr.Component(evt.Entity, "Position").(*game.Position)

	switch actor.Size {
	case game.MEDIUM:
		h := c.field.At4(int(position.Center.X), int(position.Center.Y))
		obstacle.M = h.M
		obstacle.N = h.N
	case game.LARGE:
		h := c.field.At7(int(position.Center.X), int(position.Center.Y))
		obstacle.M = h.M
		obstacle.N = h.N
	default:
		h := c.field.At(int(position.Center.X), int(position.Center.Y))
		obstacle.M = h.M
		obstacle.N = h.N
	}

}

func (c *Combat) handleMovementConcluded(t event.Typer) {
	// FIXME: Should Obstacle movement be handled by an "obstacle" system instead?
	c.syncActorObstacle(t.(event.ActorMovementConcluded))

	c.state = AwaitingInputState
	c.MousePosition(c.x, c.y)
}

func (c *Combat) addGrass() {
	M, N := c.field.Dimensions()
	for n := 0; n < N; n++ {
		for m := 0; m < M; m++ {
			h := c.field.Get(m, n)
			e := c.mgr.NewEntity()

			c.mgr.AddComponent(e, &game.Sprite{
				Texture: "terrain.png",
				X:       0,
				Y:       0,
				W:       24,
				H:       16,
			})

			c.mgr.AddComponent(e, &game.Position{
				Center: game.Center{
					X: h.X(),
					Y: h.Y(),
				},
				Layer: 1,
			})
		}
	}
}

func (c *Combat) addTrees() {
	M, N := c.field.Dimensions()
	for n := 0; n < N; n++ {
		for m := 0; m < M; m++ {
			i := m + n*M
			h := c.field.Get(m, n)
			if i == 1 || i%17 == 1 || i%13 == 1 {
				e := c.mgr.NewEntity()
				c.mgr.AddComponent(e, &game.Sprite{
					Texture: "trees.png",
					X:       0,
					Y:       0,
					W:       24,
					H:       48,
				})
				c.mgr.AddComponent(e, &game.SpriteOffset{
					Y: -16,
				})
				c.mgr.AddComponent(e, &game.Position{
					Center: game.Center{
						X: h.X(),
						Y: h.Y(),
					},
					Layer: 10,
				})
				c.mgr.AddComponent(e, &game.Obstacle{
					M:            h.M,
					N:            h.N,
					ObstacleType: game.Tree,
				})
			}
		}
	}
}

func (c *Combat) constructHUD() {
	e := c.mgr.NewEntity()
	c.mgr.AddComponent(e, &game.Sprite{
		Texture: "hud.png",
		X:       16,
		Y:       0,
		W:       46,
		H:       14,
	})
	c.mgr.AddComponent(e, &game.Scale{
		X: 3,
		Y: 3,
	})
	c.mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: -80,
			Y: 32,
		},
		Layer:    100,
		Absolute: true,
	})

	c.mgr.AddComponent(e, &ui.Interactive{
		Trigger: func() {
			// Remove TurnToken from all actors.
			for _, e := range c.mgr.Get([]string{"Actor", "TurnToken"}) {
				c.mgr.RemoveComponent(e, c.mgr.Component(e, "TurnToken"))
			}

			// // Add turntoken to any actor.
			// c.mgr.AddComponent(c.mgr.Get([]string{"Actor"})[0], &game.TurnToken{})
			c.state = PreparingState
		},
	})
}
