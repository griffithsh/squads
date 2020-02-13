package overworld

import (
	"bytes"
	"testing"

	"github.com/griffithsh/squads/geom"
)

func TestParseRecipe(t *testing.T) {
	reference := []byte(`label: Reference Recipe

terrain: // m,n tileid
0 0 1, 0 1 1,
0 2 0

	`)

	got, err := ParseRecipe(bytes.NewReader(reference))
	if err != nil {
		t.Fatalf("parse reference recipe: %v", err)
	}
	if got.Label != "Reference Recipe" {
		t.Errorf("want label %s, got label %s", "Reference Recipe", got.Label)
	}

	if len(got.Terrain) != 3 {
		t.Fatalf("want 3 terrain tiles")
	}
	if got.Terrain[geom.Key{M: 0, N: 0}] != TileID(1) {
		t.Errorf("want 0,0 to equal 1")
	}
	if got.Terrain[geom.Key{M: 0, N: 1}] != TileID(1) {
		t.Errorf("want 0,1 to equal 1")
	}
	if got.Terrain[geom.Key{M: 0, N: 2}] != TileID(0) {
		t.Errorf("want 0,2 to equal 0")
	}
}
