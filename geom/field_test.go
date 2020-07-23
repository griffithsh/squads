package geom

import (
	"fmt"
	"strconv"
	"testing"
)

func TestKeyToWorldtoKey(t *testing.T) {
	tests := []struct {
		name                         string
		bodyWidth, wingWidth, height int
	}{
		{"Legacy-Hex-Size", 10, 7, 16},
		{"Overworld-Hex-Size", 50, 47, 96},
		{"New-Hex-Size", 30, 21, 44},
	}
	for m := -2; m < 2; m++ {
		for n := -2; n < 2; n++ {
			k := Key{m, n}
			for _, tc := range tests {
				f := NewField(tc.bodyWidth, tc.wingWidth, tc.height)
				t.Run(fmt.Sprintf("%s-(%d,%d)", tc.name, m, n), func(t *testing.T) {
					x, y := f.Ktow(k)
					got := f.Wtok(x, y)
					if got != k {
						t.Errorf("got %f,%f and %v", x, y, got)
					}
				})

			}
		}
	}
}
func TestKeyToWorld(t *testing.T) {
	t.Run("Origin-Hex", func(t *testing.T) {
		tests := []struct {
			bodyWidth, wingWidth, height int
		}{
			{1, 1, 1},
			{10, 10, 10},
			{100, 100, 100},
			{100, 50, 2},
			{127, 49, 7},
		}

		for i, tc := range tests {
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				f := NewField(tc.bodyWidth, tc.wingWidth, tc.height)

				x, y := f.Ktow(Key{0, 0})

				if x != 0 || y != 0 {
					t.Errorf("Field(%d,%d,%d) origin hex was %f,%f", tc.bodyWidth, tc.wingWidth, tc.height, x, y)
				}
			})
		}
	})

	t.Run("Neighbors", func(t *testing.T) {
		tests := []struct {
			name                         string
			bodyWidth, wingWidth, height int
		}{
			{"Legacy-Hex-Size", 10, 7, 16},
			{"Overworld-Hex-Size", 50, 47, 96},
			{"New-Hex-Size", 30, 21, 44},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				f := NewField(tc.bodyWidth, tc.wingWidth, tc.height)

				tests := []struct {
					dir          DirectionType
					k            Key
					wantX, wantY float64
				}{
					{N, Key{0, -1}, 0, -float64(tc.height)},
					{NE, Key{1, -1}, float64(tc.wingWidth + tc.bodyWidth), -float64(tc.height / 2)},
					{SE, Key{1, 0}, float64(tc.wingWidth + tc.bodyWidth), float64(tc.height / 2)},
					{S, Key{0, 1}, 0, float64(tc.height)},
					{SW, Key{-1, 0}, -float64(tc.wingWidth + tc.bodyWidth), float64(tc.height / 2)},
					{NW, Key{-1, -1}, -float64(tc.wingWidth + tc.bodyWidth), -float64(tc.height / 2)},
				}
				for _, tc := range tests {
					t.Run(tc.dir.String(), func(t *testing.T) {
						x, y := f.Ktow(tc.k)
						if x != tc.wantX || y != tc.wantY {
							t.Errorf("want %f,%f, got %f,%f", tc.wantX, tc.wantY, x, y)
						}
					})
				}
			})
		}
	})
}

func TestWorldToKey(t *testing.T) {
	sizes := []struct {
		desc                         string
		bodyWidth, wingWidth, height int
	}{
		{"Legacy-Hex-Size", 10, 7, 16},
		{"Overworld-Hex-Size", 50, 47, 96},
		{"New-Hex-Size", 30, 21, 44},
	}

	for _, size := range sizes {
		t.Run(size.desc, func(t *testing.T) {
			halfHeight := float64(size.height) / 2
			totalWidth := float64(size.wingWidth + size.bodyWidth + size.wingWidth)
			halfWidth := totalWidth / 2
			f := NewField(size.bodyWidth, size.wingWidth, size.height)
			for n := -1; n <= 1; n++ {
				for m := -1; m <= 1; m++ {
					neighbors := map[string]Key{
						"origin": Key{m, n},
						"N":      Key{m, n - 1},
						"S":      Key{m, n + 1},
					}
					if m%2 == 0 {
						neighbors["NE"] = Key{m + 1, n - 1}
						neighbors["SE"] = Key{m + 1, n}
						neighbors["SW"] = Key{m - 1, n}
						neighbors["NW"] = Key{m - 1, n - 1}
					} else {
						// odd row
						neighbors["NE"] = Key{m + 1, n}
						neighbors["SE"] = Key{m + 1, n + 1}
						neighbors["SW"] = Key{m - 1, n + 1}
						neighbors["NW"] = Key{m - 1, n}
					}

					locations := []struct {
						desc             string
						xOffset, yOffset float64
						want             string
					}{
						{"Hex-Center", 0, 0, "origin"},
						{"NW", -halfWidth, -halfHeight, "NW"},
						{"SW", -halfWidth, halfHeight, "SW"},
						{"SE", halfWidth, halfHeight, "SE"},
						{"NE", halfWidth, -halfHeight, "NE"},

						// {"NW-Third", float64(size.wingWidth)/3 - halfWidth, float64(size.height)/6 - halfHeight, "NW"},
						// {"SW-Third", float64(size.wingWidth)/3 - halfWidth, -float64(size.height)/6 + halfHeight, "SW"},
						{"SE-Third", -float64(size.wingWidth)/3 + halfWidth, -float64(size.height)/6 + halfHeight, "SE"},
						{"NE-Third", -float64(size.wingWidth)/3 + halfWidth, float64(size.height)/6 - halfHeight, "NE"},

						{"W-Point", -halfWidth, 0, "origin"},
						{"NW-Point", 1 - float64(size.bodyWidth/2), 1 - halfHeight, "origin"},
						// {"NE-Point", float64(size.bodyWidth/2) - 1, 1 - halfHeight, "origin"},
						{"E-Point", halfWidth - 1, 0, "origin"},
						// {"SE-Point", float64(size.bodyWidth/2) - 1, halfHeight - 1, "origin"},
						{"SW-Point", 1 - float64(size.bodyWidth/2), halfHeight - 1, "origin"},

						{"NW-Edge", float64(size.wingWidth/2) + 1 - halfWidth, -halfHeight / 2, "origin"},
						{"SW-Edge", float64(size.wingWidth/2) + 1 - halfWidth, halfHeight / 2, "origin"},
						// {"NE-Edge", halfWidth - 3 - float64(size.wingWidth/2), 3 - halfHeight/2, "origin"},
						// {"SE-Edge", halfWidth - float64(size.wingWidth/2) - 1, halfHeight / 2, "origin"},
					}

					t.Run(fmt.Sprintf("(%d,%d)", m, n), func(t *testing.T) {
						for _, loc := range locations {
							t.Run(loc.desc, func(t *testing.T) {
								centerX, centerY := f.Ktow(Key{m, n})
								got := f.Wtok(centerX+loc.xOffset, centerY+loc.yOffset)
								if got != neighbors[loc.want] {
									t.Errorf("(%f+%f,%f+%f): want %v, got %v", centerX, loc.xOffset, centerY, loc.yOffset, neighbors[loc.want], got)
								}
							})
						}
					})
				}
			}
		})
	}
}
