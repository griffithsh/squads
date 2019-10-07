package combat

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
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
	// SelectingTargetState is when the local player is picking a hex to use a skill on.
	SelectingTargetState
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
	mgr     *ecs.World
	bus     *event.Bus
	field   *geom.Field
	nav     *Navigator
	camera  *game.Camera
	hud     *HUD
	cursors *CursorManager

	turnToken ecs.Entity // Whose turn is it?
	state     State

	incrementAccumulator float64

	x, y             int       // where the mouse last was in screen coordinates
	wx, wy           float64   // where the mouse last was in world coordinates
	screenW, screenH float64   // most recent dimensions of the window
	selectedHex      *geom.Key // most recent hex selected

	intents      *IntentSystem
	performances *PerformanceSystem

	dormant bool
}

// NewManager should accept two opposing squads of characters, a list of
// M,N,Terrain, and a tileset, and a set of environmental effects.
// Or it could accept a *ecs.World, which would by convention contain two squads
// of characters, and some other object that contains terrain, environmental
// effects etc.
func NewManager(mgr *ecs.World, camera *game.Camera, bus *event.Bus) *Manager {
	f, _ := geom.NewField(8, 24)

	cm := Manager{
		mgr:    mgr,
		bus:    bus,
		field:  f,
		nav:    NewNavigator(bus),
		camera: camera,
		// state:   TODO: some-uninitialised-state
		hud:          NewHUD(mgr, bus, camera.GetW(), camera.GetH()),
		cursors:      NewCursorManager(mgr, bus, f),
		intents:      NewIntentSystem(mgr, bus, f),
		performances: NewPerformanceSystem(mgr, bus),

		dormant: false,
	}
	cm.setState(PreparingState)

	cm.bus.Subscribe(game.CombatActorMovementConcluded{}.Type(), cm.handleMovementConcluded)
	cm.bus.Subscribe(EndTurnRequested{}.Type(), cm.handleEndTurnRequested)
	cm.bus.Subscribe(MoveModeRequested{}.Type(), cm.handleMoveModeRequested)
	cm.bus.Subscribe(CancelSkillRequested{}.Type(), cm.handleCancelSkillRequested)

	return &cm
}

