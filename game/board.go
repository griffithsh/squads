package game

import (
	"fmt"

	"github.com/griffithsh/squads/ecs"
)

const (
	xStride = 34
	yStride = 8
)

// Hex is a hexagon tile that an Actor can occupy.
type Hex struct {
	M, N      int
	neighbors []*Hex
}

// Type of this Component.
func (h *Hex) Type() string {
	return "Hex"
}

// X coordinate of the center of this hexagon.
func (h *Hex) X() float64 {
	oddXOffset := 17
	return 12 + float64((xStride*h.M)+(h.N%2*oddXOffset))
}

// Y coordinate of the center of this hexagon.
func (h *Hex) Y() float64 {
	return 8 + float64(yStride*h.N)
}

// Neighbors of this Hex.
func (h *Hex) Neighbors() []*Hex {
	return h.neighbors
}

type key struct {
	M, N int
}

// Board to play out encounters on. Collection of Hexes.
type Board struct {
	stride int // how many hexes are in a row
	hexes  map[key]*Hex
}

// NewBoard creates a new game board.
func NewBoard(mgr *ecs.World, w, h int) (*Board, error) {
	m := make(map[key]*Hex, w)

	for i := 0; i < w*h; i++ {
		e := mgr.NewEntity()
		h := Hex{
			M: i % w,
			N: i / w,
		}
		m[key{i % w, i / w}] = &h
		mgr.AddComponent(e, &h)
		mgr.AddComponent(e, &Sprite{
			Texture: "texture.png",
			X:       24,
			Y:       0,
			W:       24,
			H:       16,
		})

		mgr.AddComponent(e, &Position{
			Center: Center{
				X: h.X(),
				Y: h.Y(),
			},
			Layer: 1,
		})

		if i == 1 || i%11 == 1 || i%17 == 1 || i%13 == 1 {
			// Scatter trees about
			e := mgr.NewEntity()
			mgr.AddComponent(e, &Sprite{
				Texture: "Untitled.png",
				X:       0,
				Y:       0,
				W:       24,
				H:       48,
			})
			mgr.AddComponent(e, &Position{
				Center: Center{
					X: h.X(),
					Y: h.Y() - 16,
				},
				Layer: 10,
			})
			mgr.AddComponent(e, &Obstacle{
				M: h.M,
				N: h.N,
			})
		}
	}

	board := Board{
		stride: w,
		hexes:  m,
	}
	for _, hex := range board.hexes {
		m, n := hex.M, hex.N
		// Find neighbor candidates.
		candidates := []struct{ m, n int }{
			{m, n - 2}, // N
			{m, n + 2}, // S
		}
		if n%2 == 0 {
			// then the E ones have the same M, and the W ones are -1 M
			candidates = append(candidates, []struct{ m, n int }{
				{m - 1, n - 1}, // NW
				{m - 1, n + 1}, // SW
				{m, n + 1},     // SE
				{m, n - 1},     // NE
			}...)
		} else {
			// then the E ones are +1 M, and the W ones have the same M
			candidates = append(candidates, []struct{ m, n int }{
				{m, n - 1},     // NW
				{m, n + 1},     // SW
				{m + 1, n + 1}, // SE
				{m + 1, n - 1}, // NE
			}...)
		}

		// Attach as neighbors only the ones that appear in the board.
		for _, candidate := range candidates {
			neighbor := board.Get(candidate.m, candidate.n)
			if neighbor == nil {
				continue
			}
			hex.neighbors = append(hex.neighbors, neighbor)
		}
	}

	return &board, nil
}

// Width of the Board in pixels.
func (b *Board) Width() float64 {
	return float64(b.stride * xStride)
}

// Height of the Board in pixels.
func (b *Board) Height() float64 {
	return float64(len(b.hexes) / b.stride * yStride)
}

// relative coordinates are global x,y coordinates translated to the roughMN
// rectangle coordinates. Will not return negative numbers, or values that
// exceed 16,15.
func relative(x, y int) (int, int) {
	rx, ry := x, y
	if isOddN(x) {
		ry -= 8
	}

	for {
		if rx >= 0 {
			break
		}
		rx += 17
	}
	rx = rx % 17

	for {
		if ry >= 0 {
			break
		}
		ry += 16
	}
	ry = ry % 16
	return rx, ry
}

// Get accepts M,N coordinates and returns the Hex with those coordinates if it
// exists in the board.
func (b *Board) Get(m, n int) *Hex {
	hex, ok := b.hexes[key{m, n}]
	if !ok {
		return nil
	}
	return hex
}

