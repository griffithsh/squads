package mathx

// MinI returns the smaller of the two passed ints.
func MinI(a, b int) int {
	if a > b {
		return b
	}
	return a
}

// MinF64 returns the smaller of the two passed values.
func MinF64(a, b float64) float64 {
	if a > b {
		return b
	}
	return a
}
