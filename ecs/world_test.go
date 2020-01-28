package ecs

import (
	"reflect"
	"testing"
)

func TestWorld(t *testing.T) {
	t.Run("HasTag", func(t *testing.T) {
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
	})
	t.Run("RemoveTag", func(t *testing.T) {
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
	})
	t.Run("ListComponents", func(t *testing.T) {
		mgr := NewWorld()
		e := mgr.NewEntity()
		mgr.AddComponent(e, &Children{})
		mgr.Tag(e, "Anything")

		got := mgr.ListComponents(e)

		want := []string{"Children", "Tags"}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("want %v, got %v", want, got)
		}
	})

	t.Run("Tagged", func(t *testing.T) {
		mgr := NewWorld()
		for i := 0; i < 100; i++ {
			e := mgr.NewEntity()
			mgr.Tag(e, "prefix-tag")
			if i%7 == 0 {
				mgr.Tag(e, "mod7iszero")
			}
		}

		prefixed := len(mgr.Tagged("prefix-tag"))
		if prefixed != 100 {
			t.Errorf("want 100 tagged with prefix, got %d", prefixed)
		}

		mod7ed := len(mgr.Tagged("mod7iszero"))
		if mod7ed != 15 {
			t.Errorf("want 15 tagged with mod7, got %d", mod7ed)
		}
	})
}
