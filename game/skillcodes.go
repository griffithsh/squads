package game

//go:generate stringer -type=SkillCode

// SkillCode identifies skills.
type SkillCode int

const (
	// BasicMovement is the skill code for moving a Character.
	BasicMovement SkillCode = iota
	// BasicAttack is a generic attack to get attacking going.
	BasicAttack
	// MageLightning is a skill that can taret anywhere on the field.
	MageLightning
)

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

// TargetingForSkill maps from SkillCodes to TargetingRules.
var TargetingForSkill = map[SkillCode]TargetingRule{
	BasicMovement: TargetAnywhere,
	BasicAttack:   TargetAdjacent,
	MageLightning: TargetAnywhere,
}

// BrushForSkill maps from SkillCodes to BrushForSkill.
var BrushForSkill = map[SkillCode]TargetingBrush{
	BasicMovement: Pathfinding,
	BasicAttack:   SingleHex,
	MageLightning: SingleHex,
}
