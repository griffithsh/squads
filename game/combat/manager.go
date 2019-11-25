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
	// Uninitialised is that default state of the combat.Manager.
	Uninitialised State = iota
	// AwaitingInputState is when the combat is waiting for the local, human player to make a move.
	AwaitingInputState
	// SelectingTargetState is when the local player is picking a hex to use a skill on.
	SelectingTargetState
	// ExecutingState is when a move or action is being played out by a character.
	ExecutingState
	// ThinkingState is when an AI-controller player is waiting to get command.
	ThinkingState
	// PreparingState is when no characters is prepared enough to make a move.
	PreparingState
	// Celebration occurs when there is only one team left, and they are
	// celebrating their victory.
	Celebration
	// FadingIn is when the combat is first starting, or returning from a menu,
	// and the curtain that obscures the scene change is disappearing.
	FadingIn
	//FadingOut is when the combat is going to another scene, and the curtain
	//that obscures the scene change is appearing.
	FadingOut
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

	turnToken            ecs.Entity // Whose turn is it? References an existing Entity.
	selectingInteractive ecs.Entity // catches clicks on the field.

	// Manager has both a state and a paused flag, so that state transition
	// logic can be isolated from pausing.
	paused bool
	state  State

	incrementAccumulator float64

	x, y             int       // where the mouse last was in screen coordinates
	wx, wy           float64   // where the mouse last was in world coordinates
	screenW, screenH float64   // most recent dimensions of the window
	selectedHex      *geom.Key // most recent hex selected

	intents      *IntentSystem
	performances *PerformanceSystem

	celebrations time.Duration

	squads []ecs.Entity
}

// NewManager creates a new combat Manager.
func NewManager(mgr *ecs.World, camera *game.Camera, bus *event.Bus) *Manager {
	f := geom.NewField()

	cm := Manager{
		mgr:                  mgr,
		bus:                  bus,
		field:                f,
		nav:                  NewNavigator(bus),
		camera:               camera,
		state:                Uninitialised,
		hud:                  NewHUD(mgr, bus, camera.GetW(), camera.GetH()),
		cursors:              NewCursorManager(mgr, bus, f),
		selectingInteractive: mgr.NewEntity(),
		intents:              NewIntentSystem(mgr, bus, f),
		performances:         NewPerformanceSystem(mgr, bus),

		paused: false,
	}

	cm.bus.Subscribe(ParticipantMovementConcluded{}.Type(), cm.handleMovementConcluded)
	cm.bus.Subscribe(EndTurnRequested{}.Type(), cm.handleEndTurnRequested)
	cm.bus.Subscribe(MoveModeRequested{}.Type(), cm.handleMoveModeRequested)
	cm.bus.Subscribe(CancelSkillRequested{}.Type(), cm.handleCancelSkillRequested)
	cm.bus.Subscribe(AttemptingEscape{}.Type(), cm.handleAttemptingEscape)

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

	// When entering Selecting Target State, we need to add an Interactive to
	// cover all areas of the field, so that we can convert those clicks to
	// MoveIntents.
	if state == SelectingTargetState {
		// Using the max float value as the size and a position of 0,0 should
		// work in all cases, and it's a lot faster than figuring out the actual
		// dimensions of the field. The goal here is to catch *any* clicks in
		// the world, remember?
		cm.mgr.AddComponent(cm.selectingInteractive, &game.Position{})
		cm.mgr.AddComponent(cm.selectingInteractive, &ui.Interactive{
			W: math.MaxFloat64, H: math.MaxFloat64,
			Trigger: func(x, y float64) {
				cm.mgr.AddComponent(cm.turnToken, &game.MoveIntent{X: x, Y: y})
				cm.setState(ExecutingState)
			},
		})
	} else if ev.Old == SelectingTargetState {
		cm.mgr.RemoveComponent(cm.selectingInteractive, &ui.Interactive{})
	}

	cm.bus.Publish(&ev)
}

