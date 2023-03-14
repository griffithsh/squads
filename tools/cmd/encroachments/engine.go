package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/output"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type engine struct {
	mgr *ecs.World
	bus *event.Bus
	vis *output.Visualizer

	animation *game.AnimationSystem
}

var errExitGame = errors.New("game has completed")

func (g *engine) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return errExitGame
	}

	// Zoom controls
	if inpututil.IsKeyJustPressed(ebiten.KeyMinus) {
		zoom /= 2
	} else if inpututil.IsKeyJustPressed(ebiten.KeyEqual) {
		zoom *= 2
	}

	// panning
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		focusY -= 42
	} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		focusY += 42
	} else if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		focusX -= 42
	} else if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		focusX += 42
	}

	g.animation.Update(g.mgr, time.Second/60)
	return nil
}

const screenWidth, screenHeight = 1800, 1050

var zoom = 0.5
var focusX = 0.0
var focusY = 0.0

func (g *engine) Draw(screen *ebiten.Image) {
	err := g.vis.Render(screen, g.mgr, focusX, focusY, zoom, screenWidth, screenHeight)
	if err != nil {
		fmt.Printf("render: %v\n", err)
	}
}

func (g *engine) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