// At accepts world coordinates and returns the Hex there if there is one.
func (b *Board) At(x, y int) *Hex {
	rx, ry := relative(x, y)
	m, n := roughMN(x, y)

	m, n = xyToMN(rx, ry, m, n)

	return b.Get(m, n)
}

// isOddN determines whether the N coordinate will be odd or not.
func isOddN(x int) bool {
	if x < 0 {
		// -51 to -35 == true
		// -34 to -18 == false
		// -17 to -1 == true
		return (x+1)/17%2 == 0
	}
	// 0 to 16 == false
	// 17 to 33 == true
	// 34 to 50 == false
	// 51 to 67 == true
	return (x/17)%2 == 1
}

// roughMN determines which 17x16 rectangle version of a hex coordinates x,y
// fall into. Each rectangle comprises the rectangular center part of the hex,
// as well as the two triangle parts of the adjacent hexes to the top left and
// bottom left. This is only a rough guess as to the final M,N coordinates, and
// needs to be processed further before it's an accurate determination.
func roughMN(x, y int) (int, int) {

	var m, n int
	if x < 0 {
		m = x/34 - 1
	} else {
		m = x / 34
	}

	if isOddN(x) {
		if y-8 < 0 {
			// -24,-9 == -3
			//  -8, 7 == -1
			n = (y-7)/16*2 - 1
		} else {
			//  8,23 == 1
			// 24,39 == 3
			// 40,55 == 5
			n = (y-8)/16*2 + 1
		}

	} else {
		if y < 0 {
			// -64,-33 == -6
			// -32,-17 == -4
			// -16, -1 == -2
			n = (y/16 - 1) * 2
		} else {
			//   0,15  ==  0
			//  16,31  ==  2
			//  32,47  ==  4
			//  48,63  ==  6
			n = (y / 16) * 2
		}

	}

	return m, n
}

// XYToMN translates x,y coordinates relative to a 17x16 rect superimposed over
// Hex m,n to the Hex coordinates that the x,y coordinates lie inside.
func xyToMN(x, y, m, n int) (int, int) {
	if x < 0 || x > 16 || y < 0 || y > 15 {
		panic(fmt.Sprintf("x/y out of bounds %d,%d\n", x, y))
	}

	// If the x coordinate is greater or equal to 7, then it is in the
	// rectangular part of the RoughMN, so there is no special calculation
	// required.
	if x >= 7 {
		return m, n
	}

	// lookup is a map of x,y coordinates where -1 represents the hex to
	// the northwest, 0 represents this hex, and 1 represents the hex to
	// the southwest.
	lookup := map[int]map[int]int{
		0: {0: -1, 1: -1, 2: -1, 3: -1, 4: -1, 5: -1, 6: -1, 7: 0, 8: 0, 9: 1, 10: 1, 11: 1, 12: 1, 13: 1, 14: 1, 15: 1},
		1: {0: -1, 1: -1, 2: -1, 3: -1, 4: -1, 5: -1, 6: 0, 7: 0, 8: 0, 9: 0, 10: 1, 11: 1, 12: 1, 13: 1, 14: 1, 15: 1},
		2: {0: -1, 1: -1, 2: -1, 3: -1, 4: -1, 5: 0, 6: 0, 7: 0, 8: 0, 9: 0, 10: 0, 11: 1, 12: 1, 13: 1, 14: 1, 15: 1},
		3: {0: -1, 1: -1, 2: -1, 3: -1, 4: 0, 5: 0, 6: 0, 7: 0, 8: 0, 9: 0, 10: 0, 11: 0, 12: 1, 13: 1, 14: 1, 15: 1},
		4: {0: -1, 1: -1, 2: -1, 3: 0, 4: 0, 5: 0, 6: 0, 7: 0, 8: 0, 9: 0, 10: 0, 11: 0, 12: 0, 13: 1, 14: 1, 15: 1},
		5: {0: -1, 1: -1, 2: 0, 3: 0, 4: 0, 5: 0, 6: 0, 7: 0, 8: 0, 9: 0, 10: 0, 11: 0, 12: 0, 13: 0, 14: 1, 15: 1},
		6: {0: -1, 1: 0, 2: 0, 3: 0, 4: 0, 5: 0, 6: 0, 7: 0, 8: 0, 9: 0, 10: 0, 11: 0, 12: 0, 13: 0, 14: 0, 15: 1},
	}

	switch lookup[x][y] {
	case -1: // top-left triangle
		if n%2 == 0 {
			return m - 1, n - 1
		}
		return m, n - 1
	case 0: // center triangle
		return m, n
	case 1: // bottom left triangle
		if n%2 == 0 {
			return m - 1, n + 1
		}
		return m, n + 1
	default:
		panic("lookup table contained a value other than -1, 0, or 1: incoherant state not handled")
	}

}
