package main

import (
	"fmt"
	"image"
	"math/rand"
	"os"
	"time"

	"github.com/griffithsh/squads/data"
	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/game/combat"
	"github.com/griffithsh/squads/game/embark"
	"github.com/griffithsh/squads/game/overworld"
	"github.com/griffithsh/squads/output"
	"github.com/griffithsh/squads/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type squads struct {
	bus          *event.Bus
	video        *output.Visualizer
	expiry       *ecs.ExpirySystem
	anim         *game.AnimationSystem
	traversals   *overworld.TraversalSystem
	collisions   *overworld.CollisionSystem
	fonts        *game.FontSystem
	hierarchy    *ecs.ParentSystem
	leash        *game.LeashSystem
	fades        *game.FadeSystem
	wipes        *game.SceneWipeSystem
	interactives *ui.InteractiveSystem

	embark    *embark.Manager
	overworld *overworld.Manager
	combat    *combat.Manager

	mgr       *ecs.World
	camera    *game.Camera
	lastMouse image.Point
	last      time.Time
}

func newSquads(w, h int) (*squads, error) {
	bus := &event.Bus{}
	mgr := ecs.NewWorld()
	camera := game.NewCamera(w, h, bus)
	archive, err := data.NewArchive()
	if err != nil {
		return nil, fmt.Errorf("construct data archive: %v", err)
	}

	f, err := os.Open("./squads.data")
	if err != nil {
		return nil, fmt.Errorf("open: %v", err)
	}
	err = archive.Load(f)
	if err != nil {
		return nil, fmt.Errorf("load archive: %v", err)
	}
	s := squads{
		bus:        bus,
		video:      output.NewVisualizer(archive),
		anim:       &game.AnimationSystem{},
		expiry:     ecs.NewExpirySystem(mgr),
		traversals: &overworld.TraversalSystem{},
		collisions: overworld.NewCollisionSystem(mgr, bus),
		embark:     embark.NewManager(mgr, bus, archive),
		overworld:  overworld.NewManager(mgr, bus, archive),
		combat:     combat.NewManager(mgr, camera, bus, archive),

		mgr:    mgr,
		camera: camera,

		fonts:        game.NewFontSystem(mgr),
		hierarchy:    ecs.NewParentSystem(mgr),
		leash:        &game.LeashSystem{},
		fades:        &game.FadeSystem{},
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
	bus.Subscribe(embark.Embarked{}.Type(), func(t event.Typer) {
		s.embark.End()
		rng := rand.New(rand.NewSource(time.Now().UnixNano()))
		s.overworld.Begin(rng.Int63())
	})
	bus.Subscribe(overworld.Complete{}.Type(), func(t event.Typer) {
		s.overworld.End()
		rng := rand.New(rand.NewSource(time.Now().UnixNano()))
		s.overworld.Begin(rng.Int63())
	})

	s.bus.Publish(&game.WindowSizeChanged{
		OldW: 0,
		OldH: 0,
		NewW: w,
		NewH: h,
	})

	// Init combat?
	// TODO: pass a thing that has enough information to construct a Field and
	// the enemies you'll face in the combat
	// teams := []ecs.Entity{}
	// s.combat.Begin(teams)
	// s.combat.End()

	s.embark.Begin()

	s.last = time.Now()

	return &s, nil
}

func (s *squads) setScreenSize(w, h int) {
	s.bus.Publish(&game.WindowSizeChanged{
		OldW: s.camera.GetW(),
		OldH: s.camera.GetH(),
		NewW: w,
		NewH: h,
	})
	ebiten.SetWindowSize(w, h)
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

var (
	frames      = 0
	accumulated = time.Second * 0
	second      = time.Tick(time.Second)
)

func (s *squads) Update() error {
	// While developing this, it's nice to be able to kill the game quickly to get back to the code.
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return errExitGame
	}

	start := time.Now()
	defer func() {
		d := time.Since(start)
		frames++
		accumulated += d
	}()

	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		s.setScreenSize(640, 480)
	} else if inpututil.IsKeyJustPressed(ebiten.Key2) {
		s.setScreenSize(800, 600)
	} else if inpututil.IsKeyJustPressed(ebiten.Key3) {
		s.setScreenSize(1024, 768)
	}

	x, y := ebiten.CursorPosition()

	if s.lastMouse.X != x || s.lastMouse.Y != y {
		s.combat.MousePosition(x, y)
		s.lastMouse.X = x
		s.lastMouse.Y = y
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		wx, wy := s.camera.ScreenToWorld(x, y)
		s.bus.Publish(&ui.Interact{
			AbsoluteX: float64(x),
			AbsoluteY: float64(y),
			X:         wx,
			Y:         wy,
		})
	}

	elapsed := time.Since(s.last)
	s.last = time.Now()

	ctrl := controls()
	debugControlCamera(s.camera, elapsed, ctrl)

	s.combat.Run(elapsed)
	s.overworld.Run(elapsed)

	s.expiry.Update(elapsed)
	s.fonts.Update()
	s.leash.Update(s.mgr, elapsed)
	s.fades.Update(s.mgr, elapsed)
	s.anim.Update(s.mgr, elapsed)
	s.traversals.Update(s.mgr, elapsed)
	s.wipes.Update(s.mgr, elapsed)
	s.hierarchy.Update()
	s.camera.Update(elapsed)

	select {
	case <-second:
		var fps time.Duration
		if time.Duration(frames) > 0 {
			fps = time.Second / (accumulated / time.Duration(frames))
		}
		title := "Project Never"
		ebiten.SetWindowTitle(fmt.Sprintf("%s | FPS: %d | Entites: %d", title, fps, s.mgr.Len()))
	default:
	}
	return nil
}

func (s *squads) Draw(screen *ebiten.Image) {
	w, h := ebiten.WindowSize()
	if err := s.video.Render(screen, s.mgr, s.camera.GetX(), s.camera.GetY(), s.camera.GetZoom(), float64(w), float64(h)); err != nil {
		panic("Draw frame: " + err.Error())
	}
}

func (s *squads) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
