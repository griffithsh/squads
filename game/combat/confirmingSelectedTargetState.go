package combat

import (
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/geom"
)

type confirmingSelectedTargetState struct {
	Skill  game.SkillCode
	Target geom.Key
}

// Value satisfies the StateContext interface, and can always return
// ConfirmingSelectedTargetState.
func (confirmingSelectedTargetState) Value() State {
	return ConfirmingSelectedTargetState
}
