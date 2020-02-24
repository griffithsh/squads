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

// TODO: There is an item ... "type" here as well, because a class like swords
// should include some fast swords and some slow swords. The prep and AP values
// should be configured at this intermediate level between class and instance.
// This *might* be the Recipe level?

// ItemInstance is a rolled item that can be equipped.
type ItemInstance struct {
	Class     ItemClass
	Name      string
	Modifiers map[Modifier]float64 // base damage, or base armor, or any other modifier
}
