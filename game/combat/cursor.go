package combat

import (
	"fmt"
	"time"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/geom"
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
	field       *geom.Field
	selectedKey *geom.Key

	// Whose turn is it?
	turnToken        ecs.Entity
	highlightedHexes game.TargetingBrush
}

// NewCursorManager creates a new CursorManager.
func NewCursorManager(mgr *ecs.World, bus *event.Bus, f *geom.Field) *CursorManager {
	cm := CursorManager{
		mgr:   mgr,
		bus:   bus,
		field: f,
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
		cm.highlightedHexes = game.BrushForSkill[ctx.Skill]
		cm.showHighlightedHexes()
	case ConfirmingSelectedTargetState:
		cm.selectedKey = value.K
		ctx := value.Context.(*confirmingSelectedTargetState)
		cm.highlightedHexes = game.BrushForSkill[ctx.Skill]
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

func (cm *CursorManager) hideLiveParticipants() {
	for _, e := range cm.mgr.Tagged(liveParticipantsTag) {
		cm.mgr.DestroyEntity(e)
	}
}

func (cm *CursorManager) repaintLiveParticipants() {
	entities := cm.mgr.Get([]string{"Participant"})
	for i, slot := range cm.mgr.Tagged(liveParticipantsTag) {
		if i < len(entities) {
			spr := game.Sprite{
				Texture: "cursors.png",
			}
			participant := cm.mgr.Component(entities[i], "Participant").(*Participant)
			switch participant.Size {
			case game.SMALL:
				spr.X = 0
				spr.Y = 0
				spr.W = 24
				spr.H = 16
			case game.MEDIUM:
				spr.X = 58
				spr.Y = 32
				spr.W = 58
				spr.H = 32
			case game.LARGE:
				spr.X = 58
				spr.Y = 64
				spr.W = 58
				spr.H = 48
			}
			cm.mgr.AddComponent(slot, &spr)
			cm.mgr.AddComponent(slot, &game.Leash{
				Owner:       entities[i],
				LayerOffset: -1,
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
	case game.SingleHex:
		cm.paintSingleHex()
	case game.Pathfinding:
		cm.paintNavigationHighlights()
	}
}

func (cm *CursorManager) paintSingleHex() {
	for i, e := range cm.mgr.Tagged(pathNavigationTag) {
		cm.mgr.RemoveTag(e, invalidatedCursorsTag)
		if i == 0 && cm.selectedKey != nil {
			h := cm.field.Get(cm.selectedKey.M, cm.selectedKey.N)
			cm.mgr.AddComponent(e, &game.Sprite{
				Texture: "cursors.png",

				X: 0, Y: 16,
				W: 24, H: 16,
			})
			cm.mgr.AddComponent(e, &game.Position{
				Center: game.Center{
					X: h.X(),
					Y: h.Y(),
				},
				Layer: 10,
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

	f := game.AdaptField(cm.field, participant.Size)

	var start, goal geom.Key

	sh := f.At(int(pos.Center.X), int(pos.Center.Y))
	start = sh.Key()

	exists := ExistsFuncFactory(cm.field, participant.Size)
	costs := CostsFuncFactory(cm.field, cm.mgr, cm.turnToken)

	type comps struct {
		p game.Position
		s game.Sprite
	}
	c := []comps{}
	used := map[*geom.Hex]struct{}{}

	var steps []geom.NavigateStep
	var err error
	if cm.selectedKey == nil {
		goto repaintLabel
	}
	goal = *cm.selectedKey
	steps, err = geom.Navigate(start, goal, exists, costs)
	if err != nil {
		h := f.Get(goal.M, goal.N)
		if h == nil {
			goto repaintLabel
		}
		for _, h := range h.Hexes() {
			c = append(c, comps{
				s: game.Sprite{
					Texture: "cursors.png",

					X: 24, Y: 16,
					W: 24, H: 16,
				},
				p: game.Position{
					Center: game.Center{
						X: h.X(),
						Y: h.Y(),
					},
					Layer: 10,
				},
			})
		}
		goto repaintLabel
	}

	for _, step := range steps {
		h := f.Get(step.M, step.N)
		for _, h := range h.Hexes() {
			if _, ok := used[h]; ok {
				continue
			}
			used[h] = struct{}{}
			c = append(c, comps{
				s: game.Sprite{
					Texture: "cursors.png",

					X: 0, Y: 16,
					W: 24, H: 16,
				},
				p: game.Position{
					Center: game.Center{
						X: h.X(),
						Y: h.Y(),
					},
					Layer: 10,
				},
			})
			if int(step.Cost) > participant.ActionPoints.Cur {
				c[len(c)-1].s.X = 48
			}
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
