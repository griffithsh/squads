package overworld

import (
	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/geom"
)

// Token represents a squad on the overworld map.
type Token struct {
	Key   geom.Key
	Squad ecs.Entity
}

// Type of this Component.
func (Token) Type() string {
	return "Token"
}
