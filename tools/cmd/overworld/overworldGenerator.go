package main

import (
	"errors"
	"fmt"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/output"
	"github.com/hajimehoshi/ebiten/v2"
)

type overworldGenerator struct {
	mgr *ecs.World
	bus *event.Bus
	vis *output.Visualizer
}

var errExitGame = errors.New("game has completed")

func (g *overworldGenerator) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return errExitGame
	}

	moveSpeed := 10.0
	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		focusY -= moveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		focusX -= moveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		focusY += moveSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		focusX += moveSpeed
	}

	return nil
}

var zoom = 1.0
var focusX = 0.0
var focusY = 0.0

func (g *overworldGenerator) Draw(screen *ebiten.Image) {
	err := g.vis.Render(screen, g.mgr, focusX, focusY, zoom, screenWidth, screenHeight)
	if err != nil {
		fmt.Printf("render: %v\n", err)
	}
}

func (g *overworldGenerator) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
