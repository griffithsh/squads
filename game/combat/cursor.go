package combat

import (
	"fmt"
	"time"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/geom"
	"github.com/griffithsh/squads/graph"
	"github.com/griffithsh/squads/skill"
	"github.com/griffithsh/squads/targeting"
)

var (
	cursorsTag          = "CURSORS_TAG"
	liveParticipantsTag = cursorsTag + ".LIVE_PARTICIPANTS"
	pathNavigationTag   = cursorsTag + ".PATH_NAVIGATION"

	invalidatedCursorsTag = cursorsTag + ".INVALIDATED"
)

// CursorManager controls the visibility of cursors in a game combat. Cursors
// are visual highlights applied to Hexes. Examples might be permanent boundary
// markers of every Character, or temporary red/green blockouts when selecting a
// place to target a skill. All Cursor Entities are tagged with "combat".
type CursorManager struct {
	mgr         *ecs.World
	bus         *event.Bus
	archive     SkillArchive
	field       *geom.Field
	selectedKey *geom.Key

	// Whose turn is it?
	turnToken ecs.Entity

	targeting *targeting.Rule

	lastState State
}

// NewCursorManager creates a new CursorManager.
func NewCursorManager(mgr *ecs.World, bus *event.Bus, archive SkillArchive, f *geom.Field) *CursorManager {
	cm := CursorManager{
		mgr:     mgr,
		bus:     bus,
		archive: archive,
		field:   f,
	}

	// Subscribes
	bus.Subscribe(game.CombatBegan{}.Type(), cm.handleCombatBegan)
	bus.Subscribe(DifferentHexSelected{}.Type(), cm.handleDifferentHexSelected)
	bus.Subscribe(StateTransition{}.Type(), cm.handleCombatStateTransition)
	t := ParticipantTurnChanged{}.Type()
	bus.Subscribe(t, cm.handleParticipantTurnChanged)

	return &cm
}

func (cm *CursorManager) handleCombatBegan(ev event.Typer) {
	cm.showLiveParticipants()
}

func (cm *CursorManager) handleDifferentHexSelected(ev event.Typer) {
	value := ev.(*DifferentHexSelected)

	cm.lastState = value.Context.Value()
	if ctx, ok := value.Context.(*selectingTargetState); ok {
		if ctx.Skill == skill.BasicMovement {
			cm.lastState = SelectingPathState
		}
	}
	switch cm.lastState {
	case SelectingPathState:
		cm.selectedKey = value.K
		cm.showHighlightedHexes()

	case SelectingTargetState:
		ctx := value.Context.(*selectingTargetState)
		s := cm.archive.Skill(ctx.Skill)
		cm.selectedKey = value.K
		cm.targeting = &s.Targeting

		cm.showHighlightedHexes()

	case ConfirmingSelectedTargetState:
		ctx := value.Context.(*confirmingSelectedTargetState)
		s := cm.archive.Skill(ctx.Skill)
		cm.selectedKey = value.K
		cm.targeting = &s.Targeting

		cm.showHighlightedHexes()

	default:
		cm.selectedKey = nil
		cm.targeting = nil
		cm.hideHighlightedHexes()
	}
}

func (cm *CursorManager) handleCombatStateTransition(ev event.Typer) {
	cm.hideHighlightedHexes()
}

func (cm *CursorManager) handleParticipantTurnChanged(ev event.Typer) {
	atc := ev.(*ParticipantTurnChanged)
	cm.turnToken = atc.Entity
}

// Update the CursorManager, repainting invalidated Cursors.
func (cm *CursorManager) Update(elapsed time.Duration) {
	var e ecs.Entity

	e = cm.mgr.AnyTagged(liveParticipantsTag)
	if e != 0 && cm.mgr.HasTag(e, invalidatedCursorsTag) {
		cm.repaintLiveParticipants()
	}
	e = cm.mgr.AnyTagged(pathNavigationTag)
	if e != 0 && cm.mgr.HasTag(e, invalidatedCursorsTag) {
		cm.repaintHighlightedHexes()
	}
}

const maxLiveParticipants int = 25

func (cm *CursorManager) showLiveParticipants() {
	for _, e := range cm.mgr.Tagged(liveParticipantsTag) {
		cm.mgr.DestroyEntity(e)
	}

	for i := 0; i < maxLiveParticipants; i++ {
		e := cm.mgr.NewEntity()
		cm.mgr.Tag(e, "combat")

		cm.mgr.Tag(e, liveParticipantsTag)
		cm.mgr.Tag(e, invalidatedCursorsTag)
	}
}

func (cm *CursorManager) repaintLiveParticipants() {
	entities := cm.mgr.Get([]string{"Participant"})

	// liveEntity determines if an there is an alive entity at index i of entities.
	liveEntity := func(i int) bool {
		if i < len(entities) {
			participant := cm.mgr.Component(entities[i], "Participant").(*Participant)
			if participant.Status == Alive {
				return true
			}
		}
		return false
	}

	for i, slot := range cm.mgr.Tagged(liveParticipantsTag) {
		if liveEntity(i) {
			spr := game.Sprite{
				Texture: "cursors.png",

				X: 0, Y: 0,
				W: hexagonTileWidth, H: hexagonHeight,
			}
			cm.mgr.AddComponent(slot, &spr)
			cm.mgr.AddComponent(slot, &game.Leash{
				Owner:       entities[i],
				LayerOffset: -5,
			})
		} else {
			// hide cursor
			cm.mgr.RemoveComponent(slot, &game.Sprite{})
			cm.mgr.RemoveComponent(slot, &game.Position{})
			cm.mgr.RemoveComponent(slot, &game.Leash{})
		}
		cm.mgr.RemoveTag(slot, invalidatedCursorsTag)
	}
}

