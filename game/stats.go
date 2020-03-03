package game

func MaxHealth(rawHealth, vitality int) int {
	const healthPerVit float64 = 9.5
	// TODO: Also include vitality and base health from equipped items.

	return rawHealth + int(healthPerVit*float64(vitality)) + 25

}
