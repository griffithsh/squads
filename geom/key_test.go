package geom

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestNeighbors(t *testing.T) {
	for m := -2; m < 2; m++ {
		for n := -2; n < 2; n++ {
			k := Key{m, n}

			dirsByKey := k.Neighbors()
			keysByDir := k.Adjacent()

			t.Run(fmt.Sprintf("Key:%d,%d", k.M, k.N), func(t *testing.T) {
				for k, d := range dirsByKey {
					if keysByDir[d] != k {
						t.Errorf("%v: got %v", d, keysByDir[d])
					}
				}
				for d, k := range keysByDir {
					if dirsByKey[k] != d {
						t.Errorf("%v: got %v", k, dirsByKey[k])
					}
				}
			})
		}
	}
}

func TestHexesFrom(t *testing.T) {
	t.Run("curated-list", func(t *testing.T) {
		tests := []struct {
			A, B Key
			want int
		}{
			{Key{0, 0}, Key{0, 0}, 0},
			{Key{0, 0}, Key{3, 1}, 3},
			{Key{0, 0}, Key{-2, -2}, 3},
			{Key{-2, 2}, Key{3, -2}, 6},
			{Key{-2, -2}, Key{-2, 2}, 4},
			{Key{3, -2}, Key{3, 2}, 4},
		}

		for _, tc := range tests {
			t.Run(fmt.Sprintf("A(%d,%d)-to-B(%d,%d)", tc.A.M, tc.A.N, tc.B.M, tc.B.N), func(t *testing.T) {
				forward := tc.A.HexesFrom(tc.B)
				backward := tc.B.HexesFrom(tc.A)

				if forward != tc.want {
					t.Errorf("want %d forward, got %d", tc.want, forward)
				}
				if backward != tc.want {
					t.Errorf("want %d backward, got %d", tc.want, backward)
				}
			})
		}
	})
	t.Run("pseudo-random", func(t *testing.T) {
		// Start with a static seed for stable results.
		prng := rand.New(rand.NewSource(0))

		for i := 0; i < 128; i++ {
			key := Key{prng.Intn(1024) - 512, prng.Intn(1024) - 512}
			t.Run(fmt.Sprintf("{%d,%d}", key.M, key.N), func(t *testing.T) {
				t.Run("two-down", func(t *testing.T) {
					goal := key.ToS().ToS()
					want := 2

					forward := key.HexesFrom(goal)
					backward := goal.HexesFrom(key)

					if forward != want || backward != want {
						t.Errorf("want %d, got %d and %d", want, forward, backward)
					}
				})

				t.Run("three-uppish", func(t *testing.T) {
					goal := key.ToNW().ToNW().ToNW()
					want := 3

					forward := key.HexesFrom(goal)
					backward := goal.HexesFrom(key)

					if forward != want || backward != want {
						t.Errorf("want %d, got %d and %d", want, forward, backward)
					}
				})

				t.Run("around", func(t *testing.T) {
					goal := key.ToSE().ToNE().ToNW()
					want := 1

					forward := key.HexesFrom(goal)
					backward := goal.HexesFrom(key)

					if forward != want || backward != want {
						t.Errorf("want %d, got %d and %d", want, forward, backward)
					}
				})

				t.Run("snakey", func(t *testing.T) {
					goal := key.ToNW().ToNE().ToNW().ToNE().ToNW().ToNE()
					want := 3

					forward := key.HexesFrom(goal)
					backward := goal.HexesFrom(key)

					if forward != want || backward != want {
						t.Errorf("want %d, got %d and %d", want, forward, backward)
					}
				})

				t.Run("longer", func(t *testing.T) {
					goal := key.ToN().ToN().ToN().ToN().ToNE().ToNE()
					want := 6

					forward := key.HexesFrom(goal)
					backward := goal.HexesFrom(key)

					if forward != want || backward != want {
						t.Errorf("want %d, got %d and %d", want, forward, backward)
					}
				})

				t.Run("quite-long", func(t *testing.T) {
					goal := key.ToSE().ToSE().ToSE().ToSE().ToSE().ToSE().ToSE().ToSE().ToSE().ToSE()
					want := 10

					forward := key.HexesFrom(goal)
					backward := goal.HexesFrom(key)

					if forward != want || backward != want {
						t.Errorf("want %d, got %d and %d", want, forward, backward)
					}
				})
			})
		}
	})
}

func BenchmarkHexesFrom(b *testing.B) {
	dimension := 16
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		am, an, bm, bn :=
			rand.Intn(dimension*2)-dimension,
			rand.Intn(dimension*2)-dimension,
			rand.Intn(dimension*2)-dimension,
			rand.Intn(dimension*2)-dimension
		A, B := Key{am, an}, Key{bm, bn}
		b.StartTimer()

		distance := A.HexesFrom(B)

		b.StopTimer()
		if B.HexesFrom(A) != distance {
			b.Fatalf("A%v-B%v didn't match reverse", A, B)
		}
		b.StartTimer()
	}
}