const maxPathNavigationCursors int = 100 // FIXME: could be dynamic based on number of hexes in field?

func (cm *CursorManager) showHighlightedHexes() {
	for _, e := range cm.mgr.Tagged(pathNavigationTag) {
		cm.mgr.DestroyEntity(e)
	}

	for i := 0; i < maxPathNavigationCursors; i++ {
		e := cm.mgr.NewEntity()
		cm.mgr.Tag(e, "combat")

		cm.mgr.Tag(e, pathNavigationTag)
		cm.mgr.Tag(e, invalidatedCursorsTag)
	}
}

func (cm *CursorManager) hideHighlightedHexes() {
	for _, e := range cm.mgr.Tagged(pathNavigationTag) {
		cm.mgr.DestroyEntity(e)
	}
}

func (cm *CursorManager) repaintHighlightedHexes() {
	type cursorSprite struct {
		// s is the sprite to use
		s game.Sprite
		// p is where to paint it
		p game.Position
	}

	paints := []cursorSprite{}

	if cm.selectedKey == nil {
		// When nothing is selected, then we have nothing to repaint, skip
		// straight to the Component removal.
	} else if cm.lastState == SelectingPathState {
		participant := cm.mgr.Component(cm.turnToken, "Participant").(*Participant)
		pos := cm.mgr.Component(cm.turnToken, "Position").(*game.Position)

		sh := cm.field.At(pos.Center.X, pos.Center.Y)
		start, goal := sh.Key(), *cm.selectedKey

		cost := CostsFuncFactory(cm.field, cm.mgr, cm.turnToken)
		edges := EdgeFuncFactory(cm.field)
		guess := HeuristicFactory(cm.field)
		steps := graph.NewSearcher(cost, edges, guess).Search(start, goal)
		goalHex := cm.field.Get(goal)
		if goalHex == nil {
			// TODO: wait, what?
		} else if steps == nil {
			x, y := cm.field.Ktow(goal)
			paints = append(paints, cursorSprite{
				s: game.Sprite{
					Texture: "cursors.png",

					X: hexagonTileWidth * 1, Y: hexagonHeight * 1,
					W: hexagonTileWidth, H: hexagonHeight,
				},
				p: game.Position{
					Center: game.Center{
						X: x,
						Y: y,
					},
					Layer: cursorLayer,
				},
			})
		} else {
			for _, step := range steps {
				h := cm.field.Get(step.V.(geom.Key))
				x, y := h.Center()
				paints = append(paints, cursorSprite{
					s: game.Sprite{
						Texture: "cursors.png",

						X: 0, Y: hexagonHeight * 1,
						W: hexagonTileWidth, H: hexagonHeight,
					},
					p: game.Position{
						Center: game.Center{
							X: x,
							Y: y,
						},
						Layer: cursorLayer,
					},
				})

				// If we don't have enough AP to get all the way, then swap the cursor
				// sprite to the red one.
				if step.Cost > float64(participant.ActionPoints.Cur) {
					paints[len(paints)-1].s.X = 0
					paints[len(paints)-1].s.Y = hexagonHeight * 2
				}
			}
		}

	} else if cm.lastState == SelectingTargetState || cm.lastState == ConfirmingSelectedTargetState {
		obstacle := cm.mgr.Component(cm.turnToken, "Obstacle").(*game.Obstacle)
		ok, highlighted := cm.targeting.Execute(*cm.selectedKey, geom.Key{M: obstacle.M, N: obstacle.N})
		if !ok {
			// Add a single red cursor on selected hex.
			x, y := cm.field.Ktow(*cm.selectedKey)
			paints = append(paints, cursorSprite{
				s: game.Sprite{
					Texture: "cursors.png",

					X: 0, Y: hexagonHeight * 2,
					W: hexagonTileWidth, H: hexagonHeight,
				},
				p: game.Position{
					Center: game.Center{
						X: x,
						Y: y,
					},
					Layer: cursorLayer,
				},
			})

		} else {
			for _, k := range highlighted {
				x, y := cm.field.Ktow(k)
				paints = append(paints, cursorSprite{
					s: game.Sprite{
						Texture: "cursors.png",

						X: 0, Y: hexagonHeight,
						W: hexagonTileWidth, H: hexagonHeight,
					},
					p: game.Position{
						Center: game.Center{
							X: x,
							Y: y,
						},
						Layer: cursorLayer,
					},
				})

			}
		}

	}

	if len(paints) > len(cm.mgr.Tagged(pathNavigationTag)) {
		fmt.Println("not enough Entity slots for this path!")
	}

	for i, e := range cm.mgr.Tagged(pathNavigationTag) {
		cm.mgr.RemoveTag(e, invalidatedCursorsTag)
		if i < len(paints) {
			cm.mgr.AddComponent(e, &paints[i].s)
			cm.mgr.AddComponent(e, &paints[i].p)
			continue
		}
		cm.mgr.RemoveComponent(e, &game.Position{})
		cm.mgr.RemoveComponent(e, &game.Sprite{})
	}
}
