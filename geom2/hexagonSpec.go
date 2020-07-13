package geom

// HexagonSpec defines the geometric layout of a hexagon.
type HexagonSpec struct {
	BodyWidth int
	WingWidth int
	Height    int
	// TotalWidth = WingWidth + BodyWidth + WingWidth
	// XStride = (WingWidth + BodyWidth) * 2
	// YStride = Height / 2
}

// TotalWidth of the Hexagons.
func (s *HexagonSpec) TotalWidth() int {
	return s.WingWidth + s.BodyWidth + s.WingWidth
}

// XStride is the difference between adjacent hexagons in the X dimension.
func (s *HexagonSpec) XStride() int {}

// YStride is the difference between adjacent hexagons in the Y dimension.
func (s *HexagonSpec) YStride() int {}

// graph.Searcher(costProvider, heuristic, connectionProvider).Search(a,b Vertex)
