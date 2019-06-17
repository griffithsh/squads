package game

import (
	"bytes"
	"encoding/binary"
	"hash/fnv"
	"math"

	"github.com/griffithsh/squads/ecs"
)

// Font is renderable text.
type Font struct {
	Text string
}

// Type of this Component.
func (*Font) Type() string {
	return "Font"
}

// NewFontSystem constructs a new FontSystem.
func NewFontSystem(mgr *ecs.World) *FontSystem {
	return &FontSystem{
		mgr:    mgr,
		hashes: map[ecs.Entity][]byte{},
	}
}

// FontSystem manages the synchronization of Font Components to their composing
// child entities.
type FontSystem struct {
	mgr *ecs.World

	// map of Font entity to their hashed values of Font, Position, and Offset?
	hashes map[ecs.Entity][]byte
}

func (*FontSystem) hash(f *Font, p *Position) []byte {
	h := fnv.New128()

	h.Write([]byte(f.Text))

	fby := func(float float64) []byte {
		bits := math.Float64bits(float)
		bytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(bytes, bits)
		return bytes
	}
	h.Write(fby(p.Center.X))
	h.Write(fby(p.Center.Y))

	iby := func(i int) []byte {
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, uint64(p.Layer))

		return b
	}
	h.Write(iby(p.Layer))

	bby := func(x bool) []byte {
		b := []byte{0}
		if x {
			b[0] = 1
		}

		return b
	}
	h.Write(bby(p.Absolute))

	sum := h.Sum([]byte{})
	return sum
}

// reset by removing all Entities created as children of this Font.
func (s *FontSystem) reset(e ecs.Entity) {
	children := s.mgr.Component(e, "Children").(*ecs.Children)
	for _, e := range children.Value {
		s.mgr.DestroyEntity(e)
	}
}

