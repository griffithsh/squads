package geom

import (
	"fmt"
)

/*
 A Hexagon that has N, S, NE, SE, NW, and SW faces could be comprised of a
 triangular section and a rectangular section, and then another triangular section
 going left to right.
    _________
   /|       |\
  / |       | \
 /  |       |  \
 \  |       |  /
  \ |       | /
   \|_______|/

The width of this hexagon would be the 2*the width of the triangular section
plus the width of the rectangular section.

Stride here is used for dividing a hexagonal field into rectangular areas that
are composed of the first triangular area of a hexagon, and the rectangular
area.
*/
const (
	hexTriWidth    = 7
	hexSquareWidth = 10
	hexWidth       = hexTriWidth + hexSquareWidth + hexTriWidth
	hexHeight      = 16
	xStride        = (hexTriWidth + hexSquareWidth) * 2
	yStride        = hexHeight / 2
)

// Key is a way of referencing a Hex in a Field.
type Key struct {
	M, N int
}

// Field to play out encounters on. A collection of Hexes.
type Field struct {
	stride int // how many hexes are in a row
	hexes  map[Key]*Hex
}

// NewField creates a new game field.
func NewField(w, h int) (*Field, error) {
	hexes := make(map[Key]*Hex, w*h)

	for i := 0; i < w*h; i++ {
		h := Hex{
			M: i % w,
			N: i / w,
		}
		hexes[Key{i % w, i / w}] = &h
	}

	field := Field{
		stride: w,
		hexes:  hexes,
	}
	field.calcNeighbors()

	return &field, nil
}

func (f *Field) calcNeighbors() {
	for _, hex := range f.hexes {
		hex.neighbors = []*Hex{}
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

		// Attach as neighbors only the ones that appear in the field.
		for _, candidate := range candidates {
			neighbor := f.Get(candidate.m, candidate.n)
			if neighbor == nil {
				continue
			}
			hex.neighbors = append(hex.neighbors, neighbor)
		}
	}
}

// Width of the Field in pixels.
func (f *Field) Width() float64 {
	return float64(f.stride * xStride)
}

// Height of the Field in pixels.
func (f *Field) Height() float64 {
	return float64(len(f.hexes) / f.stride * yStride)
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
// exists in the field.
func (f *Field) Get(m, n int) *Hex {
	hex, ok := f.hexes[Key{m, n}]
	if !ok {
		return nil
	}
	return hex
}

// At accepts world coordinates and returns the Hex there if there is one.
func (f *Field) At(x, y int) *Hex {
	rx, ry := relative(x, y)
	m, n := roughMN(x, y)

	m, n = xyToMN(rx, ry, m, n)

	return f.Get(m, n)
}

// Dimensions returns the maximum extent of the field in M and N dimensions.
func (f *Field) Dimensions() (M, N int) {
	return f.stride, len(f.hexes) / f.stride
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
