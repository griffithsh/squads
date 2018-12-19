package game

type Center struct {
	X, Y float64
}
type Position struct {
	Center Center
	Layer  int
}

// Type of this component.
func (p *Position) Type() string {
	return "Position"
}
