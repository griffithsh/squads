/* tool overworld previews overworld map generation
 */
package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/game/overworld/procedural"
	"github.com/griffithsh/squads/geom"
	"github.com/griffithsh/squads/output"
	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed recipes/*
var content embed.FS

const screenWidth, screenHeight = 1024, 768

func main() {
	seed := time.Now().Unix()

	recipe, err := content.ReadFile("recipes/shore.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "read recipe: %v", err)
		os.Exit(1)
	}

	var generator procedural.Generator
	err = json.Unmarshal(recipe, &generator)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unmarshal generator: %v\n", err)
		os.Exit(1)
	}

	generated := generator.Generate(seed, 0)
	var field = geom.NewField(36, 16, 34)

	terrainSprites := map[procedural.Code]*game.Sprite{
		"WATER":  {Texture: "temporary.png", X: 136, Y: 0, W: 68, H: 34},
		"SAND":   {Texture: "temporary.png", X: 136, Y: 34, W: 68, H: 34},
		"GRASS":  {Texture: "temporary.png", X: 136, Y: 68, W: 68, H: 34},
		"FOREST": {Texture: "temporary.png", X: 136, Y: 102, W: 68, H: 34},
		"ROCK":   {Texture: "temporary.png", X: 136, Y: 136, W: 68, H: 34},
	}

	// Add terrain!
	mgr := ecs.NewWorld()
	for key, code := range generated.Terrain {
		e := mgr.NewEntity()
		x, y := field.Ktow(key)
		mgr.AddComponent(e, &game.Position{
			Center: game.Center{
				X: x,
				Y: y,
			},
			Layer: 0,
		})
		mgr.AddComponent(e, terrainSprites[code])
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
	for key, placement := range generated.Paths {
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
			e := mgr.NewEntity()
			mgr.AddComponent(e, &pos)
			mgr.AddComponent(e, roadSprites[i][dir])
		}
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

		e := mgr.NewEntity()
		mgr.AddComponent(e, dirSprites[dir])
		mgr.AddComponent(e, &game.Position{
			Center: game.Center{
				X: x,
				Y: y,
			},
			Layer: 100,
		})
	}

	// Create an instance of an ebiten "Game"
	bus := &event.Bus{}
	g := &overworldGenerator{
		mgr: mgr,
		bus: bus,
		vis: output.NewVisualizer(imageGetter{}),
	}
	// Start ebiten looping
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("overworld")
	bus.Publish(&game.WindowSizeChanged{
		OldW: 0,
		OldH: 0,
		NewW: screenWidth,
		NewH: screenHeight,
	})
	if err := ebiten.RunGame(g); err != nil && err != errExitGame {
		fmt.Fprintf(os.Stderr, "RunGame: %v\n", err)
		os.Exit(1)
	}
}
