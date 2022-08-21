package procedural

import (
	"math/rand"

	"github.com/griffithsh/squads/geom"
)

type TerrainChances []struct {
	Value  Code
	Chance int
}

func (tc TerrainChances) roll(prng *rand.Rand) Code {
	sum := 0
	for _, c := range tc {
		sum += c.Chance
	}
	got := prng.Intn(sum)
	running := 0
	for _, chance := range tc {
		if got < chance.Chance+running {
			return chance.Value
		}
		running += chance.Chance
	}

	return tc[0].Value
}

type NoiseTerrainStrategy struct {
	Smoothing  int
	Outside    Code
	Components TerrainChances
}

func (ts *NoiseTerrainStrategy) Build(prng *rand.Rand, paths Paths) map[geom.Key]Code {
	bloated := map[geom.Key]struct{}{}
	for key := range paths.Nodes {
		bloated[key] = struct{}{}
		neighbors := key.ExpandBy(0, 2)
		for _, neighbor := range neighbors {
			bloated[neighbor] = struct{}{}
		}
	}

	result := map[geom.Key]Code{}
	for _, k := range shuffledGeomKeys(prng, bloated) {
		result[k] = ts.Components.roll(prng)
	}

	// Apply smoothing.
	// If more than 3 neighbors of a hex have the same code, the hex becomes that code.
	// If exactly 3 neighbors of a hex have the same code, and the other three
	// neighbors are spread across 2 or more codes, the hex becomes the code of
	// the three neighbors.
	nextGen := map[geom.Key]Code{}
	for i := 0; i < ts.Smoothing; i++ {
		for k, code := range result {
			m := map[geom.Key]Code{}
			for _, key := range k.Adjacent() {
				if code, ok := result[key]; ok {
					m[key] = code
				} else {
					m[key] = ts.Outside
				}
			}

			grouped := counts(m)
			nextGen[k] = code
			for v, count := range grouped {
				if (count == 3 && len(grouped) > 3) || count > 3 {
					nextGen[k] = v
				}
			}
		}
		nextGen, result = result, nextGen
	}

	return result
}
