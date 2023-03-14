package hbg

import (
	"math/rand"
	"sort"
	"time"

	"github.com/griffithsh/squads/game/overworld/procedural"
	"github.com/griffithsh/squads/geom"
)

type Layer int

const (
	BaseLayer = iota
	CornersLayer
	EncroachmentsLayer
	AdornmentsLayer // AdornmentsLayer has trees, bushes, rocks, anything with a bit of three-dimensionality.
)

type ComponentHandlerCall struct {
	Frames []ComponentHandlerFrame
	X, Y   int
	Z      Layer
}
type ComponentHandlerFrame struct {
	Texture    string
	Duration   time.Duration
	L, T, W, H int
}

type ComponentHandler func(frames ComponentHandlerCall)

// ConstructBackground takes terrain and translates it into handler calls that define the hexes that compose a background.
func ConstructBackground(terrain map[geom.Key]procedural.Code, baseTiles map[procedural.Code]BaseTile, encroachments EncroachmentsCollection, field *geom.Field, prng *rand.Rand, handler ComponentHandler) {
	hexHeight := field.HexHeight()
	hexWidth := field.HexWidth()

	// Start sorted I guess
	keys := make([]geom.Key, 0, len(terrain))
	for k := range terrain {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].N == keys[j].N {
			return keys[i].M < keys[j].M
		}
		return keys[i].N < keys[j].N
	})

	// For each hex...
	for _, key := range keys {
		keyCode := terrain[key]
		keyX, keyY := field.Ktow(key)
		// layer on a draw call for the base tile artwork
		baseTile := baseTiles[keyCode]
		option := baseTile.Variations.Roll(prng.Intn)
		if option != nil {
			baseFrames := option.Frames
			w := hexWidth
			h := hexHeight + option.ExtraHeight + option.ExtraDepth
			// Centered on ...
			x := int(keyX) - w/2
			// y := int(keyY) - h/2 + yOffset
			y := int(keyY) - hexHeight/2 - option.ExtraHeight/2 + option.ExtraDepth/2

			temp := ComponentHandlerCall{
				Frames: make([]ComponentHandlerFrame, len(baseFrames)),
				X:      x,
				Y:      y,
				Z:      BaseLayer,
			}

			for i, frame := range baseFrames {
				temp.Frames[i].Texture = baseTile.Texture
				temp.Frames[i].Duration = frame.Duration
				temp.Frames[i].L = frame.Left
				temp.Frames[i].T = frame.Top
				temp.Frames[i].W = w
				temp.Frames[i].H = h
			}
			handler(temp)
		}

		neighbors := map[geom.DirectionType]procedural.Code{
			geom.N:  terrain[key.ToDirection(geom.N)],
			geom.NE: terrain[key.ToDirection(geom.NE)],
			geom.SE: terrain[key.ToDirection(geom.SE)],
			geom.S:  terrain[key.ToDirection(geom.S)],
			geom.SW: terrain[key.ToDirection(geom.SW)],
			geom.NW: terrain[key.ToDirection(geom.NW)],
		}

		// When we need to know if we need a corner encroachment, this defines
		// which neighbor has the responsiblity of which corner.
		nextNeighbors := map[geom.DirectionType]geom.DirectionType{
			geom.N:  geom.NE,
			geom.NE: geom.SE,
			geom.SE: geom.S,
			geom.S:  geom.SW,
			geom.SW: geom.NW,
			geom.NW: geom.N,
		}

		for _, dir := range []geom.DirectionType{geom.N, geom.NE, geom.SE, geom.S, geom.SW, geom.NW} {
			// For each of its neighbors...
			neighborCode := neighbors[dir]
			e := encroachments.Get(keyCode, neighborCode)
			if e == nil {
				// There is nothing to do if this hex's code does not encroach its neighbor's code.
				continue
			}

			options, ok := e.Edges.Options[dir]
			if !ok {
				// If there's nothing configured when encroaching in this direction, that's weird, but the best way forward is to skip drawing anything.
				continue
			}
			option := options.Roll(prng.Intn)
			if option == nil {
				// When no Options are configured, then there's nothing to draw.
				continue
			}

			// Layer on some draw calls for the edge encroachment
			edgeFrames := option.Frames
			w := hexWidth
			h := hexHeight + option.ExtraHeight + option.ExtraDepth
			// Centered on ...
			x := int(keyX) - w/2
			// y := int(keyY) - h/2 + yOffset
			y := int(keyY) - hexHeight/2 - option.ExtraHeight/2 + option.ExtraDepth/2

			temp := ComponentHandlerCall{
				Frames: make([]ComponentHandlerFrame, len(edgeFrames)),
				X:      x,
				Y:      y,
				Z:      EncroachmentsLayer,
			}

			for i, frame := range edgeFrames {
				temp.Frames[i].Texture = e.Edges.Texture
				temp.Frames[i].Duration = frame.Duration
				temp.Frames[i].L = frame.Left
				temp.Frames[i].T = frame.Top
				temp.Frames[i].W = w
				temp.Frames[i].H = h
			}
			handler(temp)

			// Handle the corner encroachments...

			// We need to only deal with corners once, so we'll check the
			// mapping between directions and their responsibilities.
			nextNeighbor := nextNeighbors[dir]
			nextNeighborCode := terrain[key.ToDirection(nextNeighbor)]
			if nextNeighborCode != keyCode {
				continue
			}

			corner, ok := e.Corners.Corners[dir]
			if !ok {
				continue
			}

			w = e.Corners.Corners[dir].W
			h = e.Corners.Corners[dir].H

			option = corner.Options.Roll(prng.Intn)
			if option == nil {
				// When no Options are configured, then there's nothing to draw.
				continue
			}

			x = int(keyX) - w/2
			y = int(keyY) - h/2 - option.ExtraHeight/2 + option.ExtraDepth/2

			cornerFrames := option.Frames
			temp = ComponentHandlerCall{
				Frames: make([]ComponentHandlerFrame, len(cornerFrames)),
				X:      x + corner.Offset.X,
				Y:      y + corner.Offset.Y,
				Z:      CornersLayer,
			}

			for i, frame := range cornerFrames {
				temp.Frames[i].Texture = e.Corners.Texture
				temp.Frames[i].Duration = frame.Duration
				temp.Frames[i].L = frame.Left
				temp.Frames[i].T = frame.Top
				temp.Frames[i].W = w
				temp.Frames[i].H = h
			}

			handler(temp)
		}
	}
}
