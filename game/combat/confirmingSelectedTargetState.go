package combat

import (
	"github.com/griffithsh/squads/geom"
	"github.com/griffithsh/squads/skill"
)

type confirmingSelectedTargetState struct {
	Skill  skill.ID
	Target geom.Key
}

// Value satisfies the StateContext interface, and can always return
// ConfirmingSelectedTargetState.
func (confirmingSelectedTargetState) Value() State {
	return ConfirmingSelectedTargetState
}
