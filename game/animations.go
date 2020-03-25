package game

import (
	"math"
	"math/rand"
	"time"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/res"
)

// AnimationEndBehavior is an enum that describes what happens when a
// FrameAnimation Component finishes all of its frames.
type AnimationEndBehavior int

const (
	// AnimationLoops means that the Animation will restart from the first frame
	// after the last frame expires.
	AnimationLoops AnimationEndBehavior = iota
	// HoldLastFrame means that the Animation will run through all frames, and
	// then stick on the last frame, removing the FrameAnimation Component.
	HoldLastFrame
	// DestroyEntity means that after the Animation is complete, the Entity that
	// owns it is destroyed.
	DestroyEntity
)

// FrameAnimation sets the Sprite of the Entity based on how far through the
// whole Animation is indicated by Pointer.
type FrameAnimation struct {
	Frames  []Sprite
	Timings []time.Duration
	Pointer time.Duration

	// by default, loops forever
	// on end - maybe another animation, maybe a single frame, maybe some behaviour like trigger event
	EndBehavior AnimationEndBehavior
}

// NewFrameAnimation creates a new FrameAnimation Component from a res.Animation.
func NewFrameAnimation(a res.Animation) FrameAnimation {
	fa := FrameAnimation{}
	for _, frame := range a.Frames {
		fa.Frames = append(fa.Frames, Sprite{
			Texture: frame.Texture,
			X:       frame.X,
			Y:       frame.Y,
			W:       frame.W,
			H:       frame.H,
			OffsetX: frame.OffsetX,
			OffsetY: frame.OffsetY,
		})
		fa.Timings = append(fa.Timings, frame.Duration)
	}
	return fa
}

// Duration of the entire Animation.
func (fa *FrameAnimation) Duration() time.Duration {
	var result time.Duration
	for _, t := range fa.Timings {
		result += t
	}
	return result
}

// Index returns the index of the Sprite that Pointer is currently pointing at.
// It returns -1 if Pointer is negative or the index of one past the last
// element when Pointer is outside the range of Timings.
func (fa *FrameAnimation) Index() int {
	if fa.Pointer < 0 {
		return -1
	}

	var accumulated time.Duration
	for i, t := range fa.Timings {
		accumulated += t
		if accumulated > fa.Pointer {
			return i
		}
	}

	return len(fa.Frames)
}

// Randomise the starting position of the Pointer.
func (fa *FrameAnimation) Randomise() *FrameAnimation {
	if len(fa.Frames) <= 1 {
		fa.Pointer = 0
		return fa
	}
	fa.Pointer = time.Duration(rand.Int63n(int64(fa.Duration())))
	return fa
}

// Type of this Component.
func (*FrameAnimation) Type() string {
	return "FrameAnimation"
}

// HoverAnimation makes the entity float up and down.
type HoverAnimation struct {
	YVelocity  float64
	YTranslate float64
	Force      float64
}

// NewHoverAnimation creates a new HoverAnimation.
func NewHoverAnimation() *HoverAnimation {
	return &HoverAnimation{
		Force:      92.5,
		YTranslate: -5,
		YVelocity:  0.0,
	}
}

// Type of this Component.
func (*HoverAnimation) Type() string {
	return "HoverAnimation"
}

// FloatAwayAnimation will offset the Entity upwards based on Rate as if it is a
// helium balloon floating away.
type FloatAwayAnimation struct {
	// Rate of floating away per second
	Rate float64
	age  time.Duration
}

// Type of this Component.
func (*FloatAwayAnimation) Type() string {
	return "FloatAwayAnimation"
}

// TakeDamageAnimation adds an offset that makes the Entity pendulum back and
// forth on the X axis as if they have been struck.
type TakeDamageAnimation struct {
	age time.Duration
}

// Type of this Component.
func (*TakeDamageAnimation) Type() string {
	return "TakeDamageAnimation"
}

// AnimationSpeed changes the rate at which Animtions are animated. A value of
// 1.0 is normal speed, and a value of 0.0 would mean that the animation never
// progresses from the first frame. 0.5 would animate at half speed.
type AnimationSpeed struct {
	Speed float64
}