// semiSort provides the list of Hexes in the field roughly sorted by their
// distance from m,n. It intends to provide randomish starting locations.
func semiSort(m, n int, f *geom.Field) []*geom.Hex {
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

// isBlocked determines if a Character with a CharacterSize of sz can be placed at m,n.
func isBlocked(field *geom.Field, m, n int, sz game.CharacterSize, mgr *ecs.World) bool {
	// blockages is a set of Keys that are taken by other things
	blockages := map[geom.Key]struct{}{}
	for _, e := range mgr.Get([]string{"Obstacle"}) {
		o := mgr.Component(e, "Obstacle").(*game.Obstacle)

		h := game.AdaptFieldObstacle(field, o.ObstacleType).Get(o.M, o.N)
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

	hex := game.AdaptField(field, sz).Get(m, n)
	if hex == nil {
		return true
	}

	// occupy is the list of Hexes a Character with sz and m,n will occupy.
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

// startProvider provides randomised starting hexes for combat Participants.
// In game, something like this process should occur when additional
// Participants are summoned. Necromancers summon Skeletons (this could be
// ground targeted with a range) Gemini auto-summon their twin at the start of
// combat (this sounds more like what's happening here) Druids summon beasts
// (ground targeted again)
type startProvider struct {
	starts []geom.Key
	used   map[int64][]*geom.Hex
}

func newStartProvider(starts []geom.Key) *startProvider {
	rand.Shuffle(len(starts), func(i, j int) {
		starts[i], starts[j] = starts[j], starts[i]
	})
	return &startProvider{
		starts: starts,
		used:   map[int64][]*geom.Hex{},
	}
}

func (sp *startProvider) getNearby(team *game.Team, f *geom.Field) []*geom.Hex {
	if _, ok := sp.used[team.ID]; !ok {
		s := sp.starts[len(sp.used)]
		sp.used[team.ID] = semiSort(s.M, s.N, f)
	}
	return sp.used[team.ID]
}

func (cm *Manager) getStart(sz game.CharacterSize, nearbys []*geom.Hex) *geom.Hex {
	for _, h := range nearbys {
		if isBlocked(cm.field, h.M, h.N, sz, cm.mgr) {
			continue
		}

		return h
	}
	return nil
}

// CToP converts a Character to a Participant.
func CToP(char *game.Character) *Participant {
	return &Participant{
		Name:       char.Name,
		Level:      char.Level,
		SmallIcon:  char.SmallIcon,
		BigIcon:    char.BigIcon,
		Size:       char.Size,
		Profession: char.Profession,
		Sex:        char.Sex,
		PreparationThreshold: CurMax{
			Max: char.PreparationThreshold,
		},
		ActionPoints: CurMax{
			Max: char.ActionPoints,
		},
		// Health: CurMax{
		// 	Cur: char.CurrentHealth,
		// 	Max: char.MaxHealth,
		// },
		// Strength:     int(char.StrengthPerLevel * float64(char.Level)),
		// Dexterity:    0,
		// Intelligence: 0,
		// Vitality:     0,
	}
}

// createParticipation adds a new Entity to participate in combat based on a Character.
func (cm *Manager) createParticipation(charEntity ecs.Entity, char *game.Character, team *game.Team, h *geom.Hex) {
	e := cm.mgr.NewEntity()
	cm.mgr.Tag(e, "combat")

	// Add Participant Component.
	participant := CToP(char)
	participant.Character = charEntity
	participant.ActionPoints.Cur = participant.ActionPoints.Max
	participant.Status = Alive
	cm.mgr.AddComponent(e, participant)

	// Add Team.
	cm.mgr.AddComponent(e, team)

	// Add Position.
	f := game.AdaptField(cm.field, char.Size)
	start := f.Get(h.M, h.N)
	cm.mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: start.X(),
			Y: start.Y(),
		},
		Layer: 10,
	})

	// Add Obstacle.
	o := game.Obstacle{
		M:            h.M,
		N:            h.N,
		ObstacleType: game.SmallCharacter,
	}
	switch char.Size {
	case game.MEDIUM:
		o.ObstacleType = game.MediumCharacter
	case game.LARGE:
		o.ObstacleType = game.LargeCharacter
	}
	cm.mgr.AddComponent(e, &o)

	// Add Facer Component.
	cm.mgr.AddComponent(e, &game.Facer{Face: geom.S})
}

