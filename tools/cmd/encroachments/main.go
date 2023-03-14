package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/griffithsh/squads/data"
	"github.com/griffithsh/squads/ecs"
	"github.com/griffithsh/squads/event"
	"github.com/griffithsh/squads/game"
	"github.com/griffithsh/squads/game/overworld/hbg"
	"github.com/griffithsh/squads/game/overworld/procedural"
	"github.com/griffithsh/squads/geom"
	"github.com/griffithsh/squads/output"
	"github.com/hajimehoshi/ebiten/v2"
)

var terrainFile string

func init() {
	flag.StringVar(&terrainFile, "generated", "", "specify a json-marshaled procedural.Generated")
}

// Accepts a generated overworld level as input, and renders an overworld as
// the main squads binary would, loading encroachment and tile data as defined by the recipe.
func main() {
	flag.Parse()

	mgr := ecs.NewWorld()
	bus := &event.Bus{}
	archive, err := data.NewArchive()
	if err != nil {
		fmt.Fprintf(os.Stderr, "construct data archive: %v\n", err)
		os.Exit(1)
	}
	f, err := os.Open("./squads.data")
	if err != nil {
		fmt.Fprintf(os.Stderr, "open: %v\n", err)
		os.Exit(1)
	}
	err = archive.Load(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load archive: %v\n", err)
		os.Exit(1)
	}
	e := &engine{
		mgr:       mgr,
		bus:       bus,
		vis:       output.NewVisualizer(archive),
		animation: &game.AnimationSystem{},
	}

	terrain := map[geom.Key]procedural.Code{}
	if _, err := os.Stat(terrainFile); errors.Is(err, os.ErrNotExist) {
		fmt.Fprintf(os.Stderr, "file %q does not exist. Specify a file containing a procedural.Generated marshaled to json\n", terrainFile)
		return
	}
	b, err := os.ReadFile(terrainFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "open generated file: %v\n", err)
		os.Exit(1)
	}
	var generated = procedural.Generated{}
	if err := json.Unmarshal(b, &generated); err != nil {
		fmt.Fprintf(os.Stderr, "unmarshal generated file: %v\n", err)
		os.Exit(1)
	}
	if len(generated.Terrain) == 0 {
		fmt.Fprintf(os.Stderr, "no terrain data in generated file\n")
		os.Exit(1)
	}

	for key, code := range generated.Terrain {
		terrain[key] = code
	}

	hbg.ConstructBackground(terrain, archive.GetOverworldBaseTiles(), archive.GetOverworldEncroachments(), geom.NewField(66, 31, 64), rand.New(rand.NewSource(time.Now().UnixMilli())), func(call hbg.ComponentHandlerCall) {
		var layer int
		frames := call.Frames
		switch call.Z {
		case hbg.BaseLayer:
			layer = 0
		case hbg.CornersLayer:
			layer = 100
		case hbg.EncroachmentsLayer:
			layer = 20
		default:
			panic(fmt.Sprintf("unhandled hbg.Layer %v", call.Z))
		}
		p := &game.Position{
			Center: game.Center{
				X: float64(call.X),
				Y: float64(call.Y),
			},
			Layer: layer,
		}
		var c ecs.Component
		if len(frames) == 1 {
			frame := frames[0]
			c = &game.Sprite{
				Texture: frame.Texture,
				X:       frame.L,
				Y:       frame.T,
				W:       frame.W,
				H:       frame.H,
			}
		} else {
			sprites := make([]game.Sprite, len(frames))
			timings := make([]time.Duration, len(frames))
			for i, frame := range frames {
				sprites[i] = game.Sprite{
					Texture: frame.Texture,
					X:       frame.L,
					Y:       frame.T,
					W:       frame.W,
					H:       frame.H,
				}
				timings[i] = frame.Duration
			}

			c = &game.FrameAnimation{
				Frames:  sprites,
				Timings: timings,
			}
		}
		e := mgr.NewEntity()
		mgr.AddComponent(e, p)
		mgr.AddComponent(e, c)
	})

	// Start ebiten looping
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("squads: tools/cmd/encroachments")
	bus.Publish(&game.WindowSizeChanged{
		OldW: 0,
		OldH: 0,
		NewW: screenWidth,
		NewH: screenHeight,
	})
	if err := ebiten.RunGame(e); err != nil && err != errExitGame {
		fmt.Fprintf(os.Stderr, "RunGame: %v\n", err)
		os.Exit(1)
	}
}
