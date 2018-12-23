package game

import (
	"fmt"
	"testing"

	"github.com/griffithsh/squads/ecs"
)

func TestNavigate(t *testing.T) {
	mgr := ecs.NewWorld()
	board, _ := NewBoard(mgr, 20, 60)

	start := board.Get(0, 0)
	goal := board.Get(0, 58)
	steps, err := Navigate(start, goal, []ContextualObstacle{})
	if err != nil {
		t.Fatal(err)
	}
	if len(steps) != 58/2+1 {
		t.Errorf("want %d steps, got %d steps", 58/2+1, len(steps))
	}
	for _, step := range steps {
		fmt.Printf("To: %d,%d\n", step.M, step.N)
	}
}
