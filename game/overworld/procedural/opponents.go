package procedural

import (
	"math/rand"

	"github.com/griffithsh/squads/geom"
	"github.com/griffithsh/squads/squad"
)

type OpponentSquad struct {
	ID           squad.RecipeID
	Chance       int
	Min          int
	Max          int
	OnlySpawnsOn []Code
}

type OpponentSquads []OpponentSquad

func (s OpponentSquads) Generate(prng *rand.Rand, paths Paths, terrainCodes map[geom.Key]Code) map[geom.Key]squad.RecipeID {
	result := map[geom.Key]squad.RecipeID{}
	for _, k := range shuffledGeomKeys(prng, paths.Nodes) {
		if prng.Float64() > 0.25 {
			continue
		}

		// Don't spawn baddies on the start or the goal.
		if k == paths.Goal || k == paths.Start {
			continue
		}

		// Which terrain are we dealing with? Some opponents will only spawn on
		// certain terrain.
		terrain := terrainCodes[k]

		// Calculate counts of baddies already placed.
		counts := counts(result)

		generalContenders := OpponentSquads{}
		unsatisfiedMins := OpponentSquads{}

		for _, b := range s {
			if len(b.OnlySpawnsOn) > 0 {
				// Then terrain must appear in OnlySpawnsOn.
				disallowed := true
				for _, code := range b.OnlySpawnsOn {
					if terrain == code {
						disallowed = false
						break
					}
				}
				if disallowed {
					continue
				}
			}

			if b.Min > counts[b.ID] {
				unsatisfiedMins = append(unsatisfiedMins, b)
			}

			if b.Max == 0 || b.Max > counts[b.ID] {
				generalContenders = append(generalContenders, b)
			}
		}

		contenders := generalContenders
		if len(unsatisfiedMins) > 0 {
			contenders = unsatisfiedMins
		}

		sum := 0
		for _, baddy := range contenders {
			sum += baddy.Chance
		}

		if sum <= 0 {
			continue
		}
		roll := prng.Intn(sum)
		running := 0
		for _, baddy := range contenders {
			if roll < baddy.Chance+running {
				// Got it.
				result[k] = baddy.ID
				break
			}
			running += baddy.Chance
		}
	}
	return result
}
