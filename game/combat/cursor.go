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
	turnToken        ecs.Entity
	highlightedHexes skill.TargetingBrush
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

	// if we're navigating
	switch value.Context.Value() {
	case SelectingTargetState:
		cm.selectedKey = value.K
		ctx := value.Context.(*selectingTargetState)
		s := cm.archive.Skill(ctx.Skill)
		cm.highlightedHexes = s.TargetingBrush
		cm.showHighlightedHexes()
	case ConfirmingSelectedTargetState:
		cm.selectedKey = value.K
		ctx := value.Context.(*confirmingSelectedTargetState)
		s := cm.archive.Skill(ctx.Skill)
		cm.highlightedHexes = s.TargetingBrush
		cm.showHighlightedHexes()
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

const maxPathNavigationCursors int = 100

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
	switch cm.highlightedHexes {
	case skill.SingleHex:
		cm.paintSingleHex()
	case skill.Pathfinding:
		cm.paintNavigationHighlights()
	}
}

func (cm *CursorManager) paintSingleHex() {
	for i, e := range cm.mgr.Tagged(pathNavigationTag) {
		cm.mgr.RemoveTag(e, invalidatedCursorsTag)
		if i == 0 && cm.selectedKey != nil {
			x, y := cm.field.Ktow(*cm.selectedKey)
			cm.mgr.AddComponent(e, &game.Sprite{
				Texture: "cursors.png",

				X: 0, Y: hexagonHeight * 1,
				W: hexagonTileWidth, H: hexagonHeight,
			})
			cm.mgr.AddComponent(e, &game.Position{
				Center: game.Center{
					X: x,
					Y: y,
				},
				Layer: cursorLayer,
			})
			continue
		}
		cm.mgr.RemoveComponent(e, &game.Position{})
		cm.mgr.RemoveComponent(e, &game.Sprite{})
	}
}

func (cm *CursorManager) paintNavigationHighlights() {
	participant := cm.mgr.Component(cm.turnToken, "Participant").(*Participant)
	pos := cm.mgr.Component(cm.turnToken, "Position").(*game.Position)

	var start, goal geom.Key

	sh := cm.field.At(pos.Center.X, pos.Center.Y)
	start = sh.Key()

	cost := CostsFuncFactory(cm.field, cm.mgr, cm.turnToken)
	edges := EdgeFuncFactory(cm.field)
	guess := HeuristicFactory(cm.field)

	type comps struct {
		p game.Position
		s game.Sprite
	}
	c := []comps{}
	used := map[*geom.Hex]struct{}{}

	var steps []graph.Step
	if cm.selectedKey == nil {
		goto repaintLabel
	}
	goal = *cm.selectedKey
	steps = graph.NewSearcher(cost, edges, guess).Search(start, goal)
	if steps == nil {
		h := cm.field.Get(goal)
		if h == nil {
			goto repaintLabel
		}
		x, y := cm.field.Ktow(goal)
		c = append(c, comps{
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
		goto repaintLabel
	}

	for _, step := range steps {
		h := cm.field.Get(step.V.(geom.Key))
		if _, ok := used[h]; ok {
			continue
		}
		used[h] = struct{}{}
		x, y := h.Center()
		c = append(c, comps{
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
		if step.Cost > float64(participant.ActionPoints.Cur) {
			c[len(c)-1].s.X = 0
			c[len(c)-1].s.Y = hexagonHeight * 2
		}
	}

repaintLabel:
	if len(c) > len(cm.mgr.Tagged(pathNavigationTag)) {
		fmt.Println("not enough Entity slots for this path!")
	}

	for i, e := range cm.mgr.Tagged(pathNavigationTag) {
		cm.mgr.RemoveTag(e, invalidatedCursorsTag)
		if i >= len(c) {
			cm.mgr.RemoveComponent(e, &game.Position{})
			cm.mgr.RemoveComponent(e, &game.Sprite{})
			continue
		}

		cm.mgr.AddComponent(e, &c[i].s)
		cm.mgr.AddComponent(e, &c[i].p)
	}
}
