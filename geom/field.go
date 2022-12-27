package geom

import (
	"fmt"
	"math"
)

/*
 A Hexagon that has N, S, NE, SE, NW, and SW faces could be comprised of a
 triangular section (or wing), a rectangular section (or body), and then
 another triangular section going left to right.

 wing|  body  |wing

     __________
    /|        |\
   / |        | \
  /  |        |  \
  \  |        |  /
   \ |        | /
    \|________|/

The width of this hexagon would be the 2*the width of the triangular section
plus the width of the rectangular section.


M,N coordinates are sequenced like this:
    _________               _________               _________
   /         \             /         \             /         \
  /           \           /           \           /           \
 /     M,N     \_________/     M,N     \_________/     M,N     \_________
 \    -2,-2    /         \     0,-2    /         \     2,-2    /         \
  \           /           \           /           \           /           \
   \_________/    -1,-2    \_________/     1,-2    \_________/     3,-2    \
   /         \     M,N     /         \     M,N     /         \     M,N     /
  /           \           /           \           /           \           /
 /     M,N     \_________/     M,N     \_________/     M,N     \_________/
 \    -2,-1    /         \     0,-1    /         \     2,-1    /         \
  \           /           \           /           \           /           \
   \_________/    -1,-1    \_________/     1,-1    \_________/     3,-1    \
   /         \     M,N     /         \     M,N     /         \     M,N     /
  /           \           /           \           /           \           /
 /     M,N     \_________/     M,N     \_________/     M,N     \_________/
 \    -2,0     /         \     0,0     /         \     2,0     /         \
  \           /           \           /           \           /           \
   \_________/    -1,0     \_________/     1,0     \_________/     3,0     \
   /         \     M,N     /         \     M,N     /         \     M,N     /
  /           \           /           \           /           \           /
 /     M,N     \_________/     M,N     \_________/     M,N     \_________/
 \    -2,1     /         \     0,1     /         \     2,1     /         \
  \           /           \           /           \           /           \
   \_________/    -1,1     \_________/     1,1     \_________/     3,1     \
   /         \     M,N     /         \     M,N     /         \     M,N     /
  /           \           /           \           /           \           /
 /     M,N     \_________/     M,N     \_________/     M,N     \_________/
 \    -2,2     /         \     0,2     /         \     2,2     /         \
  \           /           \           /           \           /           \
   \_________/    -1,2     \_________/     1,2     \_________/     3,2     \
			 \     M,N     /         \     M,N     /         \     M,N     /
			  \           /           \           /           \           /
			   \_________/             \_________/             \_________/
This diagram shows a six by five Field of Hexagons centered on the 0,0
Hexagon. Hexagons have their flat sides oriented North-South, not East-West.
*/

// Field of hexagons. Provides information about world coordinates.
type Field struct {
	bodyWidth  int
	wingWidth  int
	totalWidth int
	halfWidth  float64
	halfHeight float64
	height     int

	hexes    map[Key]*Hex
	clickMap []int
}

// Strides here are used for dividing a hexagonal field into rectangular areas
// that are composed of the left wing of a hexagon, and the body.

func (f Field) xStride() float64 {
	return float64(f.wingWidth + f.bodyWidth)
}
func (f Field) yStride() float64 {
	return float64(f.height)
}

// Get the Hexagon at Key. Returns nil if there is no Hexagon for that key.
func (f *Field) Get(k Key) *Hex {
	hex, ok := f.hexes[k]
	if !ok {
		return nil
	}
	return hex
}

// At looks for a Hexagon that is located at the world coordinates x,y. Returns
// nil if there is no Hexagon for those coordinates.
func (f *Field) At(x, y float64) *Hex {
	k := f.Wtok(x, y)
	h, _ := f.hexes[k]
	return h
}

// clickMap produces an int slice containing a pattern like this.
// 11111111111111111111
// 11111111111111110000
// 11111111111100000000
// 11111111000000000000
// 11110000000000000000
// 00000000000000000000
// 00000000000000000000
// 22220000000000000000
// 22222222000000000000
// 22222222222200000000
// 22222222222222220000
// 22222222222222222222
// This box is the collision set for western-side wing of a hexagon. 1's
// represent the Hexagon to the NW, 2's represent the Hxagon to the SW, and 0's
// represent the current Hexagon.
func clickMap(w, h int) []int {
	rpt := func(ch int, count int) []int {
		result := make([]int, count)
		for i := 0; i < count; i++ {
			result[i] = ch
		}
		return result
	}
	incr := float64(w) / float64(h/2-1)
	numMids := h/2 - 2
	clicks := make([]int, 0, w*h)
	clicks = append(clicks, rpt(1, w)...)
	for i := 1; i <= numMids; i++ {
		through := int(math.Round(float64(i) * incr))
		clicks = append(clicks, rpt(1, w-through)...)
		clicks = append(clicks, rpt(0, through)...)
	}

	clicks = append(clicks, rpt(0, w)...)
	clicks = append(clicks, rpt(0, w)...)

	for i := 1; i <= numMids; i++ {
		through := int(math.Round(float64(i) * incr))
		clicks = append(clicks, rpt(2, through)...)
		clicks = append(clicks, rpt(0, w-through)...)
	}
	clicks = append(clicks, rpt(2, w)...)

	return clicks
}

