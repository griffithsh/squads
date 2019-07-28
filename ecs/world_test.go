package ecs

import "testing"

func TestHasTag(t *testing.T) {
	t.Run("nonexistant-entity", func(t *testing.T) {
		mgr := NewWorld()

		if mgr.HasTag(123456, "RedHerring") {
			t.Error("want not tagged for nonexistant entity")
		}
	})

	t.Run("tagged-not-tagged", func(t *testing.T) {
		mgr := NewWorld()
		e := mgr.NewEntity()

		mgr.Tag(e, "Present")
		mgr.Tag(e, "RedHerring")

		if !mgr.HasTag(e, "Present") {
			t.Error("want tagged with Present, got not tagged")
		}
		if mgr.HasTag(e, "NotPresent") {
			t.Error("want not tagged with NotPresent, got tagged")
		}

	})
}

func TestRemoveTag(t *testing.T) {
	mgr := NewWorld()
	e := mgr.NewEntity()
	mgr.Tag(e, "B")
	mgr.Tag(e, "RedHerring")

	e = mgr.NewEntity()
	mgr.Tag(e, "A")
	mgr.Tag(e, "B")

	if len(mgr.Tagged("A")) != 1 {
		t.Errorf("want one Entity tagged with A, got %d", len(mgr.Tagged("A")))
	}

	mgr.RemoveTag(e, "A")
	if len(mgr.Tagged("A")) != 0 {
		t.Errorf("want no Entities tagged with A, got %d", len(mgr.Tagged("A")))
	}

	if len(mgr.Tagged("B")) != 2 {
		t.Errorf("want one Entity tagged with B, got %d", len(mgr.Tagged("B")))
	}
}
