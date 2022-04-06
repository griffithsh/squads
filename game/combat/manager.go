package combat

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/griffithsh/squads/baddy"
	"github.com/griffithsh/squads/data"
	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/game/item"
	"github.com/griffithsh/squads/geom"
	"github.com/griffithsh/squads/skill"
	"github.com/griffithsh/squads/ui"
)

// SkillArchive is what is required by combat of any skill definition provider.
type SkillArchive interface {
	Skill(skill.ID) *skill.Description
	SkillsByProfession(string) []*skill.Description
	SkillsByWeaponClass(item.Class) []*skill.Description
	Appearance(profession string, sex game.CharacterSex, hair string, skin string) *game.Appearance
	Profession(profession string) *game.ProfessionDetails
}

// Manager is a game-mode. It processes turns-based Combat until one or the other
// team is knocked out.
type Manager struct {
	// Manager should own systems that are only relevant to Combat. A Turns coordinator, a preparation timer
	mgr     *ecs.World
	bus     *event.Bus
	archive *data.Archive
	field   *geom.Field
	nav     *Navigator
	camera  *game.Camera
	hud     *HUD
	cursors *CursorManager
	se      *skillExecutor
	ds      *damageSystem

	turnToken            ecs.Entity // Whose turn is it? References an existing Entity.
	selectingInteractive ecs.Entity // catches clicks on the field.

	// Manager has both a state and a paused flag, so that state transition
	// logic can be isolated from pausing.
	paused bool
	state  StateContext

	incrementAccumulator float64

	x, y             int       // where the mouse last was in screen coordinates
	screenW, screenH int       // most recent dimensions of the window
	selectedHex      *geom.Key // most recent hex selected

	intents      *IntentSystem
	performances *PerformanceSystem

	celebrations time.Duration

	squads []ecs.Entity

	// vanishers is the list of entities that will vanish when a participant
	// goes behind them.
	vanishers map[ecs.Entity][]geom.Key
}

// NewManager creates a new combat Manager.
func NewManager(mgr *ecs.World, camera *game.Camera, bus *event.Bus, archive *data.Archive) *Manager {
	f := geom.NewField(hexagonBodyWidth, hexagonWingWidth, hexagonHeight)

	cm := Manager{
		mgr:                  mgr,
		bus:                  bus,
		archive:              archive,
		field:                f,
		nav:                  NewNavigator(bus),
		camera:               camera,
		state:                Uninitialised,
		hud:                  NewHUD(mgr, bus, camera.GetW(), camera.GetH(), archive),
		cursors:              NewCursorManager(mgr, bus, archive, f),
		se:                   newSkillExecutor(mgr, bus, f, archive),
		ds:                   newDamageSystem(mgr, bus),
		selectingInteractive: mgr.NewEntity(),
		intents:              NewIntentSystem(mgr, bus, f),
		performances:         NewPerformanceSystem(mgr, bus, archive),

		paused: false,
	}

	cm.bus.Subscribe(ParticipantMovementConcluded{}.Type(), cm.handleMovementConcluded)
	cm.bus.Subscribe(EndTurnRequested{}.Type(), cm.handleEndTurnRequested)
	cm.bus.Subscribe(CancelSkillRequested{}.Type(), cm.handleCancelSkillRequested)
	cm.bus.Subscribe(AttemptingEscape{}.Type(), cm.handleAttemptingEscape)
	cm.bus.Subscribe(game.WindowSizeChanged{}.Type(), cm.handleWindowSizeChanged)
	cm.bus.Subscribe(SkillRequested{}.Type(), cm.handleSkillRequested)
	cm.bus.Subscribe(SkillUseConcluded{}.Type(), cm.handleSkillUseConcluded)
	cm.bus.Subscribe(CharacterEnteredCombat{}.Type(), cm.handleCharacterEnteredCombat)
	cm.bus.Subscribe(ParticipantDefiled{}.Type(), cm.handleParticipantDefiled)
	cm.bus.Subscribe(ParticipantMoving{}.Type(), cm.handleParticipantMoving)

	return &cm
}

