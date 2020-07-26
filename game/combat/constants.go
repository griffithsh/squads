package combat

// define the size of hexagons used by combat.
const (
	hexagonWingWidth int = 19
	hexagonBodyWidth int = 34
	hexagonTileWidth int = hexagonWingWidth + hexagonBodyWidth + hexagonWingWidth
	hexagonHeight    int = 40
)

// define the z-ordering render layers used by combat
const (
	terrainLayer     = 10
	cursorLayer      = 90
	participantLayer = 100
)
