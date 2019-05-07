package game

type Cursor struct{}

// Type of this Component.
func (*Cursor) Type() string {
	return "Cursor"
}
