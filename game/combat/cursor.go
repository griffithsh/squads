package combat

import (
	"time"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/geom"
)

var (
	cursorsTag    = "CURSORS_TAG"
	liveActorsTag = cursorsTag + ".LIVE_ACTORS"

	invalidatedCursorsTag = cursorsTag + ".INVALIDATED"
)

// CursorManager controls the visibility of cursors in a game combat. Cursors
// are visual highlights applied to Hexes. Examples might be permanent boundary
// markers of every Actor, or temporary red/green blockouts when selecting a
// place to target a skill.
type CursorManager struct {
	mgr   *ecs.World
	bus   *event.Bus
	field *geom.Field
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

	return &cm
}

func (cm *CursorManager) handleCombatBegan(ev event.Typer) {
	cm.showLiveActors()
}

// Update the CursorManager, repainting invalidated Cursors.
func (cm *CursorManager) Update(elapsed time.Duration) {
	var e ecs.Entity

	e = cm.mgr.AnyTagged(liveActorsTag)
	if e != 0 && cm.mgr.HasTag(e, invalidatedCursorsTag) {
		cm.repaintLiveActors()
	}
}

const maxLiveActors int = 25

func (cm *CursorManager) showLiveActors() {
	for _, e := range cm.mgr.Tagged(liveActorsTag) {
		cm.mgr.DestroyEntity(e)
	}

	for i := 0; i < maxLiveActors; i++ {
		e := cm.mgr.NewEntity()

		cm.mgr.Tag(e, liveActorsTag)
		cm.mgr.Tag(e, invalidatedCursorsTag)
	}
}

func (cm *CursorManager) hideLiveActors() {
	for _, e := range cm.mgr.Tagged(liveActorsTag) {
		cm.mgr.DestroyEntity(e)
	}
}

func (cm *CursorManager) repaintLiveActors() {
	entities := cm.mgr.Get([]string{"Actor"})
	for i, slot := range cm.mgr.Tagged(liveActorsTag) {
		if i < len(entities) {
			spr := game.Sprite{
				Texture: "cursors.png",
			}
			actor := cm.mgr.Component(entities[i], "Actor").(*game.Actor)
			switch actor.Size {
			case game.SMALL:
				spr.X = 0
				spr.Y = 0
				spr.W = 24
				spr.H = 16
			case game.MEDIUM:
				spr.X = 0
				spr.Y = 32
				spr.W = 58
				spr.H = 32
			case game.LARGE:
				spr.X = 0
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