// setState is the canonical way to change the CombatState.
func (cm *Manager) setState(state State) {
	if state == cm.state {
		return
	}
	ev := StateTransition{
		Old: cm.state,
		New: state,
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

	// TODO:
	// There is some entity which stores info about a "level", and produces artifacts that can be used by the combat Manager.
	// It should produce the shape of the level, and the terrain of each hex (grass, water, blocked by tree etc).
	// It should also produce starting positions for teams...
	// Some other entity should produce an opponent team for the player's squad to fight _on_ this level.

	// semiSort provides the list of Hexes in the field roughly sorted by their
	// distance from m,n. It intends to provide randomish starting locations.
	semiSort := func(m, n int, f *geom.Field) []*geom.Hex {
		type s struct {
			distance float64
			h        *geom.Hex
		}
		start := geom.Hex{M: m, N: n}
		distances := make([]s, len(f.Hexes()))

		for i, h := range f.Hexes() {
			distances[i] = s{math.Pow(math.Abs(h.X()-start.X()), 2) + math.Pow(math.Abs(h.Y()-start.Y()), 2), h}
		}
		sort.Slice(distances, func(i, j int) bool {
			return distances[i].distance < distances[j].distance
		})

		// bucket the hexes into small groups, and shuffle the hexes within
		// each group. This is going to keep the nearest together, but still
		// not always pick the same places every time.
		bucket := 25
		gi := 0 // global index
		for {
			rand.Shuffle(bucket, func(i, j int) {
				distances[i+gi], distances[j+gi] = distances[j+gi], distances[i+gi]
			})
			gi += bucket
			if gi+bucket >= len(distances) {
				break
			}
		}

		result := make([]*geom.Hex, len(distances))
		for i := range distances {
			result[i] = distances[i].h
		}
		return result
	}

	// List of start locations.
	levelStarts := []geom.Key{
		{M: 6, N: 18},
		{M: 2, N: 8},
	}
	rand.Shuffle(len(levelStarts), func(i, j int) {
		levelStarts[i], levelStarts[j] = levelStarts[j], levelStarts[i]
	})
	usedStarts := map[int64][]*geom.Hex{}

	// Then each team takes a turn placing an Actor from largest to smallest,
	// working through the semi-shuffled list of results.

	// isBlocked determines if an Actor with an ActorSize of sz can be placed at m,n.
	isBlocked := func(m, n int, sz game.CharacterSize, mgr *ecs.World) bool {
		// blockages is a set of Keys that are taken by other things
		blockages := map[geom.Key]struct{}{}
		for _, e := range mgr.Get([]string{"Obstacle"}) {
			o := mgr.Component(e, "Obstacle").(*game.Obstacle)

			h := game.AdaptFieldObstacle(cm.field, o.ObstacleType).Get(o.M, o.N)
			if h == nil {
				panic(fmt.Sprintf("there is no hex where Obstacle(%d,%s) is present (%d,%d)", e, o.ObstacleType, o.M, o.N))
			}

			// FIXME: We're making the assumption again here that all obstacles
			// are total obstacles. Even conceptually things like shallow water
			// or bushes that should only impede movement slightly.
			for _, h := range h.Hexes() {
				blockages[h.Key()] = struct{}{}
			}
		}

		hex := game.AdaptField(cm.field, sz).Get(m, n)
		if hex == nil {
			return true
		}

		// occupy is the list of Hexes an Actor with sz and m,n will occupy.
		occupy := hex.Hexes()

		for _, h := range occupy {
			if h == nil {
				return true
			}
			if _, blocked := blockages[geom.Key{M: h.M, N: h.N}]; blocked {
				return true
			}
		}

		return false
	}

	// In game, something like this process should occur when additional Actors are summoned.
	// Necromancers summon Skeletons (this could be ground targeted with a range)
	// Gemini auto-summon their twin at the start of combat (this sounds more like what's happening here)
	// Druids summon beasts (ground targeted again)

	// Upgrade all Actors with components for combat.
	entities := cm.mgr.Get([]string{"Actor"})
	for _, e := range entities {
		// This shows something that is perhaps an architectural flaw - who
		// does the Actor really belong to here? Having a list of Actors is
		// relevant when we want to show a combat board and when we want to
		// serialise a savegame. Is it appropriate that the combat really
		// "owns" the Actors?
		cm.mgr.Tag(e, "combat")

		actor := cm.mgr.Component(e, "Actor").(*Actor)
		team := cm.mgr.Component(e, "Team").(*Team)

		if _, ok := usedStarts[team.ID]; !ok {
			s := levelStarts[len(usedStarts)]
			usedStarts[team.ID] = semiSort(s.M, s.N, cm.field)
		}
		nearbys := usedStarts[team.ID]

		// Pick positions for the Actors.
		for _, h := range nearbys {
			if isBlocked(h.M, h.N, actor.Size, cm.mgr) {
				continue
			}

			f := game.AdaptField(cm.field, actor.Size)
			start := f.Get(h.M, h.N)
			cm.mgr.AddComponent(e, &game.Position{
				Center: game.Center{
					X: start.X(),
					Y: start.Y(),
				},
				Layer: 10,
			})

			o := game.Obstacle{
				M:            h.M,
				N:            h.N,
				ObstacleType: game.SmallActor,
			}
			switch actor.Size {
			case game.MEDIUM:
				o.ObstacleType = game.MediumActor
			case game.LARGE:
				o.ObstacleType = game.LargeActor
			}
			cm.mgr.AddComponent(e, &o)

			break
		}

		actor.PreparationThreshold.Cur = 0
		actor.ActionPoints.Cur = actor.ActionPoints.Max
		cm.mgr.AddComponent(e, &game.Facer{Face: geom.S})
	}

	// Announce that the Combat has begun.
	cm.bus.Publish(&game.CombatBegan{})
}

// Enable the combat Manager, responding to input and rendering the state of the
// combat.
func (cm *Manager) Enable() {
	if cm.dormant {
		for _, e := range cm.mgr.Tagged("combat") {
			cm.mgr.RemoveComponent(e, &game.Hidden{})
			cm.hud.Enable()
		}
		cm.dormant = false
	}
}

// Disable the combat Manager, ignoring input and not rendering the state of the
// combat.
func (cm *Manager) Disable() {
	if !cm.dormant {
		for _, e := range cm.mgr.Tagged("combat") {
			cm.mgr.AddComponent(e, &game.Hidden{})
			cm.hud.Disable()
		}
		cm.dormant = true
	}
}

// End should be called at the resolution of a combat encounter. It removes
// combat-specific Components and Entities.
func (cm *Manager) End() {
	// TODO: When there are summoned units, then in Manager.End(), we will need
	// to remove them. Potentialy a new Component called "Impermanent", or
	// "Summoned" could be added to these, and we would need to Destroy these
	// Entities before removing combat-only Components from the other Actors.

	removals := []string{
		"Sprite",
		"RenderOffset",
		"Scale",
		"Position",
		"Obstacle",
		"Facer",
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
	if cm.dormant {
		return
	}

	switch cm.state {
	case PreparingState:
		// Use the elapsed time as a base for the preparation increment.
		const prepPerSec float64 = 500
		cm.incrementAccumulator += elapsed.Seconds() * prepPerSec
		increment := int(cm.incrementAccumulator)
		cm.incrementAccumulator -= float64(increment)

		// But if any Actor requires less than that, then only use that amount
		// instead, so that no actor overshoots its PreparationThreshold.
		for _, e := range cm.mgr.Get([]string{"Actor"}) {
			actor := cm.mgr.Component(e, "Actor").(*Actor)

			if actor.PreparationThreshold.Max-actor.PreparationThreshold.Cur < increment {
				increment = actor.PreparationThreshold.Max - actor.PreparationThreshold.Cur
			}
		}

		// prepared captures all Actors who are fully prepared to take their
		// turn now.
		prepared := []ecs.Entity{}

		// Now that we know the increment, we can apply it with confidence that
		// we will not over-prepare.
		for _, e := range cm.mgr.Get([]string{"Actor"}) {
			actor := cm.mgr.Component(e, "Actor").(*Actor)

			actor.PreparationThreshold.Cur += increment
			cm.bus.Publish(&StatModified{
				Entity: e,
				Stat:   game.PrepStat,
				Amount: increment,
			})

			if actor.PreparationThreshold.Cur >= actor.PreparationThreshold.Max {
				prepared = append(prepared, e)
			}
		}

		// N.B. It's non-deterministic whose turn it is when multiple Actors
		// finish preparing at the same time.
		if len(prepared) > 0 {
			e := prepared[0]
			actor := cm.mgr.Component(e, "Actor").(*Actor)

			ev := &StatModified{
				Entity: e,
				Stat:   game.PrepStat,
				Amount: -actor.PreparationThreshold.Cur,
			}
			actor.PreparationThreshold.Cur = 0
			cm.bus.Publish(ev)

			cm.turnToken = e
			cm.bus.Publish(&ActorTurnChanged{Entity: cm.turnToken})
			cm.setState(AwaitingInputState)
		}

	case ExecutingState:
		cm.nav.Update(cm.mgr, elapsed)
	}

	cm.intents.Update()
	cm.performances.Update(elapsed)
	cm.hud.Update(elapsed)
	cm.cursors.Update(elapsed)
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
	if cm.dormant {
		return
	}
	if cm.state == AwaitingInputState {
		cm.checkHUD(x, y)
	} else if cm.state == SelectingTargetState {
		if handled := cm.checkHUD(x, y); handled {
			return
		}

		wx, wy := cm.camera.ScreenToWorld(x, y)

		i := game.MoveIntent{X: wx, Y: wy}
		cm.mgr.AddComponent(cm.turnToken, &i)

		cm.setState(ExecutingState)
	}
}

// MousePosition is the way to notify the Combat that the mouse has a new
// position.
func (cm *Manager) MousePosition(x, y int) {
	wx, wy := cm.camera.ScreenToWorld(x, y)

	if cm.state == SelectingTargetState {
		// When we're selecting a target, we need to highlight some hexes to
		// show where we're targeting.
		// If the change in position means we're positioned over a new hex,
		// then publish a DifferentHexSelected event.

		// The consumer needs to make a decision about what to repaint now that
		// the hex that the mouse is hovering over has changed. It might be a
		// path of hexes because we're selecting a place to move to, or it might
		// be a glob of hexes because we're targeting an AoE fireball spell etc.
		actor := cm.mgr.Component(cm.turnToken, "Actor").(*Actor)
		var newSelected *geom.Key

		f := game.AdaptField(cm.field, actor.Size)
		h := f.At(int(wx), int(wy))
		if h != nil {
			k := h.Key()
			newSelected = &geom.Key{
				M: k.M,
				N: k.N,
			}
		}

		if newSelected != nil && cm.selectedHex != nil {
			if *newSelected != *cm.selectedHex {
				cm.selectedHex = newSelected

				cm.bus.Publish(&DifferentHexSelected{
					K: cm.selectedHex,
				})

			}
		} else if newSelected != cm.selectedHex {
			cm.selectedHex = newSelected

			cm.bus.Publish(&DifferentHexSelected{
				K: cm.selectedHex,
			})
		}
	}

	// Update local cached values
	cm.x = x
	cm.y = y
	cm.wx = wx
	cm.wy = wy
}

// syncActorObstacle updates the an Actor's Obstacle to be synchronised with its
// position. It should be called when an Actor has completed a move.
func (cm *Manager) syncActorObstacle(evt *game.CombatActorMovementConcluded) {
	actor := cm.mgr.Component(evt.Entity, "Actor").(*Actor)
	obstacle := cm.mgr.Component(evt.Entity, "Obstacle").(*game.Obstacle)
	position := cm.mgr.Component(evt.Entity, "Position").(*game.Position)

	h := game.AdaptField(cm.field, actor.Size).At(int(position.Center.X), int(position.Center.Y))
	k := h.Key()
	obstacle.M = k.M
	obstacle.N = k.N
}

func (cm *Manager) handleMovementConcluded(t event.Typer) {
	// FIXME: Should Obstacle movement be handled by an "obstacle" system instead?
	cm.syncActorObstacle(t.(*game.CombatActorMovementConcluded))

	cm.setState(AwaitingInputState)
	cm.MousePosition(cm.x, cm.y)
}

func (cm *Manager) handleEndTurnRequested(event.Typer) {
	// Reset to maximum AP.
	actor := cm.mgr.Component(cm.turnToken, "Actor").(*Actor)
	actor.ActionPoints.Cur = actor.ActionPoints.Max

	// Remove turnToken
	cm.turnToken = 0
	cm.bus.Publish(&ActorTurnChanged{Entity: cm.turnToken})

	cm.setState(PreparingState)
}

func (cm *Manager) handleMoveModeRequested(event.Typer) {
	cm.setState(SelectingTargetState)
}

func (cm *Manager) handleCancelSkillRequested(event.Typer) {
	cm.setState(AwaitingInputState)
}

func (cm *Manager) addGrass() {
	M, N := cm.field.Dimensions()
	for n := 0; n < N; n++ {
		for m := 0; m < M; m++ {
			h := cm.field.Get(m, n)
			e := cm.mgr.NewEntity()
			cm.mgr.Tag(e, "combat")

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
			if m == 4 && n == 14 {
				continue
			}
			i := m + n*M
			h := cm.field.Get(m, n)
			if i%17 == 1 || i%23 == 1 {
				e := cm.mgr.NewEntity()
				cm.mgr.Tag(e, "combat")
				cm.mgr.AddComponent(e, &game.Sprite{
					Texture: "trees.png",
					X:       0,
					Y:       0,
					W:       24,
					H:       48,
					OffsetY: -16,
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
