package ui

// Interactive tags Components that can be interacted with in some way.
type Interactive struct {
	Trigger func()
	Hover   func()
}

// Type of this Component.
func (*Interactive) Type() string {
	return "Interactive"
}