// Type of this Component.
func (*AnimationSpeed) Type() string {
	return "AnimationSpeed"
}

// getAnimationElapsed is a way of slowing an animation based on its AnimationSpeed.
func getAnimationElapsed(mgr *ecs.World, e ecs.Entity, elapsed time.Duration) time.Duration {
	c := mgr.Component(e, "AnimationSpeed")
	if c == nil {
		return elapsed
	}
	speed := c.(*AnimationSpeed)

	// FIXME: There has to be a better way of multiplying elapsed by speed?
	return time.Nanosecond * time.Duration(float64(elapsed.Nanoseconds())*speed.Speed)
}

// AnimationSystem animates the visual Components of Entities. It's not
// responsible for translating or mapping game concepts like "casting a spell"
// to the assignment of appropriate animation Components for that Entity.
type AnimationSystem struct{}

// Update all Animated Entities.
func (as *AnimationSystem) Update(mgr *ecs.World, elapsed time.Duration) {
	for _, e := range mgr.Get([]string{"FrameAnimation"}) {
		anim := mgr.Component(e, "FrameAnimation").(*FrameAnimation)
		elapsed := getAnimationElapsed(mgr, e, elapsed)
		anim.Pointer += elapsed

		i := anim.Index()
		complete := i >= len(anim.Frames)

		// If the animation is not complete, assign the current frame as the
		// Entity's Sprite, and continue to the next Entity.
		if !complete {
			mgr.AddComponent(e, &anim.Frames[i])
			continue
		}

		// If the Animation _is_ complete, examine the FrameAnimation's
		// EndBehavior to figure how to handle it.
		switch anim.EndBehavior {
		case AnimationLoops:
			if anim.Duration() == 0 {
				i = 0
			} else {
				anim.Pointer = anim.Pointer % anim.Duration()
				i = anim.Index()
			}

		case HoldLastFrame:
			mgr.AddComponent(e, &anim.Frames[len(anim.Frames)-1])
			mgr.RemoveComponent(e, &FrameAnimation{})
		case DestroyEntity:
			mgr.DestroyEntity(e)
		}
	}

	for _, e := range mgr.Get([]string{"HoverAnimation"}) {
		elapsed := getAnimationElapsed(mgr, e, elapsed)
		anim := mgr.Component(e, "HoverAnimation").(*HoverAnimation)

		if anim.YTranslate > 0 {
			anim.YVelocity -= anim.Force * elapsed.Seconds()
		} else {
			anim.YVelocity += anim.Force * elapsed.Seconds()
		}

		// Apply velocity to offset.
		anim.YTranslate += anim.YVelocity * elapsed.Seconds()

		// Save offset.
		mgr.AddComponent(e, &RenderOffset{
			Y: int(anim.YTranslate),
		})
	}

	for _, e := range mgr.Get([]string{"FloatAwayAnimation"}) {
		faa := mgr.Component(e, "FloatAwayAnimation").(*FloatAwayAnimation)
		elapsed := getAnimationElapsed(mgr, e, elapsed)

		faa.age += elapsed

		mgr.AddComponent(e, &RenderOffset{
			Y: -int(faa.age.Seconds() * faa.Rate),
		})
	}

	for _, e := range mgr.Get([]string{"TakeDamageAnimation"}) {
		tda := mgr.Component(e, "TakeDamageAnimation").(*TakeDamageAnimation)
		elapsed := getAnimationElapsed(mgr, e, elapsed)

		offset, ok := mgr.Component(e, "RenderOffset").(*RenderOffset)
		if !ok {
			offset = &RenderOffset{}
			mgr.AddComponent(e, offset)
		}

		// maxAge in seconds of this animation.
		maxAge := 0.65
		tda.age += elapsed
		if tda.age.Seconds() >= maxAge {
			mgr.RemoveComponent(e, tda)
			offset.X = 0
			continue
		}

		// spd is how fast the wobble happens.
		spd := 22.5
		sec := tda.age.Seconds()
		raw := math.Sin(sec * spd)
		through := raw * (1 - sec/maxAge)

		// amp amplifies the wobble.
		amp := 7.5
		offset.X = int(through * amp)
	}
}
