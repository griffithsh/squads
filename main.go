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

	"github.com/griffithsh/squads/ui"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/game/combat"
	"github.com/griffithsh/squads/game/embark"
	"github.com/griffithsh/squads/game/overworld"
	"github.com/hajimehoshi/ebiten"
)

type system struct {
	bus          *event.Bus
	render       *game.Renderer
	anim         *game.AnimationSystem
	traversals   *overworld.TraversalSystem
	collisions   *overworld.CollisionSystem
	fonts        *game.FontSystem
	hierarchy    *ecs.ParentSystem
	leash        *game.LeashSystem
	wipes        *game.SceneWipeSystem
	interactives *ui.InteractiveSystem

	embark    *embark.Manager
	overworld *overworld.Manager
	combat    *combat.Manager

	mgr          *ecs.World
	camera       *game.Camera
	lastMouse    image.Point
	wasMouseDown bool
}

func main() {
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
	w, h := 1024, 768
	s, _ := setup(w, h)
	if err := ebiten.Run(s.run, w, h, 1, "Squads"); err == errExitGame {
		fmt.Println("See you next time.")
	} else if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

}

// Controls represent debug controls
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

func debugControlCamera(c *game.Camera, t time.Duration, ctrl Controls) {
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
	bus := &event.Bus{}
	mgr := ecs.NewWorld()
	camera := game.NewCamera(w, h, bus)
	s := system{
		bus:        bus,
		render:     game.NewRenderer(),
		anim:       &game.AnimationSystem{},
		traversals: &overworld.TraversalSystem{},
		collisions: overworld.NewCollisionSystem(mgr, bus),
		embark:     embark.NewManager(mgr, bus),
		overworld:  overworld.NewManager(mgr, bus),
		combat:     combat.NewManager(mgr, camera, bus),

		mgr:    mgr,
		camera: camera,

		fonts:        game.NewFontSystem(mgr),
		hierarchy:    ecs.NewParentSystem(mgr),
		leash:        &game.LeashSystem{},
		wipes:        game.NewSceneWipeSystem(),
		interactives: ui.NewInteractiveSystem(mgr, bus),
	}
	bus.Subscribe(game.CombatConcluded{}.Type(), func(et event.Typer) {
		s.combat.End()

		// Handle results of combat.
		ev := et.(*game.CombatConcluded)
		for e, result := range ev.Results {
			if mgr.HasTag(e, "player") {
				switch result {
				case game.Escaped:
					// player escaped, others are removed
					for otherEntity := range ev.Results {
						if e == otherEntity {
							continue
						}
						mgr.DestroyEntity(otherEntity)

					}
				case game.Defeated:
					// game is over
					mgr.DestroyEntity(e)
				}
			} else if result != game.Victorious {
				// baddy squad goes away.
				mgr.DestroyEntity(e)
			}
		}

		// force cascade of deleted components
		s.hierarchy.Update()

		s.overworld.Enable()
	})
	bus.Subscribe(overworld.CombatInitiated{}.Type(), func(t event.Typer) {
		s.overworld.Disable()
		ev := t.(*overworld.CombatInitiated)
		s.combat.Begin(ev.Squads)
	})
	bus.Subscribe(embark.SquadSelected{}.Type(), func(t event.Typer) {
		s.embark.End()
		rng := rand.New(rand.NewSource(time.Now().UnixNano()))
		s.overworld.Begin(data(rng))
	})

	// Init combat?
	// TODO: pass a thing that has enough information to construct a Field and
	// the enemies you'll face in the combat
	// teams := []ecs.Entity{}
	// s.combat.Begin(teams)
	// s.combat.End()

	s.embark.Begin()

	last = time.Now()

	return &s, nil
}

var errExitGame = errors.New("game has completed")

var (
	frames      = 0
	accumulated = time.Second * 0
	second      = time.Tick(time.Second)
)

func (s *system) setScreenSize(w, h int) {
	s.bus.Publish(&game.WindowSizeChanged{
		OldW: s.camera.GetW(),
		OldH: s.camera.GetH(),
		NewW: w,
		NewH: h,
	})
	ebiten.SetScreenSize(w, h)
}

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
	if ebiten.IsKeyPressed(ebiten.Key1) {
		s.setScreenSize(640, 480)
	} else if ebiten.IsKeyPressed(ebiten.Key2) {
		s.setScreenSize(1024, 768)
	}

	// Debug code to swap to overworld and back with Tab, Ctrl-Tab
	// if ebiten.IsKeyPressed(ebiten.KeyTab) {
	// 	if ebiten.IsKeyPressed(ebiten.KeyControl) {
	// 		s.combat.End()
	// 		s.overworld.Enable()
	// 	} else {
	// 		s.combat.()
	// 		s.overworld.Disable()
	// 	}
	// }

	x, y := ebiten.CursorPosition()

	if s.lastMouse.X != x || s.lastMouse.Y != y {
		s.combat.MousePosition(x, y)
		s.lastMouse.X = x
		s.lastMouse.Y = y
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		s.wasMouseDown = true
	} else if s.wasMouseDown {
		s.wasMouseDown = false

		wx, wy := s.camera.ScreenToWorld(x, y)
		s.bus.Publish(&ui.Interact{
			AbsoluteX: float64(x),
			AbsoluteY: float64(y),
			X:         wx,
			Y:         wy,
		})
	}

	elapsed := time.Since(last)
	last = time.Now()

	ctrl := controls()
	debugControlCamera(s.camera, elapsed, ctrl)

	s.combat.Run(elapsed)
	s.overworld.Run(elapsed)

	s.fonts.Update()
	s.leash.Update(s.mgr, elapsed)
	s.anim.Update(s.mgr, elapsed)
	s.traversals.Update(s.mgr, elapsed)
	s.wipes.Update(s.mgr, elapsed)
	s.hierarchy.Update()
	s.camera.Update(elapsed)

	w, h := float64(screen.Bounds().Max.X-screen.Bounds().Min.X), float64(screen.Bounds().Max.Y-screen.Bounds().Min.Y)

	// Render all entities in the World.
	if err := s.render.Render(screen, s.mgr, s.camera.GetX(), s.camera.GetY(), s.camera.GetZoom(), w, h); err != nil {
		return fmt.Errorf("Render: %v", err)
	}

	select {
	case <-second:
		var fps time.Duration
		if time.Duration(frames) > 0 {
			fps = time.Second / (accumulated / time.Duration(frames))
		}
		ebiten.SetWindowTitle(fmt.Sprintf("%s | FPS: %d", "Hexagons, Strategy, Entities, Components, and Systems, Oh my!", fps))
	default:
	}

	return nil
}
