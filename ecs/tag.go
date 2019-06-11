package ecs

// Tags is a Component that represents arbitrary strings that are related to an
// Entity.
type Tags []string

// Type of this Component.
func (*Tags) Type() string {
	return "Tags"
}
