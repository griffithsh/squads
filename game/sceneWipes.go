package game

import (
	"math"
	"math/rand"
	"time"

	"github.com/griffithsh/squads/ecs"
)

type row struct {
	block, anim1, anim2 ecs.Entity
}

// DiagonalMatrixWipe obscures the scene by animating in squares of curtain
// in a diagonal pattern.
type DiagonalMatrixWipe struct {
	W, H       int
	Obscuring  bool
	OnComplete func()

	age time.Duration
}

// Type of this Component.
func (*DiagonalMatrixWipe) Type() string {
	return "DiagonalMatrixWipe"
}

// SceneWipeSystem handles scene wipe Components that hide abrupt changes to the
// scene by fading in and out. Think of these as curtains that can be drawn to
// hide stage hands moving pieces of the set on and off the stage.
type SceneWipeSystem struct{}

// NewSceneWipeSystem constructs a SceneWipeSystem.
func NewSceneWipeSystem() *SceneWipeSystem {
	return &SceneWipeSystem{}
}

// Update all scene wipe Components in the passed ecs.World.
func (cs *SceneWipeSystem) Update(mgr *ecs.World, elapsed time.Duration) {
	// complete is the amount of time that this wipe takes to obscure or reveal the screen.
	const complete = time.Millisecond * 1100

	for _, e := range mgr.Get([]string{"DiagonalMatrixWipe"}) {
		dmw := mgr.Component(e, "DiagonalMatrixWipe").(*DiagonalMatrixWipe)
		numRows := (dmw.H + 63) / 64
		numCols := (dmw.W + 47) / 48
		// If the age is zero then we need to initialise the Component.
		if dmw.age == 0 {
			children := make([]ecs.Entity, 0, numRows*3)
			for i := 0; i < numRows; i++ {
				block := mgr.NewEntity()
				mgr.AddComponent(block, &Sprite{
					Texture: "diagonal-matrix-wipe.png",
					X:       0,
					Y:       0,
					W:       16,
					H:       16,
				})
				mgr.AddComponent(block, &Position{
					Center: Center{
						X: 0,
						Y: 32 + float64(i*64),
					},
					Layer:    1000,
					Absolute: true,
				})
				mgr.AddComponent(block, &Scale{
					X: 3.0, // * 16 = 48
					Y: 4.0, // * 16 = 64
				})
				if !dmw.Obscuring {
					mgr.AddComponent(block, &Scale{
						X: 3.0 + float64(numCols), // * 16 = 48
						Y: 4.0 + float64(numRows), // * 16 = 64
					})

				}

				anim1 := mgr.NewEntity()
				mgr.AddComponent(anim1, &Sprite{
					Texture: "diagonal-matrix-wipe.png",
					X:       192,
					Y:       0,
					W:       48,
					H:       64,
				})
				mgr.AddComponent(anim1, &Position{
					Center: Center{
						X: 0,
						Y: 32 + float64(i*64),
					},
					Layer:    1000,
					Absolute: true,
				})

				anim2 := mgr.NewEntity()
				mgr.AddComponent(anim2, &Sprite{
					Texture: "diagonal-matrix-wipe.png",
					X:       48,
					Y:       0,
					W:       48,
					H:       64,
				})
				mgr.AddComponent(anim2, &Position{
					Center: Center{
						X: 0,
						Y: 32 + float64(i*64),
					},
					Layer:    1000,
					Absolute: true,
				})
				children = append(children, block, anim1, anim2)
			}
			mgr.AddComponent(e, &ecs.Children{Value: children})
		}

		percentage := float64(dmw.age) / float64(complete)
		if !dmw.Obscuring {
			percentage = 1.0 - percentage
		}

		children := mgr.Component(e, "Children").(*ecs.Children)
		for i, ce := range children.Value {
			switch i % 3 {
			case 0:
				s := mgr.Component(ce, "Scale").(*Scale)
				s.X = 3.0 * math.Round(float64(numCols+numRows)*percentage)
				p := mgr.Component(ce, "Position").(*Position)
				p.Center.X = math.Round(float64(numCols+numRows)*percentage) * 24
				p.Center.X -= float64(i / 3 * 48)

			case 1:
				p := mgr.Component(ce, "Position").(*Position)
				newX := math.Round(float64(numCols+numRows)*percentage) * 48
				newX -= float64(i / 3 * 48)
				newX += 24
				if newX != p.Center.X {
					p.Center.X = newX
					spr := mgr.Component(ce, "Sprite").(*Sprite)
					spr.Y = rand.Intn(3) * 64
					spr.X += 48
					if spr.X == 48*5 {
						spr.X = 96
					}
				}

			case 2:
				// Second animating
				p := mgr.Component(ce, "Position").(*Position)
				newX := math.Round(float64(numCols+numRows)*percentage) * 48
				newX -= float64(i / 3 * 48)
				newX += 24 + 48
				if newX != p.Center.X {
					p.Center.X = newX
					spr := mgr.Component(ce, "Sprite").(*Sprite)
					spr.Y = rand.Intn(3) * 64
					spr.X += 48
					if spr.X == 48*4 {
						spr.X = 0
					}
				}

			}
		}

		dmw.age += elapsed

		// If the age of this wipe has passed the total duration, then any
		// configured OnComplete handler should be called, and the scene wipe
		// then destroyed.
		if dmw.age > complete {
			if dmw.OnComplete != nil {
				dmw.OnComplete()
			}
			for _, ce := range children.Value {
				mgr.DestroyEntity(ce)
			}
			mgr.DestroyEntity(e)
		}
	}
}
