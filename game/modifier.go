package game

//go:generate stringer -type=Modifier
type Modifier int

const (
	BaseMinDamageModifier Modifier = iota
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
)
