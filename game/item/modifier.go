package item

//go:generate stringer -type=Modifier

// Modifier enumerates the stat modifiers that appear on items and effects in
// the game.
type Modifier int

const (
	// BaseMinDamageModifier should typically appear only on weapons. It is a
	// value, not a multiplier
	BaseMinDamageModifier Modifier = iota
	// BaseMaxDamageModifier should typically appear only on weapons. It is a
	// value, not a multiplier
	BaseMaxDamageModifier

	// BaseDamageModifier multiplies the base damage before it is multiplied by
	// DamageMultiplierModifier. Should typically only appear on Weapons. Maybe
	// should never go past 99%?
	BaseDamageModifier

	// DamageMultiplierModifier is accumulated from modifiers present on
	// non-weapon items and also core stats.
	DamageMultiplierModifier

	// AdditionalMinDamageModifier
	// AdditionalMaxDamageModifier

	// PreparationModifier is added to the Character's PreparationThreshold.
	PreparationModifier
	//ActionPointModifier is added to the Character's ActionPoint maximum.
	ActionPointModifier

	// ChanceToHitModifier multiplies the base chance to hit of Attacks.
	// A value of zero does not modify the chance to hit. A value of 0.1
	// improves the chance to hit by 10%. A value of -0.5 halves the chance to
	// hit.
	ChanceToHitModifier
)
