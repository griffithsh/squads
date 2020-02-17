package overworld

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/griffithsh/squads/geom"
)

func TestParseRecipe(t *testing.T) {
	reference := []byte(`label: Reference Recipe

terrain: // m,n tileid
1 4 1, 1 2 1,
1 3 0,

paths: // link hexes
1 4 1 2, 1 4 1 3, 1 3 1 2

interesting: // pick 2 of (0,0 1,1 and 2,2) to be part of the map.
1 (-1 -1, 0 0),
2 (2 2, 3 3,4 4)

	`)

	got, err := ParseRecipe(bytes.NewReader(reference))
	if err != nil {
		t.Fatalf("parse reference recipe: %v", err)
	}
	if got.Label != "Reference Recipe" {
		t.Errorf("want label %s, got label %s", "Reference Recipe", got.Label)
	}

	// Test Terrain is correct
	if len(got.Terrain) != 3 {
		t.Fatalf("want 3 terrain tiles")
	}
	if got.Terrain[geom.Key{M: 1, N: 4}] != TileID(1) {
		t.Errorf("want 1,4 to equal 1")
	}
	if got.Terrain[geom.Key{M: 1, N: 2}] != TileID(1) {
		t.Errorf("want 1,2 to equal 1")
	}
	if got.Terrain[geom.Key{M: 1, N: 3}] != TileID(0) {
		t.Errorf("want 1,3 to equal 0")
	}

	// Test Paths are correct
	if len(got.Paths) != 3 {
		t.Fatalf("want 3 paths")
	}
	want := KeyPair{
		First:  geom.Key{M: 1, N: 4},
		Second: geom.Key{M: 1, N: 2},
	}
	if got.Paths[0] != want {
		t.Errorf("want 1 4 1 2, got %v", got.Paths[0])
	}

	want = KeyPair{
		First:  geom.Key{M: 1, N: 4},
		Second: geom.Key{M: 1, N: 3},
	}
	if got.Paths[1] != want {
		t.Errorf("want 1 4 1 3, got %v", got.Paths[1])
	}

	want = KeyPair{
		First:  geom.Key{M: 1, N: 3},
		Second: geom.Key{M: 1, N: 2},
	}
	if got.Paths[2] != want {
		t.Errorf("want 1 4 1 3, got %v", got.Paths[2])
	}

	// Test interesting loads.
	t.Run("Interesting", func(t *testing.T) {
		if len(got.Interesting) != 2 {
			t.Fatalf("want 2 Interesting, got %d (%v)", len(got.Interesting), got.Interesting)
		}

		interesting := []struct {
			want string
			got  InterestRoll
		}{
			{"1 (-1 -1, 0 0)", got.Interesting[0]},
			{"2 (2 2, 3 3, 4 4)", got.Interesting[1]},
		}

		for i, tc := range interesting {
			t.Run(strconv.Itoa(i), func(t *testing.T) {
				var members []string
				for _, opt := range tc.got.Options {
					members = append(members, fmt.Sprintf("%d %d", opt.M, opt.N))
				}
				got := fmt.Sprintf("%d (%s)", tc.got.Pick, strings.Join(members, ", "))
				if got != tc.want {
					t.Fatalf("want %s, got %s", tc.want, got)
				}
			})
		}
	})
}
