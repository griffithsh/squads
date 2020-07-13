package geom

// Field of hexagons. Provides information about world coordinates.
type Field struct{}

// Get the Hexagon at Key. Returns nil if there is no Hexagon for that key.
func (f *Field) Get(Key) *Hexagon {}

// At looks for a Hexagon that is located at the world coordinates x,y. Returns
// nil if there is no Hexagon for those coordinates.
func (f *Field) At(x, y float64) *Hexagon {}

// NewField creates an empty Field configured with a HexagonSpec.
func NewField(spec HexagonSpec) *Field {}

// Load a collection of Keys into the Field.
func (f *Field) Load(keys []Key) error {}

// MedianCenter returns a world coordinate that half the loaded hexagons are
// above, half below, half to the left, and half to the right.
func (f *Field) MedianCenter() (float64, float64) {
	return 0, 0
}
