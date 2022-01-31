package game

//go:generate stringer -type=ItemClass

// ItemClass is a low-specificity delineation of item. Within each Class is a
// collection of more specific types of that class captured as a Code. BowClass
// might contain a long bow and a short bow. The SwordClass might contain a fast
// sword and a slow sword.
type ItemClass int

const (
	UnarmedClass ItemClass = iota // FIXME?
	SwordClass
	AxeClass
	ClubClass
	DaggerClass
	SlingClass
	BowClass
	SpearClass
	PolearmClass
	StaffClass
	WandClass

	HatClass
	BodyArmorClass
	AmuletClass
	RingClass
	GloveClass
	BootClass
	BeltClass
)

func (c ItemClass) IsWeapon() bool {
	return c >= UnarmedClass && c <= WandClass
}

// ItemInstance is a rolled item that can be equipped.
type ItemInstance struct {
	Class ItemClass

	// Code represents a more specific implementation of the ItemClass. If an item's
	// name in English is Short Sword, the Code might be "short_sword". Codes
	// should be configured in separate game data, and loaded into the game at
	// runtime.
	Code string

	// Name is rendered at recipe execution time and might be something like "Deadly
	// Axe of Iciness"
	Name      string
	Modifiers map[Modifier]float64 // base damage, or base armor, or any other modifier

	// Skills []skill.ID
}
