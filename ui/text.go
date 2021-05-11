package ui

const (
	letterSpacing = 1
	lineSpacing   = 2 // FIXME: It's 2 for normal and 1 for small.
	spaceWidth    = 4
)

type TextSize int

const (
	TextSizeNormal TextSize = iota
	TextSizeSmall
)

type TextLayout int

const (
	TextLayoutLeft TextLayout = iota
	TextLayoutRight
	TextLayoutJustify
	TextLayoutCenter
)

type Word struct {
	Characters []Character
}

type Character struct {
	Raw                 rune
	Width, Height, X, Y int
}

// Text is a metadata enriched version of a standard string, that has not been laid out.
type Text struct {
	Value             []Word
	Size              TextSize
	BitmapFontTexture string
}

func NewText(value string, size TextSize) *Text {
	switchRune := switchRuneNormal
	if size == TextSizeSmall {
		switchRune = switchRuneSmall
	}

	result := Text{}
	for _, r := range value {
		switch r {
		case ' ':
			// new word
		case '\n':
			// new word
		default:
			w, h, x, y := switchRune(r)
			char := Character{
				Raw:    r,
				Width:  w,
				Height: h,
				X:      x,
				Y:      y,
			}
			result.Value = append(result.Value, Word{
				Characters: []Character{char},
			})
		}
	}
	return &result
}

func (t *Text) Width() int {
	if len(t.Value) == 0 {
		return 0
	}

	width := 0
	for _, w := range t.Value {
		width += w.Width()
	}
	width += (len(t.Value) - 1) * letterSpacing

	return width
}

func (w *Word) Width() int {
	if len(w.Characters) == 0 {
		return 0
	}

	width := 0
	for _, c := range w.Characters {
		width += c.Width
	}
	width += (len(w.Characters) - 1) * spaceWidth

	return width
}

func switchRuneSmall(r rune) (w, h, x, y int) {
	switch r {
	// Alpha
	case 'A', 'a':
		return 3, 5, 0, 75
	case 'B', 'b':
		return 3, 5, 4, 75
	case 'C', 'c':
		return 3, 5, 8, 75
	case 'D', 'd':
		return 3, 5, 12, 75
	case 'E', 'e':
		return 3, 5, 16, 75
	case 'F', 'f':
		return 3, 5, 20, 75
	case 'G', 'g':
		return 3, 5, 24, 75
	case 'H', 'h':
		return 3, 5, 28, 75
	case 'I', 'i':
		return 3, 5, 32, 75
	case 'J', 'j':
		return 4, 5, 36, 75
	case 'K', 'k':
		return 4, 5, 41, 75
	case 'L', 'l':
		return 3, 5, 46, 75
	case 'M', 'm':
		return 5, 5, 50, 75
	case 'N', 'n':
		return 4, 5, 56, 75
	case 'O', 'o':
		return 4, 5, 39, 70 // NB same as zero.
	case 'P', 'p':
		return 3, 5, 0, 80
	case 'Q', 'q':
		return 4, 5, 4, 80
	case 'R', 'r':
		return 4, 5, 9, 80
	case 'S', 's':
		return 3, 5, 14, 80
	case 'T', 't':
		return 3, 5, 18, 80
	case 'U', 'u':
		return 4, 5, 22, 80
	case 'V', 'v':
		return 5, 5, 27, 80
	case 'W', 'w':
		return 7, 5, 33, 80
	case 'X', 'x':
		return 5, 5, 41, 80
	case 'Y', 'y':
		return 5, 5, 47, 80
	case 'Z', 'z':
		return 5, 5, 53, 80

	// Numeric
	case '1':
		return 2, 5, 0, 70
	case '2':
		return 4, 5, 3, 70
	case '3':
		return 3, 5, 8, 70
	case '4':
		return 3, 5, 12, 70
	case '5':
		return 3, 5, 16, 70
	case '6':
		return 4, 5, 20, 70
	case '7':
		return 3, 5, 25, 70
	case '8':
		return 4, 5, 29, 70
	case '9':
		return 4, 5, 34, 70
	case '0':
		return 4, 5, 39, 70

	// Other
	case '!':
		return 1, 5, 59, 80
	case '.':
		return 1, 5, 61, 80
	case ',':
		return 2, 5, 0, 85
	case ';':
		return 2, 5, 3, 85
	case ':':
		return 1, 5, 6, 85
	case '-':
		return 2, 5, 8, 85
	case '_':
		return 3, 5, 11, 85
	case '/':
		return 3, 5, 15, 85
	case '\\':
		return 3, 5, 19, 85

	// Default
	case '?':
		fallthrough
	default:
		return 3, 5, 61, 85
	}
}

