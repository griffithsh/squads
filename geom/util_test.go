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

		{Key{4, 16}, Key{0, 14}, false},
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

func TestXY(t *testing.T) {
	tests := []struct {
		m, n         int
		wantX, wantY float64
	}{
		// even M, even N
		{0, -4, 72, 48 - 96 - 96},
		{0, -2, 72, 48 - 96},
		{0, 0, 72, 48},
		{0, 2, 72, 48 + 96},
		{0, 4, 72, 48 + 96 + 96},

		// odd M, even N
		{1, -4, 266, 48 - 96 - 96},
		{1, -2, 266, 48 - 96},
		{1, 0, 266, 48},
		{1, 2, 266, 48 + 96},
		{1, 4, 266, 48 + 96 + 96},

		// even M, odd N
		{0, -3, 169, 96 - 96 - 96},
		{0, -1, 169, 96 - 96},
		{0, 1, 169, 96},
		{0, 3, 169, 96 + 96},
		{0, 5, 169, 96 + 96 + 96},

		// odd M, odd N
		{1, -3, 363, 96 - 96 - 96},
		{1, -1, 363, 96 - 96},
		{1, 1, 363, 96},
		{1, 3, 363, 96 + 96},
		{1, 5, 363, 96 + 96 + 96},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("for(%d,%d)", tc.m, tc.n), func(t *testing.T) {
			gotX, gotY := XY(tc.m, tc.n, 144, 96)

			if gotX != tc.wantX || gotY != tc.wantY {
				t.Errorf("want %f,%f, got %f,%f", tc.wantX, tc.wantY, gotX, gotY)
			}
		})
	}
}
