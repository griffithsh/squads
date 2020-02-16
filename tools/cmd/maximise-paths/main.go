// command maximise-terrain calculates all possible connections between
// overworld nodes.
package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/griffithsh/squads/game/overworld"
)

func main() {
	recipe, err := overworld.ParseRecipe(os.Stdin)
	if err != nil {
		fmt.Printf("ParseRecipe: %v\n", err)
		os.Exit(1)
	}

	mapped := map[overworld.KeyPair]struct{}{}
	for k := range recipe.Terrain {
		for _, adj := range k.Adjacent() {
			if _, ok := recipe.Terrain[adj]; !ok {
				continue
			}

			// Does the inverse already exist?
			inverse := overworld.KeyPair{
				First:  adj,
				Second: k,
			}
			if _, ok := mapped[inverse]; ok {
				continue
			}

			mapped[overworld.KeyPair{First: k, Second: adj}] = struct{}{}
		}
	}

	sliced := make([]overworld.KeyPair, 0, len(mapped))

	for k := range mapped {
		sliced = append(sliced, k)
	}

	sort.Slice(sliced, func(i, j int) bool {
		if sliced[i].First.M != sliced[j].First.M {
			return sliced[i].First.M < sliced[j].First.M
		}
		if sliced[i].First.N != sliced[j].First.N {
			return sliced[i].First.N < sliced[j].First.N
		}
		if sliced[i].Second.M != sliced[j].Second.M {
			return sliced[i].Second.M < sliced[j].Second.M
		}
		return sliced[i].Second.N < sliced[j].Second.N
	})

	ss := []string{}
	prevM := math.MinInt32
	for _, kp := range sliced {
		if kp.First.M != prevM {
			fmt.Fprintln(os.Stdout)
			ss = append(ss, "\n")
			prevM = kp.First.M
		}
		ss = append(ss, fmt.Sprintf("%d %d %d %d, ", kp.First.M, kp.First.N, kp.Second.M, kp.Second.N))
	}
	fmt.Fprintf(os.Stdout, "%s", strings.Join(ss, ""))
}
