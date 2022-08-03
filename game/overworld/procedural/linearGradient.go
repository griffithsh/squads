package procedural

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/griffithsh/squads/geom"
)

//go:generate go run github.com/dmarkham/enumer -output=./linearGradient_enumer.go -type=BlendType,StrategyTargetFilter -json

type BlendType int

const (
	Noisy BlendType = iota
	Smooth
	Spiky
)

type Blend struct {
	Value Code
	Type  BlendType
}

type LinearTerrainGradient struct {
	Portions float64 // How much of the entire map this gradient is responsible for.
	Value    Code
	Blend    *Blend
}

type LinearTerrainGradientSlice []LinearTerrainGradient

func (gradients LinearTerrainGradientSlice) SubGradient(percent float64) (int, float64) {
	if percent < 0 {
		return -1, percent
	}
	sum := 0.0
	for _, g := range gradients {
		sum += g.Portions
	}

	running := 0.0
	for i, g := range gradients {
		percentage := g.Portions / sum
		if percent <= percentage+running && percent >= running {
			// this is the selected gradient
			return i, (percent - running) / percentage
		}
		running += percentage
	}

	return len(gradients), percent
}

type StrategyTargetFilter int

const (
	AnyTarget StrategyTargetFilter = iota
	Widest
	Narrowest
	// Diagonals
	// NortheastSouthwest
	// SoutheastNorthWest
	// NorthSouth
)

func (f StrategyTargetFilter) Generate(prng *rand.Rand, extents map[geom.DirectionType]geom.Key) geom.DirectionType {
	targets := []geom.DirectionType{}
	switch f {
	case AnyTarget:
		targets = []geom.DirectionType{geom.N, geom.S, geom.SE, geom.SW, geom.NE, geom.NW}
	case Widest:
		widest := extents[geom.N].HexesFrom(extents[geom.S])
		targets = []geom.DirectionType{geom.N, geom.S}
		NESW := extents[geom.NE].HexesFrom(extents[geom.SW])
		if NESW > widest {
			widest = NESW
			targets = []geom.DirectionType{geom.NE, geom.SW}
		}
		SENW := extents[geom.NW].HexesFrom(extents[geom.SE])
		if SENW > widest {
			targets = []geom.DirectionType{geom.SE, geom.NW}
		}
	case Narrowest:
		narrowest := extents[geom.N].HexesFrom(extents[geom.S])
		targets = []geom.DirectionType{geom.N, geom.S}
		NESW := extents[geom.NE].HexesFrom(extents[geom.SW])
		if NESW < narrowest {
			narrowest = NESW
			targets = []geom.DirectionType{geom.NE, geom.SW}
		}
		SENW := extents[geom.NW].HexesFrom(extents[geom.SE])
		if SENW < narrowest {
			targets = []geom.DirectionType{geom.SE, geom.NW}
		}
	}
	return targets[DeterministicIndexOf(prng, targets)]
}

type LinearGradientTerrainStrategy struct {
	TargetFilter StrategyTargetFilter
	Overflows    Code
	Underflows   Code
	Gradients    LinearTerrainGradientSlice
}

func (ts *LinearGradientTerrainStrategy) Build(prng *rand.Rand, paths map[geom.Key]Placement) map[geom.Key]Code {
	// Expand paths ...
	bloated := map[geom.Key]struct{}{}
	for key := range paths {
		bloated[key] = struct{}{}
		neighbors := key.ExpandBy(0, 2)
		for _, neighbor := range neighbors {
			bloated[neighbor] = struct{}{}
		}
	}

	extents := extentsOf(keysOf(paths))
	target := ts.TargetFilter.Generate(prng, extents)

	// Generate a terrain code for each key.
	result := map[geom.Key]Code{}
	for _, k := range shuffledGeomKeys(prng, bloated) {
		code := ts.terrainFor(prng, k, target, extents)
		result[k] = code
	}

	return result
}

func (ts *LinearGradientTerrainStrategy) terrainFor(prng *rand.Rand, query geom.Key, target geom.DirectionType, extents map[geom.DirectionType]geom.Key) Code {
	percent, perpendicular := linearPercentThrough(query, target, extents)
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

// linearPercentThrough returns how far through the distance between target and it's
// opposite query is in the range 0.0 to 1.0. The second return value is a
// relative interpretation of the perpendicular value of the query.
func linearPercentThrough(query geom.Key, target geom.DirectionType, extents map[geom.DirectionType]geom.Key) (float64, int) {
	// Magic numbers to create a field where rotating 45 degrees lines up the
	// the two NE-SW, NW-SE vectors with E-W.
	f := geom.NewField(36, 16, 34)
	radians := math.Pi / (4.050063485188)
	if target == geom.NW || target == geom.SE {
		radians = -radians
	}
	magicPerpendicularDeNoiser := 49.0

	// test := geom.Key{M: 0, N: 0}
	// for i := 0; i < 100; i++ {
	// 	test = test.ToN().ToNW()
	// }
	// x, y := f.Ktow(test)
	// sin, cos := math.Sincos(-radians)
	// y = x*sin + y*cos
	// fmt.Printf("Zero? %f\n", y)

	x0, y0 := f.Ktow(extents[target])
	x1, y1 := f.Ktow(extents[geom.Opposite[target]])
	xQ, yQ := f.Ktow(query)

	// A for actual, or answer versions. N and S use the y coordinate, and the
	// other directions use the x.
	a0, a1, aQ := 0.0, 0.0, 0.0
	perpendicularDisambiguator := 0

	switch target {
	case geom.N, geom.S:
		a0, a1, aQ = y0, y1, yQ
		perpendicularDisambiguator = int(math.Round(xQ / magicPerpendicularDeNoiser))

	default:
		sin, cos := math.Sincos(radians)
		a0 = x0*cos - y0*sin // y=x0*sin+y0*cos
		a1 = x1*cos - y1*sin // y=x1*sin+y1*cos
		aQ, yQ = xQ*cos-yQ*sin, xQ*sin+yQ*cos
		perpendicularDisambiguator = int(math.Round(yQ / magicPerpendicularDeNoiser))
	}

	return 1.0 - float64(aQ-a1)/float64(a0-a1), perpendicularDisambiguator
}
