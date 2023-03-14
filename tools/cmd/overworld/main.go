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
	"github.com/griffithsh/squads/output"
	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed recipes/*
var content embed.FS

const screenWidth, screenHeight = 1024, 768

func main() {
	// newRecipe returns a randomly selected new recipe.
	newRecipe := func() *procedural.Generator {
		recipes := []string{
			// "recipes/atoll.json",
			// "recipes/dark-forest.json",
			// "recipes/desert.json",
			// "recipes/edge-of-the-woods.json",
			"recipes/lakeside.json",
			// "recipes/shore.json",
		}
		i := 0
		if len(recipes) > 1 {
			rand.Seed(time.Now().UnixMilli())
			i = rand.Intn(len(recipes))
		}
		recipe, err := content.ReadFile(recipes[i])
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
		return &generator
	}

	// Create an instance of an ebiten "Game"
	mgr := ecs.NewWorld()
	bus := &event.Bus{}
	// var seed int64 = 5546037425800197631
	g := &overworldGenerator{
		mgr:               mgr,
		bus:               bus,
		vis:               output.NewVisualizer(imageGetter{}),
		generatorProvider: func() *procedural.Generator { return nil },
		core:              newRecipe(),
		// forceSeed:         &seed,
	}

	// Generate an overworld!
	g.Generate()

	// Start ebiten looping
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("squads: tools/cmd/overworld")
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
