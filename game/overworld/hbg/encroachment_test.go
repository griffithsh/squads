package hbg

import (
	"fmt"
	"math"
	"math/rand"
	"testing"
	"time"

	"github.com/griffithsh/squads/geom"
)

func TestMarshalEncroachment(t *testing.T) {
	if err := testMarshaling(Encroachment{
		Description: "fsdfkiu",
		Over:        "sdfgdsh",
		Adjacent:    "jtyjfhjhgf",
		Edges: EdgeCollection{
			Texture: "dsfhyrjw",
			Options: map[geom.DirectionType]Options{
				geom.S: {
					Option{},
					Option{
						Chance: 36433,
						Frames: []Frame{
							{},
							{Duration: time.Hour * 99, Left: 27374, Top: 987},
						},
						ExtraHeight: 123,
						ExtraDepth:  492,
					},
				},
			},
		},
		Corners: CornerCollection{
			Texture: "sdfgdhjkk",
			Corners: map[geom.DirectionType]Corner{
				geom.NE: {
					Options: Options{
						Option{
							Chance: 53433,
							Frames: []Frame{
								{Duration: time.Nanosecond * 349493, Left: 34980958, Top: 349892},
							},
							ExtraHeight: 245231,
							ExtraDepth:  43857,
						},
						Option{},
					},
					Offset: Offset{X: 6223, Y: 8764},
					W:      234,
					H:      6593,
				},
			},
		},
	}); err != nil {
		t.Errorf("%v", err)
	}
}

func TestEncroachmentCollection(t *testing.T) {
	t.Run("Get", func(t *testing.T) {
		var ec EncroachmentsCollection

		if ec.Get("GRASS", "WATER") != nil {
			t.Errorf("should be nothing in an empty EncroachmentCollection")
		}
	})
	t.Run("Put", func(t *testing.T) {
		ec := EncroachmentsCollection{}
		e := Encroachment{
			Over:     "GRASS",
			Adjacent: "WATER",
		}
		ec.Put(e)

		if ec.Get("GRASS", "WATER") == nil {
			t.Errorf("expected to Get what I Put")
		}
	})
}

func TestOptions(t *testing.T) {
	t.Run("Roll", func(t *testing.T) {
		t.Run("Within Reason", func(t *testing.T) {
			options := Options{
				Option{Chance: 5},
				Option{Chance: 25},
				Option{Chance: 70},
			}

			results := map[int]int{}
			for i := 0; i < 1000; i++ {
				results[options.Roll(rand.Intn).Chance]++
			}
			for k, v := range results {
				results[k] = int(math.Round(float64(v) / 10))
			}
			for k, v := range results {
				if k+2 < v || k-2 > v {
					t.Errorf("want reasonable distribution, but got: %d=%d", k, v)
				}
			}

		})
		t.Run("Weighted Die", func(t *testing.T) {
			options := Options{
				Option{Chance: 2},
				Option{Chance: 1},
				Option{Chance: 3},
				Option{Chance: 1},
			}

			tests := []struct {
				roll int
				want int
			}{
				{roll: 0, want: 2},
				{roll: 1, want: 2},
				{roll: 2, want: 1},
				{roll: 3, want: 3},
				{roll: 4, want: 3},
				{roll: 5, want: 3},
				{roll: 6, want: 1},
			}

			for _, tc := range tests {
				t.Run(fmt.Sprintf("roll-%d", tc.roll), func(t *testing.T) {

					got := options.Roll(func(int) int { return tc.roll })
					if got.Chance != tc.want {
						t.Errorf("got %d, want %d", got.Chance, tc.want)
					}
				})
			}
		})
	})
}
