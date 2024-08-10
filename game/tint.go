package game

// Tint sets the color of the entity when it is rendered.
type Tint struct {
	R, G, B uint8
}

var TintPseudoBlack = Tint{R: 0x14, G: 0x1b, B: 0x27}
var TintPseudoWhite = Tint{R: 0xda, G: 0xe1, B: 0xe5}

// Type of this Component.
func (*Tint) Type() string {
	return "Tint"
}
