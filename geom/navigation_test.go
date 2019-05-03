package geom

import (
	"fmt"
	"testing"
)

func TestNavigate(t *testing.T) {
	field, _ := NewField(20, 60)

	start := field.Get(0, 0)
	goal := field.Get(0, 58)
	steps, err := Navigate(start, goal, []ContextualObstacle{})
	if err != nil {
		t.Fatal(err)
	}
	if len(steps) != 58/2+1 {
		t.Errorf("want %d steps, got %d steps", 58/2+1, len(steps))
	}
	for _, step := range steps {
		fmt.Printf("To: %f,%f\n", step.X(), step.Y())
	}
}

func TestNavigate4(t *testing.T) {
	field, _ := NewField(20, 60)

	start := field.Get4(1, 1)
	goal := field.Get4(1, 57)
	steps, err := Navigate4(start, goal, []ContextualObstacle{})
	if err != nil {
		t.Fatal(err)
	}
	if len(steps) != 57/2+1 {
		t.Errorf("want %d steps, got %d steps", 57/2+1, len(steps))
	}
	fmt.Printf("Start: %f,%f\n", start.X(), start.Y())
	fmt.Printf("Goal: %f,%f\n", goal.X(), goal.Y())
	for _, step := range steps {
		fmt.Printf("Waypoint: %f,%f\n", step.X(), step.Y())
	}
}
