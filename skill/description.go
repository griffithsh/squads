package skill

import (
	"time"

	"github.com/griffithsh/squads/game"
)

// Description contains all the data for a skill so that it can be displayed in
// menus or utilised in combat.
type Description struct {
	ID          ID
	Name        string
	Explanation string
	Tags        []Classification

	Icon game.FrameAnimation

	Targeting      TargetingRule
	TargetingBrush TargetingBrush

	// Effects of triggering this skill.
	Effects []Effect

	Costs map[CostType]int
}

// Effect is anything that executing a skill could trigger.
type Effect interface {
	// Schedule is the time that this effect should be triggered in the
	// skill execution's lifetime.
	Schedule() time.Duration
}
