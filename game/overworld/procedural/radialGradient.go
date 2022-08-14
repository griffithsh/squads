package procedural

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/griffithsh/squads/geom"
)

type RadialGradientTerrainStrategy struct {
	Overflows  Code
	Underflows Code
	Gradients  TerrainGradientSlice

	center   geom.Key
	nearest  float64
	furthest float64
	f        *geom.Field
}

func (ts *RadialGradientTerrainStrategy) Build(prng *rand.Rand, paths map[geom.Key]Placement) map[geom.Key]Code {
	ts.f = geom.NewField(36, 16, 48)
	ts.nearest = math.MaxFloat64
	ts.furthest = 0.0

	for k := range paths {
		dist := ts.f.DistanceBetween(k, ts.center)
		if dist < ts.nearest {
			ts.nearest = dist
		}
		if dist > ts.furthest {
			ts.furthest = dist
		}
	}

	// Expand paths ...
	bloated := map[geom.Key]struct{}{}
	for key := range paths {
		bloated[key] = struct{}{}
		neighbors := key.ExpandBy(0, 7)
		for _, neighbor := range neighbors {
			bloated[neighbor] = struct{}{}
		}
	}

	// Generate a terrain code for each key.
	result := map[geom.Key]Code{}
	for _, k := range shuffledGeomKeys(prng, bloated) {
		code := ts.terrainFor(prng, k)
		result[k] = code
	}

	return result
}

func (ts *RadialGradientTerrainStrategy) terrainFor(prng *rand.Rand, query geom.Key) Code {
	percent, perpendicular := ts.radialPercentAway(query)
	i, subPercent := ts.Gradients.SubGradient(percent)

	if i < 0 {
		return ts.Underflows
	} else if i >= len(ts.Gradients) {
		return ts.Overflows
	} else {
		g := ts.Gradients[i]
		if g.Blend != nil {
			// offset is effectively a hash of the perpendicular value of the query,
			// mapped to 0.0 - 1.0. It's a way to generate non-linear boundary lines
			// between gradient segments deterministically. The desired outcome is a
			// non-noisy separation between two regions where there are no disconnected
			// hexes of either side in the other.
			offset := rand.New(rand.NewSource(int64(perpendicular))).Float64()

			switch g.Blend.Type {
			case Noisy:
				if subPercent > prng.Float64() {
					return g.Blend.Value
				} else {
					return g.Value
				}
			case Smooth:
				// Smooth the offset by averaging it with its neighbors.
				offsetA := rand.New(rand.NewSource(int64(perpendicular + 1))).Float64()
				offsetB := rand.New(rand.NewSource(int64(perpendicular - 1))).Float64()
				offset = (offset + offsetA/4 + offsetB/4) / 1.5
				fallthrough
			case Spiky:
				if subPercent > offset {
					return g.Blend.Value
				} else {
					return g.Value
				}
			default:
				panic(fmt.Sprintf("unhandled BlendType: %v", g.Blend.Type))
			}
		} else {
			return g.Value
		}
	}
}

// radialPercentAway returns how far away from the center the queried key is, as
// a ratio of the nearest and furthest keys. It also returns a rough indicator
// of rotation around the center the queried key is.
func (ts *RadialGradientTerrainStrategy) radialPercentAway(query geom.Key) (float64, int) {
	dist := ts.f.DistanceBetween(ts.center, query)
	percent := (dist - ts.nearest) / (ts.furthest - ts.nearest)

	cx, cy := ts.f.Ktow(ts.center)
	px, py := ts.f.Ktow(query)
	magicalDisambiguator := int(rotation(cx, cy, px, py) * 10)

	return percent, magicalDisambiguator
}
