package skill

import (
	"time"

	"github.com/griffithsh/squads/game"
)

// Effect is anything that executing a skill could trigger.
type Effect interface {
	// Schedule is the time that this effect should be triggered in the
	// skill execution's lifetime.
	Schedule() time.Duration
}

type Operator int

const (
	AddOp Operator = iota
	MultOp
)

type Operation struct {
	Operator Operator
	Variable string // literal float, or substitute value like STR, INT, fire attunement
}

type Operations []Operation

func (o Operations) Calculate(dereferencer func(string) float64) int {
	var working float64
	for _, step := range o {
		val := dereferencer(step.Variable)
		switch step.Operator {
		case MultOp:
			working = working * val
		case AddOp:
			working = working + val
		}
	}
	return int(working)
}

type DamageEffect struct {
	ScheduleTime   time.Duration
	Min            Operations     // "0.5 * attunement * 0.22 * INT"
	Max            Operations     // "10 + 5 * attunement * 0.17 * INT"
	Classification Classification // Spell or Attack: can it be negated or dodged?
	DamageType     game.DamageType
}

// Schedule this effect.
func (de DamageEffect) Schedule() time.Duration {
	return de.ScheduleTime
}

// HealEffect heals the target. By default the amount is treated as a value to
// increment the current health by, but if IsPercentage is true, then the
// current health is incremented by the maximum health multiplied by Amount.
type HealEffect struct {
	ScheduleTime time.Duration
	Amount       float64
	IsPercentage bool
}

// Schedule this effect.
func (he HealEffect) Schedule() time.Duration {
	return he.ScheduleTime
}

// ReviveEffect changes the status from KnockedDown to Alive, and sets the
// target's current health to 1. Does nothing if the target is not KnockedDown.
type ReviveEffect struct {
	ScheduleTime time.Duration
}

// Schedule this effect.
func (re ReviveEffect) Schedule() time.Duration {
	return re.ScheduleTime
}
