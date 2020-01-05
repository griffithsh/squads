package geom

// Key is a way of referencing a hypothetical hexagon.
type Key struct {
	M, N int
}

// Equal determines if the M and N values of the passed pointers differ. If
// either value is nil, then it only returns true if the other is also nil.
func Equal(a, b *Key) bool {
	if a == nil && b != nil {
		return false
	} else if a != nil && b == nil {
		return false
	} else if a == nil && b == nil {
		return true
	} else {
		return a.M == b.M && a.N == b.N
	}
}

// Neighbors calculates the neighbors of a Key and returns them keyed by their Keys.
func (k Key) Neighbors() map[Key]DirectionType {
	result := map[Key]DirectionType{
		Key{k.M, k.N - 2}: N,
		Key{k.M, k.N + 2}: S,
	}

	if k.N%2 == 0 {
		// then the E ones have the same M, and the W ones are -1 M
		result[Key{k.M - 1, k.N - 1}] = NW
		result[Key{k.M - 1, k.N + 1}] = SW
		result[Key{k.M, k.N + 1}] = SE
		result[Key{k.M, k.N - 1}] = NE
	} else {
		// then the E ones are +1 M, and the W ones have the same M
		result[Key{k.M, k.N - 1}] = NW
		result[Key{k.M, k.N + 1}] = SW
		result[Key{k.M + 1, k.N + 1}] = SE
		result[Key{k.M + 1, k.N - 1}] = NE
	}
	return result
}

// Adjacent calculates the neighbors of a Key and returns them keyed by direction.
func (k Key) Adjacent() map[DirectionType]Key {
	result := map[DirectionType]Key{
		N: Key{k.M, k.N - 2},
		S: Key{k.M, k.N + 2},
	}

	if k.N%2 == 0 {
		// then the E ones have the same M, and the W ones are -1 M
		result[NW] = Key{k.M - 1, k.N - 1}
		result[SW] = Key{k.M - 1, k.N + 1}
		result[SE] = Key{k.M, k.N + 1}
		result[NE] = Key{k.M, k.N - 1}
	} else {
		// then the E ones are +1 M, and the W ones have the same M
		result[NW] = Key{k.M, k.N - 1}
		result[SW] = Key{k.M, k.N + 1}
		result[SE] = Key{k.M + 1, k.N + 1}
		result[NE] = Key{k.M + 1, k.N - 1}
	}
	return result
}