func switchRuneNormal(r rune) (w, h, x, y int) {
	switch r {
	// Alpha
	case 'a':
		return 5, 10, 0, 0
	case 'A':
		return 5, 10, 6, 0
	case 'b':
		return 6, 10, 12, 0
	case 'B':
		return 5, 10, 19, 0
	case 'c':
		return 5, 10, 25, 0
	case 'C':
		return 7, 10, 31, 0
	case 'd':
		return 5, 10, 39, 0
	case 'D':
		return 5, 10, 45, 0
	case 'e':
		return 5, 10, 51, 0
	case 'E':
		return 4, 10, 57, 0
	case 'f':
		return 4, 10, 0, 10
	case 'F':
		return 4, 10, 5, 10
	case 'g':
		return 5, 10, 10, 10
	case 'G':
		return 5, 10, 16, 10
	case 'h':
		return 6, 10, 22, 10
	case 'H':
		return 5, 10, 29, 10
	case 'i':
		return 2, 10, 35, 10
	case 'I':
		return 3, 10, 38, 10
	case 'j':
		return 3, 10, 42, 10
	case 'J':
		return 5, 10, 46, 10
	case 'k':
		return 5, 10, 52, 10
	case 'K':
		return 6, 10, 58, 10
	case 'l':
		return 4, 10, 0, 20
	case 'L':
		return 5, 10, 5, 20
	case 'm':
		return 7, 10, 11, 20
	case 'M':
		return 7, 10, 19, 20
	case 'n':
		return 4, 10, 27, 20
	case 'N':
		return 6, 10, 32, 20
	case 'o':
		return 5, 10, 39, 20
	case 'O':
		return 6, 10, 45, 20
	case 'p':
		return 5, 10, 52, 20
	case 'P':
		return 5, 10, 58, 20
	case 'q':
		return 5, 10, 1, 30
	case 'Q':
		return 6, 10, 7, 30
	case 'r':
		return 5, 10, 14, 30
	case 'R':
		return 6, 10, 20, 30
	case 's':
		return 4, 10, 27, 30
	case 'S':
		return 5, 10, 32, 30
	case 't':
		return 3, 10, 38, 30
	case 'T':
		return 7, 10, 42, 30
	case 'u':
		return 5, 10, 50, 30
	case 'U':
		return 7, 10, 56, 30
	case 'v':
		return 5, 10, 0, 40
	case 'V':
		return 7, 10, 6, 40
	case 'w':
		return 5, 10, 14, 40
	case 'W':
		return 9, 10, 20, 40
	case 'x':
		return 5, 10, 30, 40
	case 'X':
		return 8, 10, 36, 40
	case 'y':
		return 6, 10, 45, 40
	case 'Y':
		return 5, 10, 52, 40
	case 'z':
		return 5, 10, 0, 50
	case 'Z':
		return 8, 10, 6, 50

	// Numeric
	case '0':
		return 5, 10, 15, 50
	case '1':
		return 3, 10, 21, 50
	case '2':
		return 6, 10, 25, 50
	case '3':
		return 5, 10, 32, 50
	case '4':
		return 5, 10, 38, 50
	case '5':
		return 6, 10, 44, 50
	case '6':
		return 5, 10, 51, 50
	case '7':
		return 6, 10, 57, 50
	case '8':
		return 5, 10, 0, 60
	case '9':
		return 5, 10, 6, 60

	// Other
	case '!':
		return 1, 10, 12, 60
	case '.':
		return 1, 10, 14, 60
	case ',':
		return 2, 10, 16, 60
	case ';':
		return 1, 10, 19, 60
	case ':':
		return 1, 10, 22, 60
	case '-':
		return 2, 10, 24, 60
	case '_':
		return 4, 10, 27, 60
	case '/':
		return 3, 10, 32, 60
	case '\\':
		return 3, 10, 36, 60

	// Default
	case '?':
		fallthrough
	default:
		return 6, 10, 58, 110
	}
}
