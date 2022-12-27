package procedural

import (
	"fmt"
	"testing"

	"github.com/griffithsh/squads/geom"
)

func TestBuildDoubleBlobPaths(t *testing.T) {
	for i := 0; i < 250; i++ {
		seed := int64(i) + 1234567891011
		var errs []error
		t.Run(fmt.Sprintf("seed=%v", seed), func(t *testing.T) {
			paths, err := buildDoubleBlobPaths(seed, 0)
			if err != nil {
				errs = append(errs, err)
				return
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
		if len(errs) > 2 {
			t.Errorf("above 1%% failure rate: %v", errs)
		}
	}
}
