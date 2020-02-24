package combat

import "github.com/griffithsh/squads/skill"

// selectingTargetState implements the StateContext interface because it
// contains context around "what" is being selected.
type selectingTargetState struct {
	Skill skill.ID
}

// Value satisfies the StateContext interface, and can always return
// SelectingTargetState.
func (selectingTargetState) Value() State {
	return SelectingTargetState
}
