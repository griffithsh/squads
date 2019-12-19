package geom

import (
	"fmt"
	"testing"
)

func TestAdjacent(t *testing.T) {
	tests := []struct {
		a, b Key
		want bool
	}{
		{Key{69, 420}, Key{69, 420}, false},
		{Key{0, 2}, Key{0, 0}, true},
		{Key{0, 1}, Key{0, 0}, true},

		{Key{0, 3}, Key{0, 2}, true},
		{Key{0, 3}, Key{0, 1}, true},
		{Key{0, 3}, Key{1, 2}, true},
		{Key{0, 3}, Key{1, 4}, true},
		{Key{0, 3}, Key{0, 5}, true},
		{Key{0, 3}, Key{0, 4}, true},

		{Key{1, 4}, Key{0, 3}, true},
		{Key{1, 4}, Key{1, 2}, true},
		{Key{1, 4}, Key{1, 3}, true},
		{Key{1, 4}, Key{1, 5}, true},
		{Key{1, 4}, Key{1, 6}, true},
		{Key{1, 4}, Key{0, 5}, true},

		{Key{2, 1}, Key{1, 1}, false},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%d,%d_and_%d,%d", tc.a.M, tc.a.N, tc.b.M, tc.b.N), func(t *testing.T) {
			if Adjacent(tc.a, tc.b) != Adjacent(tc.b, tc.a) {
				t.Error("results depends on order of arguments")
			}
			got := Adjacent(tc.a, tc.b)
			if got != tc.want {
				words := map[bool]string{
					true:  "adjacent",
					false: "not adjacent",
				}
				t.Errorf("got %s, want %s", words[got], words[tc.want])
			}
		})
	}
}
