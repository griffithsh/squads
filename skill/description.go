package skill

import (
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/targeting"
)

// Description contains all the data for a skill so that it can be displayed in
// menus or utilised in combat.
type Description struct {
	ID          ID
	Name        string
	Explanation string

	// Tags critically includes Attack or Spell, and allows the game to select
	// an appropriate animation to use when using the skill.
	Tags []Classification

	Icon game.FrameAnimation

	Targeting targeting.Rule

	// Effects of triggering this skill.
	Effects []Effect

	Costs map[CostType]int
}

// IsAttack returns whether a skill is an attack or not.
func (d Description) IsAttack() bool {
	for _, tag := range d.Tags {
		if tag == Attack {
			return true
		}
	}
	return false
}
