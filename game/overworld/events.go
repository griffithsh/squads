package overworld

import (
	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/geom"
)

type TokenMoved struct {
	E    ecs.Entity
	From geom.Key
	To   geom.Key
}

// Type of the Event.
func (TokenMoved) Type() event.Type {
	return "overworld.TokenMoved"
}

type TokensCollided struct {
	E1, E2 ecs.Entity
	At     geom.Key
}

// Type of the Event.
func (TokensCollided) Type() event.Type {
	return "overworld.TokensCollided"
}
