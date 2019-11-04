package overworld

import "github.com/griffithsh/squads/geom"

// Token represents a squad on the overworld map.
type Token struct {
	Key geom.Key
}

// Type of this Component.
func (Token) Type() string {
	return "Token"
}
