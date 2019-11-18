package game

// Squad is a Component that marks a Squad of Characters
type Squad struct{}

// Type of this Component.
func (*Squad) Type() string {
	return "Squad"
}
