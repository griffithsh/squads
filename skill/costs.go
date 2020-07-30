package skill

// CostType defines what resources a skill costs to use. A skill might cost
// Action Points, Mana, etc, or some combination of costs.
type CostType int

const (
	CostsActionPoints CostType = iota
	CostsMana
	CostsExhaustionPercent
	CostsHealthSacrificePercent
)
