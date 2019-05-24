package game

// Scale defines a multiplier that is applied to the Entity's size.
type Scale struct {
	X, Y float64
}

// Type of this Component.
func (*Scale) Type() string {
	return "Scale"
}
