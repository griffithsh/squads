package procedural

import (
	"math"
	"math/rand"
	"sort"

	"github.com/griffithsh/squads/geom"
)

// keysOf a map in non-deterministic order.
func keysOf[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// sortKeys sorts a slice of geom.Keys in place.
func sortKeys(keys []geom.Key) {
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].M == keys[j].M {
			return keys[i].N < keys[j].N
		}
		return keys[i].M < keys[j].M
	})
}

// shuffleSlice in-place, mutating the passed slice.
func shuffleSlice[V any](prng *rand.Rand, s []V) {
	prng.Shuffle(len(s), func(i, j int) {
		s[i], s[j] = s[j], s[i]
	})
}

// shuffledGeomKeys extracts the keys of the passed map in a deterministically
// randomised way.
func shuffledGeomKeys[V any](prng *rand.Rand, m map[geom.Key]V) []geom.Key {
	keys := keysOf(m)
	if len(keys) > 1 {
		sortKeys(keys)
		prng.Shuffle(len(keys), func(i, j int) {
			keys[i], keys[j] = keys[j], keys[i]
		})
	}
	return keys
}

func deterministicKeyFrom[V any](m map[geom.Key]V) geom.Key {
	keys := keysOf(m)
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].M == keys[j].M {
			return keys[i].N < keys[j].N
		}
		return keys[i].M < keys[j].M
	})
	return keys[0]
}

func DeterministicIndexOf[V any](prng *rand.Rand, s []V) int {
	switch len(s) {
	case 0:
		return -1
	case 1:
		return 0
	default:
		return prng.Intn(len(s))
	}
}

// extentsOf of a set of keys. The most-northerly, southerly, north-westerly, etc...
func extentsOf(keys []geom.Key) map[geom.DirectionType]geom.Key {
	// Magic numbers to create a field where rotating 45 degrees lines up the
	// the two NE-SW, NW-SE vectors with E-W.
	f := geom.NewField(36, 16, 34)
	rad := math.Pi / 4

	sinSWNE, cosSWNE := math.Sincos(-rad)
	sinNWSE, cosNWSE := math.Sincos(rad)

	most := map[geom.DirectionType]int{}
	result := map[geom.DirectionType]geom.Key{}

	for _, key := range keys {
		x, y := f.Ktow(key)

		// Is it the smallest y we've seen? - new N
		if v, ok := most[geom.N]; !ok || int(y) < v {
			most[geom.N] = int(y)
			result[geom.N] = key
		}

		// Is it the largest y we've seen? - new S
		if v, ok := most[geom.S]; !ok || int(y) > v {
			most[geom.S] = int(y)
			result[geom.S] = key
		}

		// xSWNE := x*cosSWNE - y*sinSWNE
		ySWNE := x*sinSWNE + y*cosSWNE

		// Is it the smallest y we've seen? - new NE
		if v, ok := most[geom.NE]; !ok || v > int(ySWNE) {
			most[geom.NE] = int(ySWNE)
			result[geom.NE] = key
		}

		// Is it the largest y we've seen? - new SW
		if v, ok := most[geom.SW]; !ok || v < int(ySWNE) {
			most[geom.SW] = int(ySWNE)
			result[geom.SW] = key
		}

		// xNWSE := x*cosNWSE - y*sinNWSE
		yNWSE := x*sinNWSE + y*cosNWSE

		// Is it the smallest y we've seen? - new NW
		if v, ok := most[geom.NW]; !ok || v > int(yNWSE) {
			most[geom.NW] = int(yNWSE)
			result[geom.NW] = key
		}

		// Is it the largest y we've seen? - new SE
		if v, ok := most[geom.SE]; !ok || v < int(yNWSE) {
			most[geom.SE] = int(yNWSE)
			result[geom.SE] = key
		}
	}
	return result
}

func leftOfLine(ax, ay, bx, by, cx, cy float64) bool {
	return ((bx-ax)*(cy-ay) - (by-ay)*(cx-ax)) < 0
}

// rotation calculates the angle (in radians) that p has been rotated around
// center.
func rotation(cx, cy, px, py float64) float64 {
	return math.Atan2(px-cx, py-cy)
}

// counts returns the counts of times each value appears in the passed map.
func counts[K, V comparable](m map[K]V) map[V]int {
	result := map[V]int{}
	for _, v := range m {
		result[v] = result[v] + 1
	}
	return result
}
