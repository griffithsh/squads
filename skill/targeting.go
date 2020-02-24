package skill

// TargetingRule identifies where a skill can be targeted on a field.
type TargetingRule int

const (
	// TargetAnywhere means that the skill can be targeted on any hex in the
	// field.
	TargetAnywhere TargetingRule = iota

	// TargetAdjacent means that a skill can be targeted on any hex adjacent to
	// the origin.
	TargetAdjacent
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
)
