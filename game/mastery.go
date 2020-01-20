package game

//go:generate stringer -type=Mastery

// Mastery enumerates the types of Masteries or attunements that are available
// to Characters to learn on their journey.
type Mastery int

const (
	// ShortRangeMeleeMastery (Daggers, Swords, Clubs, etc)
	ShortRangeMeleeMastery Mastery = iota
	// LongRangeMeleeMastery (Spears and Polearms)
	LongRangeMeleeMastery
	// RangedCombatMastery (Bows, Slingshots, Javelins, etc)
	RangedCombatMastery
	// CraftsmanshipMastery relates to non-magical utility skills
	CraftsmanshipMastery
	// FireMastery for fire spells
	FireMastery
	// WaterMastery for water spells
	WaterMastery
	// EarthMastery for earth spells
	EarthMastery
	// AirMastery for air spells
	AirMastery
	// LightningMastery for Lightning spells
	LightningMastery
	// DarkMastery for unholy spells
	DarkMastery
	// LightMastery for holy spells
	LightMastery
)
