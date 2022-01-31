package item

import "github.com/griffithsh/squads/skill"

//go:generate stringer -type=Class

// Class is a low-specificity delineation of item. Within each Class is a
// collection of more specific types of that class captured as a Code. BowClass
// might contain a long bow and a short bow. The SwordClass might contain a fast
// sword and a slow sword.
type Class int

const (
	UnarmedClass Class = iota // FIXME?
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

func (c Class) IsWeapon() bool {
	return c >= UnarmedClass && c <= WandClass
}

// Instance is a rolled item that can be equipped.
type Instance struct {
	Class Class

	// Code represents a more specific implementation of the ItemClass. If an item's
	// name in English is Short Sword, the Code might be "short_sword". Codes
	// should be configured in separate game data, and loaded into the game at
	// runtime.
	Code string

	// Name is rendered at recipe execution time and might be something like "Deadly
	// Axe of Iciness"
	Name      string
	Modifiers map[Modifier]float64 // base damage, or base armor, or any other modifier

	Skills []skill.ID
}