// FlatField is a default Field where the pixel-projection distance between
// hexes is equal (or close to) in all directions.
var FlatField = NewField(5, 8, 15)
var FlatFieldDistance = FlatField.DistanceBetween(Key{}, Key{}.ToN())

// NewField creates an empty Field. The parameters configure the shape of the
// Hexagons in the Field.
func NewField(bodyWidth, wingWidth, height int) *Field {
	return &Field{
		bodyWidth:  bodyWidth,
		wingWidth:  wingWidth,
		totalWidth: wingWidth + bodyWidth + wingWidth,
		halfWidth:  float64(wingWidth+bodyWidth+wingWidth) / 2,
		halfHeight: float64(height) / 2,
		height:     height,
		clickMap:   clickMap(wingWidth, height),
	}
}

// Ktow calculates the world coordinates that are at the center of a Key.
func (f *Field) Ktow(k Key) (float64, float64) {
	x := float64(k.M * (f.wingWidth + f.bodyWidth))
	y := float64(k.N*f.height) + math.Abs(float64(k.M%2))*float64(f.halfHeight)
	return x, y
}

// Wtok finds the Key that these coordinates are inside.
func (f *Field) Wtok(x, y float64) Key {
	// translate the x,y coordinates, because the center of a hex is 0,0.
	x = x + float64(f.bodyWidth/2+f.wingWidth)
	y = y + float64(f.height)/2

	m := int(math.Floor(x / f.xStride()))
	if m%2 != 0 {
		y = y - float64(f.height)/2
	}
	n := int(math.Floor(y / f.yStride()))

	// abs returns the absolute value of n.
	abs := func(n int) int {
		if n < 0 {
			return -n
		}
		return n
	}
	lx, ly := abs(int(x)%(f.wingWidth+f.bodyWidth)), abs(int(y)%f.height)
	triangle := f.local(lx, ly)

	switch triangle {
	case 1:
		// NW of m,n
		return Key{m, n}.ToNW()
	case 2:
		// SW of m,n
		return Key{m, n}.ToSW()
	}

	return Key{m, n}
}

// local takes integral coordinates and translates them to either 1 (NW), 2
// (SW), or 0 (the central Hex).
func (f *Field) local(x, y int) int {
	if x >= f.wingWidth {
		return 0
	}
	if x < 0 || y < 0 || y >= f.height {
		fmt.Printf("local(%d,%d): incoherant\n", x, y)
		return 0
	}
	return f.clickMap[x+y*f.wingWidth]
}

// Load a collection of Keys into the Field.
func (f *Field) Load(keys []Key) error {
	f.hexes = map[Key]*Hex{}

	for _, k := range keys {
		x, y := f.Ktow(k)
		current := Hex{
			x:   x,
			y:   y,
			key: k,
		}

		// TODO: any use for this??
		// for neighbor := range k.Neighbors() {
		// 	if h, ok := f.hexes[neighbor]; ok {
		// 		// link current and h
		// 	}

		// }

		f.hexes[k] = &current
	}
	return nil
}

// Center returns a world coordinate that half the loaded hexagons are
// above, half below, half to the left, and half to the right.
func (f *Field) Center() (float64, float64) {
	sumX, sumY := 0.0, 0.0
	total := float64(len(f.Hexes()))
	for _, h := range f.Hexes() {
		x, y := h.Center()
		sumX += x
		sumY += y
	}
	return sumX / total, sumY / total
}

// Hexes generates a slice of the Hexes of the Field.
func (f *Field) Hexes() []*Hex {
	result := make([]*Hex, len(f.hexes))
	i := 0
	for _, v := range f.hexes {
		result[i] = v
		i++
	}
	return result
}

// DistanceBetween calculates the distance between two keys in the field.
func (f *Field) DistanceBetween(a, b Key) float64 {
	ax, ay := f.Ktow(a)
	bx, by := f.Ktow(b)
	return math.Hypot(math.Abs(ax-bx), math.Abs(ay-by))
}
