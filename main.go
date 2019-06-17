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
	"github.com/griffithsh/squads/game/combat"
	"github.com/hajimehoshi/ebiten"
)

type system struct {
	render       *game.Renderer
	fonts        *game.FontSystem
	hierarchy    *ecs.ParentSystem
	combat       *combat.Manager
	mgr          *ecs.World
	camera       *game.Camera
	lastMouse    image.Point
	wasMouseDown bool
}

func main() {
	// Exercise structure.
	mgr := ecs.NewWorld()
	e := mgr.NewEntity()
	p := &game.Position{}
	mgr.AddComponent(e, p)
	mgr.RemoveComponent(e, p)
	mgr.RemoveType(e, "BananaramaLlamaDingdong")
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

func controlCamera(c *game.Camera, t time.Duration, ctrl Controls) {
	camSpeed := 500.0 / c.GetZoom()
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

var last time.Time

// setup the game Entities.
func setup(w, h int) (*system, error) {
	mgr := ecs.NewWorld()

	camera := game.NewCamera(w, h)
	s := system{
		render: game.NewRenderer(),
		combat: combat.NewManager(mgr, camera),

		mgr:    mgr,
		camera: camera,

		fonts:     game.NewFontSystem(mgr),
		hierarchy: ecs.NewParentSystem(mgr),
	}

	// Create some Actors that are controlled by mouse clicks
	mgr.AddComponent(mgr.NewEntity(), &game.Actor{
		Size:                 game.SMALL,
		PreparationThreshold: 701,
		ActionPoints:         100,
	})
	mgr.AddComponent(mgr.NewEntity(), &game.Actor{
		Size:                 game.MEDIUM,
		PreparationThreshold: 1103,
		ActionPoints:         80,
	})
	mgr.AddComponent(mgr.NewEntity(), &game.Actor{
		Size:                 game.LARGE,
		PreparationThreshold: 1301,
		ActionPoints:         120,
	})

	// Start combat!
	s.combat.Begin()

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
		s.wasMouseDown = true
	} else if s.wasMouseDown {
		s.combat.Interaction(x, y)
		s.wasMouseDown = false
	}

	elapsed := time.Since(last)
	last = time.Now()

	ctrl := controls()
	controlCamera(s.camera, elapsed, ctrl)

	s.combat.Run(elapsed)
	s.fonts.Update()
	s.hierarchy.Update()

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
