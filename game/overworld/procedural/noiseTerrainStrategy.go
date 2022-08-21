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

	// TODO: run smoothing ...
	// 🤔 there are only six neighbors, but N possible components ......

	return result
}
