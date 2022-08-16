package procedural

import (
	"fmt"
	"testing"

	"github.com/griffithsh/squads/geom"
)

func TestBuildVinePaths(t *testing.T) {
	for i := 0; i < 1000; i++ {
		seed := int64(i) + 123456789
		t.Run(fmt.Sprintf("seed=%v", seed), func(t *testing.T) {
			paths, err := buildVinePaths(seed, 0)
			if err != nil {
				t.Fatal(err)
			}

			if geom.Equal(&paths.Start, &paths.Goal) {
				t.Errorf("goal is the start")
			}
			if _, ok := paths.Nodes[paths.Start]; !ok {
				t.Errorf("missing start key")
			}
			if _, ok := paths.Nodes[paths.Goal]; !ok {
				t.Errorf("missing goal key")
			}
		})
	}
}
