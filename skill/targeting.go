package skill

// TargetingRule identifies where a skill can be targeted on a field.
type TargetingRule int

const (
	// TargetAnywhere means that the skill can be targeted on any hex in the
	// field.
	TargetAnywhere TargetingRule = iota

	// TargetAdjacent means that a skill can be targeted on any hex adjacent to
	// the origin.
	// FIXME: should be TargetWithin, and introduce a distance field to implement
	// the size of the aoe.
	TargetAdjacent

	// TargetLinear allows selections that are in a straight line with the origin.
	// TODO: implement TargetLinear
)

// TargetingBrush enumerates the different rule sets for which hex or hexes are
// highlighted given a target and origin.
type TargetingBrush int

const (
	// SingleHex means that only the targeted Hex is highlighted.
	SingleHex TargetingBrush = iota

	// Pathfinding is a special rule set that highlights a path of hexes on the
	// way from the origin to the target.
	Pathfinding

	// AreaOfEffect indicates that the target and all hexes within a configured
	// distance are highlighted.
	// TODO: implement AreaOfEffect
)
