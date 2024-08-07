package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"

	"github.com/griffithsh/squads/data"
	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/output"
	"github.com/griffithsh/squads/ui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var errExitGame = errors.New("game has completed")

const screenWidth, screenHeight = 1024, 768

var e ecs.Entity
var uis []func(*ecs.World, *data.Archive) = []func(*ecs.World, *data.Archive){
	setupCombatUI,
	setupCombatUIPreparing,
	setupEmbarkFocusCharacter,
}

// prevUI is a thing to preview UIs
type prevUI struct {
	mgr     *ecs.World
	bus     *event.Bus
	archive *data.Archive
	vis     *output.Visualizer
	ui      *ui.UISystem
	last    time.Time
}

func (p *prevUI) Update() error {
	elapsed := time.Since(p.last)
	p.last = time.Now()

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return errExitGame
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		p.nextUI()
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		p.bus.Publish(&ui.Interact{
			AbsoluteX: float64(x),
			AbsoluteY: float64(y),
		})
	}
	return p.ui.Update(elapsed)
}

func (p *prevUI) Draw(screen *ebiten.Image) {
	err := p.vis.Render(screen, p.mgr, 0, 0, 1.0, screenWidth, screenHeight)
	if err != nil {
		fmt.Printf("render: %v\n", err)
	}
}

func (p *prevUI) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func (p *prevUI) nextUI() {
	p.mgr.DestroyEntity(e)
	uis[0](p.mgr, p.archive)

	// Matthew 19:30
	first := uis[0]
	uis = append(uis[1:], first)
}

func main() {
	cpuProfile := flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()
	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	// Seed the PRNG
	rand.Seed(time.Now().UnixNano())

	// Set up an archive for data
	archive, err := data.NewArchive()
	if err != nil {
		fmt.Fprintf(os.Stderr, "construct data archive: %v", err)
		os.Exit(1)
	}
	f, err := os.Open("./squads.data")
	if err != nil {
		fmt.Fprintf(os.Stderr, "open: %v", err)
		os.Exit(1)
	}
	err = archive.Load(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load archive: %v", err)
		os.Exit(1)
	}

	// Create an instance of an ebiten "Game"
	mgr := ecs.NewWorld()
	bus := &event.Bus{}
	p := &prevUI{
		mgr:     mgr,
		bus:     bus,
		archive: archive,
		vis:     output.NewVisualizer(archive),
		ui:      ui.NewUISystem(mgr, bus),
		last:    time.Now(),
	}

	// Set up first lot of test data
	p.nextUI()

	// Start ebiten looping
	ebiten.SetWindowSize(screenWidth, screenHeight)
	bus.Publish(&game.WindowSizeChanged{
		OldW: 0,
		OldH: 0,
		NewW: screenWidth,
		NewH: screenHeight,
	})
	if err := ebiten.RunGame(p); err != nil && err != errExitGame {
		fmt.Fprintf(os.Stderr, "RunGame: %v\n", err)
		os.Exit(1)
	}
}
