package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/game"
	"golang.org/x/image/colornames"
)

type system struct {
	render        *game.Renderer
	board         *game.Board
	choreographer *game.Choreographer
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
	pixelgl.Run(run)
}

type Controls struct {
	Up, Down, Left, Right,
	A, B, C, D,
	Start bool
}

func controls(win *pixelgl.Window) Controls {
	return Controls{
		Up:    win.Pressed(pixelgl.KeyUp),
		Down:  win.Pressed(pixelgl.KeyDown),
		Left:  win.Pressed(pixelgl.KeyLeft),
		Right: win.Pressed(pixelgl.KeyRight),

		A: win.Pressed(pixelgl.KeyZ),
		B: win.Pressed(pixelgl.KeyX),
		C: win.Pressed(pixelgl.KeyC),
		D: win.Pressed(pixelgl.KeyV),

		Start: win.Pressed(pixelgl.KeyEnter),
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
		c.SetZoom(c.GetZoom() * 1.05)
		hud.SetZoom(c.GetZoom() * 1.05)
	} else if ctrl.B {
		c.SetZoom(c.GetZoom() * 0.95)
		hud.SetZoom(c.GetZoom() * 0.95)
	}
}

func run() {
	mgr := ecs.NewWorld()
	board, err := game.NewBoard(mgr, 8, 24)
	if err != nil {
		panic(err)
	}
	hud := game.NewHUD(1024, 768)
	s := system{
		render:        game.NewRenderer(),
		board:         board,
		choreographer: &game.Choreographer{},
	}

	cfg := pixelgl.WindowConfig{
		Title:  "Hexagons, Strategy, Entities, Components, and Systems, Oh my!",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  false,
	}
	camera := NewCamera(cfg.Bounds.Size().X, cfg.Bounds.Size().Y)
	camera.Center(pixel.Vec{X: s.board.Width() / 2, Y: s.board.Height() / 2})

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	var (
		frames = 0
		second = time.Tick(time.Second)
	)
	last := time.Now()

	cursor := mgr.NewEntity()
	mgr.AddComponent(cursor, &game.Sprite{
		Texture: "texture.png",
		X:       0,
		Y:       0,
		W:       24,
		H:       16,
		Color:   &color.RGBA{150, 150, 150, 63},
	})
	lastMouse := win.MousePosition()

	actor := mgr.NewEntity()
	h := board.Get(0, 0)
	mgr.AddComponent(actor, &game.Actor{})
	mgr.AddComponent(actor, &game.Sprite{
		Texture: "Untitled.png",
		X:       24,
		Y:       0,
		W:       24,
		H:       48,
	})
	mgr.AddComponent(actor, &game.Position{
		Center: game.Center{
			X: h.X(),
			Y: h.Y() - 16,
		},
		Layer: 10,
	})

	for !win.Closed() {
		if win.JustReleased(pixelgl.KeyEscape) || win.Pressed(pixelgl.KeyEscape) {
			win.SetClosed(true)
			return
		}

		if win.JustPressed(pixelgl.MouseButtonLeft) {
			p := win.MousePosition()
			p = camera.View().Unproject(p)

			// faiface/pixel inverts the Y coordinate
			p.Y = -p.Y

			a := mgr.Component(actor, "Actor").(*game.Actor)
			oe := mgr.Get([]string{"Obstacle"})
			var obstacles []game.ContextualObstacle
			for _, e := range oe {
				obstacle := mgr.Component(e, "Obstacle").(*game.Obstacle)
				hex := s.board.Get(obstacle.M, obstacle.N)
				if hex == nil {
					continue
				}
				obstacles = append(obstacles, game.ContextualObstacle{
					Obstacle: *obstacle,
					Cost:     math.Inf(0), // just pretend these all are total obstacles for now
				})
			}

			steps, err := game.Navigate(board.Get(a.M, a.N), board.At(int(p.X), int(p.Y)), obstacles)
			if err != nil {
				fmt.Printf("no path there: %v\n", err)
			} else {
				a.Move(steps)
			}
		}

		// A Cursor that follows the mouse...
		p := win.MousePosition()
		if p != lastMouse {
			c, ok := mgr.Component(cursor, "Position").(*game.Position)
			if ok {
				mgr.RemoveComponent(cursor, c)
			}
			p = camera.View().Unproject(p)
			p.Y = -p.Y
			if h := s.board.At(int(p.X), int(p.Y)); h != nil {
				mgr.AddComponent(cursor, &game.Position{
					Center: game.Center{
						X: h.X(),
						Y: h.Y(),
					},
					Layer: 2,
				})
			}
			lastMouse = win.MousePosition()
		}

		elapsed := time.Since(last)
		last = time.Now()

		ctrl := controls(win)
		controlCamera(camera, hud, elapsed, ctrl)

		s.choreographer.Update(mgr, elapsed)

		win.Clear(colornames.Cadetblue)

		// Render all entities.
		if err := s.render.Render(win, camera.View(), mgr); err != nil {
			panic(err)
		}

		// Render a hud.
		hud.Render(win)

		win.Update()
		frames++
		select {
		case <-second:
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
		}
	}
}
