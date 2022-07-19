package procedural

import (
	"math/rand"
	"sort"

	"github.com/griffithsh/squads/geom"
)

func contains[T comparable](s []T, v T) bool {
	for _, item := range s {
		if item == v {
			return true
		}
	}
	return false
}

func sliceLengths[T any](s [][]T) []int {
	result := make([]int, len(s))
	for i, ss := range s {
		result[i] = len(ss)
	}
	return result
}

func keysOf[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func sortKeys(keys []geom.Key) {
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].M == keys[j].M {
			return keys[i].N < keys[j].N
		}
		return keys[i].M < keys[j].M
	})
}

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