// construct all child Entities necessary to compose this Font.
func (s *FontSystem) construct(parent ecs.Entity) {
	font := s.mgr.Component(parent, "Font").(*Font)
	position := s.mgr.Component(parent, "Position").(*Position)
	scale, ok := s.mgr.Component(parent, "Scale").(*Scale)
	if !ok {
		scale = &Scale{
			X: 1,
			Y: 1,
		}
	}
	s.mgr.AddComponent(parent, &ecs.Children{})
	children := s.mgr.Component(parent, "Children").(*ecs.Children)

	// px, py are the position of each rune
	var px, py float64 = 0, 0
	lineHeight := 10.0
	if scale != nil {
		lineHeight *= scale.Y
	}
	f := func(w, x, y int) {
		e := s.mgr.NewEntity()
		children.Value = append(children.Value, e)
		s.mgr.AddComponent(e, &ecs.Parent{
			Value: parent,
		})

		// letterSpace is the distance between letters.
		letterSpace := 1.0 * scale.X

		s.mgr.AddComponent(e, &Sprite{
			Texture: "font.png",
			X:       x,
			Y:       y,
			W:       w,
			H:       10,
		})
		if scale != nil {
			w = w * int(scale.X)
			s.mgr.AddComponent(e, &Scale{
				X: scale.X,
				Y: scale.Y,
			})
		}
		s.mgr.AddComponent(e, &Position{
			Center: Center{
				X: position.Center.X + px + float64(w)/2,
				Y: position.Center.Y + py + lineHeight/2,
			},
			Layer:    position.Layer,
			Absolute: position.Absolute,
		})
		px += float64(w) + letterSpace
	}
	for _, rn := range font.Text {
		switch rn {
		// whitespace
		case '\n':
			px = 0
			py += lineHeight + 2
		case '\t':
			px += lineHeight
		case ' ':
			px += 4

		// Alpha
		case 'a':
			f(5, 0, 0)
		case 'A':
			f(5, 6, 0)
		case 'b':
			f(6, 12, 0)
		case 'B':
			f(5, 19, 0)
		case 'c':
			f(5, 25, 0)
		case 'C':
			f(7, 31, 0)
		case 'd':
			f(5, 39, 0)
		case 'D':
			f(5, 45, 0)
		case 'e':
			f(5, 51, 0)
		case 'E':
			f(4, 57, 0)
		case 'f':
			f(4, 0, 10)
		case 'F':
			f(4, 5, 10)
		case 'g':
			f(5, 10, 10)
		case 'G':
			f(5, 16, 10)
		case 'h':
			f(6, 22, 10)
		case 'H':
			f(5, 29, 10)
		case 'i':
			f(2, 35, 10)
		case 'I':
			f(3, 38, 10)
		case 'j':
			f(3, 42, 10)
		case 'J':
			f(5, 46, 10)
		case 'k':
			f(5, 52, 10)
		case 'K':
			f(6, 58, 10)
		case 'l':
			f(4, 0, 20)
		case 'L':
			f(5, 5, 20)
		case 'm':
			f(7, 11, 20)
		case 'M':
			f(7, 19, 20)
		case 'n':
			f(4, 27, 20)
		case 'N':
			f(6, 32, 20)
		case 'o':
			f(5, 39, 20)
		case 'O':
			f(6, 45, 20)
		case 'p':
			f(5, 52, 20)
		case 'P':
			f(5, 58, 20)
		case 'q':
			f(5, 1, 30)
		case 'Q':
			f(6, 7, 30)
		case 'r':
			f(5, 14, 30)
		case 'R':
			f(6, 20, 30)
		case 's':
			f(4, 27, 30)
		case 'S':
			f(5, 32, 30)
		case 't':
			f(3, 38, 30)
		case 'T':
			f(7, 42, 30)
		case 'u':
			f(5, 50, 30)
		case 'U':
			f(7, 56, 30)
		case 'v':
			f(5, 0, 40)
		case 'V':
			f(7, 6, 40)
		case 'w':
			f(5, 14, 40)
		case 'W':
			f(9, 20, 40)
		case 'x':
			f(5, 30, 40)
		case 'X':
			f(8, 36, 40)
		case 'y':
			f(6, 45, 40)
		case 'Y':
			f(5, 52, 40)
		case 'z':
			f(5, 0, 50)
		case 'Z':
			f(8, 6, 50)

		// Numeric
		case '0':
			f(5, 15, 50)
		case '1':
			f(3, 21, 50)
		case '2':
			f(6, 25, 50)
		case '3':
			f(5, 32, 50)
		case '4':
			f(5, 38, 50)
		case '5':
			f(6, 44, 50)
		case '6':
			f(5, 51, 50)
		case '7':
			f(6, 57, 50)
		case '8':
			f(5, 0, 60)
		case '9':
			f(5, 6, 60)

		// Other
		case '!':
			f(1, 12, 60)
		case '.':
			f(1, 14, 60)
		case ',':
			f(2, 16, 60)
		case ';':
			f(1, 19, 60)
		case ':':
			f(1, 22, 60)
		case '-':
			f(1, 24, 60)
		case '_':
			f(1, 27, 60)

		// Default
		case '?':
			fallthrough
		default:
			f(6, 58, 110)
		}
	}
}

// Update Fonts in the World.
func (s *FontSystem) Update() {
	entities := s.mgr.Get([]string{"Font", "Position"})
	for _, e := range entities {
		font := s.mgr.Component(e, "Font").(*Font)
		pos := s.mgr.Component(e, "Position").(*Position)

		computed := s.hash(font, pos)
		current, ok := s.hashes[e]
		if !ok {
			s.hashes[e] = computed

			// Assign child entities for this font
			s.construct(e)
			continue
		}

		if bytes.Equal(current, computed) {
			continue
		}

		s.hashes[e] = computed
		// Remove all children of this entity
		s.reset(e)

		// Assign child entities for this font
		s.construct(e)
	}
}
