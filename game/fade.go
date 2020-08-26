package game

import (
	"time"

	"github.com/griffithsh/squads/ecs"
)

// FadeIn is a Component that causes an Entity to appear over time.
type FadeIn struct {
	init     bool
	age      time.Duration
	Duration time.Duration
}

// Type of this Component.
func (*FadeIn) Type() string {
	return "FadeIn"
}

// FadeOut is a Component that causes an Entity to disappear from view over
// time.
type FadeOut struct {
	init     bool
	age      time.Duration
	Duration time.Duration
}

// Type of this Component.
func (*FadeOut) Type() string {
	return "FadeOut"
}

// FadeSystem controls the Alpha Component of Entities by managing FadeOut and
// FadeIn Components.
type FadeSystem struct{}

// Update the system.
func (sys *FadeSystem) Update(mgr *ecs.World, elapsed time.Duration) {
	for _, e := range mgr.Get([]string{"FadeIn"}) {
		in := mgr.Component(e, "FadeIn").(*FadeIn)

		if !in.init {
			mgr.AddComponent(e, &Alpha{Value: 0})
			in.init = true
			continue
		}

		// TODO what happens when we have both a fadeout and a fadein?

		in.age += elapsed
		if in.age > in.Duration {
			mgr.RemoveComponent(e, in)
			mgr.RemoveComponent(e, &Alpha{})
			continue
		}

		mgr.AddComponent(e, &Alpha{Value: float64(in.age) / float64(in.Duration)})
	}

	for _, e := range mgr.Get([]string{"FadeOut"}) {
		out := mgr.Component(e, "FadeOut").(*FadeOut)
		if !out.init {
			mgr.AddComponent(e, &Alpha{Value: 1.0})
			out.init = true
			continue
		}

		// TODO what happens when we have both a fadeout and a fadein?

		out.age += elapsed
		if out.age > out.Duration {
			mgr.RemoveComponent(e, out)
			mgr.AddComponent(e, &Alpha{Value: 0.0})
			continue
		}

		mgr.AddComponent(e, &Alpha{Value: 1.0 - float64(out.age)/float64(out.Duration)})
	}
}
