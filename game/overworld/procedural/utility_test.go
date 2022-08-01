package procedural

import (
	"fmt"
	"testing"
)

func TestLeftOfLine(t *testing.T) {
	ax, ay := 0.0, 0.0
	bx, by := 12.0, 12.0

	tests := []struct {
		qx, qy float64
		want   bool
	}{
		{5.0, 5.0, false},
		{5.0, 1.0, true},
		{1.0, 2.0, false},
		{-10.0, -12.0, true},
	}
	for _, tc := range tests {
		t.Run(fmt.Sprintf("%v", tc), func(t *testing.T) {
			got := leftOfLine(ax, ay, bx, by, tc.qx, tc.qy)

			if got != tc.want {
				t.Errorf("got: %t, want %t", got, tc.want)
			}
		})
	}
}
