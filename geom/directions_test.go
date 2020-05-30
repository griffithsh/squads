package geom

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
)

func TestBestFacing(t *testing.T) {
	origin := Key{1, 4}
	tests := []struct {
		want DirectionType
		dest Key
	}{
		{N, Key{1, 2}},
		{N, Key{1, 0}},
		{NE, Key{1, 3}},
		{NE, Key{1, 1}},
		{NE, Key{2, 2}},
		{SE, Key{1, 5}},
		{SE, Key{2, 4}},
		{SE, Key{2, 6}},
		{SE, Key{1, 7}},
		{S, Key{1, 6}},
		{S, Key{1, 8}},
		{SW, Key{0, 5}},
		{SW, Key{0, 7}},
		{SW, Key{0, 6}},
		{SW, Key{0, 4}},
		{NW, Key{0, 3}},
		{NW, Key{0, 2}},
		{NW, Key{0, 1}},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf("%v to %v", origin, tc.dest), func(t *testing.T) {
			got := BestFacing(origin, tc.dest)

			if tc.want != got {
				t.Errorf("want %v got %v", tc.want, got)
			}

		})
	}
}
func BenchmarkBestFacing(b *testing.B) {
	origin := Key{0, 0}
	r := rand.New(rand.NewSource(0))
	for i := 0; i < b.N; i++ {
		BestFacing(origin, Key{r.Int() - math.MaxInt64, r.Int() - math.MaxInt64})
	}
}