func (cm *Manager) handleTargetSelected(x, y float64) {
	ctx := cm.state.(*selectingTargetState)

	h := cm.field.At(x, y)

	s := cm.archive.Skill(ctx.Skill)
	switch s.Targeting {
	case skill.TargetAnywhere:
		// No filtering is applicable.
	case skill.TargetAdjacent:
		if h == nil {
			return
		}
		obstacle := cm.mgr.Component(cm.turnToken, "Obstacle").(*game.Obstacle)
		origin := geom.Key{M: obstacle.M, N: obstacle.N}
		adjacent := false
		for key := range origin.Neighbors() {
			if key == h.Key() {
				adjacent = true
				break
			}
		}
		if !adjacent {
			return
		}
	}

	// TODO: extra validation needs to happen here - you can only raise skeleton
	// if there is a downed character here, right?

	var selected *geom.Key
	if h != nil {
		// Go to confirming state if a Hex was selected, and save the Key of the
		// selected Hex.
		cm.setState(&confirmingSelectedTargetState{
			Skill:  ctx.Skill,
			Target: h.Key(),
		})
		pselected := h.Key()
		selected = &pselected
	}
	// Publish an event whether a Hex was selected or not, passing the Key if
	// applicable.
	cm.bus.Publish(&DifferentHexSelected{
		K:       selected,
		Context: ctx,
	})
}

func (cm *Manager) handleTargetConfirmed(x, y float64) {
	ctx := cm.state.(*confirmingSelectedTargetState)

	obstacle := cm.mgr.Component(cm.turnToken, "Obstacle").(*game.Obstacle)
	origin := cm.field.Get(geom.Key{M: obstacle.M, N: obstacle.N})
	selected := cm.field.At(x, y)
	if selected == nil {
		// We cannot confirm the selection of something outside the hexes of the
		// field.
		return
	}

	// Go back to selectingTargetState.
	if selected.Key() != ctx.Target {
		cm.setState(&selectingTargetState{
			Skill: ctx.Skill,
		})
		return
	}

	// Special handling for movement.
	if ctx.Skill == skill.BasicMovement {
		cm.mgr.AddComponent(cm.turnToken, &MoveIntent{X: x, Y: y})
		cm.setState(ExecutingState)
		return
	}

	s := cm.archive.Skill(ctx.Skill)
	switch s.Targeting {
	case skill.TargetAnywhere:
		// Because we can target anywhere, there are no reasons to return out of
		// this function for this case.
	case skill.TargetAdjacent:
		adjacent := false

		for key := range origin.Key().Neighbors() {
			if key == selected.Key() {
				adjacent = true
				break
			}
		}
		if !adjacent {
			return
		}
	}

	cm.setState(ExecutingState)

	if s.IsAttack() {
		cm.mgr.AddComponent(cm.turnToken, &game.Facer{Face: geom.S})
	}
	cm.bus.Publish(&UsingSkill{
		User:     cm.turnToken,
		Skill:    s.ID,
		Selected: selected,
	})
}

