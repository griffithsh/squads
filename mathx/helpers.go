package mathx

// MinI returns the smaller of the two passed ints.
func MinI(a, b int) int {
	if a > b {
		return b
	}
	return a
}
