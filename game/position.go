package game

type Position struct {
	X, Y  float64
	Layer int
}

// Type of this component.
func (p *Position) Type() string {
	return "Position"
}
