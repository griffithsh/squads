package geom

import "strconv"

// Key is a way of referencing a Hexagon in a Field.
type Key struct {
	M, N int
}

func (k Key) String() string {
	return "(" + strconv.Itoa(k.M) + "," + strconv.Itoa(k.N) + ")"
}

// Equal determines if the M and N values of the passed pointers differ. If
// either value is nil, then it only returns true if the other is also nil.
// TODO: this may make more sense in the combat package as that is its only usage.
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

// ToN returns the Key that is to the North of that Key.
func (k Key) ToN() Key {
	return Key{k.M, k.N - 1}
}

// ToS returns the Key that is to the South of that Key.
func (k Key) ToS() Key {
	return Key{k.M, k.N + 1}
}

// ToNW returns the Key that is to the Northwest of that Key.
func (k Key) ToNW() Key {
	if k.M%2 != 0 {
		return Key{k.M - 1, k.N}
	}
	return Key{k.M - 1, k.N - 1}
}

// ToSW returns the Key that is to the Southwest of that Key.
func (k Key) ToSW() Key {
	if k.M%2 != 0 {
		return Key{k.M - 1, k.N + 1}
	}
	return Key{k.M - 1, k.N}
}

// ToNE returns the Key that is to the Northeast of that Key.
func (k Key) ToNE() Key {
	if k.M%2 != 0 {
		return Key{k.M + 1, k.N}
	}
	return Key{k.M + 1, k.N - 1}
}

// ToSE returns the Key that is to the Southeast of that Key.
func (k Key) ToSE() Key {
	if k.M%2 != 0 {
		return Key{k.M + 1, k.N + 1}
	}
	return Key{k.M + 1, k.N}
}
func (k Key) ToDirection(dir DirectionType) Key {
	switch dir {
	case N:
		return k.ToN()
	case NE:
		return k.ToNE()
	case SE:
		return k.ToSE()
	case S:
		return k.ToS()
	case SW:
		return k.ToSW()
	case NW:
		return k.ToNW()
	default:
		return k
	}
}

// Neighbors calculates the neighbors of a Key and returns them keyed by their Keys.
func (k Key) Neighbors() map[Key]DirectionType {
	result := map[Key]DirectionType{
		{k.M, k.N - 1}: N,
		{k.M, k.N + 1}: S,
	}

	if k.M%2 == 0 {
		// Then the Southern ones have the same N, and the Northern ones are -1 N.
		result[Key{k.M - 1, k.N - 1}] = NW
		result[Key{k.M - 1, k.N}] = SW
		result[Key{k.M + 1, k.N}] = SE
		result[Key{k.M + 1, k.N - 1}] = NE
	} else {
		// Then the Southern ones are +1 N, and the Northern ones have the same N.
		result[Key{k.M - 1, k.N}] = NW
		result[Key{k.M - 1, k.N + 1}] = SW
		result[Key{k.M + 1, k.N + 1}] = SE
		result[Key{k.M + 1, k.N}] = NE
	}
	return result
}

// Adjacent calculates the neighbors of a Key and returns them keyed by direction.
func (k Key) Adjacent() map[DirectionType]Key {
	return map[DirectionType]Key{
		N:  k.ToN(),
		S:  k.ToS(),
		SE: k.ToSE(),
		SW: k.ToSW(),
		NE: k.ToNE(),
		NW: k.ToNW(),
	}
}

// HexesFrom calculates how many Hexes away another Key is.
func (k Key) HexesFrom(other Key) int {
	mDiff := k.M - other.M
	// Convert diff to absolute.
	if mDiff < 0 {
		mDiff = -mDiff
	}

	// if M is odd ...
	minN := k.N - (mDiff / 2)
	maxN := k.N + ((1 + mDiff) / 2)
	// else if M is even
	if k.M%2 == 0 {
		minN = k.N - ((1 + mDiff) / 2)
		maxN = k.N + (mDiff / 2)
	}

	if other.N > maxN {
		return mDiff + other.N - maxN
	} else if other.N < minN {
		return mDiff + minN - other.N
	}
	return mDiff
}

// ExpandBy determines the Keys that are between min and max away from the Key. The order is randomised.
func (k Key) ExpandBy(min, max int) []Key {
	inRange := []Key{}

	if min == 0 {
		inRange = append(inRange, k)
	}
	opens := map[Key]struct{}{
		k: {},
	}
	closed := map[Key]struct{}{
		k: {},
	}

	for i := 0; i < max; i++ {
		// extract the keys of everything in open
		keys := make([]Key, 0, len(opens)*6)
		for open := range opens {
			for k2 := range open.Neighbors() {
				if _, ok := closed[k2]; ok {
					continue
				}
				keys = append(keys, k2)
			}
		}

		// empty opens
		opens = map[Key]struct{}{}

		for _, it := range keys {
			if _, ok := closed[it]; ok {
				continue
			}
			opens[it] = struct{}{}
			closed[it] = struct{}{}

			// Is this between min and max?
			if i+1 >= min {
				inRange = append(inRange, it)
			}
		}
	}

	return inRange
}
