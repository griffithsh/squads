package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
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

// prevUI is a thing to preview UIs
type prevUI struct {
	mgr *ecs.World
	bus *event.Bus
	vis *output.Visualizer
	ui  *ui.UISystem
}

func (p *prevUI) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return errExitGame
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		p.bus.Publish(&ui.Interact{
			AbsoluteX: float64(x),
			AbsoluteY: float64(y),
		})
	}
	p.ui.Update()
	return nil
}

func (p *prevUI) Draw(screen *ebiten.Image) {
	p.vis.Render(screen, p.mgr, 0, 0, 1.0, screenWidth, screenHeight)
}

func (p *prevUI) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
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
		mgr: mgr,
		bus: bus,
		vis: output.NewVisualizer(archive),
		ui:  ui.NewUISystem(mgr, bus),
	}

	// Set up test data
	// setupCombatUI(mgr, archive)
	setupEmbarkFocusCharacter(mgr, archive)

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
