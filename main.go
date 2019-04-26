package main

import (
	"errors"
	"fmt"
	"image"
	"log"
	"math"
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
	board     *geom.Field
	nav       *game.Navigator
	mgr       *ecs.World
	camera    *Camera
	hud       *game.HUD
	cursor    ecs.Entity
	lastMouse image.Point
	actor     ecs.Entity
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

func controlCamera(c *Camera, hud *game.HUD, t time.Duration, ctrl Controls) {
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
			if i == 1 || i%11 == 1 || i%17 == 1 || i%13 == 1 {
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

	hud := game.NewHUD(w, h)
	s := system{
		render: game.NewRenderer(),
		board:  board,
		nav:    &game.Navigator{},
		mgr:    mgr,
		hud:    hud,
	}

	s.camera = NewCamera(float64(w), float64(h))
	s.camera.Center(Vec{X: s.board.Width() / 2, Y: s.board.Height() / 2})

	s.cursor = mgr.NewEntity()
	mgr.AddComponent(s.cursor, &game.Sprite{
		Texture: "texture.png",
		X:       0,
		Y:       0,
		W:       24,
		H:       16,
	})

	// Create an Actor that is controlled by mouse clicks
	start := board.Get(0, 0)
	s.actor = mgr.NewEntity()
	mgr.AddComponent(s.actor, &game.Actor{})
	mgr.AddComponent(s.actor, &game.Facer{Face: geom.S})
	mgr.AddComponent(s.actor, &game.Sprite{
		Texture: "Untitled.png",
		X:       24,
		Y:       0,
		W:       24,
		H:       48,
	})
	mgr.AddComponent(s.actor, &game.Position{
		Center: game.Center{
			X: start.X(),
			Y: start.Y(),
		},
		Layer: 10,
	})
	mgr.AddComponent(s.actor, &game.SpriteOffset{
		Y: -16,
	})
	mgr.AddComponent(s.actor, &game.Obstacle{
		M:            0,
		N:            0,
		ObstacleType: game.ACTOR,
	})

	last = time.Now()

	return &s, nil
}

var errExitGame = errors.New("game has completed")

var (
	frames                    = 0
	accumulated time.Duration = time.Second * 0
	second                    = time.Tick(time.Second)
)

func (s *system) run(screen *ebiten.Image) error {
	start := time.Now()
	defer func() {
		d := time.Since(start)
		// fmt.Println(d)
		frames++
		accumulated += d
	}()

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return errExitGame
	}

	// Get x and y as screen coordinates.
	x, y := ebiten.CursorPosition()

	// Correct for current zoom.
	zoom := s.camera.GetZoom()
	x, y = int(float64(x)/zoom), int(float64(y)/zoom)

	// Correct for camera focus.
	x, y = x+int(s.camera.GetX()), y+int(s.camera.GetY())

	// Correct for size of screen (!?).
	x, y = int(float64(x)-1024/2/zoom), int(float64(y)-768/2/zoom)

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		a := s.mgr.Component(s.actor, "Actor").(*game.Actor)

		var obstacles []geom.ContextualObstacle
		for _, e := range s.mgr.Get([]string{"Obstacle"}) {
			obstacle := s.mgr.Component(e, "Obstacle").(*game.Obstacle)
			hex := s.board.Get(obstacle.M, obstacle.N)
			if hex == nil {
				continue
			}

			// Translate the Obstacles into ContextualObstacles based on
			// how much of an Obstacle this is to the Mover in this context.
			obstacles = append(obstacles, geom.ContextualObstacle{
				M:    obstacle.M,
				N:    obstacle.N,
				Cost: math.Inf(0), // just pretend these all are total obstacles for now
			})
		}

		steps, err := geom.Navigate(s.board.Get(a.M, a.N), s.board.At(x, y), obstacles)
		if err != nil {
			fmt.Printf("no path there: %v\n", err)
		} else {
			s.mgr.AddComponent(s.actor, &game.Mover{
				Moves: steps,
			})
		}
	}

	// A Cursor that follows the mouse...
	if s.lastMouse.X != x || s.lastMouse.Y != y {
		c, ok := s.mgr.Component(s.cursor, "Position").(*game.Position)
		if ok {
			s.mgr.RemoveComponent(s.cursor, c)
		}

		if h := s.board.At(x, y); h != nil {
			s.mgr.AddComponent(s.cursor, &game.Position{
				Center: game.Center{
					X: h.X(),
					Y: h.Y(),
				},
				Layer: 2,
			})
		}
		s.lastMouse.X, s.lastMouse.Y = ebiten.CursorPosition()
	}

	elapsed := time.Since(last)
	last = time.Now()

	ctrl := controls()
	controlCamera(s.camera, s.hud, elapsed, ctrl)

	s.nav.Update(s.mgr, elapsed)

	// Render all entities.
	if err := s.render.Render(screen, s.camera.GetX(), s.camera.GetY(), s.camera.GetZoom(), 1024, 768, s.mgr); err != nil {
		panic(err)
	}

	// Render a hud separately because the hud is not composed of ECS Entities.
	s.hud.Render(screen, s.camera.GetZoom(), 1024, 768)

	select {
	case <-second:
		fps := time.Second / (accumulated / time.Duration(frames))
		ebiten.SetWindowTitle(fmt.Sprintf("%s | FPS: %d", "Hexagons, Strategy, Entities, Components, and Systems, Oh my!", fps))
	default:
	}

	return nil
}
