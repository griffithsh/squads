package game

// Alpha is a component that controls how visible an Entity should be rendered.
type Alpha struct {
	Value float64 // Value should be between 0.0 and 1.0
}

// Type of this Component.
func (*Alpha) Type() string {
	return "Alpha"
}