// Begin should be called at the start of an engagement to set up components
// necessary for the combat.
// TODO: Begin should be provided with information about the squads and the
// terrain to fight on.
func (cm *Manager) Begin(participatingSquads []ecs.Entity) {
	var keys []geom.Key = geom.MByN(rand.Intn(3)+6, rand.Intn(7)+20)
	cm.field.Load(keys)
	cm.setState(FadingIn)
	e := cm.mgr.NewEntity()
	cm.mgr.Tag(e, "combat")
	cm.mgr.AddComponent(e, &game.DiagonalMatrixWipe{
		W: 1024, H: 768, // FIXME: need access to screen size
		Obscuring: false, // ergo revealing
		OnComplete: func() {
			cm.setState(PreparingState)
		},
		OnInitialised: func() {
			// Debug hacky code to add some terrain.
			cm.addGrass()
			cm.addTrees()

			// TODO:
			// There is some entity which stores info about a "level", and produces
			// artifacts that can be used by the combat Manager. It should produce the
			// shape of the level, and the terrain of each hex (grass, water, blocked by
			// tree etc). It should also produce starting positions for teams... Some
			// other entity should produce an opponent team for the player's squad to
			// fight _on_ this level.

			// FIXME: Hard-coded list of start locations.
			sp := newStartProvider([]geom.Key{
				{M: 6, N: 18},
				{M: 2, N: 8},
			})

			entities := []ecs.Entity{}
			for _, e := range participatingSquads {
				cm.squads = append(cm.squads, e)
				squad := cm.mgr.Component(e, "Squad").(*game.Squad)
				entities = append(entities, squad.Members...)
			}

			// Create a Participating Entity for every Character we have.
			for _, e := range entities {
				team := cm.mgr.Component(e, "Team").(*game.Team)
				char := cm.mgr.Component(e, "Character").(*game.Character)
				near := sp.getNearby(team, cm.field)
				h := cm.getStart(char.Size, near)

				cm.createParticipation(e, char, team, h)
			}

			// Announce that the Combat has begun.
			cm.bus.Publish(&game.CombatBegan{})
		},
	})
	cm.camera.Center(cm.field.Width()/2, cm.field.Height()/2)
}

// End should be called at the resolution of a combat encounter. It removes
// combat-specific Components and Entities.
func (cm *Manager) End() {
	// It is no-one's turn.
	cm.turnToken = 0
	cm.setState(Uninitialised)
	cm.squads = cm.squads[:0]

	// Destroy Entities that were added for combat.
	for _, e := range cm.mgr.Tagged("combat") {
		cm.mgr.DestroyEntity(e)
	}

	// Remove Components that are only relevant to combat.
	removals := []string{
		"Obstacle",
		"Facer",
		"Participant",
	}
	for _, e := range cm.mgr.Get([]string{"Participant"}) {
		for _, comp := range removals {
			cm.mgr.RemoveType(e, comp)
		}
	}
}

// Pause the combat Manager, ignoring input and not rendering the state of the
// combat. Pause should be called when an in-combat modal menu is entered, and a
// return to the current combat is imminent.
func (cm *Manager) Pause() {
	if !cm.paused {
		for _, e := range cm.mgr.Tagged("combat") {
			cm.mgr.AddComponent(e, &game.Hidden{})
			cm.hud.Disable()
		}
		cm.paused = true
	}
}

// Unpause the combat Manager, responding to input and rendering the state of the
// combat.
func (cm *Manager) Unpause() {
	if cm.paused {
		for _, e := range cm.mgr.Tagged("combat") {
			cm.mgr.RemoveComponent(e, &game.Hidden{})
			cm.hud.Enable()
		}
		cm.paused = false
	}
}

