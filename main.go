package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"

	"github.com/griffithsh/squads/game"

	_ "image/png"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/griffithsh/squads/ecs"
	"golang.org/x/image/colornames"
)

type system struct {
	render *game.Renderer
	board  *game.Board
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
		c.SetZoom(c.GetZoom() * 1.05)
	} else if ctrl.B {
		c.SetZoom(c.GetZoom() * 0.95)
	}
}

func run() {
	mgr := ecs.NewWorld()
	board, err := game.NewBoard(mgr, 8, 24)
	if err != nil {
		panic(err)
	}
	s := system{
		render: game.NewRenderer(),
		board:  board,
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
		X:       24,
		Y:       0,
		W:       24,
		H:       16,
	})
	lastMouse := win.MousePosition()
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

			// Add a wolf!
			e := mgr.NewEntity()
			mgr.AddComponent(e, &game.Sprite{
				Texture: "Untitled.png",
				X:       97,
				Y:       171,
				W:       34,
				H:       29,
			})
			mgr.AddComponent(e, &game.Position{
				Center: game.Center{
					X: p.X,
					Y: p.Y,
				},
				Layer: 10,
			})
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
					Layer: 10,
				})
			}
			lastMouse = win.MousePosition()
		}

		elapsed := time.Since(last)
		last = time.Now()

		ctrl := controls(win)
		controlCamera(camera, elapsed, ctrl)

		// Rendering
		win.Clear(colornames.Cadetblue)

		// Render all entities
		if err := s.render.Render(win, camera.View(), mgr); err != nil {
			panic(err)
		}

		// TODO: render a hud
		// ...

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
