package game

//go:generate stringer -type=SkillCode

// SkillCode identifies skills.
type SkillCode int

const (
	// BasicMovement is the skill code for moving a Character.
	BasicMovement SkillCode = iota
	// BasicAttack is a generic attack to get attacking going.
	BasicAttack
)
