/* tool overworld previews overworld map generation
 */
package main

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/game/overworld/procedural"
	"github.com/griffithsh/squads/geom"
	"github.com/griffithsh/squads/output"
	"github.com/hajimehoshi/ebiten/v2"
)

var errExitGame = errors.New("game has completed")

const screenWidth, screenHeight = 1024, 768

var resource image.Image

type overworldGenerator struct {
	mgr *ecs.World
	bus *event.Bus
	vis *output.Visualizer
}

func (g *overworldGenerator) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return errExitGame
	}

	return nil
}

func (g *overworldGenerator) Draw(screen *ebiten.Image) {
	err := g.vis.Render(screen, g.mgr, 0, 0, 0.5, screenWidth, screenHeight)
	if err != nil {
		fmt.Printf("render: %v\n", err)
	}
}
func (g *overworldGenerator) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func init() {
	f, err := os.Open("./temporary.png")
	if err != nil {
		panic(fmt.Errorf("couldn't open the image file %v", err))
	}
	decoded, err := png.Decode(f)
	if err != nil {
		panic(fmt.Errorf("couldn't decode the image file %v", err))
	}
	resource = decoded
}

type imageGetter struct{}

func (ig imageGetter) GetImage(name string) (val image.Image, ok bool) {
	return resource, true
}

func main() {
	generated := procedural.Generate()
	mgr := ecs.NewWorld()
	f := geom.NewField(36, 16, 34)
	for k, placement := range generated.Paths {
		x, y := f.Ktow(k)

		e := mgr.NewEntity()

		mgr.AddComponent(e, &generated.BaseTerrain)
		mgr.AddComponent(e, &game.Position{
			Center: game.Center{
				X: x,
				Y: y,
			},
			Layer: 0,
		})
		for dir := range placement.Connections {
			spr := generated.RoadSprites[dir]
			e := mgr.NewEntity()
			mgr.AddComponent(e, &spr)
			mgr.AddComponent(e, &game.Position{
				Center: game.Center{
					X: x,
					Y: y,
				},
				Layer: 10,
			})
		}

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
