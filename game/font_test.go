package game

import (
	"testing"

	"github.com/griffithsh/squads/ecs"
)

func TestFont(t *testing.T) {
	const text = "Exactamento!"
	const text2 = "Secondi!"
	mgr := ecs.NewWorld()

	e := mgr.NewEntity()
	mgr.AddComponent(e, &Font{
		Text: text,
	})
	mgr.AddComponent(e, &Position{})

	// Should have no sprites.
	sprites := mgr.Get([]string{"Sprite"})
	if len(sprites) != 0 {
		t.Errorf("setup: want no sprites, got %d", len(sprites))
	}

	fs := NewFontSystem(mgr)
	fs.Update()

	// Should have same number of sprites as test text.
	sprites = mgr.Get([]string{"Sprite"})
	if len(sprites) != len(text) {
		t.Errorf("assert: want %d sprites, got %d", len(text), len(sprites))
	}

	// Test a mutation.
	f := mgr.Component(e, "Font").(*Font)
	f.Text = text2
	fs.Update()
	sprites = mgr.Get([]string{"Sprite"})
	if len(sprites) != len(text2) {
		t.Errorf("assert: after mutating want %d sprites, got %d", len(text2), len(sprites))
	}

	// Clean up the entity
	mgr.DestroyEntity(e)
	ps := ecs.NewParentSystem(mgr)
	ps.Update()

	// Should have no sprites.
	sprites = mgr.Get([]string{"Sprite"})
	if len(sprites) != 0 {
		t.Errorf("concluded: want no sprites, got %d", len(sprites))
	}
}
