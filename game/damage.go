package game

//go:generate go run github.com/dmarkham/enumer -type=DamageType -json
type DamageType int

const (
	PhysicalDamage DamageType = iota
	MagicalDamage
	FireDamage
)
