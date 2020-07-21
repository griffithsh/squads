package geom

import (
	"fmt"
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
