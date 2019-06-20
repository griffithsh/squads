package combat

import (
	"time"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/geom"
	"github.com/griffithsh/squads/ui"
)

//go:generate stringer -type=State

// State enumerates the States that a Combat could be in.
type State int

const (
	// AwaitingInputState is when the combat is waiting for the local, human player to make a move.
	AwaitingInputState State = iota
	// ExecutingState is when a move or action is being played out by a character.
	ExecutingState
	// ThinkingState is when an AI-controller player is waiting to get command.
	ThinkingState
	// PreparingState is when no characters is prepared enough to make a move.
	PreparingState
)

// Manager is a game-mode. It processes turns-based Combat until one or the other
// team is knocked out.
type Manager struct {
	// Manager should own systems that are only relevant to Combat. A Turns coordinator, a preparation timer
	mgr    *ecs.World
	bus    *event.Bus
	field  *geom.Field
	nav    *game.Navigator
	camera *game.Camera
	state  State
	hud    *HUD

	x, y             int     // where the mouse last was in screen coordinates
	wx, wy           float64 // where the mouse last was in world coordinates
	screenW, screenH float64 // most recent dimensions of the window

	// actors ActorSystem...?
	cursors *game.CursorSystem
	intents *game.IntentSystem
}

// NewManager should accept two opposing squads of characters, a list of
// M,N,Terrain, and a tileset, and a set of environmental effects.
// Or it could accept a *ecs.World, which would by convention contain two squads
// of characters, and some other object that contains terrain, environmental
// effects etc.
func NewManager(mgr *ecs.World, camera *game.Camera /**/) *Manager {
	f, _ := geom.NewField(8, 24)
	bus := &event.Bus{}

	cm := Manager{
		mgr:     mgr,
		bus:     bus,
		field:   f,
		nav:     game.NewNavigator(bus),
		camera:  camera,
		state:   PreparingState,
		hud:     NewHUD(mgr, bus),
		cursors: game.NewCursorSystem(mgr),
		intents: game.NewIntentSystem(mgr, bus, f),
	}

	cm.bus.Subscribe(event.MovementConcluded, cm.handleMovementConcluded)
	cm.bus.Subscribe(event.EndTurnRequestedType, cm.handleEndTurnRequested)

	return &cm
}

// setState is the canonical way to change the CombatState.
func (cm *Manager) setState(state State) {
	if state == cm.state {
		return
	}
	ev := event.CombatStateTransition{
		Old: int(cm.state),
		New: int(state),
	}
	cm.state = state
	cm.bus.Publish(&ev)
}

// Begin should be called at the start of an engagement to set up components
// necessary for the combat.
func (cm *Manager) Begin() {
	/*
		At the start of Combat, we need to add a sprite and position component to
		every actor, because a Combat should be the thing responsible for deciding
		how to render an actor on the field.
	*/
	cm.camera.Center(cm.field.Width()/2, cm.field.Height()/2)

	cm.addGrass()
	cm.addTrees()

	// Upgrade all actors with components for visibility.
	entities := cm.mgr.Get([]string{"Actor"})
	for _, e := range entities {
		actor := cm.mgr.Component(e, "Actor").(*game.Actor)

		if actor.Size == game.SMALL {
			cm.mgr.AddComponent(e, &game.Sprite{
				Texture: "figure.png",
				X:       0,
				Y:       0,
				W:       24,
				H:       48,
			})
			cm.mgr.AddComponent(e, &game.SpriteOffset{
				Y: -16,
			})

			start := cm.field.Get(0, 0)
			cm.mgr.AddComponent(e, &game.Position{
				Center: game.Center{
					X: start.X(),
					Y: start.Y(),
				},
				Layer: 10,
			})

			cm.mgr.AddComponent(e, &game.Obstacle{
				M:            0,
				N:            0,
				ObstacleType: game.SmallActor,
			})
		} else if actor.Size == game.MEDIUM {
			cm.mgr.AddComponent(e, &game.Sprite{
				Texture: "wolf.png",
				X:       0,
				Y:       0,
				W:       58,
				H:       48,
			})
			cm.mgr.AddComponent(e, &game.SpriteOffset{
				Y: -4,
			})

			start := cm.field.Get4(0, 7)
			cm.mgr.AddComponent(e, &game.Position{
				Center: game.Center{
					X: start.X(),
					Y: start.Y(),
				},
				Layer: 10,
			})

			cm.mgr.AddComponent(e, &game.Obstacle{
				M:            0,
				N:            7,
				ObstacleType: game.MediumActor,
			})
		} else if actor.Size == game.LARGE {
			cm.mgr.AddComponent(e, &game.Sprite{
				Texture: "figure.png",
				X:       0,
				Y:       0,
				W:       24,
				H:       48,
			})
			cm.mgr.AddComponent(e, &game.SpriteOffset{
				Y: -32,
			})
			cm.mgr.AddComponent(e, &game.Scale{X: 2, Y: 2})

			start := cm.field.Get7(3, 8)
			cm.mgr.AddComponent(e, &game.Position{
				Center: game.Center{
					X: start.X(),
					Y: start.Y(),
				},
				Layer: 10,
			})

			cm.mgr.AddComponent(e, &game.Obstacle{
				M:            3,
				N:            8,
				ObstacleType: game.LargeActor,
			})
		}

		a := cm.mgr.Component(e, "Actor").(*game.Actor)
		cm.mgr.AddComponent(e, &game.CombatStats{
			CurrentPreparation: 0,
			ActionPoints:       a.ActionPoints,
		})
		cm.mgr.AddComponent(e, &game.Facer{Face: geom.S})
	}

	// Announce that the Combat has begun.
	cm.bus.Publish(event.CombatBegun{})
}

