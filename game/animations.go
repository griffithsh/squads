package game

import (
	"math/rand"
	"time"

	"github.com/griffithsh/squads/ecs"
)

// FrameAnimation sets the Sprite of the Entity based what how far through the
// total Animation Pointer indicates.
type FrameAnimation struct {
	Frames  []Sprite
	Timings []time.Duration
	Pointer time.Duration
	// by default, loops forever
	// on end - maybe another animation, maybe a single frame, maybe some behaviour like trigger event
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

// AnimationSystem animates the visual Components of Entities. It's not
// responsible for translating or mapping game concepts like "casting a spell"
// to the assignment of appropriate animation Components for that Entity.
type AnimationSystem struct{}

// Update all Animated Entities.
func (as *AnimationSystem) Update(mgr *ecs.World, elapsed time.Duration) {
	for _, e := range mgr.Get([]string{"FrameAnimation", "Sprite"}) {
		anim := mgr.Component(e, "FrameAnimation").(*FrameAnimation)

		anim.Pointer += elapsed

		i := anim.Index()
		if i >= len(anim.Frames) {
			anim.Pointer = anim.Pointer % anim.Duration()
			i = anim.Index()
		}

		mgr.AddComponent(e, &anim.Frames[i])
	}

	for _, e := range mgr.Get([]string{"HoverAnimation"}) {
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
}
