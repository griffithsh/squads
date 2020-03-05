package skill

type CostType int

const (
	CostsActionPoints CostType = iota
	CostsMana
	CostsExhaustionPercent
	CostsHealthSacrificePercent
)
