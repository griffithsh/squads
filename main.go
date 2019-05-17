package main

import (
	"errors"
	"fmt"
	"image"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/geom"
	"github.com/hajimehoshi/ebiten"
)

type system struct {
	render    *game.Renderer
	combat    *Combat
	mgr       *ecs.World
	camera    *Camera
	lastMouse image.Point
}

func main() {
	// Exercise structure.
	mgr := ecs.NewWorld()
	e := mgr.NewEntity()
	p := &game.Position{}
	mgr.AddComponent(e, p)
	mgr.RemoveComponent(e, p)
	mgr.DestroyEntity(e)

	// dump performance with pprof
	f, err := os.Create("pprof/cpu.prof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	if err := pprof.StartCPUProfile(f); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()

	rand.Seed(time.Now().Unix())
	s, _ := setup(1024, 768)
	ebiten.Run(s.run, 1024, 768, 1, "Squads")
}

type Controls struct {
	Up, Down, Left, Right,
	A, B, C, D,
	Start bool
}

func controls() Controls {
	return Controls{
		Up:    ebiten.IsKeyPressed(ebiten.KeyUp),
		Down:  ebiten.IsKeyPressed(ebiten.KeyDown),
		Left:  ebiten.IsKeyPressed(ebiten.KeyLeft),
		Right: ebiten.IsKeyPressed(ebiten.KeyRight),

		A: ebiten.IsKeyPressed(ebiten.KeyZ),
		B: ebiten.IsKeyPressed(ebiten.KeyX),
		C: ebiten.IsKeyPressed(ebiten.KeyC),
		D: ebiten.IsKeyPressed(ebiten.KeyV),

		Start: ebiten.IsKeyPressed(ebiten.KeyEnter),
	}
}

func controlCamera(c *Camera, t time.Duration, ctrl Controls) {
	camSpeed := 500.0 / c.zoom
	dt := t.Seconds()

	if ctrl.Left {
		c.SetX(c.GetX() - camSpeed*dt)
	} else if ctrl.Right {
		c.SetX(c.GetX() + camSpeed*dt)
	}

	if ctrl.Down {
		c.SetY(c.GetY() + camSpeed*dt)
	} else if ctrl.Up {
		c.SetY(c.GetY() - camSpeed*dt)
	}

	if ctrl.A {
		c.SetZoom(c.GetZoom() * 1.02)
	} else if ctrl.B {
		c.SetZoom(c.GetZoom() * 0.98)
	}
}

func addGrass(mgr *ecs.World, b *geom.Field) {
	M, N := b.Dimensions()
	for n := 0; n < N; n++ {
		for m := 0; m < M; m++ {
			h := b.Get(m, n)
			e := mgr.NewEntity()

			mgr.AddComponent(e, &game.Sprite{
				Texture: "texture.png",
				X:       24,
				Y:       0,
				W:       24,
				H:       16,
			})

			mgr.AddComponent(e, &game.Position{
				Center: game.Center{
					X: h.X(),
					Y: h.Y(),
				},
				Layer: 1,
			})
		}
	}
}

func addTrees(mgr *ecs.World, b *geom.Field) {
	M, N := b.Dimensions()
	for n := 0; n < N; n++ {
		for m := 0; m < M; m++ {
			i := m + n*M
			h := b.Get(m, n)
			if i == 1 || i%17 == 1 || i%13 == 1 {
				e := mgr.NewEntity()
				mgr.AddComponent(e, &game.Sprite{
					Texture: "Untitled.png",
					X:       0,
					Y:       0,
					W:       24,
					H:       48,
				})
				mgr.AddComponent(e, &game.SpriteOffset{
					Y: -16,
				})
				mgr.AddComponent(e, &game.Position{
					Center: game.Center{
						X: h.X(),
						Y: h.Y(),
					},
					Layer: 10,
				})
				mgr.AddComponent(e, &game.Obstacle{
					M: h.M,
					N: h.N,
				})
			}
		}
	}
}

func addHud(mgr *ecs.World) {
	for i := 0; i < 4; i++ {
		e := mgr.NewEntity()
		mgr.AddComponent(e, &game.Sprite{
			Texture: "hud.png",
			X:       0,
			Y:       0,
			W:       16,
			H:       24,
		})
		mgr.AddComponent(e, &game.Position{
			Center: game.Center{
				X: float64(8 + (16+8)*i),
				Y: 12,
			},
			Layer:    100,
			Absolute: true,
		})

	}
}

var last time.Time

// setup the game Entities.
func setup(w, h int) (*system, error) {
	mgr := ecs.NewWorld()
	board, err := geom.NewField(8, 24)
	if err != nil {
		return nil, fmt.Errorf("game.NewBoard: %v", err)
	}
	addGrass(mgr, board)
	addTrees(mgr, board)
	addHud(mgr)

	camera := NewCamera(w, h)
	s := system{
		render: game.NewRenderer(),
		combat: NewCombat(mgr, camera),

		mgr:    mgr,
		camera: camera,
	}
	s.combat.Begin()

	// Create an Actor that is controlled by mouse clicks
	start := board.Get(3, 8)
	actor := mgr.NewEntity()
	mgr.AddComponent(actor, &game.Actor{
		Size: game.SMALL,
	})
	mgr.AddComponent(actor, &game.Facer{Face: geom.S})
	mgr.AddComponent(actor, &game.Sprite{
		Texture: "Untitled.png",
		X:       24,
		Y:       0,
		W:       24,
		H:       48,
	})
	mgr.AddComponent(actor, &game.Position{
		Center: game.Center{
			X: start.X(),
			Y: start.Y(),
		},
		Layer: 10,
	})
	mgr.AddComponent(actor, &game.SpriteOffset{
		Y: -16,
	})

	// FIXME: actor construction should create one or more obstacles to match the Size of the actor.
	// mgr.AddComponent(actor, &game.Obstacle{
	// 	M:            3,
	// 	N:            8,
	// 	ObstacleType: game.ACTOR,
	// })

	last = time.Now()

	return &s, nil
}

var errExitGame = errors.New("game has completed")

var (
	frames      = 0
	accumulated = time.Second * 0
	second      = time.Tick(time.Second)
)

func (s *system) run(screen *ebiten.Image) error {
	start := time.Now()
	defer func() {
		d := time.Since(start)
		frames++
		accumulated += d
	}()

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return errExitGame
	}

	x, y := ebiten.CursorPosition()

	if s.lastMouse.X != x || s.lastMouse.Y != y {
		s.combat.MousePosition(x, y)
		s.lastMouse.X = x
		s.lastMouse.Y = y
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		s.combat.Interaction(x, y)
	}

	elapsed := time.Since(last)
	last = time.Now()

	ctrl := controls()
	controlCamera(s.camera, elapsed, ctrl)

	s.combat.Run(elapsed)

	w, h := float64(screen.Bounds().Max.X-screen.Bounds().Min.X), float64(screen.Bounds().Max.Y-screen.Bounds().Min.Y)

	// Render all entities in the World.
	if err := s.render.Render(screen, s.camera.GetX(), s.camera.GetY(), s.camera.GetZoom(), w, h, s.mgr); err != nil {
		panic(err)
	}

	select {
	case <-second:
		fps := time.Second / (accumulated / time.Duration(frames))
		ebiten.SetWindowTitle(fmt.Sprintf("%s | FPS: %d", "Hexagons, Strategy, Entities, Components, and Systems, Oh my!", fps))
	default:
	}

	return nil
}