// End should be called at the resolution of a combat encounter. It removes
// combat-specific Components.
func (cm *Manager) End() {
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
	for _, e := range cm.mgr.Get([]string{"Actor"}) {
		for _, comp := range removals {
			cm.mgr.RemoveType(e, comp)
		}
	}

	// TODO: publish combat ended event.
}

// Run a frame of this Combat.
func (cm *Manager) Run(elapsed time.Duration) {
	switch cm.state {
	case PreparingState:
		// Use the elapsed time as a base for the preparation increment.
		increment := int(elapsed.Seconds() * 500)

		// But if any Actor requires less than that, then only use that amount
		// instead, so that no actor overshoots its PreparationThreshold.
		for _, e := range cm.mgr.Get([]string{"Actor", "CombatStats"}) {
			s := cm.mgr.Component(e, "CombatStats").(*game.CombatStats)
			actor := cm.mgr.Component(e, "Actor").(*game.Actor)

			if actor.PreparationThreshold-s.CurrentPreparation < increment {
				increment = actor.PreparationThreshold - s.CurrentPreparation
			}
		}

		// prepared captures all Actors who are fully prepared to take their
		// turn now.
		prepared := []ecs.Entity{}

		// Now that we know the increment, we can apply it with confidence that
		// we will not over-prepare.
		for _, e := range cm.mgr.Get([]string{"Actor", "CombatStats"}) {
			s := cm.mgr.Component(e, "CombatStats").(*game.CombatStats)
			actor := cm.mgr.Component(e, "Actor").(*game.Actor)

			s.CurrentPreparation += increment
			cm.bus.Publish(&event.CombatStatModified{
				Entity: e,
				Stat:   event.PrepStat,
				Amount: increment,
			})

			if s.CurrentPreparation >= actor.PreparationThreshold {
				prepared = append(prepared, e)
			}
		}

		// N.B. It's non-deterministic whose turn it is when multiple Actors
		// finish preparing at the same time.
		if len(prepared) > 0 {
			e := prepared[0]
			s := cm.mgr.Component(e, "CombatStats").(*game.CombatStats)

			ev := &event.CombatStatModified{
				Entity: e,
				Stat:   event.PrepStat,
				Amount: -s.CurrentPreparation,
			}
			s.CurrentPreparation = 0
			cm.bus.Publish(ev)

			cm.mgr.AddComponent(e, &game.TurnToken{})
			cm.setState(AwaitingInputState)
		}

	case ExecutingState:
		cm.nav.Update(cm.mgr, elapsed)
	}

	cm.intents.Update()
}