// Run a frame of this Combat.
func (cm *Manager) Run(elapsed time.Duration) {
	if cm.paused {
		return
	}

	// Do a check for a victory condition.
	if cm.state == PreparingState || cm.state == AwaitingInputState {
		remainingTeams := map[int64]struct{}{}
		victoriousEntities := []ecs.Entity{}
		for _, e := range cm.mgr.Get([]string{"Participant", "Team"}) {
			participant := cm.mgr.Component(e, "Participant").(*Participant)
			if participant.Status != Alive {
				continue
			}
			team := cm.mgr.Component(e, "Team").(*game.Team)
			remainingTeams[team.ID] = struct{}{}
			victoriousEntities = append(victoriousEntities, e)
		}
		if len(remainingTeams) < 2 {
			// TODO: set Victory, Escape, Defeat banner
			// ...

			for _, e := range victoriousEntities {
				cm.bus.Publish(&CharacterCelebrating{Entity: e})
			}

			cm.setState(Celebration)
			return
		}
	}

	switch cm.state {
	case PreparingState:
		// Use the elapsed time as a base for the preparation increment.
		const prepPerSec float64 = 500
		cm.incrementAccumulator += elapsed.Seconds() * prepPerSec
		increment := int(cm.incrementAccumulator)
		cm.incrementAccumulator -= float64(increment)

		// But if any Character requires less than that, then only use that amount
		// instead, so that no Character overshoots its PreparationThreshold.
		for _, e := range cm.mgr.Get([]string{"Participant"}) {
			participant := cm.mgr.Component(e, "Participant").(*Participant)

			if participant.Status != Alive {
				continue
			}

			if participant.PreparationThreshold.Max-participant.PreparationThreshold.Cur < increment {
				increment = participant.PreparationThreshold.Max - participant.PreparationThreshold.Cur
			}
		}

		// prepared captures all Participants who are fully prepared to take their
		// turn now.
		prepared := []ecs.Entity{}

		// Now that we know the increment, we can apply it with confidence that
		// we will not over-prepare.
		for _, e := range cm.mgr.Get([]string{"Participant"}) {
			participant := cm.mgr.Component(e, "Participant").(*Participant)
			if participant.Status != Alive {
				continue
			}

			participant.PreparationThreshold.Cur += increment
			cm.bus.Publish(&StatModified{
				Entity: e,
				Stat:   game.PrepStat,
				Amount: increment,
			})

			if participant.PreparationThreshold.Cur >= participant.PreparationThreshold.Max {
				prepared = append(prepared, e)
			}
		}

		// N.B. It's non-deterministic whose turn it is when multiple
		// Participants finish preparing at the same time.
		if len(prepared) > 0 {
			e := prepared[0]
			participant := cm.mgr.Component(e, "Participant").(*Participant)

			ev := &StatModified{
				Entity: e,
				Stat:   game.PrepStat,
				Amount: -participant.PreparationThreshold.Cur,
			}
			participant.PreparationThreshold.Cur = 0
			cm.bus.Publish(ev)

			cm.turnToken = e
			cm.bus.Publish(&ParticipantTurnChanged{Entity: cm.turnToken})
			cm.setState(AwaitingInputState)
		}

	case ExecutingState:
		cm.nav.Update(cm.mgr, elapsed)
	case Celebration:
		// Celebrate for a time ...
		cm.celebrations += elapsed
		if cm.celebrations > time.Second*2 {
			cm.celebrations = 0
			cm.mgr.AddComponent(cm.mgr.NewEntity(), &game.DiagonalMatrixWipe{
				W: 1024, H: 768, // FIXME: need access to screen dimensions!
				Obscuring: true,
				OnComplete: func() {
					cc := game.CombatConcluded{
						Results: map[ecs.Entity]game.CombatResult{},
					}
					for _, squadEntity := range cm.squads {
						squad := cm.mgr.Component(squadEntity, "Squad").(*game.Squad)
						cc.Results[squadEntity] = game.Defeated
						for _, e1 := range cm.mgr.Get([]string{"Participant"}) {
							participant := cm.mgr.Component(e1, "Participant").(*Participant)

							for _, e2 := range squad.Members {
								if participant.Character == e2 {
									participant := cm.mgr.Component(e1, "Participant").(*Participant)
									if participant.Status == Alive {
										cc.Results[squadEntity] = game.Victorious
										goto nextSquad
									} else if participant.Status == Escaped {
										cc.Results[squadEntity] = game.Escaped
									}

								}
							}
						}
					nextSquad:
					}
					cm.bus.Publish(&cc)
					for _, e := range cm.mgr.Get([]string{"Token", "Position"}) {
						if !cm.mgr.HasTag(e, "player") {
							continue
						}
						p := cm.mgr.Component(e, "Position").(*game.Position)
						cm.bus.Publish(&game.SomethingInteresting{
							X: p.Center.X,
							Y: p.Center.Y,
						})
						break
					}
				},
			})
		}
	}

	cm.intents.Update()
	cm.performances.Update(elapsed)
	cm.hud.Update(elapsed)
	cm.cursors.Update(elapsed)
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
		participant := cm.mgr.Component(cm.turnToken, "Participant").(*Participant)
		var newSelected *geom.Key

		f := game.AdaptField(cm.field, participant.Size)
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

// syncParticipantObstacle updates the Participant's Obstacle to be synchronised
// with its position. It should be called when a Participant has completed a
// move.
func (cm *Manager) syncParticipantObstacle(evt *ParticipantMovementConcluded) {
	participant := cm.mgr.Component(evt.Entity, "Participant").(*Participant)
	obstacle := cm.mgr.Component(evt.Entity, "Obstacle").(*game.Obstacle)
	position := cm.mgr.Component(evt.Entity, "Position").(*game.Position)

	h := game.AdaptField(cm.field, participant.Size).At(int(position.Center.X), int(position.Center.Y))
	k := h.Key()
	obstacle.M = k.M
	obstacle.N = k.N
}

func (cm *Manager) handleMovementConcluded(t event.Typer) {
	// FIXME: Should Obstacle movement be handled by an "obstacle" system instead?
	cm.syncParticipantObstacle(t.(*ParticipantMovementConcluded))

	cm.setState(AwaitingInputState)
	cm.MousePosition(cm.x, cm.y)
}

func (cm *Manager) handleEndTurnRequested(event.Typer) {
	// Reset to maximum AP.
	participant := cm.mgr.Component(cm.turnToken, "Participant").(*Participant)
	participant.ActionPoints.Cur = participant.ActionPoints.Max

	// Remove turnToken
	cm.turnToken = 0
	cm.bus.Publish(&ParticipantTurnChanged{Entity: cm.turnToken})

	cm.setState(PreparingState)
}

func (cm *Manager) handleMoveModeRequested(event.Typer) {
	cm.setState(SelectingTargetState)
}

func (cm *Manager) handleCancelSkillRequested(event.Typer) {
	cm.setState(AwaitingInputState)
}

func (cm *Manager) handleAttemptingEscape(t event.Typer) {
	ev := t.(*AttemptingEscape)
	participant := cm.mgr.Component(ev.Entity, "Participant").(*Participant)
	participant.Status = Escaped

	cm.mgr.RemoveComponent(ev.Entity, &game.Sprite{})
	cm.mgr.RemoveComponent(ev.Entity, &game.Obstacle{})
	cm.mgr.RemoveComponent(ev.Entity, &game.Position{})
	cm.mgr.RemoveComponent(ev.Entity, &game.Facer{})

	cm.turnToken = 0
	cm.bus.Publish(&ParticipantTurnChanged{Entity: cm.turnToken})

	cm.setState(PreparingState)
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
