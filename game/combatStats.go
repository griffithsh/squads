package game

// CombatStats are the current stats values that are used only inside a Combat,
// and are discarded afterwards.
type CombatStats struct {
	CurrentPreparation int
	ActionPoints       int
}

// Type of this Component.
func (*CombatStats) Type() string {
	return "CombatStats"
}
