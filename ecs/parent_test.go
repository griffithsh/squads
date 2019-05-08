package ecs

import "testing"

func TestParent(t *testing.T) {
	mgr := NewWorld()

	original := mgr.NewEntity()
	child := mgr.NewEntity()

	mgr.AddComponent(child, &Parent{Value: original})

	if !mgr.Exists(child) {
		t.Error("want child to exist, got child did not exist")
	}

	mgr.DestroyEntity(original)
	system := ParentSystem{}

	system.Update(mgr)

	if mgr.Exists(child) {
		t.Error("want child to be gone, got child still exists")
	}
}
