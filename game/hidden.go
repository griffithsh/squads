package game

// Hidden is a Component that prevents the Entity from being rendered.
type Hidden struct{}

// Type of this Component.
func (f *Hidden) Type() string {
	return "Hidden"
}
