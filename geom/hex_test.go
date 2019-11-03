package geom

import (
	"fmt"
	"testing"
)

func TestHexXY(t *testing.T) {
	tests := []struct {
		m, n int
		x, y float64
	}{
		{0, 0, 12, 8},
		{0, 1, 29, 16},
		{0, 5, 29, 48},
		{1, 3, 63, 32},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("%d,%d", tc.m, tc.n), func(t *testing.T) {
			h := Hex{
				M: tc.m,
				N: tc.n,
			}

			if h.X() != tc.x || h.Y() != tc.y {
				t.Errorf("want %f,%f got %f,%f", tc.x, tc.y, h.X(), h.Y())
			}
		})
	}
}
