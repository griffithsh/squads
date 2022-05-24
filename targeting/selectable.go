package targeting

// SelectableType enumerates the styles of selection permissions.
type SelectableType int

//go:generate stringer -type=SelectableType

// SelectableTypeFromString converts strings to SelectableTypes. Mostly useful
// for Unmarshaling.
func SelectableTypeFromString(s string) *SelectableType {
	for i := 0; i < len(_SelectableType_index)-1; i++ {
		c := SelectableType(i)

		if c.String() == s {
			return &c
		}
	}
	return nil
}

const (
	// SelectAnywhere means that the skill can be targeted on any hex in the
	// field.
	SelectAnywhere SelectableType = iota

	// SelectWithin means that a skill can be targeted on any hex whose distance from the Origin is
	// not less than MinRange and does not exceed MaxRange.
	SelectWithin

	// SelectLinear allows selections that are in a straight line with the origin.
	// TODO: implement SelectLinear

	// Untargeted skills cannot be targeted on any specific hex.
	Untargeted
)

// This thing identifies where something is allowed to be targetted on a field.
// defines the hexes that are permitted to be selected.
type Selectable struct {
	Type SelectableType

	MinRange int
	MaxRange int

	// LinearLength    int
	// LinearDirection geom.DirectionType // should be a contextual direction - "Forward", "Back-Left"
}
