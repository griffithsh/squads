package main

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/game/overworld/procedural"
	"github.com/griffithsh/squads/geom"
	"github.com/griffithsh/squads/output"
	"github.com/hajimehoshi/ebiten/v2"
)

type overworldGenerator struct {
	mgr *ecs.World
	bus *event.Bus
	vis *output.Visualizer

	core procedural.Generator
}

func (g *overworldGenerator) Generate() {
	g.mgr.Clear()
	seed := time.Now().UnixMilli()

	generated := g.core.Generate(seed, 0)
	var field = geom.NewField(36, 16, 34)

	terrainSprites := map[procedural.Code]*game.Sprite{
		"WATER":  {Texture: "temporary.png", X: 136, Y: 0, W: 68, H: 34},
		"SAND":   {Texture: "temporary.png", X: 136, Y: 34, W: 68, H: 34},
		"GRASS":  {Texture: "temporary.png", X: 136, Y: 68, W: 68, H: 34},
		"FOREST": {Texture: "temporary.png", X: 136, Y: 102, W: 68, H: 34},
		"ROCK":   {Texture: "temporary.png", X: 136, Y: 136, W: 68, H: 34},
	}

	// Add terrain!
	for key, code := range generated.Terrain {
		e := g.mgr.NewEntity()
		x, y := field.Ktow(key)
		g.mgr.AddComponent(e, &game.Position{
			Center: game.Center{
				X: x,
				Y: y,
			},
			Layer: 0,
		})
		g.mgr.AddComponent(e, terrainSprites[code])
	}

	// Add roads!
	roadSprites := []map[geom.DirectionType]*game.Sprite{
		{
			geom.NE: {Texture: "temporary.png", X: 204, Y: 0, W: 68, H: 34},
			geom.NW: {Texture: "temporary.png", X: 272, Y: 0, W: 68, H: 34},
			geom.SE: {Texture: "temporary.png", X: 204, Y: 34, W: 68, H: 34},
			geom.SW: {Texture: "temporary.png", X: 272, Y: 34, W: 68, H: 34},
			geom.N:  {Texture: "temporary.png", X: 204, Y: 68, W: 68, H: 34},
			geom.S:  {Texture: "temporary.png", X: 272, Y: 68, W: 68, H: 34},
		},
		{
			geom.NE: {Texture: "temporary.png", X: 340, Y: 0, W: 68, H: 34},
			geom.NW: {Texture: "temporary.png", X: 408, Y: 0, W: 68, H: 34},
			geom.SE: {Texture: "temporary.png", X: 340, Y: 34, W: 68, H: 34},
			geom.SW: {Texture: "temporary.png", X: 408, Y: 34, W: 68, H: 34},
			geom.N:  {Texture: "temporary.png", X: 340, Y: 68, W: 68, H: 34},
			geom.S:  {Texture: "temporary.png", X: 408, Y: 68, W: 68, H: 34},
		},
		{
			geom.NE: {Texture: "temporary.png", X: 204, Y: 102, W: 68, H: 34},
			geom.NW: {Texture: "temporary.png", X: 272, Y: 102, W: 68, H: 34},
			geom.SE: {Texture: "temporary.png", X: 204, Y: 136, W: 68, H: 34},
			geom.SW: {Texture: "temporary.png", X: 272, Y: 136, W: 68, H: 34},
			geom.N:  {Texture: "temporary.png", X: 204, Y: 170, W: 68, H: 34},
			geom.S:  {Texture: "temporary.png", X: 272, Y: 170, W: 68, H: 34},
		},
		{
			geom.NE: {Texture: "temporary.png", X: 340, Y: 102, W: 68, H: 34},
			geom.NW: {Texture: "temporary.png", X: 408, Y: 102, W: 68, H: 34},
			geom.SE: {Texture: "temporary.png", X: 340, Y: 136, W: 68, H: 34},
			geom.SW: {Texture: "temporary.png", X: 408, Y: 136, W: 68, H: 34},
			geom.N:  {Texture: "temporary.png", X: 340, Y: 170, W: 68, H: 34},
			geom.S:  {Texture: "temporary.png", X: 408, Y: 170, W: 68, H: 34},
		},
	}
	prng := rand.New(rand.NewSource(0))
	for key, placement := range generated.Paths.Nodes {
		x, y := field.Ktow(key)
		pos := game.Position{
			Center: game.Center{
				X: x,
				Y: y,
			},
			Layer: 1,
		}
		i := procedural.DeterministicIndexOf(prng, roadSprites)
		for dir := range placement.Connections {
			e := g.mgr.NewEntity()
			g.mgr.AddComponent(e, &pos)
			g.mgr.AddComponent(e, roadSprites[i][dir])
		}
	}

	// add baddies
	for k := range generated.Opponents {
		x, y := field.Ktow(k)

		e := g.mgr.NewEntity()
		g.mgr.AddComponent(e, &game.Sprite{
			Texture: "temporary.png",

			X: 0, Y: 204,
			W: 68, H: 34,
		})
		g.mgr.AddComponent(e, &game.Position{
			Center: game.Center{
				X: x,
				Y: y,
			},
			Layer: 110,
		})
	}

	// dirSprites show the extents being generated - just for debugging.
	dirSprites := map[geom.DirectionType]*game.Sprite{
		geom.N:  {Texture: "temporary.png", W: 68, H: 34, X: 0, Y: 170},
		geom.NE: {Texture: "temporary.png", W: 68, H: 34, X: 0, Y: 102},
		geom.SE: {Texture: "temporary.png", W: 68, H: 34, X: 0, Y: 136},
		geom.S:  {Texture: "temporary.png", W: 68, H: 34, X: 68, Y: 170},
		geom.SW: {Texture: "temporary.png", W: 68, H: 34, X: 68, Y: 136},
		geom.NW: {Texture: "temporary.png", W: 68, H: 34, X: 68, Y: 102},
	}

	for dir, key := range generated.PathExtents {
		x, y := field.Ktow(key)

		e := g.mgr.NewEntity()
		g.mgr.AddComponent(e, dirSprites[dir])
		g.mgr.AddComponent(e, &game.Position{
			Center: game.Center{
				X: x,
				Y: y,
			},
			Layer: 100,
		})
	}

	// Add start hex
	x, y := field.Ktow(generated.Paths.Start)
	e := g.mgr.NewEntity()
	g.mgr.AddComponent(e, &game.Sprite{
		Texture: "temporary.png",
		X:       68,
		Y:       238,
		W:       68,
		H:       34,
	})
	g.mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: x, Y: y,
		},
		Layer: 1000,
	})
	// Add goal hex
	x, y = field.Ktow(generated.Paths.Goal)
	e = g.mgr.NewEntity()
	g.mgr.AddComponent(e, &game.Sprite{
		Texture: "temporary.png",
		X:       0,
		Y:       238,
		W:       68,
		H:       34,
	})
	g.mgr.AddComponent(e, &game.Position{
		Center: game.Center{
			X: x, Y: y,
		},
		Layer: 1000,
	})
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

	// Regenerate map button.
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		if !regeneratingDown {
			focusX = 0
			focusY = 0
			g.Generate()
			regeneratingDown = true
		}
	} else {
		regeneratingDown = false
	}

	return nil
}

var regeneratingDown bool
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
