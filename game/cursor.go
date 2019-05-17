package game

import (
	"github.com/griffithsh/squads/ecs"
)

type Cursor struct{}

// Type of this Component.
func (*Cursor) Type() string {
	return "Cursor"
}

// Types
/*
small cursor
small cursor blocked
medium cursor
large cursor
small path highlight
skill targeting highlight
partially blocked medium cursor
partially blocked large cursor
*/

type CursorSystem struct {
	mgr *ecs.World
}

func NewCursorSystem(mgr *ecs.World) *CursorSystem {
	return &CursorSystem{
		mgr: mgr,
	}
}

func (s *CursorSystem) Update() {
	// TODO: work out what Update means in this context.

	// Some rambling thoughts...
	// For every cursor, create a child (or children?) that represents the cursor, then remove the Cursor Component?

	// Is this just synchronisation?

	// For every cursor, create/or edit the Sprite and Position Components.

	// What if there was a Children Component that stored the Children of this Cursor

	// For every cursor, destroy all children, create new Children based on
	// context in the current Cursor Component, all the children need to have a
	// parent Component:
	//   - not particularly efficient, but should work.
	//   - could add second Component type to separate Tagging Components from Constructing Components

	// I wonder what other cursor-like things I might want later.
	//   - A Path of lit-up hexes that previews navigation maybe?
	//   - lines, shapes, etc that show which hexes a skill will target

	// The core of this is that I want to be able to create and destroy Cursor
	// Components, but have a system that handles the tedium of synchronising
	// that to a renderable representation.
}

// Clear all Cursors
func (s *CursorSystem) Clear() {
	for _, e := range s.mgr.Get([]string{"Cursor"}) {
		s.mgr.DestroyEntity(e)
	}
}

// Add a new Cursor via this System
func (s *CursorSystem) Add(x, y float64, sz ActorSize) {
	switch sz {
	case SMALL:
		e := s.mgr.NewEntity()
		for _, c := range []ecs.Component{
			&Cursor{},
			&Sprite{
				Texture: "texture.png",
				X:       0,
				Y:       0,
				W:       24,
				H:       16,
			},
			&Position{
				Center: Center{
					X: x,
					Y: y,
				},
				Layer: 2,
			},
		} {
			s.mgr.AddComponent(e, c)
		}
	case MEDIUM:
		e := s.mgr.NewEntity()
		for _, c := range []ecs.Component{
			&Cursor{},
			&Sprite{
				Texture: "texture.png",
				X:       0,
				Y:       32,
				W:       58,
				H:       32,
			},
			&Position{
				Center: Center{
					X: x,
					Y: y,
				},
				Layer: 2,
			},
		} {
			s.mgr.AddComponent(e, c)
		}
	case LARGE:
		e := s.mgr.NewEntity()
		for _, c := range []ecs.Component{
			&Cursor{},
			&Sprite{
				Texture: "texture.png",
				X:       0,
				Y:       64,
				W:       58,
				H:       48,
			},
			&Position{
				Center: Center{
					X: x,
					Y: y,
				},
				Layer: 2,
			},
		} {
			s.mgr.AddComponent(e, c)
		}
	}
}
