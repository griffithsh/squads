package geom

import (
	"testing"
)

func TestNavigate(t *testing.T) {

	t.Run("20x60", func(t *testing.T) {
		field := NewField()
		field.Load(MByN(20, 60))

		start := Key{0, 0}
		goal := Key{0, 58}
		exists := func(k Key) bool {
			return field.Get(k.M, k.N) != nil
		}
		costs := func(Key) float64 {
			return 1.0
		}
		steps, err := Navigate(start, goal, exists, costs)
		if err != nil {
			t.Fatal(err)
		}
		if len(steps) != 58/2+1 {
			t.Errorf("want %d steps, got %d steps", 58/2+1, len(steps))
		}
		prev := -1.0
		for i, step := range steps {
			if step.Cost <= prev {
				t.Errorf("want steps to increase in cost, got step %d cost (%f) is less than or equal to previous (%f)", i, step.Cost, prev)
			}
		}
	})

	t.Run("negatives", func(t *testing.T) {
		exists := func(k Key) bool {
			result := k.M >= -1 && k.M <= 0 && k.N >= -3 && k.N <= 3
			return result
		}
		costs := func(k Key) float64 {
			return 1.0
		}
		start := Key{-1, -2}
		goal := Key{0, 3}
		steps, err := Navigate(start, goal, exists, costs)
		if err != nil {
			t.Fatal(err)
		}
		if len(steps) != 6 {
			t.Errorf("want 6 steps got %d", len(steps))
		}
	})
}
