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

