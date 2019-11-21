package game

import "github.com/griffithsh/squads/ecs"

// Squad is a Component that marks a Squad of Characters
type Squad struct {
	Members []ecs.Entity
}

// Type of this Component.
func (*Squad) Type() string {
	return "Squad"
}
