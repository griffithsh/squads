package skill

import (
	"fmt"
	"time"
)

// Timing represents when a skill effect should be triggered. It allows both
// concrete time.Duration values that specify the amount of time after skill
// execution starts, and virtualised "points" in the skill execution like "end"
// or "when the blow lands".
type Timing struct {
	point *TimingPoint
	sched time.Duration
}

type TimingPoint int

const (
	AttackApexTimingPoint TimingPoint = iota
	EndTimingPoint
)

// NewTiming constructs a Timing for an effect that should be triggered the
// specified Duration after the skill execution starts.
func NewTiming(at time.Duration) Timing {
	return Timing{
		sched: at,
	}
}

// NewTimingFromPoint constructs a Timing from a TimingPoint.
func NewTimingFromPoint(point TimingPoint) Timing {
	return Timing{
		point: &point,
	}
}

// when to trigger an effect. Returns either a virtual TimingPoint, or a
// concrete time as a time.Duration. Use a realiser provided by
// NewTimingRealiser to resolve either case to a concrete time.Duration.
func (t Timing) when() interface{} {
	if t.point != nil {
		return *t.point
	}
	return t.sched
}

// NewTimingRealiser creates a function to realise when to trigger an effect.
func NewTimingRealiser(m map[TimingPoint]time.Duration) func(Timing) time.Duration {
	return func(t Timing) time.Duration {
		switch v := t.when().(type) {
		case TimingPoint:
			return m[v]
		case time.Duration:
			return v
		default:
			panic(fmt.Sprintf("misconfigured skill.Timing or code error, %T", v))
		}
	}
}
