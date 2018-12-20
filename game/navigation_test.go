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
	goal := board.Get(19, 59)
	steps, err := Navigate(start, goal)
	if err != nil {
		t.Fatal(err)
	}
	for _, step := range steps {
		fmt.Printf("To: %d,%d\n", step.M, step.N)
	}
}
