package geom

func neighbors(M, N int) []Key {
	result := []Key{
		Key{M, N - 2}, // North
		Key{M, N + 2}, // South
	}

	if N%2 == 0 {
		// then the E ones have the same M, and the W ones are -1 M
		result = append(result, []Key{
			{M - 1, N - 1}, // NW
			{M - 1, N + 1}, // SW
			{M, N + 1},     // SE
			{M, N - 1},     // NE
		}...)
	} else {
		// then the E ones are +1 M, and the W ones have the same M
		result = append(result, []Key{
			{M, N - 1},     // NW
			{M, N + 1},     // SW
			{M + 1, N + 1}, // SE
			{M + 1, N - 1}, // NE
		}...)
	}

	return result
}

// KeySet is a set of Keys.
type KeySet map[Key]struct{}

// Without creates a new Keyset without the provided Keys to exclude.
func (ks KeySet) Without(exclude KeySet) KeySet {
	result := make(KeySet)

	for k := range ks {
		if _, found := exclude[k]; !found {
			result[k] = struct{}{}
		}
	}
	return result
}

// Overlaps create a new KeySet of the set of Keys that are in the original
// KeySet and the passed KeySet.
func (ks KeySet) Overlaps(overlaps KeySet) KeySet {
	result := make(KeySet)

	for k := range ks {
		if _, ok := overlaps[k]; ok {
			result[k] = struct{}{}
		}
	}
	return result
}

// NeighborSet returns the neighbors of a Key as a unique set.
func NeighborSet(M, N int) KeySet {
	keys := neighbors(M, N)
	result := make(KeySet)
	for _, k := range keys {
		result[k] = struct{}{}
	}
	return result
}

// XY calculates the X and Y centre of a hexagon.
func XY(m, n, hexW, hexH int) (float64, float64) {
	xOffset := hexW - ((hexH - 2) / 2)
	x := float64(hexW/2) + float64(m*2*xOffset) + float64(n%2*xOffset)
	y := float64(hexH/2) + float64((hexH/2)*n)

	return x, y
}

// MByN creates a slice of Keys that fill a rectangular field of m by n.
func MByN(m, n int) []Key {
	var result []Key
	for ni := 0; ni < n; ni++ {
		for mi := 0; mi < m; mi++ {
			result = append(result, Key{M: mi, N: ni})
		}
	}
	return result
}

// Adjacent determines if two Keys are next to each other.
func Adjacent(a, b Key) bool {
	if a.M == b.M && (a.N-b.N == 2 || a.N-b.N == -2) {
		return true
	}
	if a.N%2 == 1 {
		if a.M == b.M {
			if a.N-b.N == 1 || a.N-b.N == -1 {
				return true
			}
		} else if a.M-b.M == 1 || a.M-b.M == -1 {
			if a.N-b.N == 1 || a.N-b.N == -1 {
				return true
			}
		}
	} else {
		if a.M == b.M {
			if a.N-b.N == 1 || a.N-b.N == -1 {
				return true
			}
		} else if a.M-b.M == 1 || a.M-b.M == -1 {
			if a.N-b.N == 1 || a.N-b.N == -1 {
				return true
			}
		}
	}
	return false
}
