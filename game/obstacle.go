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

// ContextualObstacle captures how much of an obstacle this is to the navigator.
// A bird can fly right over a tree, a snake is not impeded by a swamp. A horse
// runs fastest when the ground is level and clear. The Cost multiplies the
// normal traversal time. A Cost of 2 implies that taking this path is twice as
// long as it normally would be.
type ContextualObstacle struct {
	Obstacle

	Cost float64
}
