package skill

import (
	"github.com/griffithsh/squads/game"
)

// Effect is anything that executing a skill could trigger.
type Effect struct {
	When Timing
	What interface{}
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

// DamageEffect deals damage.
type DamageEffect struct {
	Min            Operations     // "0.5 * attunement * 0.22 * INT"
	Max            Operations     // "10 + 5 * attunement * 0.17 * INT"
	Classification Classification // Spell or Attack: can it be negated or dodged?
	DamageType     game.DamageType
}

// HealEffect heals the target. By default the amount is treated as a value to
// increment the current health by, but if IsPercentage is true, then the
// current health is incremented by the maximum health multiplied by Amount.
type HealEffect struct {
	Amount       float64
	IsPercentage bool
}

// ReviveEffect changes the status from KnockedDown to Alive, and sets the
// target's current health to 1. Does nothing if the target is not KnockedDown.
type ReviveEffect struct{}

// DefileEffect changes the target's state from KnockedDown to Defiled.
type DefileEffect struct{}

// SpawnParticipantEffect spawns a new participant.
type SpawnParticipantEffect struct {
	Profession string
	Level      Operations
}

//go:generate stringer -type=InjuryType

// InjuryType enumerates injuries.
type InjuryType int

const (
	BleedingInjury InjuryType = iota
)

func InjuryTypeFromString(s string) *InjuryType {
	for i := 0; i <= int(BleedingInjury); i++ {
		t := InjuryType(i)

		if t.String() == s {
			return &t
		}
	}
	return nil
}

// InjuryEffect applies an injury to the target.
type InjuryEffect struct {
	Type  InjuryType
	Value int
}
