/* tool overworld previews overworld map generation
 */
package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/game/overworld/procedural"
	"github.com/griffithsh/squads/output"
	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed recipes/*
var content embed.FS

const screenWidth, screenHeight = 1024, 768

func main() {
	// recipe, err := content.ReadFile("recipes/atoll.json")
	// recipe, err := content.ReadFile("recipes/edge-of-the-woods.json")
	recipe, err := content.ReadFile("recipes/lakeside.json")
	// recipe, err := content.ReadFile("recipes/shore.json")
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

	// Create an instance of an ebiten "Game"
	mgr := ecs.NewWorld()
	bus := &event.Bus{}
	g := &overworldGenerator{
		mgr:  mgr,
		bus:  bus,
		vis:  output.NewVisualizer(imageGetter{}),
		core: generator,
	}

	// Generate an overworld!
	g.Generate()

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
