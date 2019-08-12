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
	cursorsTag        = "CURSORS_TAG"
	liveActorsTag     = cursorsTag + ".LIVE_ACTORS"
	pathNavigationTag = cursorsTag + ".PATH_NAVIGATION"

	invalidatedCursorsTag = cursorsTag + ".INVALIDATED"
)

// CursorManager controls the visibility of cursors in a game combat. Cursors
// are visual highlights applied to Hexes. Examples might be permanent boundary
// markers of every Actor, or temporary red/green blockouts when selecting a
// place to target a skill.
type CursorManager struct {
	mgr         *ecs.World
	bus         *event.Bus
	field       *geom.Field
	selectedKey *geom.Key
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

	return &cm
}

func (cm *CursorManager) handleCombatBegan(ev event.Typer) {
	cm.showLiveActors()
}

func (cm *CursorManager) handleDifferentHexSelected(ev event.Typer) {
	value := ev.(*DifferentHexSelected)

	cm.selectedKey = value.K

	cm.showPathNavigationCursor()
}

func (cm *CursorManager) handleCombatStateTransition(ev event.Typer) {
	cm.hidePathNavigationCursor()
}

// Update the CursorManager, repainting invalidated Cursors.
func (cm *CursorManager) Update(elapsed time.Duration) {
	var e ecs.Entity

	e = cm.mgr.AnyTagged(liveActorsTag)
	if e != 0 && cm.mgr.HasTag(e, invalidatedCursorsTag) {
		cm.repaintLiveActors()
	}
	e = cm.mgr.AnyTagged(pathNavigationTag)
	if e != 0 && cm.mgr.HasTag(e, invalidatedCursorsTag) {
		cm.repaintPathNavigationCursor()
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

const maxPathNavigationCursors int = 100

func (cm *CursorManager) showPathNavigationCursor() {
	for _, e := range cm.mgr.Tagged(pathNavigationTag) {
		cm.mgr.DestroyEntity(e)
	}

	for i := 0; i < maxPathNavigationCursors; i++ {
		e := cm.mgr.NewEntity()

		cm.mgr.Tag(e, pathNavigationTag)
		cm.mgr.Tag(e, invalidatedCursorsTag)
	}
}

func (cm *CursorManager) hidePathNavigationCursor() {
	for _, e := range cm.mgr.Tagged(pathNavigationTag) {
		cm.mgr.DestroyEntity(e)
	}
}

func (cm *CursorManager) repaintPathNavigationCursor() {
	entities := cm.mgr.Get([]string{"TurnToken", "Actor", "Position", "CombatStats"})
	if len(entities) != 1 {
		// this really should not happen, right?
		return
	}
	actor := cm.mgr.Component(entities[0], "Actor").(*game.Actor)
	stats := cm.mgr.Component(entities[0], "CombatStats").(*game.CombatStats)
	pos := cm.mgr.Component(entities[0], "Position").(*game.Position)

	var start, goal geom.Key
	switch actor.Size {
	case game.MEDIUM:
		sh := cm.field.At4(int(pos.Center.X), int(pos.Center.Y))
		start = geom.Key{M: sh.M, N: sh.N}
	case game.LARGE:
		sh := cm.field.At7(int(pos.Center.X), int(pos.Center.Y))
		start = geom.Key{M: sh.M, N: sh.N}
	default:
		sh := cm.field.At(int(pos.Center.X), int(pos.Center.Y))
		start = geom.Key{M: sh.M, N: sh.N}
	}
	exists := game.ExistsFuncFactory(cm.field, actor.Size)
	costs := game.CostsFuncFactory(cm.field, cm.mgr, entities[0])

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
		switch actor.Size {
		case game.MEDIUM:
			h4 := cm.field.Get4(goal.M, goal.N)
			if h4 == nil {
				goto repaintLabel
			}
			for _, h := range h4.Hexes() {
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
		case game.LARGE:
			h7 := cm.field.Get7(goal.M, goal.N)
			if h7 == nil {
				goto repaintLabel
			}
			for _, h := range h7.Hexes() {
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
		default:
			h := cm.field.Get(goal.M, goal.N)
			if h == nil {
				goto repaintLabel
			}
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
		switch actor.Size {
		case game.MEDIUM:
			h4 := cm.field.Get4(step.M, step.N)
			for _, h := range h4.Hexes() {
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
				if int(step.Cost) > stats.ActionPoints {
					c[len(c)-1].s.X = 48
				}
			}
		case game.LARGE:
			h7 := cm.field.Get7(step.M, step.N)
			for _, h := range h7.Hexes() {
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
				if int(step.Cost) > stats.ActionPoints {
					c[len(c)-1].s.X = 48
				}
			}
		default:
			h := cm.field.Get(step.M, step.N)
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
			if int(step.Cost) > stats.ActionPoints {
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
