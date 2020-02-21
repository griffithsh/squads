package game

//go:generate stringer -type=ItemClass

// ItemClass is a low-specificity delineation of item.
type ItemClass int

const (
	UnarmedClass ItemClass = iota
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

// ItemInstance is a rolled item that can be equipped.
type ItemInstance struct {
	Class     ItemClass
	Name      string
	Modifiers map[Modifier]float64 // base damage, or base armor, or any other modifier
}
