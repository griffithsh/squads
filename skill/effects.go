package skill

import (
	"time"

	"github.com/griffithsh/squads/game"
)

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

func (de DamageEffect) Schedule() time.Duration {
	return de.ScheduleTime
}
