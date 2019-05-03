package game

// ObstacleType is an enum.
type ObstacleType int

// ObstacleTypes represent the thing that is the obstacle. These might be static
// obstacle types, like a Tree or Rock, or another Actor could be an Obstacle
// too.
const (
	ACTOR ObstacleType = iota
	CREVASSE
	TREE
)

// Obstacle is a Component that blocks a Hex.
type Obstacle struct {
	M, N int

	ObstacleType ObstacleType
}

// Type of the Component.
func (o *Obstacle) Type() string {
	return "Obstacle"
}
