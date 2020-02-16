// command expand-terrain expands the terrain of an *.overworld-recipe by one.
package main

import (
	"fmt"
	"math"
	"os"
	"sort"

	"github.com/griffithsh/squads/game/overworld"
	"github.com/griffithsh/squads/geom"
)

func main() {
	recipe, err := overworld.ParseRecipe(os.Stdin)
	if err != nil {
		fmt.Printf("ParseRecipe: %v\n", err)
		os.Exit(1)
	}

	expanded := make(map[geom.Key]overworld.TileID)
	for k := range recipe.Terrain {
		// If it's present in the original recipe's terrain, copy it.
		orig, ok := recipe.Terrain[k]
		if ok {
			expanded[k] = orig
		}

		for _, adj := range k.Adjacent() {
			// if it's already in the expanded set, try the next one.
			if _, ok := expanded[adj]; ok {
				continue
			}

			// If it's already in the recipe's terrain, propagate it to the expanded set.
			if t, ok := recipe.Terrain[adj]; ok {
				expanded[adj] = t
				continue
			}

			expanded[adj] = orig
		}
	}

	keys := make([]geom.Key, 0, len(expanded))
	for k := range expanded {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		if keys[i].M == keys[j].M {
			return keys[i].N < keys[j].N
		}
		return keys[i].M < keys[j].M
	})

	prevM := math.MinInt32
	for _, k := range keys {
		if prevM != k.M {
			fmt.Println()
			prevM = k.M
		}
		fmt.Printf("%d %d %d, ", k.M, k.N, expanded[k])
	}
}
