package skill

//go:generate stringer -type=CostType

// CostType defines what resources a skill costs to use. A skill might cost
// Action Points, Mana, etc, or some combination of costs.
type CostType int

const (
	CostsActionPoints CostType = iota
	CostsMana
	CostsExhaustionPercent
	CostsHealthSacrificePercent
)

func CostTypeFromString(s string) *CostType {
	for i := 0; i <= int(CostsHealthSacrificePercent); i++ {
		c := CostType(i)

		if c.String() == s {
			return &c
		}
	}
	return nil
}
