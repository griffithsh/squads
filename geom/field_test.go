package geom

import (
	"fmt"
	"testing"
)

func TestAt(t *testing.T) {
	field := NewField()
	field.Load(MByN(1, 1))

	tests := []struct {
		x, y    int
		wantHex bool
	}{
		{x: 0, y: 0, wantHex: false},
		{x: 3, y: 3, wantHex: false},
		{x: 2, y: 12, wantHex: false},
		{x: 3, y: 4, wantHex: true},
		{x: 7, y: 0, wantHex: true},
		{x: 16, y: 0, wantHex: true},
		{x: 23, y: 8, wantHex: true},
		{x: 19, y: 12, wantHex: true},
		{x: 17, y: 6, wantHex: true},
		{x: 23, y: 15, wantHex: false},
		{x: 17, y: 0, wantHex: false},

		{x: 100, y: 100, wantHex: false},
		{x: -100, y: 100, wantHex: false},
		{x: 24, y: 4, wantHex: false},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("At-%d-%d", tc.x, tc.y), func(t *testing.T) {
			h := field.At(tc.x, tc.y)

			if tc.wantHex && h == nil {
				t.Error("want hex, got nil\n")

			} else if !tc.wantHex && h != nil {
				t.Error("want nil, got hex\n")

			}
		})
	}
}

func TestAtMN(t *testing.T) {
	field := NewField()
	field.Load(MByN(2, 6))

	tests := []struct {
		x, y  int
		wantM int
		wantN int
	}{
		{11, 8, 0, 0},
		{45, 7, 1, 0},
		{28, 16, 0, 1},
		{66, 18, 1, 1},
		{0, 24, 0, 2},
		{47, 31, 1, 2},
		{30, 24, 0, 3},
		{74, 31, 1, 3},
		{6, 33, 0, 4},
		{57, 39, 1, 4},
		{27, 55, 0, 5},
		{70, 48, 1, 5},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("AtMN-%d-%d)", tc.x, tc.y), func(t *testing.T) {
			h := field.At(tc.x, tc.y)
			if h == nil {
				t.Fatal("no hex\n")
			}

			if h.M != tc.wantM || h.N != tc.wantN {
				t.Errorf("\nwant %d,%d\ngot  %d,%d\n", tc.wantM, tc.wantN, h.M, h.N)
			}
		})
	}
}

func TestRoughMN(t *testing.T) {
	tests := []struct {
		x, y         int
		wantM, wantN int
	}{
		// top-left hex, simple
		{0, 0, 0, 0},
		{16, 0, 0, 0},
		{16, 15, 0, 0},
		{0, 15, 0, 0},

		// second hex, just below the first
		{0, 16, 0, 2},
		{16, 31, 0, 2},

		{6, 52, 0, 6},   // long way down
		{0, -1, 0, -2},  // just above
		{0, -15, 0, -2}, // nearly two above
		{0, -16, 0, -4}, // two above

		// second column
		{17, 6, 0, -1},
		{17, 24, 0, 3},

		// some random ones
		{51, 8, 1, 1},
		{67, 55, 1, 5},

		{-6, 23, -1, 1},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%d,%d", tc.x, tc.y), func(t *testing.T) {
			m, n := roughMN(tc.x, tc.y)

			if m != tc.wantM || n != tc.wantN {
				t.Errorf("\nwant %d,%d\ngot  %d,%d", tc.wantM, tc.wantN, m, n)
			}
		})
	}
}