// checkHUD for interactions at x,y. Although this might sit better as a method
// of HUD, getting access to camera for Modulo is more awkward.
func (cm *Manager) checkHUD(x, y int) bool {
	for _, e := range cm.mgr.Get([]string{"Interactive", "Position", "Sprite"}) {
		position := cm.mgr.Component(e, "Position").(*game.Position)
		// Only going to handle Absolute Components for now I think
		if !position.Absolute {
			continue
		}
		interactive := cm.mgr.Component(e, "Interactive").(*ui.Interactive)
		sprite := cm.mgr.Component(e, "Sprite").(*game.Sprite)
		scale := cm.mgr.Component(e, "Scale").(*game.Scale)

		// Because Absolutely positioned components might have negative
		// position, we need to modulo them.
		px, py := cm.camera.Modulo(int(position.Center.X), int(position.Center.Y))

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

// Interaction is the way to notify the Combat Manager that a mouse click or
// touch event occurred.
func (cm *Manager) Interaction(x, y int) {
	if cm.state == AwaitingInputState {
		if handled := cm.checkHUD(x, y); handled {
			return
		}

		actor := cm.actorAwaitingInput()

		wx, wy := cm.camera.ScreenToWorld(x, y)

		i := game.MoveIntent{X: wx, Y: wy}
		cm.mgr.AddComponent(actor, &i)

		cm.setState(ExecutingState)

		cm.cursors.Clear()
	}
}

// MousePosition is the way to notify the Combat that the mouse has a new
// position.
func (cm *Manager) MousePosition(x, y int) {
	wx, wy := cm.camera.ScreenToWorld(x, y)

	cm.x = x
	cm.y = y
	cm.wx = wx
	cm.wy = wy

	if cm.state == AwaitingInputState {
		cm.cursors.Clear()
		e := cm.actorAwaitingInput()
		if e == 0 {
			return
		}
		a := cm.mgr.Component(e, "Actor").(*game.Actor)
		position := cm.mgr.Component(e, "Position").(*game.Position)

		switch a.Size {
		case game.SMALL:
			h := cm.field.At(int(wx), int(wy))
			if h == nil {
				break
			}
			cm.cursors.Add(h.X(), h.Y(), a.Size)
			cm.cursors.Add(position.Center.X, position.Center.Y, a.Size)

		case game.MEDIUM:
			h := cm.field.At4(int(wx), int(wy))
			if h == nil {
				break
			}
			cm.cursors.Add(h.X(), h.Y(), a.Size)
			cm.cursors.Add(position.Center.X, position.Center.Y, a.Size)

		case game.LARGE:
			h := cm.field.At7(int(wx), int(wy))
			if h == nil {
				break
			}
			cm.cursors.Add(h.X(), h.Y(), a.Size)
			cm.cursors.Add(position.Center.X, position.Center.Y, a.Size)
		}
	}

}

func (cm *Manager) actorAwaitingInput() ecs.Entity {
	entities := cm.mgr.Get([]string{"Actor", "TurnToken"})
	if len(entities) == 0 {
		// FIXME: this is a flow error - there should always be an entity
		return 0
	}
	return entities[0]
}

// syncActorObstacle updates the an Actor's Obstacle to be synchronised with its
// position. It should be called when an Actor has completed a move.
func (cm *Manager) syncActorObstacle(evt event.ActorMovementConcluded) {
	actor := cm.mgr.Component(evt.Entity, "Actor").(*game.Actor)
	obstacle := cm.mgr.Component(evt.Entity, "Obstacle").(*game.Obstacle)
	position := cm.mgr.Component(evt.Entity, "Position").(*game.Position)

	switch actor.Size {
	case game.MEDIUM:
		h := cm.field.At4(int(position.Center.X), int(position.Center.Y))
		obstacle.M = h.M
		obstacle.N = h.N
	case game.LARGE:
		h := cm.field.At7(int(position.Center.X), int(position.Center.Y))
		obstacle.M = h.M
		obstacle.N = h.N
	default:
		h := cm.field.At(int(position.Center.X), int(position.Center.Y))
		obstacle.M = h.M
		obstacle.N = h.N
	}

}

func (cm *Manager) handleMovementConcluded(t event.Typer) {
	// FIXME: Should Obstacle movement be handled by an "obstacle" system instead?
	cm.syncActorObstacle(t.(event.ActorMovementConcluded))

	cm.setState(AwaitingInputState)
	cm.MousePosition(cm.x, cm.y)
}

func (cm *Manager) handleEndTurnRequested(t event.Typer) {
	// Remove TurnToken from all actors.
	for _, e := range cm.mgr.Get([]string{"Actor", "TurnToken"}) {
		// Reset to maximum AP.
		actor := cm.mgr.Component(e, "Actor").(*game.Actor)
		stats := cm.mgr.Component(e, "CombatStats").(*game.CombatStats)
		stats.ActionPoints = actor.ActionPoints

		// And then remove TurnToken.
		cm.mgr.RemoveComponent(e, cm.mgr.Component(e, "TurnToken"))
	}

	cm.setState(PreparingState)
}

func (cm *Manager) addGrass() {
	M, N := cm.field.Dimensions()
	for n := 0; n < N; n++ {
		for m := 0; m < M; m++ {
			h := cm.field.Get(m, n)
			e := cm.mgr.NewEntity()

			cm.mgr.AddComponent(e, &game.Sprite{
				Texture: "terrain.png",
				X:       0,
				Y:       0,
				W:       24,
				H:       16,
			})

			cm.mgr.AddComponent(e, &game.Position{
				Center: game.Center{
					X: h.X(),
					Y: h.Y(),
				},
				Layer: 1,
			})
		}
	}
}

func (cm *Manager) addTrees() {
	M, N := cm.field.Dimensions()
	for n := 0; n < N; n++ {
		for m := 0; m < M; m++ {
			i := m + n*M
			h := cm.field.Get(m, n)
			if i == 1 || i%17 == 1 || i%13 == 1 {
				e := cm.mgr.NewEntity()
				cm.mgr.AddComponent(e, &game.Sprite{
					Texture: "trees.png",
					X:       0,
					Y:       0,
					W:       24,
					H:       48,
				})
				cm.mgr.AddComponent(e, &game.SpriteOffset{
					Y: -16,
				})
				cm.mgr.AddComponent(e, &game.Position{
					Center: game.Center{
						X: h.X(),
						Y: h.Y(),
					},
					Layer: 10,
				})
				cm.mgr.AddComponent(e, &game.Obstacle{
					M:            h.M,
					N:            h.N,
					ObstacleType: game.Tree,
				})
			}
		}
	}
}
