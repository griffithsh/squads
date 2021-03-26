package embark

// Embarking is a Component that stores whether a Character is embarking to go
// on an adventure, or will stay home in the village.
type Embarking struct {
	Value bool
}

// Type of this Component.
func (*Embarking) Type() string {
	return "Embarking"
}