// setState is the canonical way to change the CombatState.
func (cm *Manager) setState(stateContext StateContext) {
	state := stateContext.Value()
	if state == cm.state.Value() {
		return
	}
	ev := StateTransition{
		Old: cm.state,
		New: stateContext,
	}
	cm.state = stateContext

	// When entering Selecting Target State, we need to add an Interactive to
	// cover all areas of the field, so that we can convert those clicks to
	// MoveIntents.
	if ev.Old.Value() == SelectingTargetState || ev.Old.Value() == ConfirmingSelectedTargetState {
		cm.mgr.RemoveComponent(cm.selectingInteractive, &ui.Interactive{})
	}
	switch state.Value() {
	case SelectingTargetState:
		// Using the max float value as the size and a position of 0,0 should
		// work in all cases, and it's a lot faster than figuring out the actual
		// dimensions of the field. The goal here is to catch *any* clicks in
		// the world, remember?
		cm.mgr.AddComponent(cm.selectingInteractive, &game.Position{})
		cm.mgr.AddComponent(cm.selectingInteractive, &ui.Interactive{
			W: math.MaxFloat64, H: math.MaxFloat64,
			Trigger: cm.handleTargetSelected,
		})
	case ConfirmingSelectedTargetState:
		// Using the max float value as the size and a position of 0,0 should
		// work in all cases, and it's a lot faster than figuring out the actual
		// dimensions of the field. The goal here is to catch *any* clicks in
		// the world, remember?
		cm.mgr.AddComponent(cm.selectingInteractive, &game.Position{})
		cm.mgr.AddComponent(cm.selectingInteractive, &ui.Interactive{
			W: math.MaxFloat64, H: math.MaxFloat64,
			Trigger: cm.handleTargetConfirmed,
		})
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

	startX, startY := f.Ktow(geom.Key{M: m, N: n})
	distances := make([]s, len(f.Hexes()))

	for i, h := range f.Hexes() {
		x, y := f.Ktow(h.Key())
		distances[i] = s{math.Pow(math.Abs(x-startX), 2) + math.Pow(math.Abs(y-startY), 2), h}
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

// isBlocked determines if a Character can be placed at m,n.
func isBlocked(field *geom.Field, k geom.Key, mgr *ecs.World) bool {
	// blockages is a set of Keys that are taken by other things
	blockages := map[geom.Key]struct{}{}
	for _, e := range mgr.Get([]string{"Obstacle"}) {
		o := mgr.Component(e, "Obstacle").(*game.Obstacle)

		h := field.Get(geom.Key{M: o.M, N: o.N})
		if h == nil {
			panic(fmt.Sprintf("there is no hex where Obstacle(%d,%s) is present (%d,%d)", e, o.ObstacleType, o.M, o.N))
		}

		// FIXME: We're making the assumption again here that all obstacles
		// are total obstacles. Even conceptually things like shallow water
		// or bushes that should only impede movement slightly.
		blockages[h.Key()] = struct{}{}
	}

	hex := field.Get(k)
	if hex == nil {
		return true
	}

	if _, blocked := blockages[k]; blocked {
		return true
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

func (cm *Manager) getStart(nearbys []*geom.Hex) *geom.Hex {
	for _, h := range nearbys {
		if isBlocked(cm.field, h.Key(), cm.mgr) {
			continue
		}

		return h
	}
	return nil
}

// createParticipation adds a new Entity to participate in combat based on a Character.
func (cm *Manager) createParticipation(charEntity ecs.Entity, team *game.Team, atHex *geom.Hex) {
	e := cm.mgr.NewEntity()
	cm.mgr.Tag(e, "combat")

	// Add Participant Component.
	equipment, _ := cm.mgr.Component(charEntity, "Equipment").(*item.Equipment)
	char := cm.mgr.Component(charEntity, "Character").(*game.Character)
	prof := cm.archive.Profession(char.Profession)

	app := cm.archive.Appearance(char.Profession, char.Sex, char.Hair, char.Skin)

	// FIXME: there is a configuration layer here, because the player can select
	// a subset of their skills to use.
	configuredSkills := []skill.ID{}
	if equipment != nil && equipment.Weapon != nil {
		configuredSkills = equipment.Weapon.Skills
	}

	participant := &Participant{
		Name:               char.Name,
		Level:              char.Level,
		Hair:               char.Hair,
		Skin:               char.Skin,
		SmallPortraitBG:    game.PortraitBGSmall[char.PortraitBG],
		BigPortraitBG:      game.PortraitBGBig[char.PortraitBG],
		SmallPortraitFrame: game.PortraitFrameSmall[char.PortraitFrame],
		BigPortraitFrame:   game.PortraitFrameBig[char.PortraitFrame],
		SmallIcon:          app.SmallIcon(),
		BigIcon:            app.BigIcon(),
		Profession:         char.Profession,
		Sex:                char.Sex,
		PreparationThreshold: CurMax{
			Max: char.InherantPreparation + prof.Preparation + equipment.WeaponPreparation(),
		},
		ActionPoints: CurMax{
			Max: char.InherantActionPoints + prof.ActionPoints + equipment.WeaponActionPoints(),
		},
		BaseHealth:    char.BaseHealth,
		CurrentHealth: char.CurrentHealth,
		// Strength:     int(char.StrengthPerLevel * float64(char.Level)),
		// Agility:    0,
		// Intelligence: 0,
		// Vitality:     0,
		Disambiguator: char.Disambiguator,
		Masteries:     char.Masteries,

		EquippedWeaponClass: equipment.WeaponClass(),
		ItemStats:           equipment.SumModifiers(),
		// FIXME: Skills should come from a subset of the available skills
		// configured by the player. Available skills come from the equipped
		// items and the profession of the Character.
		Skills: append(configuredSkills, ([]skill.ID{
			"debug-lightning",
			"debug-revive",
			"raise-skeleton",
		})...),
	}

	participant.Character = charEntity
	participant.ActionPoints.Cur = participant.ActionPoints.Max
	participant.Status = Alive
	cm.mgr.AddComponent(e, participant)

	// Add Team.
	cm.mgr.AddComponent(e, team)

	// Add Position.
	cm.field.Ktow(atHex.Key())
	x, y := atHex.Center()
	cm.mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: x,
			Y: y,
		},
		Layer: participantLayer,
	})

	// Add Obstacle.
	o := game.Obstacle{
		M:            atHex.Key().M,
		N:            atHex.Key().N,
		ObstacleType: game.CharacterObstacle,
	}
	cm.mgr.AddComponent(e, &o)

	// Add Facer Component.
	cm.mgr.AddComponent(e, &game.Facer{Face: geom.S})

	leash := game.Leash{
		Owner:       e,
		LayerOffset: -1,
	}
	e = cm.mgr.NewEntity()
	cm.mgr.Tag(e, "combat")

	spr := cm.archive.GetPedestal(team.PedestalAppearance)

	cm.mgr.AddComponent(e, spr)
	cm.mgr.AddComponent(e, &leash)
}

// Begin should be called at the start of an engagement to set up components
// necessary for the combat.
// TODO: Begin should be provided with information about the squads and the
// terrain to fight on.
func (cm *Manager) Begin(participatingSquads []ecs.Entity) {
	cm.setState(FadingIn)
	e := cm.mgr.NewEntity()
	cm.mgr.Tag(e, "combat")
	cm.mgr.AddComponent(e, &game.DiagonalMatrixWipe{
		W: cm.screenW, H: cm.screenH,
		Obscuring: false, // ergo revealing
		OnComplete: func() {
			cm.setState(PreparingState)
		},
		OnInitialised: func() {
			// TODO: we should provide the terrain we're entering combat to this
			// method so that a map with an appropriate tileset can be selected.
			combatMap := cm.archive.GetCombatMap()

			keys := make([]geom.Key, len(combatMap.Hexes))
			for i, hex := range combatMap.Hexes {
				keys[i] = hex.Position
			}
			cm.field.Load(keys)
			cm.vanishers = make(map[ecs.Entity][]geom.Key)
			// TODO: center on the midpoint between selected start locations instead of
			// the middle of the map. Or maybe the start location of the player's team?
			x, y := cm.field.Center()
			cm.bus.Publish(&game.SomethingInteresting{
				X: x,
				Y: y,
			})
			for _, hex := range combatMap.Hexes {
				h := cm.field.Get(geom.Key{M: hex.Position.M, N: hex.Position.N})
				if hex.Obstacle != game.NonObstacle {
					// add obstacle
					e := cm.mgr.NewEntity()
					cm.mgr.Tag(e, "combat")
					cm.mgr.AddComponent(e, &game.Obstacle{
						M:            h.Key().M,
						N:            h.Key().N,
						ObstacleType: hex.Obstacle,
					})
				}
				for _, vis := range hex.Visuals {
					if len(vis.Frames) == 0 {
						continue
					}
					e := cm.mgr.NewEntity()
					cm.mgr.Tag(e, "combat")

					for _, o := range vis.Obscures {
						if val, ok := cm.vanishers[e]; ok {
							cm.vanishers[e] = append(val, o)
						} else {
							cm.vanishers[e] = []geom.Key{o}
						}
					}

					// add position
					x, y := h.Center()
					cm.mgr.AddComponent(e, &game.Position{
						Center: game.Center{
							X: x,
							Y: y,
						},
						Layer: vis.Layer,
					})

					sprite := func(frame game.CombatMapRecipeHexFrame) game.Sprite {
						return game.Sprite{
							Texture: frame.Texture,
							X:       frame.X,
							Y:       frame.Y,
							W:       frame.W,
							H:       frame.H,
							OffsetX: vis.XOffset + frame.W/2 - hexagonTileWidth/2,
							OffsetY: vis.YOffset - frame.H/2 + hexagonHeight/2,
						}
					}
					// If there is only one frame, add a Sprite Component.
					if len(vis.Frames) == 1 {
						spr := sprite(vis.Frames[0])
						cm.mgr.AddComponent(e, &spr)
					} else { // not zero or one, and definitely not negative
						// If there is more than one frame, add a FrameAnimationComponent.
						fa := game.FrameAnimation{}
						for _, frame := range vis.Frames {
							fa.Frames = append(fa.Frames, sprite(frame))
							fa.Timings = append(fa.Timings, frame.Duration)
						}
						fa.Randomise()
						cm.mgr.AddComponent(e, &fa)
					}

				}
			}

			// TODO:
			// There is some entity which stores info about a "level", and produces
			// artifacts that can be used by the combat Manager. It should produce the
			// shape of the level, and the terrain of each hex (grass, water, blocked by
			// tree etc). It should also produce starting positions for teams...

			sp := newStartProvider(combatMap.Starts)

			entities := []ecs.Entity{}
			for _, e := range participatingSquads {
				cm.squads = append(cm.squads, e)
				squad := cm.mgr.Component(e, "Squad").(*game.Squad)
				entities = append(entities, squad.Members...)
			}

			// Create a Participating Entity for every Character we have.
			for _, e := range entities {
				team := cm.mgr.Component(e, "Team").(*game.Team)
				near := sp.getNearby(team, cm.field)
				h := cm.getStart(near)

				cm.createParticipation(e, team, h)
			}

			// This is a bit of a hacky way call the code to hide or show trees
			// etc on initialisation.
			cm.handleParticipantMoving(nil)

			// Announce that the Combat has begun.
			cm.bus.Publish(&game.CombatBegan{})
		},
	})
	cm.hud.Enable()
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
	cm.hud.Disable()
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
		// TODO: recalc obscures if they're implemented with a Hidden Component.
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

		// FIXME: It's non-deterministic whose turn it is when multiple
		// Participants finish preparing at the same time.
		if len(prepared) > 0 {

			sort.Slice(prepared, func(i, j int) bool {
				p1 := cm.mgr.Component(prepared[i], "Participant").(*Participant)
				p2 := cm.mgr.Component(prepared[j], "Participant").(*Participant)

				return p1.Disambiguator < p2.Disambiguator
			})
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
		cm.se.Update(elapsed)
	case Celebration:
		// Celebrate for a time ...
		cm.celebrations += elapsed
		if cm.celebrations > time.Second*2 {
			cm.celebrations = 0
			cm.mgr.AddComponent(cm.mgr.NewEntity(), &game.DiagonalMatrixWipe{
				W: cm.screenW, H: cm.screenH,
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

	var skill skill.ID
	var selecting bool

	switch cm.state.Value() {
	case SelectingTargetState:
		ctx := cm.state.(*selectingTargetState)
		skill = ctx.Skill
		selecting = true
	case ConfirmingSelectedTargetState:
		ctx := cm.state.(*confirmingSelectedTargetState)
		skill = ctx.Skill
		// ctx.Target
		selecting = false
	default:
		return
	}

	h := cm.field.At(wx, wy)
	var newSelected *geom.Key
	if h != nil {
		k := h.Key()
		newSelected = &geom.Key{
			M: k.M,
			N: k.N,
		}
	}
	if geom.Equal(newSelected, cm.selectedHex) {
		return
	}

	if !selecting {
		cm.setState(&selectingTargetState{
			Skill: skill,
		})
	}
	cm.handleTargetSelected(wx, wy)
	cm.selectedHex = newSelected
}

// syncParticipantObstacle updates the Participant's Obstacle to be synchronised
// with its position. It should be called when a Participant has completed a
// move.
func (cm *Manager) syncParticipantObstacle(evt *ParticipantMovementConcluded) {
	obstacle := cm.mgr.Component(evt.Entity, "Obstacle").(*game.Obstacle)
	position := cm.mgr.Component(evt.Entity, "Position").(*game.Position)

	h := cm.field.At(position.Center.X, position.Center.Y)
	k := h.Key()
	obstacle.M = k.M
	obstacle.N = k.N
}

func (cm *Manager) handleMovementConcluded(t event.Typer) {
	ev := t.(*ParticipantMovementConcluded)
	// FIXME: Should Obstacle movement be handled by an "obstacle" system instead?
	cm.syncParticipantObstacle(ev)

	cm.mgr.RemoveComponent(ev.Entity, &game.FrameAnimation{})

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

func (cm *Manager) handleWindowSizeChanged(e event.Typer) {
	wsc := e.(*game.WindowSizeChanged)
	cm.screenW, cm.screenH = wsc.NewW, wsc.NewH
}

func (cm *Manager) handleSkillRequested(e event.Typer) {
	evt := e.(*SkillRequested)

	// Can we afford this skill?
	participant := cm.mgr.Component(cm.turnToken, "Participant").(*Participant)
	s := cm.archive.Skill(evt.Code)
	for ty, amount := range s.Costs {
		switch ty {
		case skill.CostsActionPoints:
			if participant.ActionPoints.Cur < amount {
				// TODO: Let the UI know the player can't afford it.
				return
			}
		default:
			panic(fmt.Sprintf("combat.Manager.handleSkillRequested: Costs type not implemented %T", ty))
		}
	}

	cm.setState(&selectingTargetState{
		Skill: evt.Code,
	})
}
func (cm *Manager) handleSkillUseConcluded(e event.Typer) {
	// evt := e.(*SkillUseConcluded)
	cm.setState(AwaitingInputState)
}

func (cm *Manager) handleCharacterEnteredCombat(et event.Typer) {
	evt := et.(*CharacterEnteredCombat)

	// FIXME: evt.Profession should be used to retrieve a recipe
	recipe := baddy.Recipes[baddy.Skellington]

	// FIXME: evt.Level should be used to Construct the recipe.
	char := recipe.Construct(nil)
	e := cm.mgr.NewEntity()
	cm.mgr.AddComponent(e, char)
	cm.createParticipation(e, evt.Team, cm.field.Get(evt.At))
}

func (cm *Manager) handleParticipantDefiled(et event.Typer) {
	pde := et.(*ParticipantDefiled)

	cm.mgr.RemoveComponent(pde.Entity, &game.Obstacle{})
}

func (cm *Manager) handleParticipantMoving(event.Typer) {
	// get a list of the Keys that participants currently reside in.
	participants := []geom.Key{}
	for _, e := range cm.mgr.Get([]string{"Participant", "Position"}) {
		pos := cm.mgr.Component(e, "Position").(*game.Position)
		k := cm.field.Wtok(pos.Center.X, pos.Center.Y)
		participants = append(participants, k)
	}

	// Go through every "vanisher" (a tile that obscures the view of other
	// tiles), and check if a participant is on one of the tiles it obscures. If
	// it is, then hide this vanisher, otherwise, show it.
	for vanisher, keys := range cm.vanishers {
		var shouldVanish bool
		for _, key := range keys {
			for _, participant := range participants {
				if participant == key {
					shouldVanish = true

					// break out of two layers of for loop.
					goto nextVanisher
				}
			}
		}
	nextVanisher:

		var invisible bool
		if alpha := cm.mgr.Component(vanisher, "Alpha"); alpha != nil {
			if alpha.(*game.Alpha).Value == 0 {
				invisible = true
			}
		}

		var disappearing bool
		if c := cm.mgr.Component(vanisher, "FadeOut"); c != nil {
			disappearing = true
		}
		if shouldVanish && !invisible && !disappearing {
			cm.mgr.AddComponent(vanisher, &game.FadeOut{Duration: time.Millisecond * 200})
		}

		if !shouldVanish && (invisible || disappearing) {
			cm.mgr.AddComponent(vanisher, &game.FadeIn{Duration: time.Millisecond * 100})
		}
	}
}
